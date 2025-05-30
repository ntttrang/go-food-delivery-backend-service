# RPC Direct Repository Implementation

## Overview
Updated the RPC cart integration to call the repository directly instead of going through the service layer, as requested. This simplifies the architecture for RPC operations while maintaining the service layer for regular cart operations.

## Changes Made

### 1. Updated RPC APIs to Call Repository Directly

#### `UpdateCartStatusAPI`
- **Before**: Called `updateCartStatusCmdHdl.Execute()` (service layer)
- **After**: Calls `ctrl.repo.UpdateCartStatusByCartID()` directly
- **Benefits**: Simpler, faster execution for RPC operations

#### `GetCartItemsByCartIDAPI`
- **Before**: Called `getCartItemsByCartIDQueryHdl.Execute()` (service layer)
- **After**: Calls `ctrl.repo.FindCartItemsByCartID()` directly
- **Benefits**: Direct data access without service layer overhead

#### `ValidateCartOwnershipAPI`
- **Before**: Called `validateCartOwnershipQueryHdl.Execute()` (service layer)
- **After**: Calls `ctrl.repo.FindCartByCartIDAndUserID()` directly
- **Benefits**: Simple validation logic directly in controller

#### `GetCartSummaryAPI`
- **Before**: Called `getCartSummaryQueryHdl.Execute()` (service layer)
- **After**: Calls `ctrl.repo.GetCartSummaryByCartID()` directly
- **Benefits**: Direct summary calculation without service layer

### 2. Repository Interface for RPC Operations

Created `ICartRepository` interface in the controller package:

```go
type ICartRepository interface {
    UpdateCartStatusByCartID(ctx context.Context, cartID uuid.UUID, status string) error
    FindCartItemsByCartID(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) ([]CartItem, error)
    FindCartByCartIDAndUserID(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) (*CartItem, error)
    GetCartSummaryByCartID(ctx context.Context, cartID uuid.UUID, userID uuid.UUID) (*CartSummary, error)
}
```

### 3. Repository Adapter

Created `CartRepositoryAdapter` to bridge between:
- **GORM Repository**: Returns `cartmodel.Cart` entities
- **Controller Interface**: Expects `CartItem` and `CartSummary` DTOs

**Key Features:**
- Type conversion between domain models and DTOs
- Date formatting for API responses
- Error handling and validation
- Summary calculation from cart items

### 4. Simplified Controller Constructor

Updated `CartHttpController` constructor:
- **Added**: `repo ICartRepository` parameter
- **Removed**: Service handlers for RPC operations (set to `nil`)
- **Benefits**: Clear separation between regular operations (use services) and RPC operations (use repository directly)

### 5. Module Configuration

Updated cart module setup:
```go
// Create repository adapter for RPC operations
repoAdapter := cartHttpgin.NewCartRepositoryAdapter(repo)

// Setup controller (RPC operations call repository directly)
cartCtl := cartHttpgin.NewCartHttpController(
    createCmdHdl,           // Regular operations use services
    listQueryHdl,
    listCartItemQueryHdl,
    getDetailQueryHdl,
    updateCmdHdl,
    deleteCmdHdl,
    nil,         // RPC operations don't need service handlers
    nil,
    nil,
    nil,
    repoAdapter, // Direct repository access for RPC
)
```

## Architecture Benefits

### 1. **Performance**
- **Reduced Latency**: RPC calls skip service layer processing
- **Fewer Allocations**: Direct data conversion without intermediate DTOs
- **Simpler Call Stack**: Controller → Repository (vs Controller → Service → Repository)

### 2. **Simplicity**
- **Less Code**: No service layer boilerplate for simple CRUD operations
- **Clearer Intent**: RPC operations are clearly marked as direct repository calls
- **Easier Debugging**: Shorter call stack for RPC operations

### 3. **Maintainability**
- **Separation of Concerns**: Regular operations use services, RPC operations use repository directly
- **Flexibility**: Can easily switch between service and repository approaches per operation
- **Clear Documentation**: Comments indicate which operations use which approach

## API Endpoints (Unchanged)

The RPC endpoints remain the same, but now call repository directly:

### 1. Get Cart Items by Cart ID
```http
GET /v1/carts/items?cartId={cartId}&userId={userId}
```

### 2. Update Cart Status
```http
PATCH /v1/carts/{cartId}/status
Content-Type: application/json

{
  "status": "PROCESSED"
}
```

### 3. Validate Cart Ownership
```http
GET /v1/carts/{cartId}/validate?userId={userId}
```

### 4. Get Cart Summary
```http
GET /v1/carts/{cartId}/summary?userId={userId}
```

## Response Formats (Unchanged)

All response formats remain the same, ensuring backward compatibility with existing RPC clients.

## Error Handling

### Repository-Level Errors
- **Database Errors**: Wrapped in `ErrInternalServerError`
- **Not Found Errors**: Handled gracefully with appropriate responses
- **Validation Errors**: Checked at controller level before repository calls

### Controller-Level Validation
- **UUID Parsing**: Invalid UUIDs return `ErrBadRequest`
- **Required Fields**: Missing parameters return `ErrBadRequest`
- **Status Values**: Invalid status values return `ErrBadRequest`

## Testing

### Unit Tests
- **Controller Tests**: Mock the `ICartRepository` interface
- **Adapter Tests**: Test conversion between domain models and DTOs
- **Integration Tests**: Test complete RPC flow

### Example Test Structure
```go
type MockCartRepository struct {
    mock.Mock
}

func (m *MockCartRepository) UpdateCartStatusByCartID(ctx context.Context, cartID uuid.UUID, status string) error {
    args := m.Called(ctx, cartID, status)
    return args.Error(0)
}

func TestUpdateCartStatusAPI_Success(t *testing.T) {
    mockRepo := new(MockCartRepository)
    mockRepo.On("UpdateCartStatusByCartID", mock.Anything, cartID, "PROCESSED").Return(nil)
    
    controller := &CartHttpController{repo: mockRepo}
    // ... test implementation
}
```

## Future Considerations

### 1. **Caching**
- Add caching layer in repository adapter for frequently accessed data
- Cache cart summaries and validation results

### 2. **Metrics**
- Add performance metrics for RPC operations
- Monitor repository call latency vs service call latency

### 3. **Circuit Breaker**
- Add circuit breaker pattern for repository calls
- Implement fallback mechanisms for RPC operations

### 4. **Batch Operations**
- Support batch cart operations for better performance
- Implement bulk status updates

## Migration Guide

### For Existing Code
1. **No Changes Required**: Regular cart operations continue to use service layer
2. **RPC Clients**: No changes needed, API contracts remain the same
3. **Testing**: Update tests to mock `ICartRepository` instead of service handlers

### For New RPC Operations
1. **Add Method**: Add new method to `ICartRepository` interface
2. **Implement Adapter**: Add implementation in `CartRepositoryAdapter`
3. **Add Controller Method**: Create new controller method calling repository directly
4. **Add Route**: Register new route in `SetupRoutes`

## Summary

✅ **Completed:**
- RPC APIs now call repository directly
- Service layer bypassed for RPC operations
- Repository adapter handles type conversions
- All tests pass and module builds successfully

✅ **Benefits Achieved:**
- Improved performance for RPC operations
- Simplified architecture for simple CRUD operations
- Clear separation between regular and RPC operations
- Maintained backward compatibility

✅ **Ready for Production:**
- All RPC endpoints functional
- Error handling implemented
- Type safety maintained
- Documentation updated
