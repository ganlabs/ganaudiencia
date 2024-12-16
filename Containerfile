FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY . .

RUN go build -o ganaudiencia .

FROM scratch AS runner

WORKDIR /app

COPY --from=builder /app/ganaudiencia .

CMD ["./ganaudiencia"]