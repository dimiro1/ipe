FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN apk add git gcc musl-dev
RUN go build -o ipe ./cmd
FROM alpine
USER root
RUN  mkdir -p /config
RUN adduser -S -D -H -h /app appuser
COPY ./entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh
USER appuser
WORKDIR /app
COPY --from=builder /build/ipe /app/
COPY --from=builder /build/config-example.yml /app/config-example.yml
VOLUME /config
CMD ["/bin/sh", "/app/entrypoint.sh"]
EXPOSE 4343
EXPOSE 8080