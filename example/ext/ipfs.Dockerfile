###
FROM golang:1.14.4-alpine3.12 as builder

RUN apk add --no-cache bash git gcc libc-dev make openssl-dev

WORKDIR /app

RUN git clone https://github.com/gostones/go-ipfs.git ipfs

ARG CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOPATH=

# COPY ipdr/go.mod ipdr/go.sum ipdr/
# RUN cd ipdr && go mod download

# COPY ipfs/go.mod ipfs/go.sum ipfs/
# RUN cd ipfs && go mod download

# COPY . .

# RUN cd ipfs && go build -o ../bin/ipfs ./cmd/ipfs
RUN cd ipfs && make build CGO_ENABLED=1 GOTAGS=openssl
# RUN cd ipdr && go build -o ../bin/ipdr ./cmd/ipdr

ADD ipfs-entrypoint.sh .

###
FROM alpine:3.12

RUN apk add --no-cache tini

# COPY --from=builder /app/bin/* /usr/local/bin/
COPY --from=builder /app/ipfs/cmd/ipfs/ipfs /usr/local/bin/
COPY --from=builder /app/ipfs-entrypoint.sh /

# docker v2 registry
EXPOSE 5000
# Swarm TCP; should be exposed to the public
EXPOSE 4001
# Swarm UDP; should be exposed to the public
EXPOSE 4001/udp
# Daemon API; must not be exposed publicly but to client services under you control
EXPOSE 5001
# Web Gateway; can be exposed publicly with a proxy, e.g. as https://ipfs.example.org
EXPOSE 8080
# Swarm Websockets; must be exposed publicly when the node is listening using the websocket transport (/ipX/.../tcp/8081/ws).
EXPOSE 8081

ENV USER=app
ENV DATA_PATH=/data
ENV IPFS_PATH=/data/ipfs

RUN adduser -D -h /data -u 1000 -G users $USER \
    && mkdir -p $IPFS_PATH \
    && chown -R $USER:users /data

VOLUME $DATA_PATH

##
USER $USER

# ENTRYPOINT ["/sbin/tini", "--"]
# CMD ["/ipfs-entrypoint.sh"]

ENTRYPOINT ["/sbin/tini", "--", "/ipfs-entrypoint.sh"]
CMD ["ipfs", "daemon", "--migrate=true"]