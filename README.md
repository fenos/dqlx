# Deku
---

Deku is a [DGraph](https://github.com/dgraph-io/dgraph) query builder 🦸 </br>
Compose fluently complex and dynamic DGraph queries & mutations!

---

### Status
The project is currently **🦸 Working Progress 🦸**, the API hasn't fully stabilised just yet.

### Why?
The DGraph query language is awesome! it is really powerful, and you can achieve a lot with it.
However, as you start trying to add dynamicity (like any other declarative query language) you soon starts
fiddling with a lot strings concatenations and can quickly get messy.

Deku tries to simplify the construction of DGraph Queries and mutations with a fluent API.

### Features

- [x] Schema Builder (Types, Predicates, Indexes)
- [x] Filtering - Connecting Filters (AND / OR)
- [x] Nested Selection / Filters
- [ ] Functions 90%
- [x] Pagination
- [x] Aggregation
- [x] Sorting
- [x] GroupBy
- [x] Multiple Query Block
- [x] Query Variables
- [x] Values Variables
- [ ] Facets
- [ ] Mutations

## Getting Started

---

### Installation
```bash
go get github.com/fenos/deku
```

### Quick Overview

Here is how you produce a **super simple** query:
```go
query, variables, err := dql.
    Query("bladerunner", dql.EqFn("item", "value")).
    Fields(`
        uid
        name
        initial_release_date
        netflix_id
    `).
    Filter(dql.Eq{"field1": "value1"}).
    ToDQL()


print(query)
```

Produces
```graphql
query Bladerunner($0:string, $1:string) {
    bladerunner(func: eq(item,$0)) @filter(eq(field1,$1)) {
        uid
        name
        initial_release_date
        netflix_id
    }
}
```

The true power of **Deku** shows when you start getting serious

```go
query, variables, err := dql.
    Query("bladerunner", dql.EqFn("name@en", "Blade Runner")).
    Fields(`
        uid
        name
        initial_release_date
        netflix_id
    `).
    Edge("authors", dql.Fields(`
        uid
        name
        surname
        age
    `), dql.Eq{"age": 20}).
    Edge("actors", dql.Fields(`
        uid
        surname
        age
    `), dql.Gt{"age": []int{18, 20, 30}}).
    Edge("actors->rewards"), dql.Fields(`
        uid
        points
    `), dql.Gt{"points": 3}).
    ToDQL()
```

Produces

```graphql
query Bladerunner($0:string, $1:int, $2:int, $3:int, $4:int, $5:int) {
    bladerunner(func: eq(name@en,$0)) {
        uid
        name
        initial_release_date
        netflix_id
        authors @filter(eq(age,$1)) {
            uid
            name
            surname
            age
        }
        actors @filter(gt(age,[$2,$3,$4])) {
            uid
            surname
            age
            rewards @filter(gt(points,$5)) {
                uid
                points
            }
        }
    }
}
```

Not yet convinced?

Check the [Test Cases](https://github.com/fenos/deku/blob/main/query_test.go) until the full documentation is ready

### Licence
MIT
