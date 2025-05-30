# Testing Order Flow - Complete Guide

## Overview
This guide provides comprehensive testing strategies for the order food flow implementation, including unit tests, integration tests, and API testing.

## Testing Strategy

### 1. Unit Testing
Test individual services in isolation with mocked dependencies.

### 2. Integration Testing
Test the complete order flow with real database connections.

### 3. API Testing
Test HTTP endpoints with various scenarios.

### 4. Manual Testing
Step-by-step manual testing procedures.

## Prerequisites

### Database Setup
Ensure your database has the required tables:
- `orders`
- `order_trackings`
- `order_details`
- `carts`
- `foods`
- `restaurants`
- `cards` (payment module)

### Environment Variables
```bash
export DB_DSN="user:password@tcp(localhost:3306)/food_delivery?charset=utf8mb4&parseTime=True&loc=Local"
export USER_SERVICE_URL="http://localhost:3001"
export FOOD_SERVICE_URL="http://localhost:3002"
export PORT="3000"
```

## Unit Testing

### Test Structure
```
modules/order/service/
├── convert_cart_to_order_test.go
├── process_payment_test.go
├── check_inventory_test.go
├── manage_order_state_test.go
└── create_order_test.go
```

### Running Unit Tests
```bash
# Run all order service tests
go test ./modules/order/service/... -v

# Run specific test
go test ./modules/order/service/ -run TestCreateOrderFromCart -v

# Run with coverage
go test ./modules/order/service/... -cover -v
```

## Integration Testing

### Database Test Setup
1. Create test database
2. Run migrations
3. Seed test data

### Test Data Requirements
```sql
-- Test restaurant
INSERT INTO restaurants (id, owner_id, name, addr, city_id, lat, lng, shipping_fee_per_km, status, created_at, updated_at)
VALUES ('test-restaurant-id', 'test-owner-id', 'Test Restaurant', '123 Test St', 1, 10.762622, 106.660172, 5000, 'ACTIVE', NOW(), NOW());

-- Test food items
INSERT INTO foods (id, restaurant_id, category_id, name, description, price, images, status, created_at, updated_at)
VALUES
('test-food-1', 'test-restaurant-id', 'test-category-id', 'Test Food 1', 'Delicious test food', 50000, 'image1.jpg', 'ACTIVE', NOW(), NOW()),
('test-food-2', 'test-restaurant-id', 'test-category-id', 'Test Food 2', 'Another test food', 75000, 'image2.jpg', 'ACTIVE', NOW(), NOW());

-- Test user cart
INSERT INTO carts (id, user_id, food_id, restaurant_id, quantity, status, created_at, updated_at)
VALUES
('test-cart-id', 'test-user-id', 'test-food-1', 'test-restaurant-id', 2, 'ACTIVE', NOW(), NOW()),
('test-cart-id', 'test-user-id', 'test-food-2', 'test-restaurant-id', 1, 'ACTIVE', NOW(), NOW());

-- Test payment card
INSERT INTO cards (id, method, provider, cardholder_name, card_number, card_type, expiry_month, expiry_year, cvv, user_id, status, created_at, updated_at)
VALUES ('test-card-id', 'CREDIT_CARD', 'STRIPE', 'Test User', '4111111111111111', 'VISA', '12', '2025', '123', 'test-user-id', 'ACTIVE', NOW(), NOW());
```

## API Testing

### 1. Test Order Creation from Cart

#### Request
```bash
curl -X POST http://localhost:3000/v1/orders/from-cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "cartId": "test-cart-id",
    "deliveryAddress": {
      "cityId": 1,
      "cityName": "Ho Chi Minh City",
      "addr": "123 Delivery St",
      "lat": 10.762622,
      "lng": 106.660172
    },
    "paymentMethod": "card",
    "cardId": "test-card-id"
  }'
```

#### Expected Response
```json
{
  "data": {
    "orderId": "generated-order-id",
    "message": "Order created successfully from cart"
  }
}
```

### 2. Test Order State Management

#### Update Order State
```bash
curl -X PATCH http://localhost:3000/v1/orders/{order-id}/state \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "newState": "preparing",
    "shipperId": "test-shipper-id"
  }'
```

#### Assign Shipper
```bash
curl -X PATCH http://localhost:3000/v1/orders/{order-id}/assign-shipper \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "shipperId": "test-shipper-id"
  }'
```

#### Update Payment Status
```bash
curl -X PATCH http://localhost:3000/v1/orders/{order-id}/payment-status \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "paymentStatus": "paid"
  }'
```

### 3. Test Error Scenarios

#### Test ADMIN Role Restriction for CreateOrderAPI
```bash
# Test with non-admin user (should fail)
curl -X POST http://localhost:3000/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer USER_JWT_TOKEN" \
  -d '{
    "totalPrice": 100000,
    "restaurantId": "test-restaurant-id",
    "deliveryAddress": {
      "cityId": 1,
      "cityName": "Ho Chi Minh City",
      "addr": "123 Test Street"
    },
    "paymentMethod": "cash",
    "orderDetails": [...]
  }'

# Expected Response: 403 Forbidden
# {
#   "error": "only administrators can create orders manually"
# }

# Test with admin user (should succeed)
curl -X POST http://localhost:3000/v1/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_JWT_TOKEN" \
  -d '{...same payload...}'
```

#### Invalid Cart ID
```bash
curl -X POST http://localhost:3000/v1/orders/from-cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "cartId": "invalid-cart-id",
    "deliveryAddress": {...},
    "paymentMethod": "card",
    "cardId": "test-card-id"
  }'
```

#### Missing Payment Method
```bash
curl -X POST http://localhost:3000/v1/orders/from-cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "cartId": "test-cart-id",
    "deliveryAddress": {...}
  }'
```

## Manual Testing Procedures

### Complete Order Flow Test

#### Step 1: Setup Test Data
1. Create test restaurant with active status
2. Add test food items to restaurant
3. Create test user account
4. Add test payment card for user

#### Step 2: Add Items to Cart
```bash
# Add first item to cart
curl -X POST http://localhost:3000/v1/carts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "foodId": "test-food-1",
    "quantity": 2
  }'

# Add second item to cart
curl -X POST http://localhost:3000/v1/carts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "foodId": "test-food-2",
    "quantity": 1
  }'
```

#### Step 3: Verify Cart Contents
```bash
curl -X GET http://localhost:3000/v1/carts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Step 4: Create Order from Cart
```bash
curl -X POST http://localhost:3000/v1/orders/from-cart \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "cartId": "your-cart-id",
    "deliveryAddress": {
      "cityId": 1,
      "cityName": "Ho Chi Minh City",
      "addr": "123 Delivery St",
      "lat": 10.762622,
      "lng": 106.660172
    },
    "paymentMethod": "card",
    "cardId": "your-card-id"
  }'
```

#### Step 5: Verify Order Creation
```bash
curl -X GET http://localhost:3000/v1/orders/{order-id} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Step 6: Test State Transitions
```bash
# Move to preparing
curl -X PATCH http://localhost:3000/v1/orders/{order-id}/state \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"newState": "preparing", "shipperId": "test-shipper-id"}'

# Move to on_the_way
curl -X PATCH http://localhost:3000/v1/orders/{order-id}/state \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"newState": "on_the_way"}'

# Move to delivered
curl -X PATCH http://localhost:3000/v1/orders/{order-id}/state \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"newState": "delivered"}'
```

#### Step 7: Verify Final State
```bash
curl -X GET http://localhost:3000/v1/orders/{order-id} \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Test Scenarios

### Happy Path Scenarios
1. ✅ Create order from cart with card payment
2. ✅ Create order from cart with cash payment
3. ✅ Complete order state transitions
4. ✅ Assign shipper to order
5. ✅ Update payment status

### Error Scenarios
1. ❌ Create order with empty cart
2. ❌ Create order with processed cart
3. ❌ Create order with invalid payment method
4. ❌ Create order with inactive restaurant
5. ❌ Create order with unavailable food items
6. ❌ Invalid state transitions
7. ❌ Assign shipper to delivered order

### Edge Cases
1. 🔄 Create order with mixed restaurant items
2. 🔄 Create order with insufficient inventory
3. 🔄 Payment processing failure
4. 🔄 Network timeout during order creation
5. 🔄 Concurrent cart modifications

## Performance Testing

### Load Testing
```bash
# Install hey for load testing
go install github.com/rakyll/hey@latest

# Test order creation endpoint
hey -n 100 -c 10 -m POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d @order_payload.json \
  http://localhost:3000/v1/orders/from-cart
```

### Database Performance
```sql
-- Check order creation performance
EXPLAIN SELECT * FROM orders WHERE user_id = 'test-user-id';
EXPLAIN SELECT * FROM order_trackings WHERE order_id = 'test-order-id';
EXPLAIN SELECT * FROM order_details WHERE order_id = 'test-order-id';
```

## Monitoring and Debugging

### Application Logs
```bash
# Monitor application logs
tail -f /var/log/food-delivery/app.log

# Filter order-related logs
grep "order" /var/log/food-delivery/app.log
```

### Database Monitoring
```sql
-- Monitor order creation
SELECT COUNT(*) as total_orders,
       DATE(created_at) as order_date
FROM orders
GROUP BY DATE(created_at)
ORDER BY order_date DESC;

-- Check order states
SELECT state, COUNT(*) as count
FROM order_trackings
GROUP BY state;
```

## Troubleshooting Common Issues

### Issue: Order creation fails with "cart not found"
**Solution**: Verify cart ID exists and belongs to the user

### Issue: Payment validation fails
**Solution**: Check card status and ownership

### Issue: Invalid state transition
**Solution**: Verify current order state allows the transition

### Issue: Inventory check fails
**Solution**: Verify restaurant and food items are active

## Next Steps

1. **Implement Repository Interfaces**: Create actual implementations for the service interfaces
2. **Add Comprehensive Unit Tests**: Write tests for all service methods
3. **Set up CI/CD Pipeline**: Automate testing in your deployment pipeline
4. **Add Monitoring**: Implement metrics and alerting for order flow
5. **Performance Optimization**: Optimize database queries and add caching
