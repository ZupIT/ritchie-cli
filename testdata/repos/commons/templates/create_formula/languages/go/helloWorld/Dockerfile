FROM alpine:latest

ADD . .
RUN chmod +x set_umask.sh

WORKDIR /app
ENTRYPOINT ["/set_umask.sh"]
CMD ["/linux/main"]