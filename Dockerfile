FROM golang:1.20 as builder
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN go build -o ./bin/promoter .


FROM debian:bookworm-slim
COPY --from=builder /app/bin/promoter /bin/promoter

RUN apt update -y && apt install -y git

CMD ["/bin/promoter"]
