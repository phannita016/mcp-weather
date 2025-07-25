package nws

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	nwsAPIBase = "https://api.weather.gov"
	userAgent  = "weather-app/1.0"
)

// MakeNWSRequest sends a GET request to the specified NWS API URL.
func MakeNWSRequest(url string) ([]byte, error) {
	client := http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/geo+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func GetAlertsURL(state string) string {
	return fmt.Sprintf("%s/alerts/active/area/%s", nwsAPIBase, state)
}

func GetForecastURL(latitude, longitude float64) string {
	return fmt.Sprintf("%s/points/%.4f,%.4f", nwsAPIBase, latitude, longitude)
}
