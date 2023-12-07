FROM golang:1.21-alpine as builder

# configuration for dependency from private repository
RUN apk update && apk add --no-cache git
ARG GITHUB_TOKEN
ENV CGO_ENABLED=0 GO111MODULE=on GOOS=linux
ENV GOPRIVATE="github.com/jinleejun-corp"
RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

COPY . .

RUN go build -o main

FROM alpine:3.18

COPY --from=builder /app ./app

# set default timezone
RUN apk add --no-cache tzdata
ENV TZ Asia/Bangkok
RUN ln -sf /usr/share/zoneinfo/Asia/Bangkok /etc/localtime
RUN echo "Asia/Bangkok" > /etc/timezone

EXPOSE 8080

CMD [ "./app/main" ]