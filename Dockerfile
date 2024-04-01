FROM golang:1.22.1-alpine3.19 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# CGO_ENABLED=0 確保不依賴任何 C 函式庫
# GOOS=linux 即使在非 Linux 環境生成的二進位文件也是 for Linux系统
# -a 強制重新建構不使用暫存
# -installsuffix cgo 和 -a 搭配使用 給安装的包添加一个後綴確保构是完全乾淨重新建立的二進位文件
# cgo 是後綴名，用意在於表達這是 CGO 禁止的情況下編譯的
RUN GOOS=linux go build -o main ./cmd/main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .env
COPY build/start.sh .
COPY scripts/db/migration ./migration

EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]
CMD ["/app/main"]
