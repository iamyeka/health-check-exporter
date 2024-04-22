# 第一阶段：构建 Go 应用程序
FROM golang:1.14 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# 第二阶段：创建运行镜像
FROM alpine:3.12
WORKDIR /root/
COPY --from=builder /app/app .
CMD ["./app"]
