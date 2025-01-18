package v1

// Destination represents an external-element
// It holds all the info needed to make a call to an External System
// It supports calls to HTTP server, Open Policy Agent or user-functions
// It is used in Action of type Send and can be also used (same as Lupus Element) as Next is Element spec
type Destination struct {
	// Type specifies if the external system is: HTTP server in gerneral, special type of HTTP server as Open Policy Agent or internal, user-defined Go function
	Type string `json:"type" kubebuilder:"validation:Enum=http;opa;gofunc"`
	// One of these fields is not null depending on a type, it has specifiaction specific to types
	HTTP   *HTTPDestination   `json:"http,omitempty" kubebuilder:"validation:Optional"`
	Opa    *OpaDestination    `json:"opa,omitempty" kubebuilder:"validation:Optional"`
	GoFunc *GoFuncDestination `json:"gofunc,omitempty" kubebuilder:"validation:Optional"`
}

// HTTPDestination defines fields specific to HTTP type
// This is information needed to make a HTTP request
type HTTPDestination struct {
	// Path specifies HTTP URI
	Path string `json:"path"`
	// Method specifies HTTP method
	Method string `json:"method"`
}

// OpaDestination defines fields specific to Open Policy Agent type
// This is information needed to make an Open Policy Agent request
// Call to Opa is actually a special type of HTTP call
type OpaDestination struct {
	// Path specifies HTTP URI, since method is known
	Path string `json:"path"`
}

// GoFuncDestination defines fields specific to GoFunc type
// This is information needed to call an user-function
type GoFuncDestination struct {
	// Name specifies the name of the function
	Name string `json:"name"`
}
