package v1

// It specifies the of next node in loop workflow, it may be other Lupus element or some external Destination
// It allows not to forward the whole Data final form, but also parts of it
type Next struct {
	// Type specifies class of next node in loop workflow, it may be other Lupus element or some external Destination
	Type string `json:"type" kubebuilder:"validation:Enum=element,destination"`
	// List of input keys (Data fields) that have to be forwarded
	// Pass array with single element '*' to forward the whole input
	Keys        []string     `json:"keys"`
	Element     *NextElement `json:"element,omitempty" kubebuilder:"validation:Optional"`
	Destination *Destination `json:"destination,omitempty" kubebuilder:"validation:Optional"`
}

type NextElement struct {
	// Kubernetes name of the API Object
	// This is the name that you give in Master CR spec
	Name string `json:"name"`
}
