FROM golang:1.14-alpine AS builder
WORKDIR /app
COPY . .

# Build the Go app
RUN go build -o ./bin/sysl-playground .

FROM alpine:latest

ENV PORT=80
RUN apk add --no-cache go maven openjdk8 graphviz font-noto-cjk bash
RUN go version
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV SYSL_PLANTUML=http://localhost:8080/plantuml

#install sysl
RUN GOPATH=$(go env GOPATH)
RUN GO111MODULE=on go get -u github.com/anz-bank/sysl/cmd/sysl
RUN sysl --version

WORKDIR /src
COPY --from=builder /app .
EXPOSE 80
RUN chmod +x ./script/start.sh
CMD ["./script/start.sh"]