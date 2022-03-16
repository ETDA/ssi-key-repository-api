FROM golang:1.16.4-alpine3.13

RUN apk update && apk upgrade && \
apk add --no-cache bash git openssh
WORKDIR /app
RUN git config --global url."https://idin-core-deploy:Zr2TnbF6X9oMLAQKxvvX@gitlab.finema.co".insteadOf "https://gitlab.finema.co"
ADD go.mod go.sum /app/
RUN go mod download
ADD . /app/
ADD .env /app/seeds/
CMD cd /app/ && go run indexes/index.go
