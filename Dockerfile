FROM ubuntu:22.04

RUN apt update
RUN apt upgrade -y

RUN apt install -y curl 

COPY otelcol /otelcol
RUN chmod +x /otelcol

ENTRYPOINT ["/otelcol"]
EXPOSE 4317 55680 55679
