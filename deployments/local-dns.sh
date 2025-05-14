#!/bin/bash

# This script provides a local DNS resolution solution without modifying /etc/hosts
# It uses dnsmasq in a Docker container to provide DNS resolution for food-delivery.local

# Exit when any command fails
set -e

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
  echo "Error: Docker is not running. Please start Docker and try again."
  exit 1
fi

# Check if the dnsmasq container is already running
if docker ps | grep -q "food-delivery-dns"; then
  echo "DNS container is already running."
  echo "To restart it, run: docker restart food-delivery-dns"
  echo "To stop it, run: docker stop food-delivery-dns"
  exit 0
fi

# Create a Docker network for DNS if it doesn't exist
if ! docker network ls | grep -q "food-delivery-net"; then
  echo "Creating Docker network for DNS..."
  docker network create food-delivery-net
fi

# Start dnsmasq container
echo "Starting local DNS container..."
docker run --name food-delivery-dns \
  --restart=always \
  -d \
  -p 53:53/udp \
  -p 53:53/tcp \
  --cap-add=NET_ADMIN \
  --network food-delivery-net \
  -v "$(pwd)/deployments/dnsmasq.conf:/etc/dnsmasq.conf" \
  alpine:latest \
  sh -c "apk add --no-cache dnsmasq && echo 'address=/food-delivery.local/127.0.0.1' > /etc/dnsmasq.d/local.conf && dnsmasq -k"

echo ""
echo "Local DNS server is running."
echo "To use it, configure your system to use 127.0.0.1 as your DNS server."
echo ""
echo "For macOS:"
echo "  1. Go to System Preferences > Network"
echo "  2. Select your active network connection"
echo "  3. Click 'Advanced' > 'DNS'"
echo "  4. Add 127.0.0.1 as the first DNS server"
echo ""
echo "For temporary DNS resolution without changing system settings, you can use:"
echo "  - For curl: curl -H 'Host: food-delivery.local' http://localhost"
echo "  - For browser: Use the port-forward.sh script which provides direct URLs"
echo ""
echo "To stop the DNS server: docker stop food-delivery-dns"
echo "To remove the DNS server: docker rm food-delivery-dns"
