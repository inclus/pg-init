
## pg-init

![](https://github.com/inclus/pg-init/workflows/Go/badge.svg)

[![Go Report Card](https://goreportcard.com/badge/github.com/inclus/pg-init)](https://goreportcard.com/report/github.com/inclus/pg-init)


Simple application that tries to establish a connection with a postgresql database and then configures the postgis extension.

It was built to be run as an init container to make sure the application starts when the database is ready.

You can download the docker image [here](https://github.com/inclus/pg-init/packages/31983). 
