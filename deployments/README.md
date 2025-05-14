# Food Delivery Microservices Kubernetes Deployment

This directory contains Kubernetes deployment configurations for the Food Delivery microservices application.

## Prerequisites

- Docker Desktop with Kubernetes enabled
- kubectl command-line tool
- Ingress controller installed (NGINX Ingress Controller recommended)

## Deployment Structure

The deployment is organized into the following files:

- `00-namespace.yaml`: Creates the food-delivery namespace
- `01-configmap.yaml`: ConfigMap for application configuration
- `02-secrets.yaml`: Secrets for sensitive information
- `03-persistent-volumes.yaml`: PersistentVolumes and PersistentVolumeClaims for stateful services
- `04-mysql.yaml`: MySQL database deployment
- `05-redis.yaml`: Redis cache deployment
- `06-minio.yaml`: MinIO object storage deployment
- `07-elasticsearch-kibana.yaml`: Elasticsearch and Kibana deployments
- `08-restaurant-service.yaml`: Restaurant service deployment
- `09-user-service.yaml`: User service deployment
- `10-food-service.yaml`: Food service deployment
- `11-category-service.yaml`: Category service deployment
- `12-media-service.yaml`: Media service deployment
- `13-cart-service.yaml`: Cart service deployment
- `14-order-service.yaml`: Order service deployment
- `15-payment-service.yaml`: Payment service deployment
- `16-ingress.yaml`: Ingress configuration for external access
- `18-filebeat.yaml`: Filebeat DaemonSet for log collection
- `filebeat-config.yaml`: Configuration for Filebeat
- `kibana-dashboards.yaml`: Kibana dashboards for log visualization

## Deployment Scripts

- `deploy.sh`: Main deployment script with various options
- `local-dns.sh`: Script to set up local DNS resolution without modifying /etc/hosts
- `cleanup.sh`: Script to clean up all resources

## Deployment Options

The `deploy.sh` script supports several options:

```
Usage: ./deployments/deploy.sh [options]

Options:
  -h, --help                 Show help message
  -s, --skip-build           Skip building Docker images
  -i, --skip-infra           Skip deploying infrastructure components
  -d, --setup-dns            Setup local DNS resolution (no /etc/hosts modification needed)
  --services <service-list>  Deploy only specific services (comma-separated)
                             Available: restaurant,user,food,category,media,cart,order,payment

Examples:
  ./deployments/deploy.sh                         Deploy everything with default settings
  ./deployments/deploy.sh --skip-build            Deploy without rebuilding Docker images
  ./deployments/deploy.sh --services restaurant,user  Deploy only restaurant and user services
  ./deployments/deploy.sh --skip-infra --services cart  Deploy only cart service, skip infrastructure
```

## Deployment Steps

1. Make sure Docker Desktop with Kubernetes is running
2. Install NGINX Ingress Controller if not already installed:
   ```
   kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.2/deploy/static/provider/cloud/deploy.yaml
   ```
3. Run the deployment script with your preferred options:
   ```
   ./deployments/deploy.sh
   ```
4. Access the application using one of the following methods:

### Method 1: Using Local DNS (No /etc/hosts modification needed)

Run the deployment with DNS option:
```
./deployments/deploy.sh --setup-dns
```

Or set up DNS separately:
```
./deployments/local-dns.sh
```

Then access the application at: http://food-delivery.local

### Method 2: Using /etc/hosts (Traditional method)

Since you've already added the entry to your `/etc/hosts` file:
```
127.0.0.1 food-delivery.local
```

You can access the application at: http://food-delivery.local

## Service Endpoints

### API Services (via Ingress)
- Restaurant Service: http://food-delivery.local/api/restaurant
- User Service: http://food-delivery.local/api/user
- Food Service: http://food-delivery.local/api/food
- Category Service: http://food-delivery.local/api/category
- Media Service: http://food-delivery.local/api/media
- Cart Service: http://food-delivery.local/api/cart
- Order Service: http://food-delivery.local/api/order
- Payment Service: http://food-delivery.local/api/payment

### Infrastructure Services (via Ingress)
- Kibana: http://food-delivery.local/kibana
- MinIO Console: http://food-delivery.local/minio

### Direct Access via NodePort
For more reliable access, especially for Kibana and Elasticsearch, you can use these NodePort services:

- Kibana: http://localhost:30601
- Elasticsearch: http://localhost:30920

These NodePort services bypass the ingress and connect directly to the services, which can resolve browser compatibility issues.

## Monitoring and Management

You can monitor the deployment using:

```
kubectl get all -n food-delivery
```

To view logs for a specific service directly:

```
kubectl logs -n food-delivery deployment/[service-name]
```

For example:

```
kubectl logs -n food-delivery deployment/restaurant-service
```

## Centralized Logging

The deployment includes a centralized logging system using the ELK stack (Elasticsearch, Logstash, Kibana) with Filebeat for log collection:

1. **Filebeat**: Deployed as a DaemonSet on all nodes to collect container logs
2. **Elasticsearch**: Stores and indexes all logs
3. **Kibana**: Provides visualization and search capabilities for logs

### Accessing Logs in Kibana

1. Access Kibana at http://food-delivery.local/kibana
2. Navigate to the "Discover" section to search and filter logs
3. Use the pre-configured dashboards to visualize log data:
   - "Food Delivery Overview" dashboard shows log counts by service and recent logs
   - Filter logs by service using the kubernetes.container.name field
   - Search for errors using keywords like "error", "fail", or "exception"

### Log Filtering Examples

In Kibana's search bar, you can use queries like:

- `kubernetes.container.name: "restaurant-service"` - Show only restaurant service logs
- `message: *error*` - Show logs containing the word "error"
- `kubernetes.namespace: "food-delivery" AND message: *exception*` - Show exception logs in the food-delivery namespace

## Cleanup

To remove the deployment, you can use the cleanup script:

```
./deployments/cleanup.sh
```

This script will:
- Stop the local DNS container if running
- Delete all Kubernetes resources in the food-delivery namespace
- Optionally delete all Docker images

Alternatively, you can manually remove the namespace:

```
kubectl delete namespace food-delivery
```

## Troubleshooting



### DNS Resolution Issues

If you're using the local DNS method and having issues:

1. Check if the DNS container is running:
   ```
   docker ps | grep food-delivery-dns
   ```

2. Restart the DNS container:
   ```
   docker restart food-delivery-dns
   ```

3. Verify your DNS settings are correctly pointing to 127.0.0.1

### Service Connectivity Issues

If services can't communicate with each other:

1. Check if all services are running:
   ```
   kubectl get pods -n food-delivery
   ```

2. Check service logs:
   ```
   kubectl logs -n food-delivery deployment/<service-name>
   ```

3. Verify the ConfigMap has the correct service URLs:
   ```
   kubectl get configmap -n food-delivery app-config -o yaml
   ```
