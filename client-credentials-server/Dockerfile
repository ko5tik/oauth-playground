FROM golang:alpine
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o client-credentials-server .
CMD ["/app/client-credentials-server"]
