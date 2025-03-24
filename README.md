# go-food-delivery-backend-service
# food-delivery-backend-service
**# UI Analysis of Food Delivery App**

## **Overview**
The food delivery app's UI is designed with a dark theme and orange highlights, providing a modern and sleek look. The UI flows smoothly across multiple functionalities, including onboarding, authentication, menu browsing, order tracking, and payment processing.

## **Key UI Components**
### **Onboarding Screens:**
- Welcome screen introducing the app.
- User guidance on menu browsing, ordering, and delivery options.

### **Authentication Screens:**
- Sign-up and login with email/password and social authentication.
- Verification and password reset functionality.

### **Home & Menu Screens:**
- Restaurant listings with filtering and sorting options.
- Food details with price, ingredients, and add-to-cart functionality.

### **Cart & Checkout:**
- Order summary and payment gateway integration.
- Multiple payment methods including card scanning.

### **Order Tracking & Delivery:**
- Real-time order tracking with map view.
- Estimated delivery time and contact options.

### **User Profile & History:**
- User profile with order history.
- Review and rating system for restaurants and orders.

## **UI Strengths:**
- Well-structured navigation with clear CTA buttons.
- Aesthetic and intuitive design improving user experience.
- Secure and seamless payment options.

## **UI Improvement Areas:**
- Additional filtering options (e.g., dietary restrictions, discounts).
- More personalization features, such as order recommendations.

---

**# User Stories**

## **User Roles:**
1. **Customer** – Places and tracks orders.
2. **Restaurant Owner** – Manages menu and orders.
3. **Delivery Partner** – Delivers orders.
4. **Admin** – Manages platform operations.

### **Customer User Stories:**
1. **As a customer**, I want to sign up and log in so that I can place orders securely.
2. **As a customer**, I want to browse restaurants and menus so that I can choose my meal.
3. **As a customer**, I want to filter restaurants based on cuisine, price, and ratings so that I can find the best options.
4. **As a customer**, I want to add items to my cart so that I can place an order.
5. **As a customer**, I want to make secure payments so that my transaction is safe.
6. **As a customer**, I want to track my order in real-time so that I know when it will arrive.
7. **As a customer**, I want to leave a review so that I can share my experience.

### **Restaurant Owner User Stories:**
1. **As a restaurant owner**, I want to list my menu so that customers can place orders.
2. **As a restaurant owner**, I want to receive and manage orders so that I can prepare them on time.

### **Delivery Partner User Stories:**
1. **As a delivery partner**, I want to accept delivery requests so that I can deliver food.
2. **As a delivery partner**, I want to update order statuses so that customers stay informed.

### **Admin User Stories:**
1. **As an admin**, I want to manage user accounts so that only valid users can access the app.
2. **As an admin**, I want to oversee transactions to ensure secure payments.

---

**# Database Design (ERD)**

## **Entities & Relationships**

1. **User (UserID, Name, Email, Phone, Role, Password)**
   - Roles: Customer, Restaurant Owner, Delivery Partner, Admin.
2. **Restaurant (RestaurantID, Name, Location, Rating, OwnerID)**
   - A restaurant is managed by one owner.
3. **MenuItem (ItemID, RestaurantID, Name, Price, Category, Description, Availability)**
   - A restaurant has multiple menu items.
4. **Order (OrderID, UserID, RestaurantID, TotalAmount, Status, OrderTime, DeliveryTime)**
   - A user places multiple orders.
5. **OrderDetails (OrderDetailID, OrderID, ItemID, Quantity, SubTotal)**
   - An order consists of multiple order details.
6. **Payment (PaymentID, OrderID, UserID, PaymentMethod, PaymentStatus, Amount, Timestamp)**
   - Each order has a payment transaction.
7. **Delivery (DeliveryID, OrderID, DeliveryPartnerID, Status, Location, ETA)**
   - A delivery partner is assigned to an order.
8. **Review (ReviewID, UserID, RestaurantID, Rating, Comment, Timestamp)**
   - A user can review a restaurant.

## **Relationships:**
- One **User** can place multiple **Orders**.
- One **Order** can contain multiple **OrderDetails**.
- One **Restaurant** has multiple **MenuItems**.
- One **Delivery Partner** delivers multiple **Orders**.
- One **Order** has a single **Payment** record.

## Micro services
1. User Service
Responsibilities:

User registration and authentication

Profile management

Endpoints:

POST /register – Register a new user

POST /login – Authenticate user

GET /users/{id} – Get user details

Database: Users table

Technology: JWT for authentication

2. Restaurant Service
Responsibilities:

Managing restaurants and menus

Restaurant listings and filtering

Endpoints:

GET /restaurants – List all restaurants

POST /restaurants – Add a new restaurant

GET /restaurants/{id} – Get restaurant details

Database: Restaurants and MenuItems tables

3. Order Service
Responsibilities:

Handling customer orders

Order status updates

Order history tracking

Endpoints:

POST /order – Place an order

GET /order/{id} – Get order details

PATCH /order/{id} – Update order status

Database: Orders and OrderDetails tables

4. Payment Service
Responsibilities:

Secure payment processing

Payment verification

Endpoints:

POST /payment – Process payment

GET /payment/{order_id} – Get payment status

Database: Payments table

Integration: Connect with third-party payment gateways (e.g., Stripe, PayPal)

5. Delivery Service
Responsibilities:

Assigning orders to delivery partners

Tracking delivery status

Endpoints:

POST /assign-delivery – Assign a delivery request

GET /track-delivery/{order_id} – Track order location

Database: Deliveries table

Integration: Google Maps API for tracking
