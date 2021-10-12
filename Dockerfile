FROM golang:1.16.0-alpine3.13

ADD ./intraStatusBot /intraStatusBot/
WORKDIR /intraStatusBot

RUN go install .

CMD intraStatusBot ./dev.env