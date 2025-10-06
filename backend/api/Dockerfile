# --- ビルドステージ ---
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

# --- 実行ステージ ---
FROM alpine

#FIXME シークレットを.envなどに移動
ENV POSTGRES_DSN=
ENV JWT_SECRET=
ENV JANUS_WS=
ENV API_PORT=60001
ENV GOOGLE_CRED_PATH=

WORKDIR /app
COPY --from=builder /app/server /app/server

EXPOSE 60001

# 実行
ENTRYPOINT ["/app/server"]
