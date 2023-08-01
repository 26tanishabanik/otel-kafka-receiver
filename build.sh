docker build -t otelcollector-test -f ./Dockerfile .
docker tag otelcollector-test:latest tanishabanik/otelcollector-test:customkafka0.2
docker push tanishabanik/otelcollector-test:customkafka0.2
