package handlers

import (
	"math"
	"testing"
)

func TestHaversine_SamePoint(t *testing.T) {
	p := Point{Lat: 28.0612, Lng: -82.4123} // Tampa
	distance := Haversine(p, p)
	if distance != 0 {
		t.Errorf("expected 0 meters, got %f", distance)
	}
}

func TestHaversine_ShortDistance(t *testing.T) {
	// two points about 30 meters apart (classroom scale)
	p1 := Point{Lat: 28.06120, Lng: -82.41230}
	p2 := Point{Lat: 28.06140, Lng: -82.41210}
	distance := Haversine(p1, p2)

	// should be roughly 28-32 meters
	if distance < 20 || distance > 40 {
		t.Errorf("expected ~30 meters, got %f", distance)
	}
}

func TestHaversine_MediumDistance(t *testing.T) {
	// USF campus to downtown Tampa — roughly 15km
	usf := Point{Lat: 28.0587, Lng: -82.4139}
	downtown := Point{Lat: 27.9506, Lng: -82.4572}
	distance := Haversine(usf, downtown)

	// should be roughly 12-16 km
	if distance < 12000 || distance > 16000 {
		t.Errorf("expected ~14km, got %f meters", distance)
	}
}

func TestHaversine_LongDistance(t *testing.T) {
	// Tampa to London — roughly 7,100 km
	tampa := Point{Lat: 28.0612, Lng: -82.4123}
	london := Point{Lat: 51.5074, Lng: -0.1278}
	distance := Haversine(tampa, london)

	// should be roughly 6900-7300 km
	if distance < 6900000 || distance > 7300000 {
		t.Errorf("expected ~7100km, got %f meters", distance)
	}
}

func TestHaversine_Symmetrical(t *testing.T) {
	p1 := Point{Lat: 28.0612, Lng: -82.4123}
	p2 := Point{Lat: 51.5074, Lng: -0.1278}

	d1 := Haversine(p1, p2)
	d2 := Haversine(p2, p1)

	if math.Abs(d1-d2) > 0.001 {
		t.Errorf("expected symmetrical distance, got %f and %f", d1, d2)
	}
}

func TestHaversine_WithinClassroomRadius(t *testing.T) {
	// simulate: classroom center and student 40 meters away
	classroom := Point{Lat: 28.06120, Lng: -82.41230}
	student := Point{Lat: 28.06150, Lng: -82.41200}
	distance := Haversine(classroom, student)
	radius := 50.0

	if distance > radius {
		t.Errorf("expected student within %f meter radius, got %f meters", radius, distance)
	}
}

func TestHaversine_OutsideClassroomRadius(t *testing.T) {
	// simulate: student 200 meters away
	classroom := Point{Lat: 28.06120, Lng: -82.41230}
	student := Point{Lat: 28.06300, Lng: -82.41050}
	distance := Haversine(classroom, student)
	radius := 50.0

	if distance < radius {
		t.Errorf("expected student outside %f meter radius, got %f meters", radius, distance)
	}
}
