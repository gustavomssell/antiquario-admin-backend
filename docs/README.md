# AntiquaOS Admin Backend

AntiquaOS is an advanced SaaS operations and catalog management system tailored for antique shops. This backend project is built in Go (Golang) using Clean Architecture principles, Gin framework for routing, GORM for database operations, and Gorilla WebSockets for real-time capabilities.

## Domain Modules

- **Auth & User**: Roles, Authentication, JWT Token management.
- **Catalog**: Core taxonomy (Categories, Periods, Styles, Tags, Materials).
- **Product**: Antique item catalog, dimensions, specific characteristics, media.
- **Commercial**: Ecosystem entities (Suppliers, Customers) and transactions (Sales, Acquisitions, Payments).
- **Operations & Restoration**: Specialized domains for Reservations, Appraisals, Auctions (via WebSockets), and Restoration Orders.

## Technology Stack
- **Go 1.24+**
- **Gin Web Framework**
- **GORM (PostgreSQL)**
- **Gorilla WebSocket**
- **Zap Logger**
- **JWT Authentication**

## Running the Application
```bash
go mod tidy
go run main.go
```