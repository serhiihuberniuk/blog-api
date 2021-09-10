FROM golang:1.16-buster AS build

WORKDIR /app

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog-api

#FROM scratch
FROM busybox

COPY --from=build /app/blog-api /
COPY --from=build /app/migrations / migrations/

HEALTHCHECK --interval=5s --timeout=5s --retries=3 \
    CMD wget -nv -t1 --spider 'http://localhost:8083/health' || exit 1

ENTRYPOINT ["/blog-api"]