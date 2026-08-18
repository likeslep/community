# syntax=docker/dockerfile:1

# ---- 构建阶段 ----
FROM golang:1.26-alpine AS builder
ARG SERVICE
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/service ./cmd/${SERVICE}

# ---- 运行阶段 ----
FROM alpine:3.20
RUN adduser -D -u 10001 app
USER app
COPY --from=builder /app/service /service
EXPOSE 8080
ENTRYPOINT ["/service"]
