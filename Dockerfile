FROM node:22-alpine AS frontend_builder

WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
RUN npm run build

FROM nginx:alpine
COPY --from=frontend_builder /app/dist /usr/share/nginx/html
EXPOSE 8083
CMD ["nginx", "-g", "daemon off;"]

FROM golang:tip-alpine3.22 AS backend_builder
RUN apk add --no-cache git ca-certificates
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o grubzo ./cmd/api

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=backend_builder /app/grubzo .
EXPOSE 8082
CMD ["./grubzo serve"]
