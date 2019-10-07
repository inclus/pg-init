FROM golang:1.13 AS build-env
ADD . /src
RUN cd /src && CGO_ENABLED=0 go build -o pg-init

FROM alpine
WORKDIR /app
COPY --from=build-env /src/pg-init /app/
ENTRYPOINT ./pg-init