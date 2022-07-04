# go-cqrs-microservices


### Golang CQRS Kafka gRPC Postgresql MongoDB Redis microservices with clean architecture example 👋

#### 👨‍💻 Full list what has been used:
[Kafka](https://github.com/segmentio/kafka-go) as messages broker<br/>
[gRPC](https://github.com/grpc/grpc-go) Go implementation of gRPC<br/>
[PostgreSQL](https://github.com/jackc/pgx) as database<br/>
[Prometheus](https://prometheus.io/) monitoring and alerting<br/>
[Grafana](https://grafana.com/) for to compose observability dashboards with everything from Prometheus<br/>
[MongoDB](https://github.com/mongodb/mongo-go-driver) Web and API based SMTP testing<br/>
[Redis](https://github.com/go-redis/redis) Type-safe Redis client for Golang<br/>
[swag](https://github.com/swaggo/swag) Swagger for Go<br/>
[Beego](https://github.com/beego/beego) framework fro Go<br/>

### CQRS Architecture
![golang clean architecture](https://github.com/radyatamaa/loyalti-go-echo/blob/dev/CQRS-architecture-2.png)

### Clean Architecture
This project has  4 Domain layer :

 * Models Layer
 * Repository Layer
 * Usecase Layer  
 * Delivery Layer

#### The diagram:

![golang clean architecture](https://github.com/bxcodec/go-clean-arch/raw/master/clean-arch.png)

The explanation about this project's structure  can read from this medium's post : https://medium.com/@imantumorang/golang-clean-archithecture-efd6d7c43047

### How To Run This Project

```bash
#move to directory
cd $GOPATH/src/github.com/radyatamaa

# Clone into YOUR $GOPATH/src
git clone https://github.com/radyatamaa/go-cqrs-microservices.git

#move to project
cd go-cqrs-microservices

# Run app writer service
go run write_service/cmd/main.go

# Run app reader service
go run reader_service/cmd/main.go

# Run app api gateway service
go run api_gateway_service/cmd/main.go

```

Or with `docker-compose`

```bash
#move to directory
cd $GOPATH/src/github.com/radyatamaa

# Clone into YOUR $GOPATH/src
git clone https://github.com/radyatamaa/go-cqrs-microservices.git

#move to project
cd go-cqrs-microservices

# Run the application
make run  OR  docker compose -f "docker-compose.yml" up -d --build


# Open at browser this url
http://localhost:8082/swagger/index.html

```

### Prometheus UI:

http://localhost:9090

### Grafana UI:

http://localhost:3000

### Swagger UI:

http://localhost:8082/swagger/index.html
