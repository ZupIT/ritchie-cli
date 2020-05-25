FROM stefanscherer/chocolatey

MAINTAINER ritchie-cli

RUN choco install -y golang
RUN choco install -y make
RUN choco install -y grep
RUN go version

RUN mkdir source
WORKDIR /source

COPY . /source
ADD . /source