package v1

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

// NextElement indicates the next loop-element in loop-workflow of type lupus-element by its name
type NextElement struct {
	// Name is the lupus-name of lupus-element (the one specified in Element struct)
	Name string `json:"name"`
}
