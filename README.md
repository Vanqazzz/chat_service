# chat_service

gRPC-service authentication and messaging microservice. 

## Features

- User registration
- User login with JWT token
- Token validating with `app_id` and `app_secret`
- PostgreSQL storage
- Message sending and receiving with a Kafka
- unit-tests
- Database migrations

## Stack
- Go 1.24+
- PostgreSQL
- Kafka
- gRPC 
- Docker / Docker Compose


## Quick start

### Requirements 

- Go 1.24+
- PostgreSQL
- Kafka
- protoc
- Docker

### Instalattion

Clone repository in your folder
```cmd
git clone https://github.com/Vanqazzz/chat_service 
cd chat_service
```

## Build container
```
docker-compose up --build
```
