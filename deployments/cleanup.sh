#!/bin/bash

# This script cleans up all resources created by the deployment
# It removes all Kubernetes resources and Docker containers/images

# Exit when any command fails
set -e

echo "Starting cleanup process..."

# No port forwarding to stop since we're not using it

# Stop local DNS if running
if docker ps | grep -q "food-delivery-dns"; then
  echo "Stopping local DNS container..."
  docker stop food-delivery-dns
  docker rm food-delivery-dns
fi

# Ask for confirmation before deleting Kubernetes resources
read -p "Are you sure you want to delete all Kubernetes resources in the food-delivery namespace? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
  echo "Deleting all resources in the food-delivery namespace..."

  # Delete ingress first to avoid connection issues
  kubectl delete -f deployments/16-ingress.yaml --ignore-not-found=true

  # Delete services
  echo "Deleting microservices..."
  kubectl delete -f deployments/15-payment-service.yaml --ignore-not-found=true
  kubectl delete -f deployments/14-order-service.yaml --ignore-not-found=true
  kubectl delete -f deployments/13-cart-service.yaml --ignore-not-found=true
  kubectl delete -f deployments/12-media-service.yaml --ignore-not-found=true
  kubectl delete -f deployments/11-category-service.yaml --ignore-not-found=true
  kubectl delete -f deployments/10-food-service.yaml --ignore-not-found=true
  kubectl delete -f deployments/09-user-service.yaml --ignore-not-found=true
  kubectl delete -f deployments/08-restaurant-service.yaml --ignore-not-found=true

  # Delete infrastructure components
  echo "Deleting infrastructure components..."
  kubectl delete -f deployments/18-filebeat.yaml --ignore-not-found=true
  kubectl delete -f deployments/filebeat-config.yaml --ignore-not-found=true
  kubectl delete -f deployments/07-elasticsearch-kibana.yaml --ignore-not-found=true
  kubectl delete -f deployments/06-minio.yaml --ignore-not-found=true
  kubectl delete -f deployments/05-redis.yaml --ignore-not-found=true
  kubectl delete -f deployments/04-mysql.yaml --ignore-not-found=true

  # Delete persistent volumes
  echo "Deleting persistent volumes..."
  kubectl delete -f deployments/03-persistent-volumes.yaml --ignore-not-found=true

  # Delete config and secrets
  echo "Deleting ConfigMaps and Secrets..."
  kubectl delete -f deployments/02-secrets.yaml --ignore-not-found=true
  kubectl delete -f deployments/01-configmap.yaml --ignore-not-found=true

  # Delete namespace
  echo "Deleting namespace..."
  kubectl delete -f deployments/00-namespace.yaml --ignore-not-found=true

  echo "All Kubernetes resources have been deleted."
else
  echo "Kubernetes resources deletion cancelled."
fi

# Ask for confirmation before deleting Docker images
read -p "Do you want to delete all food-delivery Docker images? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
  echo "Deleting Docker images..."
  docker images | grep "food-delivery/" | awk '{print $1":"$2}' | xargs -r docker rmi
  echo "Docker images deleted."
else
  echo "Docker images deletion cancelled."
fi

# Remove Docker network if it exists
if docker network ls | grep -q "food-delivery-net"; then
  echo "Removing Docker network..."
  docker network rm food-delivery-net
fi

echo "Cleanup completed successfully!"
