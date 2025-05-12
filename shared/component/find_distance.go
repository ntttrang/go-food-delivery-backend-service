package shareComponent

import (
	"math"
)

// toRadians converts degrees to radians.
func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

// haversine returns the distance between two points specified by
// latitude and longitude in kilometers.
func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth’s radius in kilometers

	dLat := toRadians(lat2 - lat1)
	dLon := toRadians(lon2 - lon1)
	lat1Rad := toRadians(lat1)
	lat2Rad := toRadians(lat2)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}
