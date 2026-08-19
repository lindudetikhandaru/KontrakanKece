FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY projectBiru/go.mod projectBiru/go.sum ./
RUN go mod download || true

COPY projectBiru/ .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/*.html ./
COPY --from=builder /app/*.css ./
COPY --from=builder /app/*.png ./

EXPOSE 8080

CMD ["./main"]
