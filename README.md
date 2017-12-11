![https://github.com/ewikum/go-postgres-redis-s3-dev](https://cdn-images-1.medium.com/max/600/1*bnyJ9a-DxAKV-qJXKDAZkQ.png)

# Golang docker Env (Minio, Postgre, Redis) 
In this setup we will have four docker containers as follows.

* 1 Go container : This is the main container which will host our Go app.
* 1 Minio container : Object storage server with Amazon S3 compatible API.
* 1 Postgres container : Postgres DB server.
* 1 Redis container : Key-value store that functions as a data structure server.


## Prerequisites

You have to have docker and docker composer installed on your local machine.
If you haven't got Docker Compose installed, 
[follow these installation instructions](https://docs.docker.com/compose/install/).

## Download, Build, Run, Code and Test 

### Download and build

On your local machine, clone this repo and spin up our four containers issuing following 3 commands:

```
git clone https://github.com/DaisukeHirata/go-test-rest-api-with-postgres.git
cd go-test-rest-api-with-postgres
docker-compose build
```
It will take some time to complete these three commands as it has to pull all the images needed at first run. 

### Run

#### Foreground mode

```
docker-compose up
```

By Issuing above command, we can have our containers up and running in the foreground. Hence we can see all the logs in our terminal window 

We can use `CTRL + C` to stop containers and get the access back to our terminal. 

#### Background mode

```
docker-compose up -d
```

Using the flag -d, we can have our containers running in the background. However we cannot see logs in our terminal window 


In either mode we can now access our go application on [http://localhost:8080/](http://localhost:8080/)

#### Stop and remove

Isuue the following command to stop and remove all running containers, networks, volumes, and images created by `docker-compose up` 

```
docker-compose up -d
```

### Code

When we start up our containers, our go app is builded from `go-app` directory 
and starts to serve. 

`go-app` is the directory where all of our go source files live.

Each time a go souce file changes, Our app will incrementally be rebuilded. So we dont need to worry about buidling our app again and again. 

We can just refresh our browser to see the changes we make.


Additinally, We can use following command to compile and run Go program.

```
docker-compose run go go run abc.go
```
abc.go is the go file path realative to the go-app directory.

### Test

The following command will run all the go test files in our go-app directory. 

```
docker-compose run go go test ./...
```

## Container Names

* 1 Go container : go-cont
* 1 Minio container : s3-cont
* 1 Postgres container : db-cont
* 1 Redis container : redis-cont

## Usefull commands

### Access shell of a container

Issue the Following command when the container is running.

replace go-cont with any of the container name above.

```
docker exec -it go-cont sh
```

## Files Structure

### `.S3`

This directory is mounted to the data directory of the minio container so that our buckets and objects wont get lost as the minio container stops.

ie : This directory is there to keep all the data for minio service and hence make our object storage persistent.


### `.DB`

Just like the .S3 directory , This directory is mounted to the data directory of the postgres container to make our postgres database persistent.

### `Dockerfile`

This file is responsible to create a new docker image by copying the Go source in `go-app` directory, build and run it in the
container. 

### `docker-compose.yml`

This file is usefull for defining and running multi-container Docker applications like ours. We use this file to configure our application’s services. Then, by `docker-compose build` command, we can create and start all the services. 


## Inside go-app/

This is the root directory where our go application is served from.
So we have to place our go source files inside this directory.

### `main.go`

This runs a http server that listens on port `:8080` inside go container, so that we can access it via web browser on `http://localhost:8080`. Since this is the starting point of our application, this file is reponsible to route a url to the correct httphandler function.

### `redis.go`

This file inclues a sample code to undestand how to deal with redis storage in this envioronment. 
[Redigo](https://github.com/garyburd/redigo/) is used in this sample code, 
  which is a Go client for the Redis database.


### `db.go`

This file inclues a sample code to undestand how to connect to postgres database in this envioronment. 

This file is responsible for the followings. 

* `http://localhost:8080/db/`: lists all the records in a table called dummytable.
* `http://localhost:8080/db/add/{some text}`: inserts a new records to a table called dummytable with the value of {some text}.

Check whats happening meanwhile inside `.DB` directory. :)
 
### `s3.go`

This file inclues a sample code to undestand how to deal with minio object storage in this envioronment. 

This file is responsible for the followings. 

* `http://localhost:8080/s3/`: lists all the objects in a bucket called testbucket.
* `http://localhost:8080/s3/triggeraput/`: uploads the `s3_upload_test_file.txt` to the testbucket.

If we check the `.S3` directory, We can see our objects are placed inside.
Additionally we can access the minio service from the mino container directly via `http://localhost:9000`

Read [official minio guide](https://docs.minio.io/docs/golang-client-quickstart-guide
) to leran more. 


### Happy coding...!