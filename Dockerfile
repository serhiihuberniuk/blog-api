FROM golang:1.16-buster AS build

WORKDIR /app

COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o blog-api

FROM scratch

COPY --from=build /app/blog-api /
COPY --from=build /app/init.sql /

ENTRYPOINT ["/blog-api"]