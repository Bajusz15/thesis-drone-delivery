FROM golang:1.16 AS builder

#RUN apk --update add \
#	git openssl \
#	&& rm /var/cache/apk/*

RUN mkdir -p /go/src/drone-delivery/drone-swarm
WORKDIR /go/src/drone-delivery/drone-swarm

COPY drone-swarm/go.mod /go/src/drone-delivery/drone-swarm
COPY drone-swarm/go.sum /go/src/drone-delivery/drone-swarm
COPY ./server/ /go/src/drone-delivery/server
COPY ./ /go/src/drone-delivery/
WORKDIR /go/src/drone-delivery/drone-swarm/cmd/drone-delivery-swarm
RUN go mod download
RUN go install

#Second stage
#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#RUN mkdir -p /go/src/drone-delivery/
#WORKDIR /go/src/drone-delivery/
#COPY --from=builder /go/src/drone-delivery/ /go/src/drone-delivery/
#COPY --from=builder /go/bin/drone-delivery-swarm /go/bin/drone-delivery-swarm
ENTRYPOINT ["/go/bin/drone-delivery-swarm"]
EXPOSE 2000