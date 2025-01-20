# Data

One of the requirements for Lupus was for it to be [data-driven](../defs.md#data-driven). [Data](../defs.md#data) is the heart and the core of fullfillment of this requirement. The runtime of [lupus-element](../defs.md#lupus-element) controller is driven by [data](../defs.md#data) contents and [actions](../defs.md#action) chain specified in [LupN](../defs.md#lupn). Data does not impose any reconciliation logic.e

Data is the way in which user can:
- retrive information from [current-state](defs.md#current-state)
- store auxiliary information (such as responses from [external-elements](../defs.md#external-element))
- store logging/debuggin information
- save information needed to formulate [control-action](defs.md#control-action)

during a single [loop-iteration](../defs.md#loop-iteration).

In each iteration [data](../defs.md#data) resets.

Data is an information carrier. Let's discuss how it stores this information.

## Data format

Data stores information in a JSON object, which is convertible to golang's `map[string]interface{}`. It means that the top-level fields of JSON have to have string keys. Their values can be anything since they are represented by `interface{}`.

[lupus-element](../defs.md#lupus-element) [controller](../defs.md#controller) has access to [data](../defs.md) by the means of [actions](../defs.md#action). Operations performed on [data](../defs.md#data) are specified in a [Lupn-file](../defs.md#lupn-file).

The keystone concept of data are data-fields. Data is a key-value store, where the keys are `string` and values can be anything (especially again a key-value object, in this case data supports field-nesting). Data-fields are indentified by their keys.

Action can:
- Get given data-field by specyfying its key
- Set given data-field by specyfing its key and value

## Implementation

The specifics of Data implementation results from the [communication-between-lupus-elements](../com-bet-lup-ele.md). Since they communicate by updating [custom-resource](../defs.md#custom-resources) [`status`](../defs.md#status), [data](../defs.md#data) had to be part of it. So, this is the first requirement that was imposed on data by the implementation layer.

The second one, was imposed by [data-driven](../defs.md#data-driven) design requirement of Lupus. Data has to be versatile and universal. The first idea that come up was json. Json is able to represent any structured data. 

But how to represent json in golang? What type of field should `ElementStatus` have?

The response is `RawExtension`. It is a type in Kubernetes used to handle **arbitrary raw JSON or YAML data**. It is part of the k8s.io/apimachinery/pkg/runtime package and is commonly used when a resource needs to embed or work with structured, yet flexible, data. `RawExtension` fits in with our requirements.

```go
// ElementStatus defines the observed state of Element
type ElementStatus struct {
	// Input contains operational data
	Input runtime.RawExtension `json:"input"`
	// Timestamp of the last update
	LastUpdated metav1.Time `json:"lastUpdated"`
}
```

`RawExtension` is the type which we will use to carry information (e.g. [current-state](../defs.md#current-state)) in Kubernetes platform. But how about inside Lupus? Can we operate easily on this type?

The `RawExtension` definition is down below:

```go
type RawExtension struct {
    Raw []byte `json:"-"` // Serialized JSON or YAML data
    Object Object         // A runtime.Object representation
}
```

The intended purpose of `RawExtension` by the kubernetes developers is so that it can then be deserialized into some known structure. But due to the [data-driven](../defs.md#data-driven) requirement, there is no such structure in our case. 

We need a structure that represents ANY json. 

The first idea that cames upon is the golang type - `interface{}` since it can represent any information. The problem with `interface{}` is that it can not be operated on, it provides no iterface to interact with. It is a primary type. 

The second idea of representing json was to use - `map[string]interface{}`, since most of json instances are indeed a key-value stores. Keys in this case are of type `string` and values can be anything (hence represented by the `interface{}`)  in golang. In most* cases json objects have several root fields and it ideally fits to the `map[string]interface{}` representation. 

This how Data was born. Data is actually a wrapper structure for the map mentioned above.

```go
type Data struct {
	Body map[string]interface{}
}
```

This struct has a plethora of functions (methods) defined which act as an interface to work with data. These methods are called during the execution of action and typically, except of methods `Get` and `Set`, one method corresponds to exactly one action.

A keystone concept of data is a [data-field](../defs.md#data-field). It is the same as field in JSON. Each field is identified by its key and stores a value. With the method `Get()` we can obtain a value that resides under a certain key and with `Set()` we can set a new value for field referred with certain key.

When using keys, either as input or output in actions, it is possible to:
- access nested fields by the `.` delimiter
- use a `"*"` wildcard to indicate all fields


*But it does not cover all json objects that world has seen. JSON allows for the top level element to be an array itself. It imposes some constrains in the loop design. Especially json representing the [current-state](../defs.md#current-state) of [managed-system](../defs.md#managed-system) sent via [lupin-interface](../defs.md#lupin-interface) needs to be serializable into `map[string]interface{}`. So it can't be:
- a primitive type
- an array
- JSON object with non string keys

The same rules apply to the response from [external-element](../defs.md#external-element) when send action specifies `OutputKey` as `["*"]`, which means that is has to replace the whole body of data.

