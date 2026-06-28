FROM node:22-alpine AS frontend_builder

WORKDIR /app

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./

RUN npm run build

FROM nginx:alpine AS frontend

COPY --from=frontend_builder /app/dist/ /usr/share/nginx/html/
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]

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

# docker build -t rohana001/grubzo-frontend:v0.1.1 --target frontend .
# docker build -t rohana001/grubzo-backend:v0.1.1 --target backend .
