FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .

RUN go build -o /app/jelly-metrics

FROM scratch AS live

COPY --from=builder /app/jelly-metrics /jelly-metrics

EXPOSE 8097
ENTRYPOINT ["./jelly-metrics"]
