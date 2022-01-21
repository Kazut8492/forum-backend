FROM golang:1.16-alpine

WORKDIR /app

COPY . ./
RUN go mod download
RUN apk add --update gcc musl-dev
RUN go build -o /forum

EXPOSE 8888

CMD [ "/forum" ]