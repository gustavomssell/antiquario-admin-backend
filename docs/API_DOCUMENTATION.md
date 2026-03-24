# AntiquaOS API Documentation

The REST API is structured under the `/v1` prefix.

## Authentication
- `POST /v1/auth/login`
- `POST /v1/auth/register`

## Core Modules (CRUDs)
Each domain module exposes standard CRUD endpoints.

For example, for the `Category` module (`/v1/categories`):
- `GET /` - Retrieve all records
- `GET /:id` - Retrieve a single record by ID
- `POST /` - Create a new record
- `PUT /:id` - Update an existing record
- `DELETE /:id` - Delete a record
- `GET /search` - Paginated and filtered search

**Available Endpoints:**
- `/v1/users`
- `/v1/products`
- `/v1/categories`, `/v1/periods`, `/v1/styles`, `/v1/tags`, `/v1/materials`
- `/v1/suppliers`, `/v1/customers`, `/v1/acquisitions`, `/v1/sales`, `/v1/payments`
- `/v1/restorations`, `/v1/reservations`, `/v1/appraisals`, `/v1/auctions`

## Real-time WebSockets
- `GET /v1/ws` - Connect to the global WebSocket event hub for real-time notifications and live auction updates.