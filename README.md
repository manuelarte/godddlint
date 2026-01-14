# Go DDD Lint

[![CI](https://github.com/manuelarte/godddlint/actions/workflows/ci.yml/badge.svg)](https://github.com/manuelarte/godddlint/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/manuelarte/godddlint)](https://goreportcard.com/report/github.com/manuelarte/godddlint)
![version](https://img.shields.io/github/v/release/manuelarte/godddlint)

## ⬇️  Getting Started

To install it, run:

```bash
go install github.com/manuelarte/godddlint@latest
```

To run it in your project:

```bash
godddlint ./...
```

## 🚀 Features

### Entities

An [entity][entity] is an object defined not by its attributes, but its identity.

#### Entities Rules

##### E001: ID is the first embedded field

##### E002: Pointer Receivers

##### E003: Custom Types Over Primitives

##### E004: Using Custom Errors

##### E005: Unexported Fields

### Value Objects

[Value Objects][value-object] are objects that are equal due to the value of their properties.

#### Value Objects Rules

##### VO001: Non Pointer Receivers

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

[entity]: https://en.wikipedia.org/wiki/Entity#In_computer_science
[value-object]: https://en.wikipedia.org/wiki/Value_object
