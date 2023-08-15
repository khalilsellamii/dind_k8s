# Two-stage build:
#    first  FROM prepares a binary file in full environment ~780MB
#    second FROM takes only binary file ~10MB

FROM golang:1.18.1-alpine AS builder

RUN go version

COPY . "/"
WORKDIR "/"


RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
RUN go get -v -t  .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -buildvcs=false  -o /go_script

#CMD ["/go_script"]



#########
# second stage to obtain a very small image
FROM scratch
WORKDIR "/"

COPY --from=builder /go_script .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/go_script"]