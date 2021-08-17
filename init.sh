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

export API_POSTGRESQL_URL=postgres://serhii:serhii@localhost:5432/api
echo "export POSTGRESQL_URL executed"

export API_POSTGRES_INIT_FILE=init.sql
echo "export POSTGRES_INIT_FILE executed"

export API_PRIVATE_KEY_FILE=app.rsa
echo "export API_PRIVATE_KEY_PATH executed"

#API_HTTP_PORT=8080;API_GRPC_PORT=8081;API_GRAPHQL_PORT=8082;API_HEALTHCHECK_PORT=8083;API_POSTGRESQL_URL=postgres://serhii:serhii@localhost:5432/api;API_POSTGRES_INIT_FILE=init.sql;API_PRIVATE_KEY_FILE=app.rsa;


