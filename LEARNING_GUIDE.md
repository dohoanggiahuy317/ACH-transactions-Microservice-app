# Simple ACH Microservice – Learning Guide

This project demonstrates how to build a simple Automated Clearing House (ACH) microservice in Go, using PostgreSQL, sqlc, and Gin. The guide below walks you through the implementation flow, from database setup to API development and testing.

---

## 1. Database Schema Initialization

### 1.1. Overview

**Database migration** refers to the process of updating your database schema safely and consistently over time.

- **`up` migration**: Applies new changes to the database schema (e.g., creating tables, modifying columns).
- **`down` migration**: Reverts changes made by the corresponding `up` migration, restoring the previous schema version.

To create migration scripts:
- Write the appropriate SQL commands to define or update the database schema.

### 1.2. Applying Migrations in This Project

1. **Understand the Schema**
   - **Tables**: This project uses three core tables: `account`, `entries`, and `transfer`.
   - **Schema Design**: Refer to the ER diagram located in the `img/` directory.
   - **Migration Scripts**: Initial schema migrations are located in `db/migration/000001_init_schema.up.sql`.

2. **Generate Migration Boilerplate**  
   Use the following command to create a new versioned migration file:
   ```sh
   migrate create -ext sql -dir db/migration -seq init_schema
   ```

3. Add the schema SQL to the generated migration file (`up` for applying, `down` for rollback).




## 2. Spin up the Database with Docker

1. **Pull and run PostgreSQL:**
```sh
docker pull postgres:17-alpine
docker run --name postgres17 -p 55432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine
```

2. **Access the database:**
```sh
docker exec -it postgres17 psql -U root
docker logs postgres17
```

## 3. Generate Query Functions

### 3.1. Overview

Implementing query functions allows the application to perform transactions and interact with the database using Go functions.

### 3.2. How It Works in This Project

This project uses **sqlc** to generate type-safe Go code from raw SQL queries, improving development speed and reducing human error.

1.  Define SQL queries for each table in the `db/query/` directory.
2.  Run `sqlc` to generate Go code based on these SQL files.
3. The generated files (e.g., `account.sql.go`, `entry.sql.go`, `transfer.sql.go`, `models.go`) are saved in the `db/sqlc/` directory.


## 4. Implement the Store Layer

### 4.1 Why We Need the Store Interface

The Store interface serves as a crucial abstraction layer that addresses several important design challenges:

**1. Transaction Management:**
- Database operations like money transfers require multiple SQL queries to execute atomically
- Without proper transaction handling, we risk data inconsistency (e.g., money deducted but not credited)
- The Store interface provides a clean way to wrap multiple operations in a single transaction

**2. Testability:**
- Direct database calls make unit testing difficult and slow
- The Store interface allows us to create mock implementations for testing
- We can test business logic without hitting the actual database

**3. Code Organization:**
- Separates database logic from business logic
- Provides a single point of entry for all database operations
- Makes the codebase more maintainable and easier to understand

### 4.2 How It's Implemented in This Project

**Step 1: Define the Interface**
```go
// db/sqlc/store.go
type Store interface {
    Querier                    // Includes all basic CRUD operations
    TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}
```

**Step 2: Create the Implementation**
```go
type SQLStore struct {
    db *sql.DB
    *Queries
}

func NewStore(db *sql.DB) Store {
    return &SQLStore{
        db:      db,
        Queries: New(db),
    }
}
```

**Step 3: Implement Complex Transactions**
```go
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
    var result TransferTxResult
    
    err := store.execTx(ctx, func(q *Queries) error {
        // 1. Create transfer record
        // 2. Create debit entry
        // 3. Create credit entry  
        // 4. Update account balances
        return nil
    })
    
    return result, err
}
```

1. **SQLStore Implementation:**
   - `SQLStore` struct implements the `Store` interface
   - Contains both the database connection and the generated Queries
   - Provides the concrete implementation for all database operations

**Usage:**
```go
store := NewStore(db)
result, err := store.TransferTx(ctx, TransferTxParams{
    FromAccountID: 1,
    ToAccountID:   2,
    Amount:        100,
})
```

This pattern provides atomic operations, testability, and clean separation of concerns.

## 5. Prevent Deadlocks in Account Updates

- To avoid deadlocks when updating account balances, use `FOR UPDATE` in your SQL queries (see migration scripts).
- This ensures row-level locking and safe concurrent transactions.

1. **Transaction Method (`TransferTx`):**
   - Handles complex money transfer operations
   - Creates transfer record, account entries, and updates balances atomically
   - Uses row-level locking (`FOR UPDATE`) to prevent deadlocks
   - Implements proper error handling and rollback on failure

2. **Deadlock Prevention:**
   - Uses account ID ordering to ensure consistent lock acquisition
   - Prevents circular wait conditions when multiple transfers involve the same accounts

## 6. Testing

### 6.1 Unit Testing Database Layer

Unit testing the database layer ensures that our SQL queries and transaction logic work correctly in isolation. This involves testing individual database operations, complex transactions, and edge cases without external dependencies.

- Create test files in `db/sqlc/` (e.g., `account_test.go`, `transfer_test.go`)
- Test each generated query function with various scenarios


### 6.3 Mock Testing for API Layer

Mock testing allows us to test HTTP endpoints without hitting the actual database. We create mock implementations of our database interface that return predefined responses, enabling fast and reliable API testing.

```sh
mockgen -package mockdb -destination db/mock/store.go github.com/yourusername/simple-ach/db/sqlc Store
```

**Key Testing Patterns:**
- **Mock Expectations**: Define what the mock should return
- **HTTP Testing**: Use `httptest` for endpoint testing
- **Response Validation**: Verify HTTP status codes and response bodies
- **Error Scenarios**: Test both success and failure cases


## 9. Implementing the HTTP API

The HTTP API serves as the interface between clients and the backend service, allowing external applications or users to interact with the system over HTTP. It provides RESTful endpoints that handle HTTP requests and responses, connecting the web layer to our business logic.


1. **Route Definition:**
   - Uses Gin framework for HTTP routing
   - RESTful endpoints for account operations:
     - `POST /accounts` - Create new account
     - `GET /accounts/:id` - Get account by ID
     - `GET /accounts` - List all accounts

2. **Error Handling:**
   - Centralized error response function
   - Consistent JSON error format
   - Proper HTTP status codes

3. **Server Startup:**
   - `Start()` method runs the HTTP server
   - Configurable address parameter
   - Graceful error handling

**Usage Example:**
```go
// Create server with database store
store := NewStore(db)
server := NewServer(store)

// Start server on port 8080
server.Start(":8080")
```
