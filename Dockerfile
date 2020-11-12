ARG BUILD_IMAGE=golang:1.13-alpine
FROM $BUILD_IMAGE

RUN echo "7224aa97-1f94-4679-917d-9c7bb10074ab" > /etc/machine-id
RUN apk update && apk add \
    alpine-sdk \
    gettext \
    curl \
    python3 \
    py-pip \
    jq \
    git \
    && pip install --no-cache-dir awscli==1.16.310 \
    && apk del py-pip \
    && rm -rf /var/cache/apk/* /root/.cache/pip/* /usr/lib/python2.7/site-packages/awscli/examples
