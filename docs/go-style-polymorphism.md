# Go-Style Polymorphism
## Polymorphysm by interfaces
Go is a statically typed language that doesn’t have inheritance or traditional object-oriented polymorphism. Instead, polymorphism in Go is typically achieved through interfaces as in the code below:

```go
package main

import (
	"fmt"
)

type Forwarder interface {
	Forward() string
}

type NextElement struct {
	Name string
}

func (e *NextElement) Forward() string {
	return fmt.Sprintf("Forwarding to element: %s", e.Name)
}

type Destination struct {
	URL string
}

func (d *Destination) Forward() string {
	return fmt.Sprintf("Forwarding to destination: %s", d.URL)
}

func ProcessForwarder(f Forwarder) {
	fmt.Println(f.Forward())
}

func main() {
	element := &NextElement{Name: "Element1"}
	destination := &Destination{URL: "https://example.com"}

	ProcessForwarder(element)
	ProcessForwarder(destination)
}
```
### How It Works

1. Interface Definition:
- Forwarder defines a method Forward() that types must implement.
2. Concrete Implementations:
- NextElement and Destination both implement the Forward() method.
3. Usage:
- Any type that satisfies the Forwarder interface can be passed to functions expecting a Forwarder.

## Polymorphysm by Pointers and Discriminator Field

Another powerful and idiomatic pattern in Go is Go-Style Polymorphism with Pointers, where a struct has optional pointer fields, and a "type" field (tag) determines which of those fields is relevant at runtime.

A **tagged union** is a design where a "tag" field specifies which of several possible data representations the object is using. In Go, this is implemented using a comination of:
- A **type discriminator field** (e.g. `Type String`)
- **Pointer fields** for the possible variants. If a field is not present in current data representations its value is simply null
- During runtime we can **validate** which data representation object is using and act accordingly

```go
type Next struct {
	// Type specifies the type of next loop-element, lupus-element (element) or external-element (destination)
	Type string `json:"type"`
	// List of input keys (Data fields) that have to be forwarded
	// Pass array with single element '*' to forward the whole input
	Keys []string `json:"keys"`
	// One of the fields below is not null
	Element     *NextElement `json:"element,omitempty"`
	Destination *Destination `json:"destination,omitempty"`
}

type NextElement struct {
	Name string `json:"name"`
}

type Destination struct {
	URL string `json:"url"`
}

func (n *Next) Validate() error {
	if n.Type == "element" && n.Element == nil {
		return fmt.Errorf("Element must be set for type 'element'")
	}
	if n.Type == "destination" && n.Destination == nil {
		return fmt.Errorf("Destination must be set for type 'destination'")
	}
	if n.Element != nil && n.Destination != nil {
		return fmt.Errorf("Only one of Element or Destination can be set")
	}
	return nil
}
```

## Comparison
| Feature              | Polymorphism by Pointers                         | Polymorphism by Interfaces                  |
| -------------------- | ------------------------------------------------ | ------------------------------------------- |
| **Runtime Behavior** | Uses a discriminator (`Type`) and pointer fields | Uses method implementation for polymorphism |
| **Type Safety**      | Explicit validation required                     | Enforced at compile time via interfaces     |
| **Serialization**    | Seamless with JSON                               | May require custom marshaling               |
| **Extensibility**    | Add new pointer fields and update `Type` enum    | Add new types implementing the interface    |
| **Ease of Use**      | Straightforward, but requires manual validation  | Clean and idiomatic in Go                   |
| **Shared Behavior**  | Requires external logic                          | Encapsulated within interface methods       |

Polymorphism with Pointers offers clear data representation, which can easily serialized/deserialized (e.g to JSON or YAML). Also, the it supports "type" being the part of data model, so it can be also stored or transmitted. On the other hand, Polymorphism by Interfaces is native to go, more clean and guarantees strong compile-time. 

Summaryzing:
- Polymorphism by Pointers is preferred in data-focused applications.
- Polymorphism by Intercaes is preffered in behavior-focused applications.

## Lupus

In Lupus a Polymorphism was needed to represent different types (varieties) of some [Lupn Objects](defs.md#lupn-object) such as [actions](defs.md#action) or [next](defs.md#next).
