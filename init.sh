# run source ./init.sh

echo "Initialisation of env variables"
export HTTP_PORT=8080
echo "export HTTP_PORT executed"

export GRPC_PORT=8081
echo "export GRPC_PORT executed"

export GRAPHQL_PORT=8082
echo "export GRAPHQL_PORT executed"

export POSTGRESQL_URL=postgres://serhii:serhii@localhost:5432/api
echo "export POSTGRESQL_URL executed"


