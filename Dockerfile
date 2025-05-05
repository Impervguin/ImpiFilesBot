FROM golang:1.24-alpine
LABEL AUTHOR="Impervguin"

RUN apk add --no-cache make

RUN mkdir /build
WORKDIR /build


COPY go.* .
RUN go mod download
COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/
COPY ./bot/ ./bot/
COPY ./config/*.yaml ./config/

CMD ["go", "run", "./cmd/bot"]

