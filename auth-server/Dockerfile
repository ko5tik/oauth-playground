FROM golang:latest
WORKDIR /src
COPY . .
RUN go build -o /auth-server .
ENTRYPOINT ["./auth-server"]
