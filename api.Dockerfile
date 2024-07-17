FROM golang:1.22.5 as base

FROM base as dev

EXPOSE 8010

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN > .env

COPY . .

ENTRYPOINT [ "go", "run" ]
CMD [ "." ]

