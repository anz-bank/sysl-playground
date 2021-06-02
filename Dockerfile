FROM anzbank/sysl as sysl

FROM golang:1.14-alpine AS builder
WORKDIR /app
COPY . .

# Build the Go app
RUN go build -o ./bin/sysl-playground .

# Build the plantuml server
FROM maven:3-jdk-11 AS builderjetty
WORKDIR /app

COPY plantuml .

RUN mvn --batch-mode --define java.net.useSystemProxies=true package

# final image
FROM golang:1.14-alpine

ENV PORT=80
RUN apk add --no-cache openjdk8 graphviz font-noto-cjk bash
ENV SYSL_PLANTUML=http://localhost:8080/plantuml

COPY --from=sysl /sysl /usr/local/bin

RUN sysl --version

WORKDIR /src
COPY --from=builder /app .
COPY --from=builderjetty /app/target ./target
EXPOSE 80
RUN chmod +x ./script/start.sh
CMD ["./script/start.sh"]
