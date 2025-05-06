# base image
FROM golang:1.24-bookworm AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o chirpy .


# application image
FROM scratch

COPY --from=builder /app/chirpy .

CMD ["./chirpy"]
