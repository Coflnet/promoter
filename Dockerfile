FROM golang:1.21 as builder
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN go build -o ./bin/promoter .


FROM gcr.io/distroless/static
COPY --from=builder /app/bin/promoter /bin/promoter

RUN apt update -y && apt install -y git

CMD ["/bin/promoter"]
