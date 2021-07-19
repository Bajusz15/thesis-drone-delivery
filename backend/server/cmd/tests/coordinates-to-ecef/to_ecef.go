package main

import (
	"fmt"
	"github.com/StefanSchroeder/Golang-Ellipsoid/ellipsoid"
)

func main() {
	var lat float64
	var lon float64
	var alt float64
	geo := ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Meter, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)
	lat = 48.080922
	lon = 20.766208
	alt = 45
	x, y, z := geo.ToECEF(lat, lon, alt)
	fmt.Printf("x = %v \ny = %v \nz = %v\n", x, y, z)
	fmt.Printf("x = %f \ny = %f \nz = %f\n", x, y, z)

	x, z, y = geo.ToECEF(48.05, 20.75, 50)
	fmt.Printf("x = %f \ny = %f \nz = %f\n", x, y, z)
	x, z, y = geo.ToECEF(48.10, 20.78, 50)
	fmt.Printf("x = %f \ny = %f \nz = %f\n", x, y, z)

}
