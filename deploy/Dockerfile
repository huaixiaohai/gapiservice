FROM alpine:3.9
RUN mkdir -p /app
WORKDIR /app

COPY ./bin/main /app/main
#COPY ./builder/deploy/entrypoint.sh /app/entrypoint.sh
COPY ./config/config.yaml /app/config/config.yaml

#RUN chmod a+x /app/entrypoint.sh
RUN chmod +x /app/main

CMD ["./main"]
