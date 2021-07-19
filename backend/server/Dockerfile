FROM golang:1.16 AS builder

#RUN apk --update add \
#	git openssl \
#	&& rm /var/cache/apk/*

RUN mkdir -p /go/src/drone-delivery/server
WORKDIR /go/src/drone-delivery/server

COPY go.mod /go/src/drone-delivery/server
COPY go.sum /go/src/drone-delivery/server
RUN go mod download
COPY ./ /go/src/drone-delivery/server
WORKDIR /go/src/drone-delivery/server/cmd/drone-delivery-server
RUN go install

#Second stage
#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#RUN mkdir -p /go/src/drone-delivery/server
#WORKDIR /go/src/drone-delivery/server
#COPY --from=builder /go/src/drone-delivery/server /go/src/drone-delivery/server
#COPY --from=builder /go/bin/drone-delivery-server /go/bin/drone-delivery-server
ENTRYPOINT ["/go/bin/drone-delivery-server"]
EXPOSE 5000
EXPOSE 50051