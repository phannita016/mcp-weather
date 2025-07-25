package dtos

type (
	Tool struct {
		Name        string         `json:"name"`
		Description string         `json:"description"`
		InputSchema map[string]any `json:"input_schema"`
	}

	Message struct {
		Role    string      `json:"role"`
		Content interface{} `json:"content"`
	}

	Toolcall struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}
)
