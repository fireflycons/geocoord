# geocoord

`geocoord` is a small Go package for representing geographic coordinates and
performing common calculations between them. It has no dependencies outside of
the Go standard library.

Coordinates use decimal degrees, with latitude first and longitude second.
Valid latitudes are in `[-90, 90]`; valid longitudes are in `[-180, 180]`.
Distances and radii are in nautical miles. Headings are initial bearings in
degrees clockwise from north, normalized to `[0, 360)`.

## Install

```sh
go get github.com/fireflycons/gocoord
```

## Examples

### Distance and heading

Create coordinates with error handling, then calculate the distance and initial
heading from New York to London:

```go
package main

import (
	"fmt"
	"log"

	"github.com/fireflycons/gocoord"
)

func main() {
	newYork, err := geocoord.NewCoordinate(40.7128, -74.0060)
	if err != nil {
		log.Fatal(err)
	}

	london, err := geocoord.NewCoordinate(51.5074, -0.1278)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Distance: %.1f nautical miles\n", newYork.DistanceTo(london))
	fmt.Printf("Initial heading: %.1f degrees\n", newYork.HeadingTo(london))
}
```

### Check a radius

Use `IsWithinRadius` to check whether another coordinate is within a radius in
nautical miles. A negative radius always returns `false`.

```go
package main

import (
	"fmt"

	"github.com/fireflycons/gocoord"
)

func main() {
	origin := geocoord.MustNewCoordinate(0, 0)
	point := geocoord.MustNewCoordinate(0, 0.1)

	fmt.Println(origin.IsWithinRadius(point, 10)) // true
}
```

`NewCoordinate` returns an error for out-of-range values. `MustNewCoordinate`
is convenient for fixed, known-valid coordinates and panics if a value is
invalid.