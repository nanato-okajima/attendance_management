FROM golang:1.24

RUN apt install -y \
    git

ARG DIR=/go/app/
WORKDIR ${DIR}

COPY ./go.mod ${DIR}

RUN go install github.com/cosmtrek/air@v1.40.2
RUN go install go.uber.org/mock/mockgen@v0.6.0

# Install golangci-lint from source for Go 1.24 compatibility
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.62.2

CMD [ "air","-c", ".air.toml" ]
