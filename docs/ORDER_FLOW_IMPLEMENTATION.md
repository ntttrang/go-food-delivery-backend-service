# Order Food Flow Implementation

## Overview
This document describes the implementation of the complete order food flow based on the provided sequence diagram. The implementation follows hexagonal architecture principles and integrates with existing cart, payment, restaurant, and food modules.

## Architecture

### Core Services Implemented

#### 1. Cart to Order Conversion Service (`convert_cart_to_order.go`)
- **Purpose**: Converts cart items to order data
- **Key Features**:
  - Validates cart items are from the same restaurant
  - Checks cart is not already processed
  - Converts cart items to order details with food information
  - Calculates total price
  - Marks cart as processed after successful order creation

#### 2. Payment Processing Service (`process_payment.go`)
- **Purpose**: Handles payment validation and processing
- **Key Features**:
  - Validates payment methods (cash/card)
  - Verifies card ownership and status
  - Processes payments through gateway (simulated)
  - Returns payment status based on method

#### 3. Inventory Checking Service (`check_inventory.go`)
- **Purpose**: Validates restaurant and food availability
- **Key Features**:
  - Checks restaurant is active and accepting orders
  - Validates food items are available
  - Verifies sufficient quantity
  - Ensures all items belong to the same restaurant

#### 4. Order State Management Service (`manage_order_state.go`)
- **Purpose**: Manages order state transitions
- **Key Features**:
  - Validates state transitions (waiting_for_shipper → preparing → on_the_way → delivered)
  - Handles shipper assignment
  - Updates payment status
  - Sends notifications for state changes

#### 5. Notification Service (`notification_service.go`)
- **Purpose**: Handles order-related notifications
- **Key Features**:
  - Order creation notifications
  - State change notifications
  - Shipper assignment notifications
  - Payment status change notifications
  - Extensible for email, SMS, and push notifications

### Enhanced Order Creation

#### Original Order Creation (ADMIN Only)
- **Access Control**: Restricted to ADMIN role only
- Enhanced with payment method validation
- Added card ID requirement for card payments
- Added delivery address validation
- **Use Case**: Administrative order creation, support scenarios

#### New Order Creation from Cart
- **Endpoint**: `POST /v1/orders/from-cart`
- **Flow**:
  1. Validate cart and convert to order data
  2. Check inventory availability
  3. Validate payment method
  4. Create order
  5. Process payment
  6. Mark cart as processed
  7. Send notifications

### New API Endpoints

#### Order Creation from Cart
```
POST /v1/orders/from-cart
Authorization: Bearer <token>
Content-Type: application/json

{
  "cartId": "uuid",
  "deliveryAddress": {
    "cityId": 1,
    "cityName": "Ho Chi Minh City",
    "addr": "123 Main St",
    "lat": 10.762622,
    "lng": 106.660172
  },
  "paymentMethod": "card",
  "cardId": "uuid" // required for card payments
}
```

#### Order State Management
```
PATCH /v1/orders/{id}/state
Authorization: Bearer <token>
Content-Type: application/json

{
  "newState": "preparing",
  "shipperId": "uuid", // optional
  "paymentStatus": "paid" // optional
}
```

#### Shipper Assignment
```
PATCH /v1/orders/{id}/assign-shipper
Authorization: Bearer <token>
Content-Type: application/json

{
  "shipperId": "uuid"
}
```

#### Payment Status Update
```
PATCH /v1/orders/{id}/payment-status
Authorization: Bearer <token>
Content-Type: application/json

{
  "paymentStatus": "paid"
}
```

## Order States

### State Flow
1. **waiting_for_shipper** → preparing, cancel
2. **preparing** → on_the_way, cancel
3. **on_the_way** → delivered, cancel
4. **delivered** (terminal state)
5. **cancel** (terminal state)

### Payment Status
- **pending**: For cash payments or failed card payments
- **paid**: For successful card payments or cash payments upon delivery

## Error Handling

### New Error Types Added
- `ErrPaymentMethodRequired`
- `ErrCardIdRequired`
- `ErrDeliveryAddressRequired`
- `ErrCartIdRequired`
- `ErrCartNotFound`
- `ErrCartEmpty`
- `ErrCartAlreadyProcessed`
- `ErrInvalidPaymentMethod`
- `ErrPaymentFailed`
- `ErrInventoryInsufficient`
- `ErrRestaurantNotAvailable`
- `ErrFoodNotAvailable`
- `ErrInvalidOrderState`
- `ErrShipperRequired`

## Integration Points

### Cart Module Integration
- Reads cart items and food details
- Updates cart status to "PROCESSED"
- Validates cart ownership

### Payment Module Integration
- Validates card ownership and status
- Processes card payments
- Handles payment method validation

### Restaurant Module Integration
- Checks restaurant availability
- Validates restaurant is accepting orders
- Gets restaurant information

### Food Module Integration
- Validates food availability
- Checks food inventory
- Gets food details for order

## Future Enhancements

### Notification System
- Email notifications for order updates
- SMS notifications for delivery updates
- Push notifications for mobile apps
- Restaurant notifications for new orders

### Payment Gateway Integration
- Stripe integration for card payments
- PayPal integration
- Refund processing for cancelled orders

### Inventory Management
- Real-time inventory tracking
- Automatic inventory updates
- Low stock alerts

### Order Tracking
- Real-time order tracking
- GPS integration for delivery tracking
- Estimated delivery time calculations

## Testing Recommendations

### Unit Tests
- Test each service independently
- Mock external dependencies
- Test error scenarios

### Integration Tests
- Test complete order flow
- Test cart to order conversion
- Test payment processing
- Test state transitions

### API Tests
- Test all new endpoints
- Test authentication and authorization
- Test error responses

## Deployment Notes

### Environment Variables
- `USER_SERVICE_URL`: For user authentication
- `FOOD_SERVICE_URL`: For food data
- `PAYMENT_GATEWAY_URL`: For payment processing
- Database connection strings

### Database Migrations
- Ensure order_trackings table has correct ENUM values
- Verify foreign key constraints
- Check indexes for performance

### Monitoring
- Monitor order creation success rates
- Track payment processing times
- Monitor notification delivery
- Alert on failed state transitions

## Conclusion

The order food flow implementation provides a comprehensive solution that:
- Follows the sequence diagram requirements
- Maintains hexagonal architecture principles
- Integrates seamlessly with existing modules
- Provides extensibility for future enhancements
- Includes proper error handling and validation
- Supports both manual order creation and cart-based ordering
