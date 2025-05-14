#!/bin/bash

# Food Delivery Microservices Deployment Script
# This script deploys the Food Delivery microservices to Kubernetes

# Exit on error
set -e

# Default values
BUILD_IMAGES=true
SKIP_INFRA=false
SETUP_DNS=false
SERVICES_TO_DEPLOY="all"

# Display help message
show_help() {
  echo "Food Delivery Microservices Deployment Script"
  echo ""
  echo "Usage: $0 [options]"
  echo ""
  echo "Options:"
  echo "  -h, --help                 Show this help message"
  echo "  -s, --skip-build           Skip building Docker images"
  echo "  -i, --skip-infra           Skip deploying infrastructure components"
  echo "  -d, --setup-dns            Setup local DNS resolution (no /etc/hosts modification needed)"
  echo "  --services <service-list>  Deploy only specific services (comma-separated)"
  echo "                             Available: restaurant,user,food,category,media,cart,order,payment"
  echo ""
  echo "Examples:"
  echo "  $0                         Deploy everything with default settings"
  echo "  $0 --skip-build            Deploy without rebuilding Docker images"
  echo "  $0 --services restaurant,user  Deploy only restaurant and user services"
  echo "  $0 --skip-infra --services cart  Deploy only cart service, skip infrastructure"
  echo ""
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--help)
      show_help
      exit 0
      ;;
    -s|--skip-build)
      BUILD_IMAGES=false
      shift
      ;;
    -i|--skip-infra)
      SKIP_INFRA=true
      shift
      ;;
    -d|--setup-dns)
      SETUP_DNS=true
      shift
      ;;
    --services)
      if [[ -z "$2" || "$2" == -* ]]; then
        echo "Error: --services requires an argument"
        exit 1
      fi
      SERVICES_TO_DEPLOY="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      show_help
      exit 1
      ;;
  esac
done

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
  echo "Error: Docker is not running. Please start Docker and try again."
  exit 1
fi

# Check if kubectl is available
if ! command -v kubectl &> /dev/null; then
  echo "Error: kubectl is not installed or not in PATH"
  exit 1
fi

# Check if Kubernetes is running
if ! kubectl get nodes &> /dev/null; then
  echo "Error: Kubernetes is not running. Please start Kubernetes and try again."
  exit 1
fi

# Build Docker images if not skipped
if [ "$BUILD_IMAGES" = true ]; then
  echo "Building Docker images for services..."

  # Determine which services to build
  if [ "$SERVICES_TO_DEPLOY" = "all" ]; then
    # Build all services
    docker build -t food-delivery/restaurant-service:latest -f modules/restaurant/Dockerfile .
    docker build -t food-delivery/user-service:latest -f modules/user/Dockerfile .
    docker build -t food-delivery/food-service:latest -f modules/food/Dockerfile .
    docker build -t food-delivery/category-service:latest -f modules/category/Dockerfile .
    docker build -t food-delivery/media-service:latest -f modules/media/Dockerfile .
    docker build -t food-delivery/cart-service:latest -f modules/cart/Dockerfile .
    docker build -t food-delivery/order-service:latest -f modules/order/Dockerfile .
    docker build -t food-delivery/payment-service:latest -f modules/payment/Dockerfile .
  else
    # Build only specified services
    IFS=',' read -ra SERVICES <<< "$SERVICES_TO_DEPLOY"
    for service in "${SERVICES[@]}"; do
      if [ -d "modules/$service" ]; then
        echo "Building $service service..."
        docker build -t food-delivery/$service-service:latest -f modules/$service/Dockerfile .
      else
        echo "Warning: Service '$service' not found, skipping build"
      fi
    done
  fi

  echo "Docker images built successfully."
else
  echo "Skipping Docker image build as requested."
fi

echo "Deploying to Kubernetes..."

# Create namespace
kubectl apply -f deployments/00-namespace.yaml

# Apply ConfigMaps and Secrets
kubectl apply -f deployments/01-configmap.yaml
kubectl apply -f deployments/02-secrets.yaml

# Apply PersistentVolumes and PersistentVolumeClaims
kubectl apply -f deployments/03-persistent-volumes.yaml

# Deploy infrastructure components if not skipped
if [ "$SKIP_INFRA" = false ]; then
  echo "Deploying infrastructure components..."
  kubectl apply -f deployments/04-mysql.yaml
  kubectl apply -f deployments/05-redis.yaml
  kubectl apply -f deployments/06-minio.yaml
  kubectl apply -f deployments/07-elasticsearch-kibana.yaml

  # Apply Filebeat for log collection
  kubectl apply -f deployments/filebeat-config.yaml
  kubectl apply -f deployments/18-filebeat.yaml

  # Wait for infrastructure to be ready
  echo "Waiting for infrastructure components to be ready..."
  kubectl wait --namespace food-delivery --for=condition=ready pod --selector=app=mysql --timeout=300s || echo "Warning: MySQL not ready in time, continuing anyway"
  kubectl wait --namespace food-delivery --for=condition=ready pod --selector=app=redis --timeout=300s || echo "Warning: Redis not ready in time, continuing anyway"
  kubectl wait --namespace food-delivery --for=condition=ready pod --selector=app=minio --timeout=300s || echo "Warning: MinIO not ready in time, continuing anyway"
  kubectl wait --namespace food-delivery --for=condition=ready pod --selector=app=elasticsearch --timeout=300s || echo "Warning: Elasticsearch not ready in time, continuing anyway"
else
  echo "Skipping infrastructure deployment as requested."
fi

# Deploy microservices
echo "Deploying microservices..."
if [ "$SERVICES_TO_DEPLOY" = "all" ]; then
  # Deploy all services
  kubectl apply -f deployments/08-restaurant-service.yaml
  kubectl apply -f deployments/09-user-service.yaml
  kubectl apply -f deployments/10-food-service.yaml
  kubectl apply -f deployments/11-category-service.yaml
  kubectl apply -f deployments/12-media-service.yaml
  kubectl apply -f deployments/13-cart-service.yaml
  kubectl apply -f deployments/14-order-service.yaml
  kubectl apply -f deployments/15-payment-service.yaml
else
  # Deploy only specified services
  IFS=',' read -ra SERVICES <<< "$SERVICES_TO_DEPLOY"
  for service in "${SERVICES[@]}"; do
    service_file=""
    case "$service" in
      restaurant)
        service_file="deployments/08-restaurant-service.yaml"
        ;;
      user)
        service_file="deployments/09-user-service.yaml"
        ;;
      food)
        service_file="deployments/10-food-service.yaml"
        ;;
      category)
        service_file="deployments/11-category-service.yaml"
        ;;
      media)
        service_file="deployments/12-media-service.yaml"
        ;;
      cart)
        service_file="deployments/13-cart-service.yaml"
        ;;
      order)
        service_file="deployments/14-order-service.yaml"
        ;;
      payment)
        service_file="deployments/15-payment-service.yaml"
        ;;
      *)
        echo "Warning: Unknown service '$service', skipping deployment"
        continue
        ;;
    esac

    if [ -n "$service_file" ]; then
      echo "Deploying $service service..."
      kubectl apply -f "$service_file"
    fi
  done
fi

# Apply Ingress
kubectl apply -f deployments/16-ingress.yaml

echo "Deployment completed successfully!"

# Setup local DNS if requested
if [ "$SETUP_DNS" = true ]; then
  echo "Setting up local DNS resolution..."
  ./deployments/local-dns.sh
  echo ""
  echo "You can now access the application at: http://food-delivery.local"
  echo "Access Kibana at: http://food-delivery.local/kibana"
  echo "Access MinIO Console at: http://food-delivery.local/minio"
else
  echo ""
  echo "Access the application at: http://food-delivery.local"
  echo "Access Kibana at: http://food-delivery.local/kibana"
  echo "Access MinIO Console at: http://food-delivery.local/minio"
fi
