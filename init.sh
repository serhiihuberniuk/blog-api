# run source ./init.sh

echo "Initialisation of env variables"
export API_HTTP_PORT=8080
echo "export HTTP_PORT executed"

export API_GRPC_PORT=8081
echo "export GRPC_PORT executed"

export API_GRAPHQL_PORT=8082
echo "export GRAPHQL_PORT executed"

export API_HEALTHCHECK_PORT=8083
echo "export HEALTHCHECK_PORT executed"

export API_POSTGRESQL_URL=postgres://serhii:serhii@localhost:5432/api?sslmode=disable
echo "export POSTGRESQL_URL executed"

export API_POSTGRES_MIGRATIONS_PATH=file://migrations
echo "export POSTGRES_INIT_FILE executed"

export API_POSTGRES_MIGRATIONS_VERSION=3
echo "export API_POSTGRES_DATABASE_VERSION executed"

export API_PRIVATE_KEY_FILE=app.rsa
echo "export API_PRIVATE_KEY_PATH executed"

export API_REDIS_ADDRESS=localhost:6379
echo "export API_REDIS_ADDRESS executed"

#API_HTTP_PORT=8080;API_GRPC_PORT=8081;API_GRAPHQL_PORT=8082;API_HEALTHCHECK_PORT=8083;API_POSTGRESQL_URL=postgres://serhii:serhii@localhost:5432/api?sslmode=disable;API_POSTGRES_MIGRATIONS_PATH=file://migrations;API_POSTGRES_MIGRATIONS_VERSION=3;API_PRIVATE_KEY_FILE=app.rsa;API_REDIS_ADDRESS=localhost:6379;


