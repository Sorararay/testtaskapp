FROM golang:1.21.6

WORKDIR /migrator

RUN apt update
RUN apt install -y postgresql-common 
RUN /usr/share/postgresql-common/pgdg/apt.postgresql.org.sh -y
RUN apt install -y postgresql-client-16
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0


COPY . .
