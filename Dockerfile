# syntax=docker/dockerfile:1

FROM golang:1.19
WORKDIR /welsh/
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build ./cmd/api
RUN CGO_ENABLED=0 go build ./cmd/cli/migrate

FROM alpine:3
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /welsh/api ./
COPY --from=0 /welsh/migrate ./
COPY --from=0 /welsh/migrations ./migrations
RUN echo './migrate -dns=${POSTGRES_DNS} -up' >> startup.sh
RUN echo 'rm -r ./migrations' >> startup.sh
RUN echo 'rm ./migrate' >> startup.sh
RUN echo './api -db-dns=${POSTGRES_DNS}' >> startup.sh
RUN chmod u+x startup.sh
CMD ["/bin/sh", "./startup.sh"]