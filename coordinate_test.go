package geocoord

import (
	"encoding/json"
	"math"
	"testing"
)

func TestCoordinateJSONRoundTrip(t *testing.T) {
	want := MustNewCoordinate(51.5, -0.12)

	data, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if got, expected := string(data), `{"lat":51.5,"lon":-0.12}`; got != expected {
		t.Fatalf("Marshal() = %s, want %s", got, expected)
	}

	var got Coordinate
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got.Latitude() != want.Latitude() || got.Longitude() != want.Longitude() {
		t.Errorf("Unmarshal() = (%v, %v), want (%v, %v)", got.Latitude(), got.Longitude(), want.Latitude(), want.Longitude())
	}
}

func TestCoordinateJSONUnmarshalValidation(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{name: "latitude out of range", json: `{"lat":90.1,"lon":0}`},
		{name: "longitude out of range", json: `{"lat":0,"lon":180.1}`},
		{name: "missing latitude", json: `{"lon":0}`},
		{name: "missing longitude", json: `{"lat":0}`},
		{name: "null latitude", json: `{"lat":null,"lon":0}`},
		{name: "non-numeric latitude", json: `{"lat":"0","lon":0}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			coordinate := MustNewCoordinate(12, 34)
			if err := json.Unmarshal([]byte(test.json), &coordinate); err == nil {
				t.Fatal("Unmarshal() error = nil, want an error")
			}
			if coordinate.Latitude() != 12 || coordinate.Longitude() != 34 {
				t.Errorf("coordinate changed after failed Unmarshal(): (%v, %v)", coordinate.Latitude(), coordinate.Longitude())
			}
		})
	}
}

func TestNewCoordinateFromString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		lat   float64
		lon   float64
	}{
		{name: "plain values", input: "51.5,-0.12", lat: 51.5, lon: -0.12},
		{name: "surrounding whitespace", input: "  51.5 , \t-0.12  ", lat: 51.5, lon: -0.12},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NewCoordinateFromString(test.input)
			if err != nil {
				t.Fatalf("NewCoordinateFromString() error = %v", err)
			}
			if got.Latitude() != test.lat || got.Longitude() != test.lon {
				t.Errorf("NewCoordinateFromString() = (%v, %v), want (%v, %v)", got.Latitude(), got.Longitude(), test.lat, test.lon)
			}
		})
	}
}

func TestNewCoordinateFromStringErrors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "missing longitude", input: "51.5"},
		{name: "extra component", input: "51.5,-0.12,1"},
		{name: "invalid latitude", input: "north,-0.12"},
		{name: "invalid longitude", input: "51.5,west"},
		{name: "latitude out of range", input: "90.1,0"},
		{name: "longitude out of range", input: "0,180.1"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewCoordinateFromString(test.input); err == nil {
				t.Fatal("NewCoordinateFromString() error = nil, want an error")
			}
		})
	}
}

func TestCoordinateStringRoundTrip(t *testing.T) {
	want := MustNewCoordinate(51.5, -0.12)
	var stringer interface{ String() string } = want
	if got, expected := stringer.String(), "51.5,-0.12"; got != expected {
		t.Fatalf("String() = %q, want %q", got, expected)
	}

	got, err := NewCoordinateFromString(want.String())
	if err != nil {
		t.Fatalf("NewCoordinateFromString(String()) error = %v", err)
	}
	if got.Latitude() != want.Latitude() || got.Longitude() != want.Longitude() {
		t.Errorf("round trip = (%v, %v), want (%v, %v)", got.Latitude(), got.Longitude(), want.Latitude(), want.Longitude())
	}
}

func TestCoordinateDistanceTo(t *testing.T) {
	tests := []struct {
		name string
		from Coordinate
		to   Coordinate
		want float64
	}{
		{name: "same point", from: MustNewCoordinate(0, 0), to: MustNewCoordinate(0, 0), want: 0},
		{name: "one degree latitude", from: MustNewCoordinate(0, 0), to: MustNewCoordinate(1, 0), want: 60.04046},
		{name: "one degree longitude at equator", from: MustNewCoordinate(0, 0), to: MustNewCoordinate(0, 1), want: 60.04046},
		{name: "antipodal points", from: MustNewCoordinate(0, 0), to: MustNewCoordinate(0, 180), want: 10807.2971},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.from.DistanceTo(test.to); math.Abs(got-test.want) > 0.001 {
				t.Errorf("DistanceTo(%v) = %v, want %v", test.to, got, test.want)
			}
		})
	}
}

func TestCoordinateHeadingTo(t *testing.T) {
	tests := []struct {
		name string
		from Coordinate
		to   Coordinate
		want float64
	}{
		{name: "north", from: MustNewCoordinate(0, 0), to: MustNewCoordinate(1, 0), want: 0},
		{name: "east", from: MustNewCoordinate(0, 0), to: MustNewCoordinate(0, 1), want: 90},
		{name: "south", from: MustNewCoordinate(0, 0), to: MustNewCoordinate(-1, 0), want: 180},
		{name: "west", from: MustNewCoordinate(0, 0), to: MustNewCoordinate(0, -1), want: 270},
		{name: "east across antimeridian", from: MustNewCoordinate(0, 179), to: MustNewCoordinate(0, -179), want: 90},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.from.HeadingTo(test.to); math.Abs(got-test.want) > 0.001 {
				t.Errorf("HeadingTo(%v) = %v, want %v", test.to, got, test.want)
			}
		})
	}
}

func TestCoordinateIsWithinRadius(t *testing.T) {
	origin := MustNewCoordinate(0, 0)

	tests := []struct {
		name   string
		other  Coordinate
		radius float64
		want   bool
	}{
		{name: "inside", other: MustNewCoordinate(0.1, 0), radius: 10, want: true},
		{name: "outside", other: MustNewCoordinate(0.2, 0), radius: 10, want: false},
		{name: "same point zero radius", other: origin, radius: 0, want: true},
		{name: "negative radius", other: origin, radius: -1, want: false},
		{name: "antimeridian", other: MustNewCoordinate(0, 179.95), radius: 10, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := origin.IsWithinRadius(test.other, test.radius); got != test.want {
				t.Errorf("IsWithinRadius(%v, %v) = %v, want %v", test.other, test.radius, got, test.want)
			}
		})
	}
}

func TestCoordinateIsWithinRadiusAcrossAntimeridian(t *testing.T) {
	origin := MustNewCoordinate(0, 179.9)
	other := MustNewCoordinate(0, -179.9)

	if !origin.IsWithinRadius(other, 15) {
		t.Error("expected coordinates across the antimeridian to be within 15 nautical miles")
	}
}

func TestCoordinateIsWithinRadiusNearPoles(t *testing.T) {
	tests := []struct {
		name   string
		origin Coordinate
		other  Coordinate
		radius float64
	}{
		{
			name:   "north pole",
			origin: MustNewCoordinate(90, 0),
			other:  MustNewCoordinate(90, 180),
			radius: 1,
		},
		{
			name:   "south pole",
			origin: MustNewCoordinate(-90, 0),
			other:  MustNewCoordinate(-90, 180),
			radius: 1,
		},
		{
			name:   "radius reaches north pole",
			origin: MustNewCoordinate(89, 0),
			other:  MustNewCoordinate(89.5, 90),
			radius: 91,
		},
		{
			name:   "radius reaches south pole",
			origin: MustNewCoordinate(-89, 0),
			other:  MustNewCoordinate(-89.5, 90),
			radius: 91,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if !test.origin.IsWithinRadius(test.other, test.radius) {
				t.Errorf("expected %v to be within %v nautical miles of %v", test.other, test.radius, test.origin)
			}
		})
	}
}
