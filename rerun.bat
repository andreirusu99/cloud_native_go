docker stop kvs
docker rm kvs
docker build --tag kvs:multipart .
docker run -d -p 5555:5555 --name kvs --add-host=host.docker.internal:host-gateway --env-file .env kvs:multipart