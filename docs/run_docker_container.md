docker run -d --name media-service --network fd-net \
  -e DB_DSN="fduser:fduser@123@tcp(your-mysql-container-name:3306)/fddb?charset=utf8&parseTime=True&loc=Local" \
  -e JWT_SECRET_KEY="your_key" \
  -e USER_SERVICE_URL="http://localhost:8085/v1/rpc/users" \
  -e NATS_URL="nats://your-nat-container-name:4222" \
  -e JAEGER_ENDPOINT="your-jaeger-container-name:4317" \
  -e MINIO_ACCESS_KEY="your_key" \
  -e MINIO_BUCKET_NAME="your_bucket_name" \
  -e MINIO_DOMAIN="your_container_name:9000" \
  -e MINIO_REGION="your_region" \
  -e MINIO_SECRET_KEY="your_key" \
  -e VIRTUAL_HOST="fooddelivery.localhost" \
  -e VIRTUAL_PORT="3000"  \
  -e VIRTUAL_PATH="/v1/users" \
  -p 8085:3000 \
  food-delivery-backend:1.0.0

<!---
nginx proxy ( API Gateway)
VIRTUAL_HOST
VIRTUAL_PORT
VIRTUAL_PATH
-->
docker network ls
docker network create fd-net
docker network connect fd-net mysql-container
docker network connect fd-net nats-container
docker network connect fd-net minio-container
docker network connect fd-net jaeger-container
docker network connect fd-net nginx-container

docker logs media-service

curl http://localhost:8085/ping

