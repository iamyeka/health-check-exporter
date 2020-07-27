FROM alpine:3.12
COPY ./app /app
ENTRYPOINT /app
