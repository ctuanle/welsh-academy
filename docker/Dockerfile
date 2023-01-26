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
# api binary
COPY --from=0 /welsh/api ./
# migrate cli binary
COPY --from=0 /welsh/migrate ./
# migration files
COPY --from=0 /welsh/migrations ./migrations

# create a startup.sh script that
# make POSGRES_DNS from POSTGRES_PASSWORD and POSTGRES_DB,
# run migrations,
# remove migrate binary and migration-files 
# and start the api
RUN echo 'POSTGRES_DNS=postgres://postgres:${POSTGRES_PASSWORD}@postgresql/${POSTGRES_DB}?sslmode=disable' >> startup.sh
RUN echo './migrate -dns=${POSTGRES_DNS} -up' >> startup.sh
RUN echo 'rm -r ./migrations' >> startup.sh
RUN echo 'rm ./migrate' >> startup.sh
RUN echo './api -db-dns=${POSTGRES_DNS}' >> startup.sh
RUN chmod u+x startup.sh

CMD ["/bin/sh", "./startup.sh"]