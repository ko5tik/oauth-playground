FROM golang:1.16.0-alpine3.13 AS build
WORKDIR /src
COPY . .
RUN go build -o /auth-server .
FROM scratch AS bin
COPY --from=build /auth-server /