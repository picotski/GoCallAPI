FROM golang:1.22.5

EXPOSE 8010

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

RUN > .env

COPY . .

ENTRYPOINT [ "go", "run" ]
CMD [ "." ]
