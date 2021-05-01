# dqlx

dqlx is a fully featured [DGraph](https://github.com/dgraph-io/dgraph) Schema and Query Builder for Go.
It aims to simplify the interaction with the awesome Dgraph database allowing you to fluently compose any queries and mutations of any complexity. It also comes with a rich Schema builder to easily develop and maintain your Dgraph schema.

---

### Status
The project is getting close to its first official release v1.0

### Why?
The DGraph query language is awesome! it is really powerful, and you can achieve a lot with it.
However, as you start trying to add dynamicity (like any other declarative query language) you soon starts
fiddling with a lot strings concatenations and can quickly get messy.

dqlx tries to simplify the interaction with DGraph by helping to construct Queries and mutations with a fluent API.

### Features

- [x] Schema Builder (Types, Predicates, Indexes)
- [x] Filtering - Connecting Filters (AND / OR)
- [x] Nested Selection / Filters
- [x] Functions
- [x] Pagination
- [x] Aggregation
- [x] Sorting
- [x] GroupBy
- [x] Multiple Query Block
- [x] Query Variables
- [x] Values Variables
- [x] Facets
- [x] Mutations

## Documentation

You can find the documentation here: https://fenos.github.io/dqlx

---

### Installation
```bash
go get github.com/fenos/dqlx
```

### Quick Overview

```go
func main() {
    // Connect to Dgraph cluster
    db, err := dqlx.Connect("localhost:9080")

    if err != nil {
        log.Fatal()
    }

    ctx := context.Background()

    var animals []map[string]interface{}

    // Query for animals
    _, err = db.
        QueryType("Animal").
        Filter(
            dqlx.Eq{"species": "Cat"},
            dqlx.Lt{"age": 5},
        ).
        Fields(`
            uid
            name
            species
            age
        `).
        UnmarshalInto(&animals).
        Execute(ctx)

    if err != nil { panic(err) }

    println(fmt.Sprintf("The animals are: %v", animals))
}
```

### Licence
MIT
