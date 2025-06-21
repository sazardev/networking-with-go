# Handling JSON and XML over HTTP: Concepts and Go Implementation 📄

> "Imagine you’re sending a package (data) through the mail (HTTP). JSON and XML are like the forms you fill out to describe what’s inside—structured, standardized, and understood by everyone."

---

## 🚦 Why Use JSON and XML in HTTP?
- **Data Interchange:** JSON and XML are the most common formats for exchanging data between clients and servers.
- **Human-Readable:** Both formats are text-based and easy to debug.
- **Widely Supported:** Every major programming language can read and write JSON and XML.
- **Analogy:** Like using a customs form—everyone knows how to read it, no matter where it’s from.

---

## 🧩 JSON vs XML: Quick Comparison
| Feature         | JSON                | XML                 |
|----------------|---------------------|---------------------|
| Syntax         | { "key": "value" } | <key>value</key>    |
| Readability    | Very high           | Medium              |
| Verbosity      | Low                 | High                |
| Data Types     | Native (numbers, etc.) | All as text      |
| Go Support     | encoding/json       | encoding/xml        |

---

## 🛠️ Go in Action: Serving and Consuming JSON (with Nested Data)

Let’s build a server that returns a list of users (with nested fields) as JSON, and a client that fetches and posts users.

### Server Example: Serve and Accept JSON

See: [`json_server.go`](../../exercises/part2/11-json-xml-examples/json_server.go)

```go
// ...see exercises/part2/11-json-xml-examples/json_server.go for full code and comments...
```

### Client Example: Fetch and Post JSON

See: [`json_client.go`](../../exercises/part2/11-json-xml-examples/json_client.go)

```go
// ...see exercises/part2/11-json-xml-examples/json_client.go for full code and comments...
```

**Key Details:**
- Use `json.NewEncoder(w).Encode(data)` to send JSON from a handler.
- Use `json.NewDecoder(r.Body).Decode(&target)` to parse JSON from a request.
- For nested/child fields, use Go slices or structs (see the `User` struct for examples).
- Always set the `Content-Type: application/json` header.

**Testing:**
- Use `curl` to test GET: `curl http://localhost:8080/users`
- Use `curl` to test POST: `curl -X POST -H "Content-Type: application/json" -d '{"name":"Test","email":"test@example.com"}' http://localhost:8080/users`

---

## 🛠️ Go in Action: Serving and Consuming XML (with Nested Elements)

Let’s build a server that returns a list of products (with nested tags) as XML, and a client that fetches and posts products.

### Server Example: Serve and Accept XML

See: [`xml_server.go`](../../exercises/part2/11-json-xml-examples/xml_server.go)

```go
// ...see exercises/part2/11-json-xml-examples/xml_server.go for full code and comments...
```

### Client Example: Fetch and Post XML

See: [`xml_client.go`](../../exercises/part2/11-json-xml-examples/xml_client.go)

```go
// ...see exercises/part2/11-json-xml-examples/xml_client.go for full code and comments...
```

**Key Details:**
- Use `xml.NewEncoder(w).Encode(data)` to send XML from a handler.
- Use `xml.NewDecoder(r.Body).Decode(&target)` to parse XML from a request.
- For nested/child elements, use Go slices with struct tags (see the `Product` struct and `tags>tag`).
- Always set the `Content-Type: application/xml` header.

**Testing:**
- Use `curl` to test GET: `curl http://localhost:8081/products`
- Use `curl` to test POST: `curl -X POST -H "Content-Type: application/xml" -d '<product><name>Test</name><tags><tag>foo</tag></tags></product>' http://localhost:8081/products`

---

## 🧪 Practical Exercise Files
- [JSON Server Example](../../exercises/part2/11-json-xml-examples/json_server.go)
- [JSON Client Example](../../exercises/part2/11-json-xml-examples/json_client.go)
- [XML Server Example](../../exercises/part2/11-json-xml-examples/xml_server.go)
- [XML Client Example](../../exercises/part2/11-json-xml-examples/xml_client.go)

---

## 🧠 Key Takeaways
- Use JSON for most web APIs—compact, fast, and easy in Go.
- Use XML when you need to interoperate with older systems or standards.
- Always set the correct `Content-Type` header.
- Use Go’s `encoding/json` and `encoding/xml` for safe, automatic encoding/decoding.
- For nested/child data, use slices and struct tags.
- Test with `curl` or Postman for real-world scenarios.

---

[Previous: HTTP: Protocol Theory and Go Implementation](10-http-protocol-theory-and-go-implementation.md) | [Next: WebSockets: Real-Time Communication](12-websockets-real-time-communication.md)
