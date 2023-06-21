docker build -t echo-service .
docker run -d -p  7777:7777 --name echo-service echo-service
docker image tag echo-service dauh/echo-service:v1
docker image push dauh/echo-service:v1