#!/bin/bash

# Kafka Setup Script for Food Delivery Service
# This script sets up Kafka infrastructure for local development

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
KAFKA_VERSION="7.4.0"
NETWORK_NAME="food-delivery-network"
COMPOSE_FILE="deployments/kafka/docker-compose.kafka.yml"

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if Docker is running
check_docker() {
    print_status "Checking Docker..."
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    print_success "Docker is running"
}

# Function to check if Docker Compose is available
check_docker_compose() {
    print_status "Checking Docker Compose..."
    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose and try again."
        exit 1
    fi
    print_success "Docker Compose is available"
}

# Function to create Docker network
create_network() {
    print_status "Creating Docker network: $NETWORK_NAME"
    if docker network ls | grep -q "$NETWORK_NAME"; then
        print_warning "Network $NETWORK_NAME already exists"
    else
        docker network create "$NETWORK_NAME"
        print_success "Network $NETWORK_NAME created"
    fi
}

# Function to start Kafka services
start_kafka() {
    print_status "Starting Kafka services..."
    
    if [ ! -f "$COMPOSE_FILE" ]; then
        print_error "Docker Compose file not found: $COMPOSE_FILE"
        exit 1
    fi
    
    docker-compose -f "$COMPOSE_FILE" up -d
    print_success "Kafka services started"
}

# Function to wait for Kafka to be ready
wait_for_kafka() {
    print_status "Waiting for Kafka to be ready..."
    
    local max_attempts=30
    local attempt=1
    
    while [ $attempt -le $max_attempts ]; do
        if docker exec food-delivery-kafka kafka-broker-api-versions --bootstrap-server localhost:9092 > /dev/null 2>&1; then
            print_success "Kafka is ready!"
            return 0
        fi
        
        print_status "Attempt $attempt/$max_attempts - Kafka not ready yet, waiting..."
        sleep 5
        ((attempt++))
    done
    
    print_error "Kafka failed to start within expected time"
    return 1
}

# Function to create topics
create_topics() {
    print_status "Creating Kafka topics..."
    
    local topics=(
        "dev.order-events.order.created:6:1"
        "dev.order-events.order.state_changed:6:1"
        "dev.order-events.order.payment_processed:3:1"
        "dev.order-events.order.shipper_assigned:3:1"
        "dev.order-events.order.cancelled:3:1"
        "dev.order-events.order.delivered:6:1"
        "dev.order-events.order.inventory_reserved:3:1"
        "dev.order-events.order.inventory_released:3:1"
    )
    
    for topic_config in "${topics[@]}"; do
        IFS=':' read -r topic_name partitions replication_factor <<< "$topic_config"
        
        print_status "Creating topic: $topic_name"
        
        if docker exec food-delivery-kafka kafka-topics --bootstrap-server localhost:9092 --list | grep -q "^$topic_name$"; then
            print_warning "Topic $topic_name already exists"
        else
            docker exec food-delivery-kafka kafka-topics \
                --bootstrap-server localhost:9092 \
                --create \
                --topic "$topic_name" \
                --partitions "$partitions" \
                --replication-factor "$replication_factor" \
                --config retention.ms=604800000 \
                --config compression.type=snappy
            
            print_success "Topic $topic_name created"
        fi
    done
}

# Function to verify setup
verify_setup() {
    print_status "Verifying Kafka setup..."
    
    # Check if all services are running
    local services=("food-delivery-zookeeper" "food-delivery-kafka" "food-delivery-kafka-ui")
    
    for service in "${services[@]}"; do
        if docker ps --format "table {{.Names}}" | grep -q "$service"; then
            print_success "$service is running"
        else
            print_error "$service is not running"
            return 1
        fi
    done
    
    # List topics
    print_status "Available topics:"
    docker exec food-delivery-kafka kafka-topics --bootstrap-server localhost:9092 --list
    
    # Show cluster info
    print_status "Kafka cluster info:"
    docker exec food-delivery-kafka kafka-broker-api-versions --bootstrap-server localhost:9092 | head -5
}

# Function to show access information
show_access_info() {
    print_success "Kafka setup completed successfully!"
    echo ""
    echo "Access Information:"
    echo "=================="
    echo "Kafka Broker:     localhost:9092"
    echo "Kafka UI:         http://localhost:8080"
    echo "Kafka Manager:    http://localhost:9000"
    echo "Schema Registry:  http://localhost:8081"
    echo "Kafka Connect:    http://localhost:8083"
    echo "Metrics:          http://localhost:9308/metrics"
    echo ""
    echo "Useful Commands:"
    echo "==============="
    echo "List topics:      docker exec food-delivery-kafka kafka-topics --bootstrap-server localhost:9092 --list"
    echo "Describe topic:   docker exec food-delivery-kafka kafka-topics --bootstrap-server localhost:9092 --describe --topic <topic-name>"
    echo "Console producer: docker exec -it food-delivery-kafka kafka-console-producer --bootstrap-server localhost:9092 --topic <topic-name>"
    echo "Console consumer: docker exec -it food-delivery-kafka kafka-console-consumer --bootstrap-server localhost:9092 --topic <topic-name> --from-beginning"
    echo ""
    echo "Stop Kafka:       docker-compose -f $COMPOSE_FILE down"
    echo "Stop and clean:   docker-compose -f $COMPOSE_FILE down -v"
}

# Function to stop Kafka services
stop_kafka() {
    print_status "Stopping Kafka services..."
    docker-compose -f "$COMPOSE_FILE" down
    print_success "Kafka services stopped"
}

# Function to clean up (stop and remove volumes)
cleanup_kafka() {
    print_status "Cleaning up Kafka services and data..."
    docker-compose -f "$COMPOSE_FILE" down -v
    print_success "Kafka services and data cleaned up"
}

# Function to show logs
show_logs() {
    local service=${1:-""}
    if [ -n "$service" ]; then
        print_status "Showing logs for $service..."
        docker-compose -f "$COMPOSE_FILE" logs -f "$service"
    else
        print_status "Showing logs for all services..."
        docker-compose -f "$COMPOSE_FILE" logs -f
    fi
}

# Function to show help
show_help() {
    echo "Kafka Setup Script for Food Delivery Service"
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  start     Start Kafka services (default)"
    echo "  stop      Stop Kafka services"
    echo "  restart   Restart Kafka services"
    echo "  cleanup   Stop services and remove data volumes"
    echo "  status    Show status of Kafka services"
    echo "  logs      Show logs for all services"
    echo "  logs <service>  Show logs for specific service"
    echo "  topics    List all topics"
    echo "  help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 start"
    echo "  $0 logs kafka"
    echo "  $0 cleanup"
}

# Function to show status
show_status() {
    print_status "Kafka services status:"
    docker-compose -f "$COMPOSE_FILE" ps
}

# Function to list topics
list_topics() {
    print_status "Kafka topics:"
    if docker ps --format "table {{.Names}}" | grep -q "food-delivery-kafka"; then
        docker exec food-delivery-kafka kafka-topics --bootstrap-server localhost:9092 --list
    else
        print_error "Kafka is not running"
    fi
}

# Main function
main() {
    local command=${1:-"start"}
    
    case "$command" in
        "start")
            check_docker
            check_docker_compose
            create_network
            start_kafka
            wait_for_kafka
            create_topics
            verify_setup
            show_access_info
            ;;
        "stop")
            stop_kafka
            ;;
        "restart")
            stop_kafka
            sleep 2
            main "start"
            ;;
        "cleanup")
            cleanup_kafka
            ;;
        "status")
            show_status
            ;;
        "logs")
            show_logs "$2"
            ;;
        "topics")
            list_topics
            ;;
        "help"|"-h"|"--help")
            show_help
            ;;
        *)
            print_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Run main function with all arguments
main "$@"
