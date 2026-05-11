package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jinzhu/copier"
)

func LoadJSONFromFile[T any](filePath string) (*T, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var value T

	if err := json.NewDecoder(file).Decode(&value); err != nil {
		return nil, err
	}

	return &value, nil
}

func AsType[T any](v interface{}) (T, error) {
	var out T

	if v == nil {
		return out, fmt.Errorf("nil cannot be converted to %T", out)
	}

	if val, ok := v.(T); ok {
		return val, nil
	}

	if ptr, ok := v.(*T); ok {
		if ptr == nil {
			return out, fmt.Errorf("nil pointer cannot be converted to %T", out)
		}
		return *ptr, nil
	}

	if m, ok := v.(map[string]any); ok {
		if _, isMap := any(out).(map[string]any); isMap {
			return any(m).(T), nil
		}

		jsonBytes, err := json.Marshal(m)
		if err != nil {
			return out, err
		}

		decoder := json.NewDecoder(bytes.NewReader(jsonBytes))
		decoder.UseNumber()

		if err := decoder.Decode(&out); err != nil {
			return out, err
		}
		return out, nil
	}

	if b, ok := v.([]byte); ok {
		if err := json.Unmarshal(b, &out); err != nil {
			return out, err
		}
		return out, nil
	}

	return out, fmt.Errorf("cannot convert %T to %T", v, out)
}

func Map[Dst any, Src any](dst *Dst, src *Src) error {
	return copier.Copy(dst, src)
}
