package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	maxBatchSize   = 1000
	maxBufferSize  = 10_000
	maxRetries     = 3
	initialBackoff = 500 * time.Millisecond
)

type lokiPayload struct {
	Streams []struct {
		Stream map[string]string `json:"stream"`
		Values [][2]string       `json:"values"`
	} `json:"streams"`
}

type LokiWriter struct {
	client *http.Client
	url    string
	labels map[string]string

	mu     sync.Mutex
	buf    [][2]string
	ticker *time.Ticker
	stop   chan struct{}
	closed bool
}

func NewLokiWriter(
	url string,
	labels map[string]string,
	flushInterval time.Duration,
) *LokiWriter {
	lw := &LokiWriter{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		url:    url,
		labels: labels,
		ticker: time.NewTicker(flushInterval),
		stop:   make(chan struct{}),
	}

	go lw.flushLoop()
	return lw
}

func (l *LokiWriter) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return len(p), nil
	}

	if len(l.buf) >= maxBufferSize {
		l.buf = l.buf[1:]
	}

	l.buf = append(l.buf, [2]string{
		strconv.FormatInt(time.Now().UnixNano(), 10),
		string(bytes.TrimSpace(p)),
	})

	return len(p), nil
}

func (l *LokiWriter) flushLoop() {
	for {
		select {
		case <-l.ticker.C:
			l.flush()
		case <-l.stop:
			l.flush()
			return
		}
	}
}

func (l *LokiWriter) flush() {
	for {
		batch := l.nextBatch()
		if len(batch) == 0 {
			return
		}

		if err := l.push(batch); err != nil {
			l.requeue(batch)
			return
		}
	}
}

func (l *LokiWriter) nextBatch() [][2]string {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.buf) == 0 {
		return nil
	}

	n := maxBatchSize
	if len(l.buf) < n {
		n = len(l.buf)
	}

	batch := l.buf[:n]
	l.buf = l.buf[n:]
	return batch
}

func (l *LokiWriter) requeue(batch [][2]string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if len(l.buf)+len(batch) > maxBufferSize {
		batch = batch[:maxBufferSize-len(l.buf)]
	}

	l.buf = append(batch, l.buf...)
}

func (l *LokiWriter) push(values [][2]string) error {
	payload := lokiPayload{
		Streams: []struct {
			Stream map[string]string `json:"stream"`
			Values [][2]string       `json:"values"`
		}{
			{
				Stream: l.labels,
				Values: values,
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	backoff := initialBackoff

	for i := 0; i < maxRetries; i++ {
		req, err := http.NewRequest(
			http.MethodPost,
			l.url+"/loki/api/v1/push",
			bytes.NewBuffer(data),
		)
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := l.client.Do(req)
		if err == nil && resp != nil {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()

			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				return nil
			}
		}

		time.Sleep(backoff)
		backoff *= 2
	}

	return os.ErrDeadlineExceeded
}

func (l *LokiWriter) Close() {
	l.mu.Lock()
	if l.closed {
		l.mu.Unlock()
		return
	}
	l.closed = true
	l.mu.Unlock()

	close(l.stop)
	l.ticker.Stop()
}
