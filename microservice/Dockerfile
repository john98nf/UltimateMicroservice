FROM golang:1.20-bookworm as builder

LABEL maintainer="Giannis Fakinos"

WORKDIR /ultimateMicroservice

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN cd cmd/app && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /ultimateMicroservice/cmd/app/app .

EXPOSE 8080

CMD ["./app"]
