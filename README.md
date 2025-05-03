# mileta-invoice-backend


# About

This repo is a backend for my [mileta-invoice-creator](https://github.com/milemik/mileta-invoice-creator) app. And backend will be written in go to.


# Local development

## Run server
You can run server using ```go run```:
```shell
go run .
```

## Run mongoDB
You can run mongoDB using docker [mongoDB official image](https://www.mongodb.com/docs/manual/tutorial/install-mongodb-community-with-docker/)

```shell
docker run --name mongodb-test -p 27017:27017 -v ./.data:/data/db -d mongodb/mongodb-community-server:latest
```