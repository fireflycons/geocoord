package geocoord

import (
	"math"
	"testing"
)

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
