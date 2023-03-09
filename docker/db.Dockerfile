FROM mongo:6.0.4@sha256:aac6f27030c440d738d3fe8c5e9d36a6f9ceed2c02c35c2d28dfe0628d55bb0f

RUN chmod a-w /etc                  && \
    groupadd --gid 20001 application && \
    useradd --uid 10001 --gid application --shell /bin/bash --create-home application

USER application

EXPOSE 27017