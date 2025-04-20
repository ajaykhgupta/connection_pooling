# Connection Pooling in Go

When multiple requests hit a system simultaneously—each requiring a separate database connection—it can overwhelm the database server. This often results in an error like:

```
too many clients already
```

This happens when the number of concurrent connections exceeds the database's allowed limit, causing some requests to fail and requiring clients to retry, ultimately leading to poor user experience.

---

## Solution: Connection Pooling

**Connection pooling** is a strategy that helps mitigate this issue. Instead of creating a new database connection for every request, a pool of reusable connections is maintained.

Here’s how it works:

1. A fixed number of database connections (e.g., 5) are pre-created and stored in a queue or channel.
2. When a new request comes in:
   - It attempts to fetch a connection from the pool.
   - If a connection is available, it uses it to perform the database operation.
   - Once done, the connection is returned (requeued) to the pool.
3. If all connections are in use:
   - The request waits until a connection is available instead of failing immediately.

This model ensures:
- Controlled concurrency
- Efficient use of database resources
- Prevention of connection overload errors
- Smooth request handling under load

---

## Goal of This POC

This proof of concept (POC) is implemented in **Go** to simulate a high-concurrency environment. It mimics how a real system would handle 100+ simultaneous database requests using a limited-size connection pool. The goal is to:
- Test system behavior under load
- Demonstrate safe concurrent access to the database
- Avoid overloading the database server


#### Full Disclosure

This README was rewritten with assistance from an LLM based on pointers I provided, to save time and improve clarity.