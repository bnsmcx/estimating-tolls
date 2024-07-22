package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Lane struct {
	origin      *Waypoint
	destination *Waypoint
	routes      *Routes
}

type Waypoint struct {
	address string
	lat     float64
	long    float64
}

type Routes struct {
	best  Route
	worst Route
	avg   Route
}

type Route struct {
	cost     float64
	polyLine string
}

func NewLane(start, end string) (*Lane, error) {
	var lane = Lane{}
	var err error
	if lane.origin, err = NewWaypoint(start); err != nil {
		return nil, err
	}
	if lane.destination, err = NewWaypoint(end); err != nil {
		return nil, err
	}
	if lane.routes, err = lane.CalcuateRoutes(); err != nil {
		return nil, err
	}
	return &lane, err
}

func (l *Lane) CalcuateRoutes() (*Routes, error) {

	pl, err := getPolyLine(l.origin, l.destination)
	if err != nil {
		return nil, err
	}

	//TODO: Logic to permute routes to get different costs
	var routes [3]Route
	routes[0] = Route{getTollCost(pl), pl}

	// TODO: Implement Sort interface here

	return &Routes{
		best:  routes[0],
		avg:   routes[0],
		worst: routes[0],
	}, nil
}

func getTollCost(pl string) float64 {
	// TODO: Toll Guru api call goes here
	return 13.37
}

func getPolyLine(origin *Waypoint, destination *Waypoint) (string, error) {
	baseURL := "https://maps.googleapis.com/maps/api/directions/json"
	params := url.Values{}
	params.Add("origin", fmt.Sprintf("%f,%f", origin.lat, origin.long))
	params.Add("destination", fmt.Sprintf("%f,%f", destination.lat, destination.long))
	params.Add("key", os.Getenv("GOOGLE_API"))

	// Make the request
	resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code")
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Parse the JSON response
	var directionsResponse DirectionsResponse
	if err := json.Unmarshal(body, &directionsResponse); err != nil {
		return "", err
	}

	if directionsResponse.Status != "OK" {
		return "", fmt.Errorf("received non-OK status: %s", directionsResponse.Status)
	}

	if len(directionsResponse.Routes) == 0 {
		return "", fmt.Errorf("no routes found")
	}

	overviewPolyline := directionsResponse.Routes[0].OverviewPolyline.Points
	return overviewPolyline, nil
}

type DirectionsResponse struct {
	Routes []struct {
		OverviewPolyline struct {
			Points string `json:"points"`
		} `json:"overview_polyline"`
	} `json:"routes"`
	Status string `json:"status"`
}

func (l Lane) String() string {
	var sb strings.Builder
	sb.WriteString("Lane:\n")
	sb.WriteString("\n  Origin:\n\n")
	sb.WriteString(fmt.Sprintf("    %s\n", l.origin.address))
	sb.WriteString(fmt.Sprintf("    Lat: %.2f, Long: %.2f\n", l.origin.lat, l.origin.long))
	sb.WriteString("\n  Destination:\n\n")
	sb.WriteString(fmt.Sprintf("    %s\n", l.destination.address))
	sb.WriteString(fmt.Sprintf("    Lat: %.2f, Long: %.2f\n", l.destination.lat, l.destination.long))
	sb.WriteString("\n  Routes:\n")
	sb.WriteString("\n    Best:\n\n")
	sb.WriteString(fmt.Sprintf("        Cost: $%.2f\n", l.routes.best.cost))
	sb.WriteString(fmt.Sprintf("        Poly Line: %s\n", l.routes.best.polyLine))
	sb.WriteString("\n    Worst:\n\n")
	sb.WriteString(fmt.Sprintf("        Cost: $%.2f\n", l.routes.worst.cost))
	sb.WriteString(fmt.Sprintf("        Poly Line: %s\n", l.routes.worst.polyLine))
	sb.WriteString("\n    Avg:\n\n")
	sb.WriteString(fmt.Sprintf("        Cost: $%.2f\n", l.routes.avg.cost))
	sb.WriteString(fmt.Sprintf("        Poly Line: %s\n", l.routes.avg.polyLine))
	return sb.String()
}

func NewWaypoint(address string) (*Waypoint, error) {
	lat, long, err := getLatLong(address)
	if err != nil {
		return nil, fmt.Errorf("getting lat/long for address: %s", err)
	}

	return &Waypoint{
		address: address,
		lat:     lat,
		long:    long,
	}, nil
}

func getLatLong(address string) (lat, long float64, err error) {
	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"
	params := url.Values{}
	params.Add("address", address)
	params.Add("key", os.Getenv("GOOGLE_API"))

	// Make the request
	resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, 0, fmt.Errorf("received non-200 response code")
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

  if resp.StatusCode != http.StatusOK {
    log.Println(resp.Status)
    log.Println(string(body))
  }

	// Parse the JSON response
	var geocodingResponse GeocodingResponse
	if err := json.Unmarshal(body, &geocodingResponse); err != nil {
		return 0, 0, err
	}

	if geocodingResponse.Status != "OK" {
		return 0, 0, fmt.Errorf("received non-OK status: %s", geocodingResponse.Status)
	}

	location := geocodingResponse.Results[0].Geometry.Location
	return location.Lat, location.Lng, nil
}

type GeocodingResponse struct {
	Results []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}
