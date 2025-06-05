# Order Flow Implementation

This document describes the complete order food flow implementation, including order creation, state management, and cancellation functionality

## I. Order Creation

### Original Order Creation (ADMIN Only)

- **Endpoint**: `POST /v1/orders`
- **Access Control**: Restricted to ADMIN role only
- Enhanced with payment method validation
- Added card ID requirement for card payments
- Added delivery address validation
- **Use Case**: Administrative order creation, support scenarios

### New Order Creation from Cart

- **Endpoint**: `POST /v1/orders/from-cart`
- **Flow**:
  1. Validate cart and convert to order data
  2. Check inventory availability
  3. Validate payment method
  4. Create order
  5. Process payment
  6. Mark cart as processed
  7. Send notifications

### Order States and Flow

#### State Flow

```
waiting_for_shipper → preparing → on_the_way → delivered
        ↓               ↓           ↓
      cancel         cancel      cancel
```

1. **waiting_for_shipper** → preparing, cancel
2. **preparing** → on_the_way, cancel
3. **on_the_way** → delivered, cancel
4. **delivered** (terminal state)
5. **cancel** (terminal state)

#### Payment Status

- **pending**: For cash payments or failed card payments
- **paid**: For successful card payments or cash payments upon delivery

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

## II. Order State Management API

**Endpoint**: `PATCH /v1/orders/{id}/state`

#### Request Format

```json
{
  "newState": "string",           // Required: target state
  "shipperId": "string",          // Optional: assign shipper
  "paymentStatus": "string",      // Optional: update payment
  "cancellationReason": "string"  // Required for cancellations
}
```

#### Response Format

```json
{
  "data": {
    "orderId": "order-123",
    "message": "Order state updated successfully"
  }
}
```

### Usage Examples

#### 1. Simple State Transition

Move order from "waiting_for_shipper" to "preparing":

```http
PATCH /v1/orders/order-123/state
Authorization: Bearer <token>
Content-Type: application/json

{
  "newState": "preparing"
}
```

#### 2. State Transition + Shipper Assignment

Move to "preparing" and assign shipper in one call:

```http
PATCH /v1/orders/order-123/state
Authorization: Bearer <token>
Content-Type: application/json

{
  "newState": "preparing",
  "shipperId": "shipper-456"
}
```

#### 3. Update Payment Status Only

Update payment without changing state:

```http
PATCH /v1/orders/order-123/state
Authorization: Bearer <token>
Content-Type: application/json

{
  "newState": "waiting_for_shipper",
  "paymentStatus": "paid"
}
```

#### 4. Complex State Change

Move to "on_the_way" with shipper and payment update:

```http
PATCH /v1/orders/order-123/state
Authorization: Bearer <token>
Content-Type: application/json

{
  "newState": "on_the_way",
  "shipperId": "shipper-789",
  "paymentStatus": "paid"
}
```

## III. Order Cancellation

Order cancellation functionality with proper refund processing, inventory restoration, and multi-party notifications.

Cancellation Logic Flow

1. **Validation**: Check cancellation reason is provided
2. **Refund Processing**: Process refunds for paid orders
3. **Inventory Restoration**: Return items to inventory
4. **Shipper Management**: Clear shipper assignment
5. **Timestamp Tracking**: Record cancellation time
6. **Database Update**: Save all changes in transaction
7. **Notifications**: Send notifications to all parties

### Cancellation Features

#### 1. Cancellation Validation

- **Cancellation Reason Required**: All cancellations must include a reason
- **State Transition Validation**: Orders can only be cancelled from valid states:
  - `waiting_for_shipper` → `cancel`
  - `preparing` → `cancel`
  - `on_the_way` → `cancel`

#### 2. Refund Processing

- **Automatic Refund Detection**: Paid orders automatically trigger refund processing
- **Payment Method Support**:
  - **Cash**: No refund needed (payment not yet collected)
  - **Card**: Processes refund through payment gateway
- **Refund Status Tracking**: Updates payment status to indicate refund processing

#### 3. Inventory Restoration

- **Automatic Restoration**: Returns cancelled items back to inventory
- **Item-by-Item Processing**: Handles each order item individually
- **Failure Resilience**: Continues processing even if some items fail

#### 4. Shipper Management

- **Assignment Clearing**: Removes shipper assignment when order is cancelled
- **Notification**: Notifies assigned shipper about cancellation

#### 5. Comprehensive Notifications

- **Customer**: Cancellation confirmation with refund information
- **Restaurant**: Order cancellation notification
- **Shipper**: Cancellation notification (if assigned)

### Cancel Order Request

```http
PATCH /v1/orders/{id}/state
Authorization: Bearer <token>
Content-Type: application/json

{
  "newState": "cancel",
  "cancellationReason": "Customer requested cancellation"
}
```

## Future Enhancements

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

### Refund Service

- Integration with Stripe, PayPal, and other payment gateways
- Partial refund support
- Refund retry mechanisms
- Refund status tracking and reporting

### Inventory Service

- Real-time inventory updates via food service API
- Inventory reservation and release
- Inventory audit trails
- Cache invalidation for inventory changes

### Notification Service

- SMS and push notification implementation
- Notification templates and localization
- Delivery confirmation and retry logic
- Notification preferences management
