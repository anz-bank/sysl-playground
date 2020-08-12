FROM golang:1.14-alpine
ENV PORT=80 
ENV SYSL_PLANTUML=http://www.plantuml.com/plantuml
RUN apk add --no-cache git
#install sysl
RUN GO111MODULE=on go get -u github.com/anz-bank/sysl/cmd/sysl
RUN sysl --version

# Set the Current Working Directory inside the container
WORKDIR /app

COPY . .

# Build the Go app
RUN go build -o ./bin/sysl-playground .

# This container exposes port 8080 to the outside world
EXPOSE 80

# Run the binary program produced by `go install`
CMD ["./bin/sysl-playground"]