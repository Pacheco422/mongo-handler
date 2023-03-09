# First Stage
FROM ubuntu:22.04@sha256:b2175cd4cfdd5cdb1740b0e6ec6bbb4ea4892801c0ad5101a81f694152b6c559

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y \
        curl=7.81.0-1ubuntu1.8 && \
    mkdir /application
WORKDIR /application

ENTRYPOINT [ "sleep", "3600" ]
