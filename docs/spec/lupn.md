# LupN
LupN (for Lup (Loop) Notation) is a language/notation to express a [loop-workflow]. It lacks the description of [computing-part](../defs.md#computing-part) of the [loop-logic](../defs.md#loop-logic). [Computing-part](../defs.md#computing-part) is specified outside of Lupus, in [external-elements](../defs.md#external-element).

LupN specifies then:
- [workflow](../defs.md#workflow) of [lupus-elements](../defs.md#loop-element) withing a loop
- references to [external-elements](../defs.md#external-element) expressed as [destinations](../defs.md#destination)
- [workflow](../defs.md#worokflow) of [actions](../defs.md#action) within a [lupus-element]
- reference or references to [egress-agent](../defs.md#egress-agent) expressed as [destination](../defs.md#destination)

As you can note LupN expresses some workflow on 2 levels, one global (a workflow of lupus-elements) and one inside a lupus-element (a workflow of actions). The capabilities of both are close to each other, but ultimately divergent. This document will also cover this issue.

From the implementation point of view a [LupN file](../defs.md#lupn-file) is actually a [YAML manifest file](../defs.md#yaml-manifest-file) of [Lupus-Master CR](../defs.md#lupus-master). Once [applied](../defs.md#yaml-manifest-file), a [Lupus-Master controller](../defs.md#lupus-master) spawns [lupus-elements](../defs.md#lupus-element) that deliver the expressed [loop-workflow](../defs.md#loop-workflow).

[LupN](../defs.md#lupn) expresses [loop-workflow](../defs.md#loop-workflow) by the specification of various objects in [YAML notation](https://yaml.org). Let's call these object a [LupN objects](../defs.md#lupn-object). This document will specify these objects and relation between them. Also it will indicate what usage of each one will mean in the [loop-workflow](../defs.md#loop-workflow) terminology and how [lupus-element](../defs.md#lupus-element) [controller](../defs.md#controller) will interprete them during runtime.

It occurs, that YAML object inside [YAML manifest files](../defs.md#yaml-manifest-file) are derived from Golang structs (Golang types), therefore we can describe [lupn-objects](../defs.md#lupn-object) based on these Golang structs.

It is mandatory to be familiar with [YAML](https://yaml.org) first. 

It does not matter wheter we consider `apiVersion`, `kind` and `metadata` as LupN or not. In some way it specifies the loop (e.g. metadata has object name), in some not (name is anyway repeated later in `spec`). But as for sure the `spec` objects states a Loop description.

### MasterSpec
It corresponds to the `MasterSpec` golang struct:
```go
// MasterSpec defines the desired state of Master
type MasterSpec struct {
	// Name of the Master CR (indicating the name of the loop)
	Name string `json:"name"`
	// Elements is a list of Lupus-Elements
	Elements []*ElementSpec `json:"elements"`
}
```

Each element of the `Elements` list will trigger [Lupus-Master](../defs.md#lupus-master) [controller](../defs.md#controller) to spawn an [API object](../defs.md#api-object) of type [lupus-element](../defs.md#lupus-element) with the given spec. The sequence of elements workflow is expressed in the elements itself (as next property).

### ElementSpec
It corresponds to the `ElemenetSpec` golang struct:
```go
// ElementSpec defines the desired state of Element
type ElementSpec struct {
	// Name is the name of the element, its distinct from Kubernetes API Object name, but rather serves ease of managemenet aspect for loop-designer
	Name string `json:"name"`
	// Descr is the description of the lupus-element, same as Name is serves as ease of management aspect for loop-designer
	Descr string `json:"descr"`
	// Actions is a list of Actions that lupus-element has to perform
	Actions []Action `json:"actions,omitempty"`
	// Next is a list of next objects (can be lupus-element or external-element) to which send the final-data
	Next []Next `json:"next,omitempty"`
	// Name of master element (used to use it as prefix for lupus-element name)
	Master string `json:"master,omitempty"`
}
```

### Next
```go
// It specifies the of next loop-element in loop workflow, it may be either lupus-element or reference to external-element
// It allows to forward the whole final-data, but also parts of it
type Next struct {
	// Type specifies the type of next loop-element, lupus-element (element) or external-element (destination)
	Type string `json:"type" kubebuilder:"validation:Enum=element,destination"`
	// List of input keys (Data fields) that have to be forwarded
	// Pass array with single element '*' to forward the whole input
	Keys []string `json:"keys"`
	// One of the fields below is not null
	Element     *NextElement `json:"element,omitempty" kubebuilder:"validation:Optional"`
	Destination *Destination `json:"destination,omitempty" kubebuilder:"validation:Optional"`
}
```
