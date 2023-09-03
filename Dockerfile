FROM golang:1.21 as builder
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .
RUN go build -o ./bin/promoter .


FROM alpine:3.18
COPY --from=builder /app/bin/promoter /bin/promoter

RUN apk add git

CMD ["/bin/promoter"]
