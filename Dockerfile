FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY ./config ./config
COPY ./locales ./locales

RUN mv ./config/config.staging.yaml ./config/config.yaml

COPY ./build/server-linux ./server
CMD ["./server"]