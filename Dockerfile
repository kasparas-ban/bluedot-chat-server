FROM golang:1.19.2-alpine3.16 as build
WORKDIR /server
COPY . /server
RUN go build -o /server-app

FROM alpine
COPY --from=build ./server-app ./
COPY --from=build ./server/.env ./
EXPOSE ${SERVER_PORT}
ENTRYPOINT ["/server-app"]