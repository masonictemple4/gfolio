FROM golang:1.21

WORKDIR /usr/src/masonictempl

# Note: Make sure you have the certs and creds folders in a local env folder
COPY . .

# Should be set by default on the target machine, but just incase.
ENV PORT=8080

RUN go install

EXPOSE 8080

ENTRYPOINT ["masonictempl"]


