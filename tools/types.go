package tools

// Command represents a generic command structure for all tools
type Command struct {
	Action string      `json:"action"`
	Key    string      `json:"key,omitempty"`
	Value  interface{} `json:"value,omitempty"`
	Path   string      `json:"path,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

// Registers specific types
type RegistersData struct {
	Registers map[string]interface{} `json:"registers"`
}

type RegistersCommand struct {
	Action string      `json:"action"`
	Key    string      `json:"key,omitempty"`
	Value  interface{} `json:"value,omitempty"`
}

// Math specific types
type MathCommand struct {
	Action string      `json:"action"`
	A      float64     `json:"a,omitempty"`
	B      float64     `json:"b,omitempty"`
	Key    string      `json:"key,omitempty"`
	Value  interface{} `json:"value,omitempty"`
	Expr   string      `json:"expr,omitempty"`
	Vars   []string    `json:"vars,omitempty"`
}

type MathResult struct {
	Result float64 `json:"result"`
	Key    string  `json:"key,omitempty"`
}

// Storage specific types
type StorageCommand struct {
	Action string      `json:"action"`
	Path   string      `json:"path,omitempty"`
	Key    string      `json:"key,omitempty"`
	Value  interface{} `json:"value,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

type StorageResult struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Path    string      `json:"path,omitempty"`
}

// Memory specific types
type MemoryData struct {
	Data map[string]interface{} `json:"data"`
}

type MemoryCommand struct {
	Action string      `json:"action"`
	Key    string      `json:"key,omitempty"`
	Value  interface{} `json:"value,omitempty"`
}
