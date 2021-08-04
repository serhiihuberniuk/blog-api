# run source ./init.sh

echo "Initialisation of env variables"
export API_HTTP_PORT=8080
echo "export HTTP_PORT executed"

export API_GRPC_PORT=8081
echo "export GRPC_PORT executed"

export API_GRAPHQL_PORT=8082
echo "export GRAPHQL_PORT executed"

export API_POSTGRESQL_URL=postgres://serhii:serhii@localhost:5432/api
echo "export POSTGRESQL_URL executed"

#API_HTTP_PORT=8080;API_GRPC_PORT=8081;API_GRAPHQL_PORT=8082;API_POSTGRESQL_URL=postgres://serhii:serhii@localhost:5432/api


