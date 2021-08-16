FROM golang:alpine

LABEL key="MINT"

RUN mkdir /src
WORKDIR /src

COPY . .


RUN go env -w GO111MODULE=on \
  && go env -w GOPROXY=https://goproxy.cn,direct \
  && go mod download 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.BuildEnv=prod" main.go


FROM alpine:3.12.3

RUN apk --no-cache add ca-certificates vips \
  && update-ca-certificates

ENV CGO_CFLAGS_ALLOW=-Xpreprocessor

COPY --from=0 /src/main .
COPY --from=0 /src/config.prod.json .
RUN mkdir /upload

EXPOSE 9099

ENTRYPOINT ["./main"]