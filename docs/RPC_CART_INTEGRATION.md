# RPC Cart Integration Implementation

## Overview
This document describes the implementation of RPC integration between the Order service and Cart service for the `FindCartItemsByCartID` functionality and related cart operations.

## Architecture

### RPC Client Layer
- **Location**: `shared/infras/rpc/cart_rpc_client.go`
- **Purpose**: Handles HTTP communication with the cart service
- **Methods**:
  - `FindCartItemsByCartID()` - Retrieves cart items by cart ID and user ID
  - `UpdateCartStatus()` - Updates cart status (e.g., mark as processed)
  - `ValidateCartOwnership()` - Validates cart belongs to user
  - `GetCartSummary()` - Gets cart summary information

### Cart Service Endpoints
- **Location**: `modules/cart/infras/controller/http-gin/cart_items_by_cart_id_api.go`
- **New Endpoints**:
  - `GET /carts/items?cartId=xxx&userId=xxx` - Get cart items by cart ID
  - `PATCH /carts/{cartId}/status` - Update cart status
  - `GET /carts/{cartId}/validate?userId=xxx` - Validate cart ownership
  - `GET /carts/{cartId}/summary?userId=xxx` - Get cart summary

### Service Layer
- **Location**: `modules/cart/service/`
- **New Services**:
  - `get_cart_items_by_cart_id.go` - Service for retrieving cart items by cart ID
  - `update_cart_status.go` - Service for updating cart status (extended)
  - `validate_cart_ownership.go` - Service for validating cart ownership
  - `get_cart_summary.go` - Service for getting cart summary

### Repository Layer
- **Location**: `modules/cart/infras/repository/gorm-mysql/`
- **New Methods**:
  - `FindCartItemsByCartID()` - Database query for cart items by cart ID
  - `UpdateCartStatusByCartID()` - Update cart status by cart ID
  - `FindCartByCartIDAndUserID()` - Find specific cart by ID and user
  - `GetCartSummaryByCartID()` - Get aggregated cart summary

## Configuration

### Environment Variables
Added `CART_SERVICE_URL` to the configuration:

```bash
export CART_SERVICE_URL="http://localhost:3001/v1"
```

### Config Structure
Updated `shared/datatype/config.go`:
```go
type Config struct {
    // ... existing fields
    CartServiceURL string
    // ... other fields
}
```

## Data Transfer Objects (DTOs)

### Cart Item RPC DTO
```go
type CartItemRPCDto struct {
    ID           uuid.UUID `json:"id"`
    UserID       uuid.UUID `json:"userId"`
    FoodID       uuid.UUID `json:"foodId"`
    RestaurantID uuid.UUID `json:"restaurantId"`
    Quantity     int       `json:"quantity"`
    Status       string    `json:"status"`
    CreatedAt    string    `json:"createdAt"`
    UpdatedAt    string    `json:"updatedAt"`
}
```

### Cart Summary RPC DTO
```go
type CartSummaryRPCDto struct {
    CartID       uuid.UUID `json:"cartId"`
    UserID       uuid.UUID `json:"userId"`
    RestaurantID uuid.UUID `json:"restaurantId"`
    ItemCount    int       `json:"itemCount"`
    TotalPrice   float64   `json:"totalPrice"`
    Status       string    `json:"status"`
    CreatedAt    string    `json:"createdAt"`
    UpdatedAt    string    `json:"updatedAt"`
}
```

## Order Service Integration

### RPC Adapter
- **Location**: `modules/order/service/cart_conversion_rpc_adapter.go`
- **Purpose**: Adapts the cart RPC client to the order service's cart conversion interface
- **Methods**:
  - `ValidateCartForOrder()` - Validates cart can be used for order
  - `ConvertCartToOrderData()` - Converts cart items to order data
  - `MarkCartAsProcessed()` - Marks cart as processed after order creation

### Module Integration
- **Location**: `modules/order/module.go`
- **Status**: Prepared for integration (currently commented out)
- **Usage**: Will be integrated when full order flow is implemented

## API Endpoints

### Cart Service RPC Endpoints

#### 1. Get Cart Items by Cart ID
```http
GET /v1/carts/items?cartId={cartId}&userId={userId}
```

**Response:**
```json
{
  "data": [
    {
      "id": "cart-item-id",
      "userId": "user-id",
      "foodId": "food-id",
      "restaurantId": "restaurant-id",
      "quantity": 2,
      "status": "ACTIVE",
      "createdAt": "2023-01-01T00:00:00Z",
      "updatedAt": "2023-01-01T00:00:00Z"
    }
  ]
}
```

#### 2. Update Cart Status
```http
PATCH /v1/carts/{cartId}/status
Content-Type: application/json

{
  "status": "PROCESSED"
}
```

#### 3. Validate Cart Ownership
```http
GET /v1/carts/{cartId}/validate?userId={userId}
```

**Response:**
```json
{
  "data": {
    "valid": true,
    "reason": ""
  }
}
```

#### 4. Get Cart Summary
```http
GET /v1/carts/{cartId}/summary?userId={userId}
```

**Response:**
```json
{
  "data": {
    "cartId": "cart-id",
    "userId": "user-id",
    "restaurantId": "restaurant-id",
    "itemCount": 3,
    "totalPrice": 150000.0,
    "status": "ACTIVE",
    "createdAt": "2023-01-01T00:00:00Z",
    "updatedAt": "2023-01-01T00:00:00Z"
  }
}
```

## Error Handling

### RPC Client Errors
- **Configuration Error**: Cart service URL not configured
- **HTTP Errors**: 4xx/5xx responses from cart service
- **Network Errors**: Connection timeouts, network failures

### Service Layer Errors
- **Validation Errors**: Invalid cart ID, user ID, or status values
- **Business Logic Errors**: Cart not found, cart already processed, etc.
- **Database Errors**: Connection issues, query failures

## Testing

### Unit Tests
- **RPC Client**: `shared/infras/rpc/cart_rpc_client_test.go`
- **Service Layer**: Individual service test files
- **Repository Layer**: Repository method tests

### Integration Tests
- **End-to-End**: Cart service → RPC client → Order service
- **API Tests**: Direct HTTP calls to cart service endpoints

### Test Commands
```bash
# Run RPC client tests
go test ./shared/infras/rpc/ -v

# Run cart service tests
go test ./modules/cart/service/ -v

# Run cart repository tests
go test ./modules/cart/infras/repository/gorm-mysql/ -v
```

## Usage Example

### From Order Service
```go
// Initialize RPC client
cartRPCClient := sharerpc.NewCartRPCClient("http://localhost:3001/v1")

// Get cart items
cartItems, err := cartRPCClient.FindCartItemsByCartID(ctx, cartID, userID)
if err != nil {
    return err
}

// Process cart items...

// Mark cart as processed
err = cartRPCClient.UpdateCartStatus(ctx, cartID, "PROCESSED")
if err != nil {
    return err
}
```

## Future Enhancements

### 1. Authentication
- Add JWT token passing between services
- Implement service-to-service authentication

### 2. Circuit Breaker
- Add circuit breaker pattern for RPC calls
- Implement fallback mechanisms

### 3. Caching
- Cache cart data to reduce RPC calls
- Implement cache invalidation strategies

### 4. Monitoring
- Add metrics for RPC call success/failure rates
- Implement distributed tracing

### 5. Load Balancing
- Support multiple cart service instances
- Implement client-side load balancing

## Deployment Considerations

### Environment Setup
```bash
# Cart service
export CART_SERVICE_URL="http://cart-service:3001/v1"

# Order service
export CART_SERVICE_URL="http://cart-service:3001/v1"
```

### Docker Compose
```yaml
services:
  cart-service:
    # ... cart service configuration
    ports:
      - "3001:3000"
  
  order-service:
    # ... order service configuration
    environment:
      - CART_SERVICE_URL=http://cart-service:3000/v1
    depends_on:
      - cart-service
```

### Kubernetes
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: order-service-config
data:
  CART_SERVICE_URL: "http://cart-service.default.svc.cluster.local:3000/v1"
```

## Status

### ✅ Completed
- RPC client implementation
- Cart service endpoints
- Service layer implementation
- Repository layer implementation
- Basic testing structure
- Configuration setup

### 🔄 In Progress
- Full integration with order service
- Comprehensive error handling
- Performance optimization

### 📋 TODO
- Authentication implementation
- Circuit breaker pattern
- Comprehensive integration tests
- Performance benchmarks
- Documentation updates
