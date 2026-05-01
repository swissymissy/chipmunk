package handlers

import (
	"math"
)

const (
	earthRM = 6371000.0 // r of earth in meters
)

// coordinate: latitude , longtitude
type Point struct {
	Lat float64
	Lng float64
}

// converts degrees to radians
func degreesToRad(d float64) float64 {
	return d * math.Pi/ 180
}

// calculate the shortest path between 2 points on earth
// using haversine formula
// this function returns meters
func Haversine(p1, p2 Point) float64 {
	p1Lat := degreesToRad(p1.Lat)
	p1Lng := degreesToRad(p1.Lng)

	p2Lat := degreesToRad(p2.Lat)
	p2Lng := degreesToRad(p2.Lng)

	diffLat := p2Lat - p1Lat 
	diffLng := p2Lng - p1Lng

	// haversine formula 
	// square of half the chord length
	a := math.Pow(math.Sin(diffLat/2),2) + math.Cos(p1Lat)*math.Cos(p2Lat)*math.Pow(math.Sin(diffLng/2),2)

	// calculate the central angle between 2 points
	// angular distance in radians
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	// convert to meters
	m := c * earthRM

	return m
}