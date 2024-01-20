FROM golang:1.21

WORKDIR /usr/src/masonictempl
ENV PWD=/usr/src/masonictempl

# Note: Make sure you have the certs and creds folders in a local env folder
COPY . .

RUN go mod download

RUN make build

ENV PORT=8080

# Should be set by default on the target machine, but just incase.
EXPOSE 8080

CMD ["./bin/masonictempl"]


