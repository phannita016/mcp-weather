package dtos

type (
	ForecastParams struct {
		Latitude  float64 `json:"latitude" jsonschema:"latitude of the location"`
		Longitude float64 `json:"longitude" jsonschema:"longitude of the location"`
	}

	ForecastData struct {
		Properties ForecastProperties `json:"properties"`
	}

	ForecastPeriod struct {
		Name             string `json:"name"`
		Temperature      int    `json:"temperature"`
		TemperatureUnit  string `json:"temperatureUnit"`
		WindSpeed        string `json:"windSpeed"`
		WindDirection    string `json:"windDirection"`
		DetailedForecast string `json:"detailedForecast"`
	}

	PointsData struct {
		Properties ForecastProperties `json:"properties"`
	}

	ForecastProperties struct {
		Periods     []ForecastPeriod `json:"periods"`
		ForecastURL string           `json:"forecast"`
	}
)
