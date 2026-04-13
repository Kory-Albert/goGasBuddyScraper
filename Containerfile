FROM golang:1.24-alpine AS builder
RUN apk add --no-cache tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/main .

FROM scratch
COPY --from=builder /app/main /main
ENTRYPOINT ["/main"]