# GO clean code demo

simple API server to upload file designed with clean architecture

## To run

- Start external service, mongoDB and localstack 
```shell
$ cd docker
$ docker-compose up -d
```

- Start API server
```shell
$ go run main.go
```

- POST data to API server
> I used restcli for the test
```shell
$ cd http-request
$ java -jar your.path.to.restcli.jar asset.http
```