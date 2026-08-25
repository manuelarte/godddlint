# 🧐 Go DDD Lint

[![CI](https://github.com/manuelarte/godddlint/actions/workflows/ci.yml/badge.svg)](https://github.com/manuelarte/godddlint/actions/workflows/ci.yml)
![version](https://img.shields.io/github/v/release/manuelarte/godddlint)

Go DDD Lint is an opinionated linter that checks some properties a
[Domain-Driven Design][ddd] (DDD) model should achieve.

Mark your [value object][value-object] structs with `//godddlint:valueObject` and your entities with `//godddlint:entity`
and then run this linter.

## ⬇️  Getting Started

### Run it as a standalone linter

To install it, run:

```bash
go install github.com/manuelarte/godddlint@latest
```

Then, you can run it in your project by executing:

```bash
godddlint ./...
```

### Run it as a module plugin in golangci-lint

You can integrate this linter with [golangci-lint](https://golangci-lint.run/)
by using the [module plugin](https://golangci-lint.run/docs/plugins/module-plugins/).

Example `custom-gcl.yml` with this module plugin:

```yaml
version: v2.12.2
plugins:
  - module: "github.com/manuelarte/godddlint"
    import: "github.com/manuelarte/godddlint/plugin"
    version: latest
```

### Quick demo

Annotate your domain struct and run the linter:

```go
//godddlint:entity
type User struct {
  id   string
  Name string
}
```

```bash
godddlint ./...
```

Expected diagnostics look similar to:

```text
path/to/user.go:3:3: E003: Prefer custom domain types to primitives
path/to/user.go:4:3: E005: Entity's field is exported
```

## 🚀 Features

### Entities

An [entity][entity] is an object defined not by its attributes, but its identity.

#### Entity Rules

Current entity rule codes: `E001`, `E003`, `E004`, `E005`.

##### E001: Pointer Receivers

An `Entity` can mutate, so its methods should use pointer receivers.

```go
//godddlint:entity
type User struct {
 id      UserID
 name    Name
 address Address
}

// entities must have a pointer receiver
func (c *User) ...
```

You can disable this rule at struct level, and also at method level by adding `//godddlint:disable:E001`.

##### E003: Custom Domain Types Over Primitives

An `Entity` field should have domain meaning beyond a primitive type.

<table>
<thead><tr><th>❌ Bad</th><th>✅ Good</th></tr></thead>
<tbody>
<tr><td>

```go
type User struct {
  id      string
  name    string
  address string
}
...
```

</td><td>

```go
type (
  UserID int
  Name string
  Address string
)


//godddlint:entity
type User struct {
  id      UserID
  name    Name
  address Address
}
...
```

</td></tr>

</tbody>
</table>

You can disable this rule at struct level, and also at field level by adding `//godddlint:disable:E003`.

##### E004: Use Custom Domain Errors

Business processes that can return an error should return a meaningful domain error, not a generic one.

<table>
<thead><tr><th>❌ Bad</th><th>✅ Good</th></tr></thead>
<tbody>
<tr><td>

```go
func (c *User) AddAddress(na Address) error {
  if len(c.addresses) >= 2 {
    return errors.New("max number of addresses reached")
  }
  c.addresses = append(c.addresses, na)
  return nil
}
...
```

</td><td>

```go
func (c *User) AddAddress(na Address) error {
  if len(c.addresses) >= 2 {
    return MaxNumberOfAddressesError{}
  }
  c.addresses = append(c.addresses, na)
  return nil
}
...
```

</td></tr>

</tbody>
</table>

You can disable this rule at struct level, and also at method level by adding `//godddlint:disable:E004`.

##### E005: Unexported Fields

Entity fields should be mutated through methods with business meaning,
not by changing fields directly.

<table>
<thead><tr><th>❌ Bad</th><th>✅ Good</th></tr></thead>
<tbody>
<tr><td>

```go
type User struct {
  ID      UserID
  Name    Name
  Address Address
}
...
u.Address = na
```

</td><td>

```go
//godddlint:entity
type User struct {
  id      UserID
  name    Name
  address Address
}

func (c *User) MovedTo(na Address) {
  c.address = na
}
...
```

</td></tr>

</tbody>
</table>

You can disable this rule at struct level, and also at field level by adding `//godddlint:disable:E005`.

```go
//godddlint:entity
//godddlint:disable:E005
type User struct {
  id      UserID
  name    Name
  //godddlint:disable:E005
  Address Address
}
```

### Value Objects

[Value Objects][value-object] are objects that are defined by the value of their properties.

#### Value Object Rules

Current value object rule codes: `VO001`, `VOX001`, `VOX002`.

##### VO001: Non Pointer Receivers

A value object is assumed to be immutable, so internal mutation should not be allowed.

```go
//godddlint:valueObject
type Point struct {
    x, y int
}

// valueObject must not have a pointer receiver
func (c *Point) ...
```

You can disable this rule at struct level, and also at method level by adding `//godddlint:disable:VO001`.

```go
//godddlint:valueObject
type Point struct {
    x, y int
}

//godddlint:disable:VO001
func (c *Point) ...
```

##### VOX001: Immutable

A value object makes sense when its properties are immutable.
This rule checks that a value object has a constructor and that all fields are unexported.

```go
//godddlint:valueObject
type Point struct {
    x, y int
}

func New(x, y int) Point {
    return Point{x: x, y: y}
}
```

You can disable this rule at struct level, and also at field level by adding `//godddlint:disable:VOX001`.

```go
//godddlint:valueObject
type Point struct {
  //godddlint:disable:VOX001
  x, y int
}

func New(x, y int) Point {
  return Point{x: x, y: y}
}
```

#### VOX002: Maps/Slices Not Defensive Copied

When using a `map` or a `slice` inside a value object, you should prevent external mutation.
Use defensive copy in constructors and accessors.

## 📚 Glossary

- Constructor

In Go there is no constructor keyword, but in practice a constructor is a function that returns
a struct (or struct + error) and starts with `New` or `Must`, for example `func NewMyStruct(...) MyStruct`.

- Domain Error

It is a struct that implements `Error() string` and has domain meaning.

- Domain Struct

A domain struct is a Go struct that represents a domain object.
It can be an Aggregate, Entity, or Value Object.

[ddd]: https://en.wikipedia.org/wiki/Domain-driven_design
[entity]: https://en.wikipedia.org/wiki/Entity#In_computer_science
[value-object]: https://en.wikipedia.org/wiki/Value_object
