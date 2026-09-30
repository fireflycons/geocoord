// Package geocoord represents geographic positions as latitude and longitude
// in decimal degrees and provides common calculations between coordinates.
//
// Latitude must be in [-90, 90] and longitude in [-180, 180]. Distances and
// radius values are measured in nautical miles; headings are initial bearings
// in degrees clockwise from north.
//
// Use NewCoordinate to return an error for out-of-range values, or
// MustNewCoordinate to panic when given invalid values.
//
// This package has no dependencies outside of the standard library.
package geocoord
