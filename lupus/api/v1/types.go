package v1

// Next is used in Observe Spec
// It specifies to which element forward the input
// It allows not to forward the whole input, but also parts of it
type Next struct {
	// Kubernetes name of the API Object
	// This is the name that you give in Master CR spec
	Name string `json:"name"`
	// List of input tags (json root fields) that have to be forwarded
	// Pass array with single element '*' to forward the whole input
	Tags []string `json:"tags"`
}

// Destination in Action type and in Learn Spec
// Destination is a polymorphic field that can represent different types
type Destination struct {
	// Discriminator: "HTTP", "FILE", "gRPC", etc.
	Type string `json:"type" kubebuilder:"validation:Enum=HTTP;FILE;gRPC;Opa"`

	// HTTP-specific fields
	HTTP *HTTPDestination `json:"http,omitempty" kubebuilder:"validation:Optional"`

	// File-specific fields
	File *FileDestination `json:"file,omitempty" kubebuilder:"validation:Optional"`

	// gRPC-specific fields
	GRPC *GRPCDestination `json:"grpc,omitempty" kubebuilder:"validation:Optional"`
	// Opa-specific fields
	Opa *OpaDestination `json:"opa,omitempty" kubebuilder:"validation:Optional"`
}

// HTTPDestination defines fields specific to HTTP type
type HTTPDestination struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

// FileDestination defines fields specific to FILE type
type FileDestination struct {
	Path string `json:"path"`
}

// GRPCDestination defines fields specific to gRPC type
type GRPCDestination struct {
	Host    string `json:"host"`
	Service string `json:"service"`
	Method  string `json:"method"`
}

// OpaDestination defines fields specific to Open Policy Agent type
type OpaDestination struct {
	Path string `json:"path"`
}

// Action is used in Decide spec
// It represents the Action that Decide has to perform on its input
// As for now only the Action of sending the Input somewhere (to HTTP server, Open Policy Agent, gRPC service) is supported
// Thus Action includes its name, input_tag (part of input that has to be sent) and destination
type Action struct {
	// Name of the Action, it is for designer to ease the management of the Loop
	Name string `json:"name"`
	// Specifies the root field of input json that will be send, pass * for whole input to be sent
	Input_tag string `json:"input_tag"`
	// Specifies Destination where the input has to be sent
	Destination Destination `json:"destination"`
}
