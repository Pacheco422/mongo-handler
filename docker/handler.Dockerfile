# First Stage
FROM ubuntu:22.04@sha256:b2175cd4cfdd5cdb1740b0e6ec6bbb4ea4892801c0ad5101a81f694152b6c559 as builder

ENV DEBIAN_FRONTEND=noninteractive
ENV SSL_CERT_DIR=/etc/ssl/certs
RUN apt update && \
    apt install -y \
        ca-certificates=20211016ubuntu0.22.04.1 \
        openssl=3.0.2-0ubuntu1.8 \
        golang-go=2:1.18~0ubuntu2 && \
    mkdir /application
COPY src/ /application/
WORKDIR /application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o mongo-handler .

# Second Stage
FROM alpine:3.17.1 as application
RUN chmod a-w /etc && \
    addgroup -S appgroup && \
    adduser -S appuser -G appgroup -h /home/appuser && \
    rm -rf /bin/*

COPY --from=builder /application /home/appuser/
EXPOSE 9539
USER appuser

ENTRYPOINT [ "/home/appuser/mongo-handler" ]
