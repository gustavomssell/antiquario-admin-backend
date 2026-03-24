# Paginated Search and Filtering

All CRUD modules in AntiquaOS support advanced data pagination and filtering via the `/search` endpoint (e.g., `/v1/products/search`).

## Query Parameters

- `page`: The page number (default: 1)
- `pageSize`: The number of items per page (default: 10)

### Text Filtering
To filter by a textual column, append `_like` to the field name in the query string.
The system automatically wraps the value in `%` for `ILIKE` operations in PostgreSQL.

**Example: Search Products by Title**
```
GET /v1/products/search?page=1&pageSize=20&title_like=vase
```

### JSON Response Structure
The search endpoint responds with a standardized SearchResult container:
```json
{
  "data": [ ... array of entities ... ],
  "total": 150,
  "page": 1,
  "pageSize": 20,
  "totalPages": 8,
  "filters": {
      "page": 1,
      "pageSize": 20,
      "likeFilters": {
          "title": ["vase"]
      }
  }
}
```