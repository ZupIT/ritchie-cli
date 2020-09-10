FROM chiaretto/dotnet-core-debian-buildessential:3.1 AS build-env

USER root

RUN curl -fsSL https://commons-repo.ritchiecli.io/install.sh | bash

RUN mkdir /rit
COPY . /rit
RUN sed -i 's/\r//g' /rit/set_umask.sh
RUN sed -i 's/\r//g' /rit/run.sh
RUN chmod +x /rit/set_umask.sh

WORKDIR /app
ENTRYPOINT ["/rit/set_umask.sh"]
CMD ["/rit/run.sh"]
