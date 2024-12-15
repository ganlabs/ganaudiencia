FROM golang:1.23-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o ganaudiencia .

FROM scratch as runner

WORKDIR /app

COPY --from=builder /app/ganaudiencia .

CMD ["./ganaudiencia"]