FROM alpine:3.4

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY dbheroapp_linux /app/
COPY scripts/website_entrypoint.sh /app/entrypoint.sh
COPY www /app/www/

EXPOSE 80 443

CMD ["/app/entrypoint.sh"]
