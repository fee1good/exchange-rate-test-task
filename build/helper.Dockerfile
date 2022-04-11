FROM golang:1.18

ENV GO111MODULE="on"
ENV CGO_ENABLED="0"
ENV GOOS="linux"

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar -xvz \
    && mv ./migrate.linux-amd64 /bin/migrate \
    && wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.38.0 \
    && mv ./bin/golangci-lint /bin/golangci-lint

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
