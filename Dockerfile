FROM debian:stable-slim
RUN mkdir -p /app/
COPY ./bin/go-ikea-server /app
WORKDIR /app
ENTRYPOINT [ "./go-ikea-server"]