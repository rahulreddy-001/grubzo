FROM golang:1.25-alpine AS backend_builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -o grubzo ./cmd/api

FROM alpine:3.23 AS backend

WORKDIR /app

RUN apk add --no-cache libstdc++ ca-certificates

COPY --from=backend_builder /app/grubzo ./grubzo

VOLUME ["/app"]

EXPOSE 80

CMD ["./grubzo", "serve"]
