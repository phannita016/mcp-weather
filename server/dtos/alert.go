package dtos

type (
	AlertsParams struct {
		State string `json:"state" jsonschema:"two-letter US state code"`
	}

	FeatureCollection struct {
		Features []Feature `json:"features"`
	}

	Feature struct {
		AlertProperties `json:"properties"`
	}

	AlertProperties struct {
		Event       string `json:"event"`
		AreaDesc    string `json:"areaDesc"`
		Severity    string `json:"severity"`
		Description string `json:"description"`
		Instruction string `json:"instruction"`
	}
)
