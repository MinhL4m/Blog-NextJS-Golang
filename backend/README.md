# Blogs backend

## Technologies

Language: Golang

### Dependencies

- gopkg.in/go-playground/assert.v1: For testing. [Doc API](https://pkg.go.dev/gopkg.in/go-playground/assert.v1)
- github.com/jinzhu/gorm: ORM package for go

## Learnt

### GORM for model

- [GORM DOC](https://gorm.io/docs/models.html)

### Time in model

- Use `time.Time` for data type

### Prevent csrf

- For each model, have a `Prepare()` method that will run before save into db. In this method, use `html.EscapeString(strings.TrimSpace(<the field>))` to escape everything. Visit any model to see how it got done.

### Full Text Search for post title

- Implement all 3 db queries: [](https://medium.com/@bencagri implementing-multi-table-full-text-search-with-gorm-632518257d15)

### Take only one for update and delete

- When query, the db will return a list or a single record. Even query with uid, still need to make sure only grab one record. `Take()` will make sure only get the first return record.

### JWT Claims

- [JSON Web Token Claims](https://auth0.com/docs/tokens/json-web-tokens/json-web-token-claims)

### JWT Structure

- set of claims are stored in jws payload: [JWT Structure](https://auth0.com/docs/tokens/json-web-tokens/json-web-token-structure)

```go
// Set what cryptographic algorithms used + claims -> header and payload
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// Signed -> Signature
return token.SignedString([]byte(os.Getenv("APP_API_SECRET")))
```

### Type Assertions

- This application use type assertion to check what cryptographic algorithms used. Also used in other places

```go
token.Method.(*jwt.SigningMethodHMAC);
```

### Middleware

- Is just a HOC of React.

### Location Response Header

- https://www.geeksforgeeks.org/http-headers-location/

- Used in 2 circumstances:

  - ask a browser to redirect a URL(status code 3xx)
  - provide information about the location of a newly created resource (status code of 201 http.StatusCreated)
