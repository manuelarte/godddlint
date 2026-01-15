# 🧐 Go DDD Lint

[![CI](https://github.com/manuelarte/godddlint/actions/workflows/ci.yml/badge.svg)](https://github.com/manuelarte/godddlint/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/manuelarte/godddlint)](https://goreportcard.com/report/github.com/manuelarte/godddlint)
![version](https://img.shields.io/github/v/release/manuelarte/godddlint)

Go DDD Lint is an opinionated linter that checks for some of the properties a DDD model should achieve.

## ⬇️  Getting Started

To install it, run:

```bash
go install github.com/manuelarte/godddlint@latest
```

Mark your [value object][value-object] structs with `//godddlint:valueObject` and your entities with `//godddlint:entity`
and then run this linter in your project with:

```bash
godddlint ./...
```

## 🚀 Features

### Entities

An [entity][entity] is an object defined not by its attributes, but its identity.

#### Entities Rules

##### E001: ID is the first embedded field

TODO

##### E002: Pointer Receivers

An `Entity` can mutate, so then an internal mutation is allowed.

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

##### E003: Custom Domain Types Over Primitives

An `Entity` field needs to have more meaning than just a primitive value.

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
type UserID int
type Name string
type Address string

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

##### E004: Use Custom Domain Errors

Business processes that can return an error need to return a meaningful error, not a generic one.

<table>
<thead><tr><th>❌ Bad</th><th>✅ Good</th></tr></thead>
<tbody>
<tr><td>

```go
func (c *User) AddAddress(na Address) error {
  if len(c.addresses) >= 2 {
    return errors.New("max number of moves reached")
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
    return UserNotAllowedToMoveError{}
  }
  c.addresses = append(c.addresses, na)
  return nil
}
...
```

</td></tr>

</tbody>
</table>

##### E005: Unexported Fields

Entity fields need to be mutated by a method that indicates a business process.
Not by just changing the field.

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

func (c *User) UserMoved(na Address) {
  c.address = na
}
...
```

</td></tr>

</tbody>
</table>

### Value Objects

[Value Objects][value-object] are objects that are equal due to the value of their properties.

#### Value Objects Rules

##### VO001: Non Pointer Receivers

A value object is assumed to be immutable, so then no internal mutation is allowed.

```go
//godddlint:valueObject
type Point struct {
    x, y int
}

// valueObject must not have a pointer receiver
func (c *Point) ...
```

##### VOX001: Immutable

A value object makes sense when the properties are immutable.
This rule checks that a value object can only be created using a constructor that
tries to prevent that developers mutate fields in the struct.
Also checks that all the fields are unexported.

```go
//godddlint:valueObject
type Point struct {
    x, y int
}

func New(x, y int) Point {
    return Point{x: x, y: y}
}
```

#### VOX002: Maps/Slices Not Defensive Copied

When using a `map` or a `slice` inside a value object, we should prevent that it gets mutated.
To avoid that, you can use *Defensive Copy*.

## 📚 Glossary

* Constructor

In Go there is no concept of a constructor, but normally a constructor is considered a function that returns
a struct, or a struct and an error, and starts with `New` or `Must`, e.g. `func NewMyStruct(...) MyStruct`.

* Domain Error

It is a Go error struct (struct that implements `Error() string`) but that gives a domain meaning to
the error.

* Domain Struct

A domain struct is a Go struct that represents a domain object.
It can be an Aggregate, Entity or a Value Object.

[entity]: https://en.wikipedia.org/wiki/Entity#In_computer_science
[value-object]: https://en.wikipedia.org/wiki/Value_object
