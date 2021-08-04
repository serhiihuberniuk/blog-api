FROM golang

WORKDIR /blog-api

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o main

ENV API_HTTP_PORT=8080
ENV API_GRPC_PORT=8081
ENV API_GRAPHQL_PORT=8082
ENV API_POSTGRESQL_URL=postgres://serhii:serhii@localhost:5432/api

CMD ["./main"]



