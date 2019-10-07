FROM golang:1.13 AS build-env
ADD . /src
RUN cd /src && go build -o pg-init

FROM alpine
WORKDIR /app
COPY --from=build-env /src/pg-init /app/
ENTRYPOINT ./pg-init