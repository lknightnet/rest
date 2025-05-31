FROM ubuntu

RUN apt update && \
    DEBIAN_FRONTEND=noninteractive apt install -y tzdata && \
    rm -rf /var/lib/apt/lists/*

RUN apt update && \
    apt install -y wget build-essential

WORKDIR /rest-service

COPY config/config.yaml config/
COPY rest-service /rest-service

EXPOSE 8080

CMD ["./rest-service"]