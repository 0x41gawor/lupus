package v1

// Next is used in Observe Spec
// It specifies to which element forward the input
// It allows not to forward the whole input, but also parts of it
type Next struct {
	// Type specifies the type of the element ("Observe", "Decide", "Learn", "Execute", etc.)
	Type string `json:"type" kubebuilder:"validation:Enum=Observe;Decide;Learn;Execute"`
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
// It represents the Action that Decide has to perform on Data
type Action struct {
	// Name of the Action, it is for designer to ease the management of the Loop
	Name string `json:"name"`
	// Type of Action one of send;concat,remove,rename,duplicate
	Type string `json:"type" kubebuilder:"validation:Enum=send;concat,remove,rename,duplicate"`
	// One of these fields is not null depending on a Type.
	Send      *SendAction      `json:"send,omitempty" kubebuilder:"validation:Optional"`
	Concat    *ConcatAction    `json:"concat,omitempty" kubebuilder:"validation:Optional"`
	Remove    *RemoveAction    `json:"remove,omitempty" kubebuilder:"validation:Optional"`
	Rename    *RenameAction    `json:"rename,omitempty" kubebuilder:"validation:Optional"`
	Duplicate *DuplicateAction `json:"duplicate,omitempty" kubebuilder:"validation:Optional"`
}

type SendAction struct {
	InputKey    string      `json:"inputKey"`
	Destination Destination `json:"destination"`
	OutputKey   string      `json:"outputKey"`
}

type ConcatAction struct {
	InputKeys []string `json:"inputKeys"`
	OutputKey string   `json:"outputKey"`
}

type RemoveAction struct {
	InputKeys []string `json:"inputKeys"`
}

type RenameAction struct {
	InputKey  string `json:"inputKey"`
	OutputKey string `json:"outputKey"`
}

type DuplicateAction struct {
	InputKey  string `json:"inputKey"`
	OutputKey string `json:"outputKey"`
}

// Element is a polymorphic structure that can represent different types of specs
type Element struct {
	// Name is the name of the element
	Name string `json:"name"`
	// Type specifies the type of the element ("Observe", "Decide", "Learn", "Execute", etc.)
	Type string `json:"type" kubebuilder:"validation:Enum=Observe;Decide;Learn;Execute"`

	ObserveSpec *ObserveSpec `json:"observeSpec,omitempty"`
	DecideSpec  *DecideSpec  `json:"decideSpec,omitempty"`
	LearnSpec   *LearnSpec   `json:"learnSpec,omitempty"`
	ExecuteSpec *ExecuteSpec `json:"executeSpec,omitempty"`
}
