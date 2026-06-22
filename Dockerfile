FROM golang:1.26.3-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/cmd ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder /app/cmd .

RUN mkdir /app/logs

RUN chown -R appuser:appuser /app

USER appuser

EXPOSE 8000

CMD ["./cmd", "run"]
