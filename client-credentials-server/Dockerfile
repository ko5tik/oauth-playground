FROM golang:latest
WORKDIR /src
COPY . .
RUN go build -o /client-credentials-server .
ENTRYPOINT ["./client-credentials-server"]
