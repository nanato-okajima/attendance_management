FROM golang:1.18.3

RUN apt install -y \
    git

ARG DIR=/go/app/
WORKDIR ${DIR}

COPY ./go.mod ${DIR}

RUN go install github.com/cosmtrek/air@v1.40.2

CMD [ "air","-c", ".air.toml" ]
