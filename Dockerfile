FROM golang:1.25-alpine AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY *.go ./

# 编译
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 最终镜像
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# 从 builder 复制编译好的二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY config.yaml .

# 暴露端口
EXPOSE 52143

# 设置时区（可选）
ENV TZ=Asia/Shanghai

# 运行
CMD ["./main"]