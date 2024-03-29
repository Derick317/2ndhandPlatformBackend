# Start with the Back-End

@author: Deming Chen (陈德铭)

## PostgreSQL in Docker

Suppose you are under the `Backend` directory. To start with, initialize a PostgreSQL container in docker:

```bash
 docker-compose up
```

This command reads `docker-compose.yml` and does a handful of things, such as pulling the docker image (if it does not exist), creating a container and running it. For details, see [this guide](./docker.md). It occupies a terminal window, so you could open a new one. Also, you could just press `ctrl+C` to stop the container. Pressing `ctrl+C` has the same effect as click `Stop` in docker desktop. ![](./images/docker_desktop_stop_container.png)

Afterwards, you are able to restart it by click `Start` in docker desktop or run this command again.![image](./images/docker_desktop_start_container.png)

## Connecting to PostgreSQL Database using Go

In the starting code, we connect to PostgreSQL database using Go. We need to install the `pq` package using go package manager:

```bash
# go get github.com/lib/pq
go get -u gorm.io/gorm
go get -u gorm.io/driver/postgres
```

Here, we have initialize a module called `secondHand`. Now, change your directory to it. To test connection, run

```bash
go run main.go
```

You are expected to see

```
started-service
Connected to PostgreSQL database successfully!
Initialized PostgreSQL database successfully!
```

## Constants

`src/secondHand/constants/constants.go` saves a lot of configurations. You should add your own.

### `GCS_BUCKE`

`GCS_BUCKE` is just the name of the bucket in google cloud storage:![](./images/gcs_bucket.png)

### `GCS_CREDENTIALS_FILE_PATH`

`GCS_CREDENTIALS_FILE_PATH` is the path where you locate the `json` credential file. You can download it by following [this guide](https://blog.csdn.net/leolee_0606/article/details/121654458).
