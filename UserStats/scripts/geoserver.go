package main

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/StefanSchroeder/Golang-Ellipsoid/ellipsoid"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lukeroth/gdal" // requirement libgdal v2.4.0
	"github.com/tidwall/gjson"
	//"github.com/umahmood/haversine"
	"github.com/fhs/go-netcdf/netcdf"
	"github.com/im7mortal/UTM"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// get RGB color from name NOT USED
func getRGB(colorName string) (int, int, int) {
	if colorName == "blue" {
		return 0, 0, 255
	} else if colorName == "cyan" {
		return 0, 153, 255
	} else if colorName == "green" {
		return 0, 153, 0
	} else if colorName == "darkgreen" {
		return 0, 95, 0
	} else if colorName == "yellowgreen" {
		return 0, 255, 0
	} else if colorName == "yellow" {
		return 255, 255, 0
	} else if colorName == "gold" {
		return 255, 187, 0
	} else if colorName == "orange" {
		return 255, 102, 0
	} else if colorName == "red" {
		return 255, 0, 0
	} else if colorName == "darkred" {
		return 153, 0, 0
	} else if colorName == "maroon" {
		return 84, 0, 0
	} else {
		return -1, -1, -1
	}
}

// get color NOT USED
func getColorOld(conf map[string]string, metricName string, value float64) (int, int, int) {
	if metricName == "noiseLAeq" || metricName == "noiseLA" ||
		metricName == "noiseLAmax" || metricName == "LAeq" {
		if value < 44.3 {
			return getRGB("blue")
		} else if value >= 44.3 && value <= 48.8 {
			return getRGB("cyan")
		} else if value > 48.8 && value <= 53.3 {
			return getRGB("green")
		} else if value > 53.3 && value <= 57.7 {
			return getRGB("yellowgreen")
		} else if value > 57.7 && value <= 62.1 {
			return getRGB("yellow")
		} else if value > 62.1 && value <= 66.6 {
			return getRGB("gold")
		} else if value > 66.6 && value <= 71 {
			return getRGB("orange")
		} else if value > 71 && value <= 75.5 {
			return getRGB("red")
		} else if value > 75.5 && value <= 79.9 {
			return getRGB("darkred")
		} else {
			return getRGB("maroon")
		}
	} else if metricName == "airTemperature" || metricName == "minTemperature" ||
		metricName == "maxTemperature" || metricName == "minGroundTemperature" {
		if value < -20 {
			return getRGB("blue")
		} else if value >= -20 && value <= 0 {
			return getRGB("cyan")
		} else if value > 0 && value <= 9 {
			return getRGB("green")
		} else if value > 9 && value <= 15 {
			return getRGB("yellowgreen")
		} else if value > 15 && value <= 18 {
			return getRGB("yellow")
		} else if value > 18 && value <= 21 {
			return getRGB("gold")
		} else if value > 21 && value <= 25 {
			return getRGB("orange")
		} else if value > 25 && value <= 30 {
			return getRGB("red")
		} else if value > 30 && value <= 34 {
			return getRGB("darkred")
		} else {
			return getRGB("maroon")
		}
	} else if metricName == "airHumidity" {
		if value < 40 {
			return getRGB("blue")
		} else if value >= 40 && value <= 45.5 {
			return getRGB("cyan")
		} else if value > 45.5 && value <= 51.1 {
			return getRGB("green")
		} else if value > 51.1 && value <= 56.7 {
			return getRGB("yellowgreen")
		} else if value > 56.7 && value <= 62.2 {
			return getRGB("yellow")
		} else if value > 62.2 && value <= 67.8 {
			return getRGB("gold")
		} else if value > 67.8 && value <= 73.3 {
			return getRGB("orange")
		} else if value > 73.3 && value <= 78.9 {
			return getRGB("red")
		} else if value > 78.9 && value <= 84.4 {
			return getRGB("darkred")
		} else {
			return getRGB("maroon")
		}
	} else if metricName == "windSpeed" {
		if value <= 3.9 {
			return getRGB("blue")
		} else if value > 3.9 && value <= 7.9 {
			return getRGB("cyan")
		} else if value > 7.9 && value <= 11.9 {
			return getRGB("green")
		} else if value > 11.9 && value <= 15.9 {
			return getRGB("yellowgreen")
		} else if value > 15.9 && value <= 19.9 {
			return getRGB("yellow")
		} else if value > 19.9 && value <= 23.9 {
			return getRGB("gold")
		} else if value > 23.9 && value <= 27.9 {
			return getRGB("orange")
		} else if value > 27.9 && value <= 31.9 {
			return getRGB("red")
		} else if value > 31.9 && value <= 35.9 {
			return getRGB("darkred")
		} else {
			return getRGB("maroon")
		}
	} else if metricName == "windGust" {
		if value <= 3.32 {
			return getRGB("blue")
		} else if value > 3.32 && value <= 6.66 {
			return getRGB("cyan")
		} else if value > 6.66 && value <= 9.99 {
			return getRGB("green")
		} else if value > 9.99 && value <= 13.32 {
			return getRGB("yellowgreen")
		} else if value > 13.32 && value <= 16.66 {
			return getRGB("yellow")
		} else if value > 16.66 && value <= 19.99 {
			return getRGB("gold")
		} else if value > 19.99 && value <= 23.32 {
			return getRGB("orange")
		} else if value > 23.32 && value <= 26.66 {
			return getRGB("red")
		} else if value > 26.66 && value <= 30 {
			return getRGB("darkred")
		} else {
			return getRGB("maroon")
		}
	} else if metricName == "dewPoint" {
		if value < -10 {
			return getRGB("blue")
		} else if value >= -10 && value <= -7.99 {
			return getRGB("cyan")
		} else if value > -7.99 && value <= -5.99 {
			return getRGB("green")
		} else if value > -5.99 && value <= -3.99 {
			return getRGB("yellowgreen")
		} else if value > -3.99 && value <= -1.99 {
			return getRGB("yellow")
		} else if value > -1.99 && value <= -0.01 {
			return getRGB("gold")
		} else if value > -0.01 && value <= 1.99 {
			return getRGB("orange")
		} else if value > 1.99 && value <= 3.99 {
			return getRGB("red")
		} else if value > 3.99 && value <= 5.99 {
			return getRGB("darkred")
		} else {
			return getRGB("maroon")
		}
	} else if metricName == "airQualityAQI" || metricName == "airqualityAQI" {
		if int(value) == 0 {
			return getRGB("blue")
		} else if int(value) == 1 {
			return getRGB("cyan")
		} else if int(value) == 2 {
			return getRGB("green")
		} else if int(value) == 3 {
			return getRGB("yellowgreen")
		} else if int(value) == 4 {
			return getRGB("yellow")
		} else if int(value) == 5 {
			return getRGB("gold")
		} else if int(value) == 6 {
			return getRGB("orange")
		} else if int(value) == 7 {
			return getRGB("red")
		} else if int(value) == 8 {
			return getRGB("darkred")
		} else {
			return getRGB("maroon")
		}
	} else if metricName == "airQualityPM10" || metricName == "PM10" {
		if value <= 10 {
			return getRGB("blue")
		} else if value > 10 && value <= 20 {
			return getRGB("cyan")
		} else if value > 20 && value <= 30 {
			return getRGB("green")
		} else if value > 30 && value <= 40 {
			return getRGB("yellowgreen")
		} else if value > 40 && value <= 50 {
			return getRGB("yellow")
		} else if value > 50 && value <= 60 {
			return getRGB("gold")
		} else if value > 60 && value <= 70 {
			return getRGB("orange")
		} else if value > 70 && value <= 80 {
			return getRGB("red")
		} else if value > 80 && value <= 90 {
			return getRGB("darkred")
		} else {
			return getRGB("maroon")
		}
	} else if metricName == "airQualityNO2" || metricName == "NO2" {
		if value <= 20 {
			return getRGB("blue")
		} else if value > 20 && value <= 50 {
			return getRGB("cyan")
		} else if value > 50 && value <= 70 {
			return getRGB("green")
		} else if value > 70 && value <= 120 {
			return getRGB("yellowgreen")
		} else if value > 120 && value <= 150 {
			return getRGB("yellow")
		} else if value > 150 && value <= 180 {
			return getRGB("gold")
		} else if value > 180 && value <= 200 {
			return getRGB("orange")
		} else if value > 200 && value <= 250 {
			return getRGB("red")
		} else if value > 250 && value <= 300 {
			return getRGB("darkred")
		} else {
			return getRGB("maroon")
		}
	} else if metricName == "airQualityPM2_5" || metricName == "PM2_5" {
		if value <= 5 {
			return getRGB("blue")
		} else if value > 5 && value <= 10 {
			return getRGB("cyan")
		} else if value > 10 && value <= 15 {
			return getRGB("green")
		} else if value > 15 && value <= 25 {
			return getRGB("yellowgreen")
		} else if value > 25 && value <= 35 {
			return getRGB("yellow")
		} else if value > 35 && value <= 40 {
			return getRGB("gold")
		} else if value > 40 && value <= 50 {
			return getRGB("orange")
		} else if value > 50 && value <= 60 {
			return getRGB("red")
		} else if value > 60 && value <= 70 {
			return getRGB("darkred")
		} else {
			return getRGB("maroon")
		}
	} else if metricName == "safetyOnBikeDensity" || metricName == "bikeSafety" {
		if int(value) <= -7 {
			return getRGB("red")
		} else if int(value) >= -6 && int(value) <= -4 {
			return getRGB("orange")
		} else if int(value) >= -3 && int(value) <= -1 {
			return getRGB("gold")
		} else if int(value) == 0 {
			return getRGB("yellow")
		} else if int(value) > 1 && int(value) <= 3 {
			return getRGB("yellowgreen")
		} else if int(value) > 4 && int(value) <= 6 {
			return getRGB("green")
		} else {
			return getRGB("darkgreen")
		}
	} else if metricName == "accidentDensity" {
		if value == 1 {
			return getRGB("yellowgreen")
		} else if int(value) >= 2 && int(value) <= 3 {
			return getRGB("yellow")
		} else if int(value) >= 4 && int(value) <= 5 {
			return getRGB("gold")
		} else if int(value) >= 6 && int(value) <= 7 {
			return getRGB("orange")
		} else if int(value) >= 8 && int(value) <= 9 {
			return getRGB("red")
		} else {
			return getRGB("maroon")
		}
	} else if metricName == "airQualityNOx" {
		if value <= 10 {
			return getRGB("blue")
		} else if value > 10 && value <= 25 {
			return getRGB("cyan")
		} else if value > 25 && value <= 35 {
			return getRGB("green")
		} else if value > 35 && value <= 60 {
			return getRGB("yellowgreen")
		} else if value > 60 && value <= 75 {
			return getRGB("yellow")
		} else if value > 75 && value <= 90 {
			return getRGB("gold")
		} else if value > 90 && value <= 105 {
			return getRGB("orange")
		} else if value > 105 && value <= 125 {
			return getRGB("red")
		} else if value > 125 && value <= 150 {
			return getRGB("darkred")
		} else {
			return getRGB("maroon")
		}
	} else if metricName == "CAQI" {
		if value <= 25 {
			return getRGB("yellowgreen")
		} else if value > 25 && value <= 50 {
			return getRGB("yellow")
		} else if value > 50 && value <= 75 {
			return getRGB("gold")
		} else if value > 75 && value <= 100 {
			return getRGB("orange")
		} else {
			return getRGB("darkred")
		}
	} else if metricName == "EAQI" {
		if value == 1 {
			return getRGB("green")
		} else if value == 2 {
			return getRGB("yellowgreen")
		} else if value == 3 {
			return getRGB("yellow")
		} else if value == 4 {
			return getRGB("orange")
		} else {
			return getRGB("darkred")
		}
	} else if metricName == "CO" {
		if value < 1.9 {
			return getRGB("green")
		} else if value >= 1.9 && value <= 3.9 {
			return getRGB("yellowgreen")
		} else if value > 3.9 && value <= 5.9 {
			return getRGB("yellow")
		} else if value > 5.9 && value <= 7.9 {
			return getRGB("gold")
		} else if value > 7.9 && value <= 10 {
			return getRGB("orange")
		} else {
			return getRGB("red")
		}
	} else if metricName == "Benzene" {
		if value < 0.9 {
			return getRGB("green")
		} else if value >= 0.9 && value <= 1.9 {
			return getRGB("yellowgreen")
		} else if value > 1.9 && value <= 2.9 {
			return getRGB("yellow")
		} else if value > 2.9 && value <= 3.9 {
			return getRGB("gold")
		} else if value > 3.9 && value <= 5 {
			return getRGB("orange")
		} else {
			return getRGB("red")
		}
	} else {
		fmt.Println("Color table not found for metric: " + metricName)
		return -1, -1, -1
	} /*else if metricName == "accidentDensity" {
		if value <= 0 {
			return getRGB("blue")
		} else if value > 0 && value <= 1 {
			return getRGB("cyan")
		} else if value > 1 && value <= 2 {
			return getRGB("green")
		} else if value > 2 && value <= 3 {
			return getRGB("yellowgreen")
		} else if value > 3 && value <= 4 {
			return getRGB("yellow")
		} else if value > 4 && value <= 5 {
			return getRGB("gold")
		} else if value > 5 && value <= 6 {
			return getRGB("orange")
		} else if value > 6 && value <= 7 {
			return getRGB("red")
		} else if value > 7 && value <= 8 {
			return getRGB("darkred")
		} else {
			return getRGB("maroon")
		}
	}*/
}

// get color
func getColor(colorMap map[string]map[int]map[string]interface{}, metricName string, value float64) (int, int, int) {
	var max, min float64
	var rgb []int
	for order := 1; order <= len(colorMap[metricName]); order++ {
		if _, ok := colorMap[metricName][order]; ok {
			max_bool := false
			min_bool := false
			if val, ok := colorMap[metricName][order]["max"]; ok {
				max_bool = true
				max = val.(float64)
			}
			if val, ok := colorMap[metricName][order]["min"]; ok {
				min_bool = true
				min = val.(float64)
			}
			rgb = colorMap[metricName][order]["rgb"].([]int)
			if !min_bool && max_bool {
				if value < max {
					return rgb[0], rgb[1], rgb[2]
				}
			} else if !max_bool && min_bool {
				if value >= min {
					return rgb[0], rgb[1], rgb[2]
				}
			} else {
				if value >= min && value < max {
					return rgb[0], rgb[1], rgb[2]
				}
			}
		}
	}
	return -1, -1, -1
}

// get colors' maps from MySQL
func getColorsMaps(conf map[string]string) map[string]map[int]map[string]interface{} {
	db, err := sql.Open("mysql", conf["MySQL_username"]+":"+conf["MySQL_password"]+"@tcp("+conf["MySQL_hostname"]+":"+conf["MySQL_port"]+")/"+conf["MySQL_database"])

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	defer db.Close()

	// get map's colors
	results, err := db.Query("SELECT metric_name, `min`, `max`, `order`, rgb FROM heatmap.colors")
	if err != nil {
		panic(err.Error())
	}
	var metricName, min_b, max_b, order_b, rgb_b []byte
	result := map[string]map[int]map[string]interface{}{}
	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&metricName, &min_b, &max_b, &order_b, &rgb_b)
		min, err_min := strconv.ParseFloat(string(min_b), 64)
		max, err_max := strconv.ParseFloat(string(max_b), 64)
		order, err_order := strconv.Atoi(string(order_b))
		var rgb []int
		rgb_err := json.Unmarshal(rgb_b, &rgb)
		if err_order == nil {
			if _, ok := result[string(metricName)]; !ok {
				result[string(metricName)] = map[int]map[string]interface{}{}
			}
			if _, ok := result[string(metricName)][order]; !ok {
				result[string(metricName)][order] = map[string]interface{}{}
			}
			if err_min == nil {
				result[string(metricName)][order]["min"] = min
			}
			if err_max == nil {
				result[string(metricName)][order]["max"] = max
			}
			if rgb_err == nil {
				result[string(metricName)][order]["rgb"] = rgb
			}
		}
	}
	return result
}

// get destination point given distance (m) and bearing (clockwise from north) from start point
// http://www.movable-type.co.uk/scripts/latlong.html
// http://cdn.rawgit.com/chrisveness/geodesy/v1.1.1/latlon-spherical.js
// http://williams.best.vwh.net/avform.htm#LL
// http://janmatuschek.de/LatitudeLongitudeBoundingCoordinates
func getDestinationPoint(latitude, longitude, distance, bearing float64) (float64, float64) {
	var radius float64 = 6371000

	delta := distance / radius // angular distance in radians
	theta := bearing * math.Pi / 180

	fi1 := latitude * math.Pi / 180
	lambda1 := longitude * math.Pi / 180

	sinfi1 := math.Sin(fi1)
	cosfi1 := math.Cos(fi1)
	sindelta := math.Sin(delta)
	cosdelta := math.Cos(delta)
	sintheta := math.Sin(theta)
	costheta := math.Cos(theta)

	sinfi2 := sinfi1*cosdelta + cosfi1*sindelta*costheta
	fi2 := math.Asin(sinfi2)
	y := sintheta * sindelta * cosfi1
	x := cosdelta - sinfi1*sinfi2
	lambda2 := lambda1 + math.Atan2(y, x)

	return fi2 * 180 / math.Pi, lambda2 * 180 / math.Pi // normalise to -180...+180
	//return fi2 * 180 / math.Pi, math.Mod((lambda2*180/math.Pi+540), 360) - 180 // normalise to -180...+180
}

// get bounding coordinates (decimal latitude and longitude)
// (minLat, maxLat, minLon, maxLon) from coordinates
func getBoundingCoordinates(lat_center, lon_center, bboxLengthX, bboxLengthY float64, deltaLatLon bool) (float64, float64, float64, float64) {
	// if cluster's lengths were estimating by calculating delta latitude and longitude in the data set
	if deltaLatLon {
		return lat_center - bboxLengthX/2, lat_center + bboxLengthX/2, lon_center - bboxLengthY/2, lon_center + bboxLengthY/2
	} else { // else if cluster's lengths were provided
		// alternative method for calculating the bounding box
		//minLat, _ := getDestinationPoint(lat_center, lon_center, bboxLengthY/2, 180)
		//maxLat, _ := getDestinationPoint(lat_center, lon_center, bboxLengthY/2, 0)
		//_, minLon := getDestinationPoint(lat_center, lon_center, bboxLengthX/2, -90)
		//_, maxLon := getDestinationPoint(lat_center, lon_center, bboxLengthX/2, 90)

		// alternative method for calculating the bounding box
		//maxLat, maxLon := getDestinationPoint(lat_center, lon_center, bboxLengthX*math.Sqrt(2)/2, 45)
		//minLat, minLon := getDestinationPoint(lat_center, lon_center, bboxLengthY*math.Sqrt(2)/2, -135)

		geo1 := ellipsoid.Init("WGS84", ellipsoid.Degrees, ellipsoid.Meter, ellipsoid.LongitudeIsSymmetric, ellipsoid.BearingIsSymmetric)
		minLat, _ := geo1.At(lat_center, lon_center, bboxLengthY/2, 180)
		maxLat, _ := geo1.At(lat_center, lon_center, bboxLengthY/2, 0)
		_, minLon := geo1.At(lat_center, lon_center, bboxLengthX/2, -90)
		_, maxLon := geo1.At(lat_center, lon_center, bboxLengthX/2, 90)

		return minLat, maxLat, minLon, maxLon
	}
}

// save a GeoTIFF to disk
func saveGeoTIFF(conf map[string]string, colorMap map[string]map[int]map[string]interface{}, filePath, mapName, metricName string, latitude, longitude, value float64, date string, bboxLengthX, bboxLengthY float64, deltaLatLon bool) bool {
	// format the date
	dateString := strings.Replace(date, ":", "", -1)
	dateString = strings.Replace(dateString, "-", "", -1)
	dateString = strings.Replace(dateString, " ", "T", -1) + "Z"

	// format the filename
	latitude_s := fmt.Sprintf("%v", latitude)
	longitude_s := fmt.Sprintf("%v", longitude)
	fileName := filePath + "/" + strings.Replace(latitude_s, ".", "-", -1) + "_" + strings.Replace(longitude_s, ".", "-", -1) + "_" + dateString + ".tiff"

	//  initialize the image size in pixels
	nx := 1
	ny := 1

	// get metric name's color
	r, g, b := getColor(colorMap, metricName, value)

	// if there is no color map for this metric name, then return false
	if r == -1 {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// get coordinates ranges
	var minLat, maxLat, minLon, maxLon float64

	minLat, maxLat, minLon, maxLon = getBoundingCoordinates(latitude, longitude, bboxLengthX, bboxLengthY, deltaLatLon)

	// create each channel
	r_pixels := make([]uint8, ny*nx)
	g_pixels := make([]uint8, ny*nx)
	b_pixels := make([]uint8, ny*nx)
	alpha := make([]uint8, ny*nx)
	for index, _ := range alpha {
		alpha[index] = 255
	}

	// set the pixel data
	for x := 0; x < nx; x++ {
		for y := 0; y < ny; y++ {
			// computing the values
			loc := x + y*nx

			// test for drawing black countour of squares
			// uncomment this and set nx := 10 ny := 10
			/*if (x <= 9 && y == 0) || // first line
				((x == 0 || x == 9) && y == 1) || // border
				((x == 0 || x == 9) && y == 2) || // border
				((x == 0 || x == 9) && y == 3) || // border
				((x == 0 || x == 9) && y == 4) || // border
				((x == 0 || x == 9) && y == 5) || // border
				((x == 0 || x == 9) && y == 6) || // border
				((x == 0 || x == 9) && y == 7) || // border
				((x == 0 || x == 9) && y == 8) || // border
				(x <= 9 && y == 9) || // last line
				(x >= 4 && x <= 5 && y >= 4 && y <= 5) { // center
				r_pixels[loc] = 0
				g_pixels[loc] = 0
				b_pixels[loc] = 0
			} else {
				r_pixels[loc] = uint8(r)
				g_pixels[loc] = uint8(g)
				b_pixels[loc] = uint8(b)
			}*/

			// set pixel values
			// comment this if you uncommented the code above
			r_pixels[loc] = uint8(r)
			g_pixels[loc] = uint8(g)
			b_pixels[loc] = uint8(b)
		}
	}

	// set geotransform
	xmin, ymin, xmax, ymax := minLon, minLat, maxLon, maxLat
	xres := (xmax - xmin) / float64(nx)
	yres := (ymax - ymin) / float64(ny)
	//geotransform := [6]float64{xmin, xres, 0, ymax, 0, -yres}
	geotransform := [6]float64{xmin - xres/2, xres, 0, ymax + yres/2, 0, -yres}

	driver, err := gdal.GetDriverByName("GTiff")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	// create the 4-band raster file with transparency (r, g, b, alpha)
	dst_ds := driver.Create(fileName, ny, nx, 4, gdal.Byte, nil)
	defer dst_ds.Close()

	spatialRef := gdal.CreateSpatialReference("") // establish encoding
	spatialRef.FromEPSG(4326)                     // WGS84 lat/lon
	srString, err := spatialRef.ToWKT()
	dst_ds.SetProjection(srString)       // export coords to file
	dst_ds.SetGeoTransform(geotransform) // specify coords

	r_band := dst_ds.RasterBand(1)
	r_band.IO(gdal.Write, 0, 0, nx, ny, r_pixels, nx, ny, 0, 0)
	g_band := dst_ds.RasterBand(2)
	g_band.IO(gdal.Write, 0, 0, nx, ny, g_pixels, nx, ny, 0, 0)
	b_band := dst_ds.RasterBand(3)
	b_band.IO(gdal.Write, 0, 0, nx, ny, b_pixels, nx, ny, 0, 0)
	alpha_band := dst_ds.RasterBand(4)
	alpha_band.IO(gdal.Write, 0, 0, nx, ny, alpha, nx, ny, 0, 0) // set alpha channel as opaque (nodata = transparent)

	return true
}

// save a GeoTIFF to disk from a whole dataset, reading data from MySQL
func saveGeoTIFFDataset(conf map[string]string, colorMap map[string]map[int]map[string]interface{}, filePath, mapName, metricName, date string) bool {
	// create folder if it does not exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, 0700)
	}

	// format the date
	dateString := strings.Replace(date, ":", "", -1)
	dateString = strings.Replace(dateString, "-", "", -1)
	dateString = strings.Replace(dateString, " ", "T", -1) + "Z"

	// format the filename
	fileName := filePath + "/" + dateString + ".tiff"

	db, err := sql.Open("mysql", conf["MySQL_username"]+":"+conf["MySQL_password"]+"@tcp("+conf["MySQL_hostname"]+":"+conf["MySQL_port"]+")/"+conf["MySQL_database"])

	// if there is an error opening the connection, handle it
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	// get map's metadata
	results, err := db.Query("SELECT x_length, y_length, projection FROM heatmap.metadata WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND date = '" + date + "'")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	var xLength, yLength int
	var projection int
	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&xLength, &yLength, &projection)
		if err != nil {
			return false
			//log.Fatal(err)
			//panic(err.Error())
		}
	}

	// get map's pixel size and coordinates' bounding box
	results, err = db.Query("SELECT (max(latitude)-min(latitude))/" + fmt.Sprintf("%v", yLength) + " AS ny, (max(longitude)-min(longitude))/" + fmt.Sprintf("%v", xLength) + " AS nx, max(latitude) AS ymax, min(latitude) AS ymin, max(longitude) AS xmax, min(longitude) AS xmin FROM heatmap.data WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND date = '" + date + "'")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	var nx, ny, xmax, xmin, ymax, ymin int
	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&ny, &nx, &ymax, &ymin, &xmax, &xmin)
		if err != nil {
			return false
			//log.Fatal(err)
			//panic(err.Error())
		}
	}

	nx += 1
	ny += 1

	// create each image's channel
	r_pixels := make([]uint8, ny*nx)
	g_pixels := make([]uint8, ny*nx)
	b_pixels := make([]uint8, ny*nx)
	alpha := make([]uint8, ny*nx)

	// get map's data
	results, err = db.Query("SELECT latitude AS y, longitude AS x, value FROM heatmap.data WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND date = '" + date + "'")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	var x, y int
	var value float64
	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&y, &x, &value)
		if err != nil {
			return false
			//log.Fatal(err)
			//panic(err.Error())
		}
		// get metric name's color
		r, g, b := getColor(colorMap, metricName, value)
		// if there is no color map for this metric name, then return false
		if r == -1 {
			return false
			//log.Fatal(err)
			//panic(err.Error())
		}
		// set the pixel data
		//loc := (x-xmin)/xLength + ((y-ymin)/yLength)*nx
		loc := (x-xmin)/xLength + ((ymax-y)/yLength)*nx
		r_pixels[loc] = uint8(r)
		g_pixels[loc] = uint8(g)
		b_pixels[loc] = uint8(b)
		alpha[loc] = 255
	}

	// set geotransform
	//xres := float64((xmax - xmin)) / float64(nx)
	//yres := float64((ymax - ymin)) / float64(ny)
	xres := float64((xmax - xmin + xLength)) / float64(nx)
	yres := float64((ymax - ymin + yLength)) / float64(ny)
	//geotransform := [6]float64{float64(xmin), xres, 0, float64(ymax), 0, -yres}
	geotransform := [6]float64{float64(xmin) - xres/2, xres, 0, float64(ymax) + yres/2, 0, -yres}

	driver, err := gdal.GetDriverByName("GTiff")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	// create the 4-band raster file with transparency (r, g, b, alpha)
	dst_ds := driver.Create(filePath+"/tmp.tiff", nx, ny, 4, gdal.Byte, nil)
	//defer dst_ds.Close()

	spatialRef := gdal.CreateSpatialReference("") // establish encoding
	spatialRef.FromEPSG(projection)               // EPSG projection
	srString, err := spatialRef.ToWKT()
	dst_ds.SetProjection(srString)       // export coords to file
	dst_ds.SetGeoTransform(geotransform) // specify coords

	r_band := dst_ds.RasterBand(1)
	r_band.IO(gdal.Write, 0, 0, nx, ny, r_pixels, nx, ny, 0, 0)
	g_band := dst_ds.RasterBand(2)
	g_band.IO(gdal.Write, 0, 0, nx, ny, g_pixels, nx, ny, 0, 0)
	b_band := dst_ds.RasterBand(3)
	b_band.IO(gdal.Write, 0, 0, nx, ny, b_pixels, nx, ny, 0, 0)
	alpha_band := dst_ds.RasterBand(4)
	alpha_band.IO(gdal.Write, 0, 0, nx, ny, alpha, nx, ny, 0, 0) // set alpha channel as opaque (nodata = transparent)
	dst_ds.Close()

	// convert GeoTIFF from projection to WGS84
	cmd := conf["gdalwarp_path"] + " " + filePath + "/tmp.tiff" + " " + fileName + " -t_srs \"+proj=longlat +datum=WGS84 +ellps=WGS84\""
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	/*ds, err := gdal.Open(filePath+"/tmp.tiff", gdal.ReadOnly)
	if err != nil {
		log.Fatal(err)
	}
	defer ds.Close()

	options := []string{"-t_srs \"+proj=longlat +ellps=WGS84\""}
	outputDs := gdal.GDALWarp(fileName, gdal.Dataset{}, []gdal.Dataset{ds}, options)
	defer outputDs.Close()*/

	// remove temp GeoTIFF
	cmd = "rm " + filePath + "/tmp.tiff"
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	return true
}

// save a GeoTIFF to disk from a whole dataset, reading data from a binary file (GRAL)
func saveGeoTIFFDatasetFile(conf map[string]string, colorMap map[string]map[int]map[string]interface{}, filePath, mapName, metricName, date string) bool {
	// create folder if it does not exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, 0700)
	}

	// format the date
	dateString := strings.Replace(date, ":", "", -1)
	dateString = strings.Replace(dateString, "-", "", -1)
	dateString = strings.Replace(dateString, " ", "T", -1) + "Z"

	// format the filename
	fileName := filePath + "/" + dateString + ".tiff"

	db, err := sql.Open("mysql", conf["MySQL_username"]+":"+conf["MySQL_password"]+"@tcp("+conf["MySQL_hostname"]+":"+conf["MySQL_port"]+")/"+conf["MySQL_database"])

	// if there is an error opening the connection, handle it
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	// get map's metadata
	results, err := db.Query("SELECT x_length, y_length, projection, insertOnDB FROM heatmap.metadata WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND date = '" + date + "'")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	var xLength, yLength int32
	var projection, insertOnDB int
	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&xLength, &yLength, &projection, &insertOnDB)
		if err != nil {
			return false
			//log.Fatal(err)
			//panic(err.Error())
		}
	}

	// prepared statement for inserting data in data table
	data_stmt, err := db.Prepare("INSERT IGNORE INTO heatmap.data " +
		"(map_name, metric_name, latitude, longitude, value, date) " +
		"VALUES (?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	// get bounding box coordinates from binary file
	var ymax int32 = -1 // max UTM Easting
	var ymin int32 = -1 // min UTM Easting
	var xmax int32 = -1 // max UTM Northing
	var xmin int32 = -1 // min UTM Northing
	// open GRAL binary file for reading
	file, err := os.Open(conf["gral_data"] + "/" + mapName + dateString)
	//defer file.Close()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	m := payload{}
	fileinfo, err := file.Stat()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	filesize := fileinfo.Size()
	// skip the first 4 bytes
	readNextBytes(file, 4)
	var c int64
	// read 12 bytes into struct [UTM Easting, UTM Northing, Value]
	for c = 0; c < (filesize-4)/12; c++ {
		data := readNextBytes(file, 12)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.LittleEndian, &m)
		if err != nil {
			return false
			//log.Fatal(err)
			//panic(err.Error())
			//log.Fatal("binary.Read failed", err)
		}
		if m.UTMN > ymax || ymax == -1 {
			ymax = m.UTMN
		}
		if m.UTMN < ymin || ymin == -1 {
			ymin = m.UTMN
		}
		if m.UTME > xmax || xmax == -1 {
			xmax = m.UTME
		}
		if m.UTME < xmin || xmin == -1 {
			xmin = m.UTME
		}
	}
	// close file
	file.Close()

	// get map's pixel size
	ny := int((ymax-ymin)/yLength) + 1
	nx := int((xmax-xmin)/xLength) + 1

	// create each image's channel
	r_pixels := make([]uint8, ny*nx)
	g_pixels := make([]uint8, ny*nx)
	b_pixels := make([]uint8, ny*nx)
	alpha := make([]uint8, ny*nx)

	// get map's data from GRAL binary file
	file, err = os.Open(conf["gral_data"] + "/" + mapName + dateString)
	defer file.Close()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	fileinfo, err = file.Stat()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	filesize = fileinfo.Size()
	// skip the first 4 bytes
	readNextBytes(file, 4)
	// read 12 bytes into struct [UTM Easting, UTM Northing, Value]
	for c = 0; c < (filesize-4)/12; c++ {
		data := readNextBytes(file, 12)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.LittleEndian, &m)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}
		// get metric name's color
		r, g, b := getColor(colorMap, metricName, float64(m.Value))
		// if there is no color map for this metric name, then return false
		if r == -1 {
			return false
			//log.Fatal(err)
			//panic(err.Error())
		}
		// set the pixel data
		loc := (m.UTME-xmin)/xLength + ((ymax-m.UTMN)/yLength)*int32(nx)
		r_pixels[loc] = uint8(r)
		g_pixels[loc] = uint8(g)
		b_pixels[loc] = uint8(b)
		alpha[loc] = 255

		// if insertOnDB == 1, then save data on MySQL
		if insertOnDB == 1 || strings.Contains(mapName, "GRALheatmapHelsinki") {
			// convert UTM coordinates to latitude, longitude
			latitude, longitude, _ := UTM.ToLatLon(float64(m.UTME), float64(m.UTMN), getUTMZone(projection), "", true)
			// insert only data within this bounding box
			if latitude >= 60.15500512818767 && latitude <= 60.16173190973275 && longitude >= 24.911051048489185 && longitude <= 24.92336775228557 {
				_, err = data_stmt.Exec(mapName, metricName, latitude, longitude, m.Value, date)
				if err != nil {
					return false
				}
			}
		}
	}

	// set geotransform
	//xres := float64((xmax-xmin)) / float64(nx)
	//yres := float64((ymax-ymin)) / float64(ny)
	xres := float64((xmax-xmin)+xLength) / float64(nx)
	yres := float64((ymax-ymin)+yLength) / float64(ny)
	//geotransform := [6]float64{float64(xmin), xres, 0, float64(ymax), 0, -yres}
	geotransform := [6]float64{float64(xmin) - xres/2, xres, 0, float64(ymax) + yres/2, 0, -yres}

	driver, err := gdal.GetDriverByName("GTiff")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	// create the 4-band raster file with transparency (r, g, b, alpha)
	dst_ds := driver.Create(filePath+"/tmp.tiff", nx, ny, 4, gdal.Byte, nil)
	//defer dst_ds.Close()

	spatialRef := gdal.CreateSpatialReference("") // establish encoding
	spatialRef.FromEPSG(projection)               // EPSG projection
	srString, err := spatialRef.ToWKT()
	dst_ds.SetProjection(srString)       // export coords to file
	dst_ds.SetGeoTransform(geotransform) // specify coords

	r_band := dst_ds.RasterBand(1)
	r_band.IO(gdal.Write, 0, 0, nx, ny, r_pixels, nx, ny, 0, 0)
	g_band := dst_ds.RasterBand(2)
	g_band.IO(gdal.Write, 0, 0, nx, ny, g_pixels, nx, ny, 0, 0)
	b_band := dst_ds.RasterBand(3)
	b_band.IO(gdal.Write, 0, 0, nx, ny, b_pixels, nx, ny, 0, 0)
	alpha_band := dst_ds.RasterBand(4)
	alpha_band.IO(gdal.Write, 0, 0, nx, ny, alpha, nx, ny, 0, 0) // set alpha channel as opaque (nodata = transparent)
	dst_ds.Close()

	// convert GeoTIFF from projection to WGS84
	cmd := conf["gdalwarp_path"] + " " + filePath + "/tmp.tiff" + " " + filePath + "/uncompressed.tiff -t_srs \"+proj=longlat +datum=WGS84 +ellps=WGS84\""
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// remove temp GeoTIFF
	cmd = "rm " + filePath + "/tmp.tiff"
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// compress GeoTIFF
	cmd = conf["gdal_translate_path"] + " -co compress=deflate -co zlevel=9 -co tiled=yes -co NUM_THREADS=ALL_CPUS --config GDAL_CACHEMAX 512 -of GTiff " + filePath + "/uncompressed.tiff" + " " + fileName
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// remove uncompressed GeoTIFF
	cmd = "rm " + filePath + "/uncompressed.tiff"
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// compress GRAL binary file
	cmd = "tar cJvf " + conf["gral_data"] + "/" + mapName + dateString + ".tar.xz" + " " + conf["gral_data"] + "/" + mapName + dateString
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	} else {
		// remove GRAL binary file
		cmd = "rm " + conf["gral_data"] + "/" + mapName + dateString
		_, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			return false
			//log.Fatal(err)
			//panic(err.Error())
		}
	}
	return true
}

// save a map's GeoTIFFs to disk (calculate the latitude and longitude deltas to determine the cluster size's length)
func saveGeoTIFFsDeltaLatLonOld(conf map[string]string, colorMap map[string]map[int]map[string]interface{}, filePath, mapName, metricName string, date string) bool {
	var latitude float64 = -100  // not set, define outside the range -90, 90
	var longitude float64 = -200 // not set, define outside the range -180, 180
	save := false

	// create folder if it does not exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, 0700)
	}

	db, err := sql.Open("mysql", conf["MySQL_username"]+":"+conf["MySQL_password"]+"@tcp("+conf["MySQL_hostname"]+":"+conf["MySQL_port"]+")/"+conf["MySQL_database"])

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	// calculate the minimum delta latitude for this data set
	results, err := db.Query("SELECT latitude FROM heatmap.data WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND clustered = " + conf["clustered"] + " AND date = '" + date + "' ORDER BY latitude")
	if err != nil {
		panic(err.Error())
	}

	// latitude
	var lat float64

	// minimum delta latitude between locations
	var lat_delta_min float64 = -1

	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&lat)
		if err != nil {
			panic(err.Error())
		}
		if latitude != -100 && latitude != lat {
			lat_delta := math.Abs(latitude - lat)
			// if the delta latitude is the minimum
			if lat_delta_min == -1 || lat_delta < lat_delta_min {
				lat_delta_min = lat_delta
			}
		}
		latitude = lat
	}

	// calculate the minimum delta longitude for this data set
	results, err = db.Query("SELECT longitude FROM heatmap.data WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND clustered = " + conf["clustered"] + " AND date = '" + date + "' ORDER BY longitude")
	if err != nil {
		panic(err.Error())
	}

	// longitude
	var lon float64

	// minimum delta longitude between locations
	var lon_delta_min float64 = -1

	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&lon)
		if err != nil {
			panic(err.Error())
		}
		if longitude != -200 && longitude != lon {
			lon_delta := math.Abs(longitude - lon)
			// if the delta longitude is the minimum
			if lon_delta_min == -1 || lon_delta < lon_delta_min {
				lon_delta_min = lon_delta
			}
		}
		longitude = lon
	}

	// get the map's data
	results, err = db.Query("SELECT latitude, longitude, value FROM heatmap.data WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND clustered = " + conf["clustered"] + " AND date = '" + date + "'")
	if err != nil {
		panic(err.Error())
	}

	// value
	var value float64

	// save GeoTIFF for this row
	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&lat, &lon, &value)
		if err != nil {
			panic(err.Error())
		}
		latitude = lat
		longitude = lon
		save = saveGeoTIFF(conf, colorMap, filePath, mapName, metricName, latitude, longitude, value, date, lat_delta_min, lon_delta_min, true)
	}
	// return save result, (true if success, otherwise false)
	return save
}

// save a map's GeoTIFFs to disk (calculate the latitude and longitude deltas to determine the cluster size's length)
func saveGeoTIFFsDeltaLatLon(conf map[string]string, colorMap map[string]map[int]map[string]interface{}, filePath, mapName, metricName string, date string) bool {
	var latitude float64 = -100  // not set, define outside the range -90, 90
	var longitude float64 = -200 // not set, define outside the range -180, 180

	// create folder if it does not exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, 0700)
	}

	// format the date
	dateString := strings.Replace(date, ":", "", -1)
	dateString = strings.Replace(dateString, "-", "", -1)
	dateString = strings.Replace(dateString, " ", "T", -1) + "Z"

	// format the filename
	fileName := filePath + "/" + dateString + ".tiff"

	db, err := sql.Open("mysql", conf["MySQL_username"]+":"+conf["MySQL_password"]+"@tcp("+conf["MySQL_hostname"]+":"+conf["MySQL_port"]+")/"+conf["MySQL_database"])

	// if there is an error opening the connection, handle it
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	// calculate the bounding box
	results, err := db.Query("SELECT max(latitude), min(latitude), max(longitude), min(longitude) FROM heatmap.data WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND clustered = " + conf["clustered"] + " AND date = '" + date + "' ORDER BY latitude")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	var maxLat, minLat, maxLon, minLon float64

	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&maxLat, &minLat, &maxLon, &minLon)
	}
	xmin, ymin, xmax, ymax := minLon, minLat, maxLon, maxLat

	// calculate the map's projection
	results, err = db.Query("SELECT projection FROM heatmap.metadata WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND clustered = " + conf["clustered"] + " AND date = '" + date + "'")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	var projection int

	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&projection)
	}
	if projection == 0 {
		projection = 4326
	}

	// calculate the minimum delta latitude for this data set
	results, err = db.Query("SELECT DISTINCT(latitude) FROM heatmap.data WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND clustered = " + conf["clustered"] + " AND date = '" + date + "' ORDER BY latitude")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// latitude
	var lat float64

	// minimum delta latitude between locations
	var lat_delta_min float64 = -1

	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&lat)
		if err != nil {
			return false
			//log.Fatal(err)
			//panic(err.Error())
		}
		if latitude != -100 && latitude != lat {
			lat_delta := math.Abs(latitude - lat)
			// if the delta latitude is the minimum
			if lat_delta_min == -1 || lat_delta < lat_delta_min {
				lat_delta_min = lat_delta
			}
		}
		latitude = lat
	}

	// calculate the minimum delta longitude for this data set
	results, err = db.Query("SELECT DISTINCT(longitude) FROM heatmap.data WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND clustered = " + conf["clustered"] + " AND date = '" + date + "' ORDER BY longitude")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// longitude
	var lon float64

	// minimum delta longitude between locations
	var lon_delta_min float64 = -1

	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&lon)
		if err != nil {
			return false
			//log.Fatal(err)
			//panic(err.Error())
		}
		if longitude != -200 && longitude != lon {
			lon_delta := math.Abs(longitude - lon)
			// if the delta longitude is the minimum
			if lon_delta_min == -1 || lon_delta < lon_delta_min {
				lon_delta_min = lon_delta
			}
		}
		longitude = lon
	}

	// if this map is of type enfuser, increase delta lat and delta lon
	if strings.Contains(mapName, "Enfuser") {
		lat_delta_min *= 1.05
		lon_delta_min *= 1.05
	}

	//fmt.Println(lat_delta_min, lon_delta_min)

	// get map's pixel size
	ny := int((maxLat-minLat)/lat_delta_min) + 1
	nx := int((maxLon-minLon)/lon_delta_min) + 1

	// create each image's channel
	r_pixels := make([]uint8, ny*nx)
	g_pixels := make([]uint8, ny*nx)
	b_pixels := make([]uint8, ny*nx)
	alpha := make([]uint8, ny*nx)

	// get the map's data
	results, err = db.Query("SELECT latitude, longitude, value FROM heatmap.data WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND clustered = " + conf["clustered"] + " AND date = '" + date + "'")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// value
	var value float64

	// save GeoTIFF for this row
	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&lat, &lon, &value)
		if err != nil {
			return false
			//log.Fatal(err)
			//panic(err.Error())
		}
		// get metric name's color
		r, g, b := getColor(colorMap, metricName, value)
		// if there is no color map for this metric name, then return false
		if r == -1 {
			return false
		}
		// set the pixel data
		loc := int32((lon-minLon)/lon_delta_min) + int32((maxLat-lat)/lat_delta_min)*int32(nx)
		r_pixels[loc] = uint8(r)
		g_pixels[loc] = uint8(g)
		b_pixels[loc] = uint8(b)
		alpha[loc] = 255
	}

	// set geotransform
	//xres := (xmax - xmin) / float64(nx)
	//yres := (ymax - ymin) / float64(ny)
	xres := (xmax - xmin + lon_delta_min) / float64(nx)
	yres := (ymax - ymin + lat_delta_min) / float64(ny)
	//geotransform := [6]float64{float64(xmin), xres, 0, float64(ymax), 0, -yres}
	geotransform := [6]float64{float64(xmin) - xres/2, xres, 0, float64(ymax) + yres/2, 0, -yres}

	driver, err := gdal.GetDriverByName("GTiff")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	// create the 4-band raster file with transparency (r, g, b, alpha)
	dst_ds := driver.Create(filePath+"/uncompressed.tiff", nx, ny, 4, gdal.Byte, nil)
	//defer dst_ds.Close()

	spatialRef := gdal.CreateSpatialReference("") // establish encoding
	spatialRef.FromEPSG(projection)               // EPSG projection
	srString, err := spatialRef.ToWKT()
	dst_ds.SetProjection(srString)       // export coords to file
	dst_ds.SetGeoTransform(geotransform) // specify coords

	r_band := dst_ds.RasterBand(1)
	r_band.IO(gdal.Write, 0, 0, nx, ny, r_pixels, nx, ny, 0, 0)
	g_band := dst_ds.RasterBand(2)
	g_band.IO(gdal.Write, 0, 0, nx, ny, g_pixels, nx, ny, 0, 0)
	b_band := dst_ds.RasterBand(3)
	b_band.IO(gdal.Write, 0, 0, nx, ny, b_pixels, nx, ny, 0, 0)
	alpha_band := dst_ds.RasterBand(4)
	alpha_band.IO(gdal.Write, 0, 0, nx, ny, alpha, nx, ny, 0, 0) // set alpha channel as opaque (nodata = transparent)
	dst_ds.Close()

	// compress GeoTIFF
	cmd := conf["gdal_translate_path"] + " -co compress=deflate -co zlevel=9 -co tiled=yes -co NUM_THREADS=ALL_CPUS --config GDAL_CACHEMAX 512 -of GTiff " + filePath + "/uncompressed.tiff" + " " + fileName
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// remove uncompressed GeoTIFF
	cmd = "rm " + filePath + "/uncompressed.tiff"
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// convert GeoTIFF from projection to WGS84
	/*cmd := conf["gdalwarp_path"] + " " + filePath + "/tmp.tiff" + " " + fileName + " -t_srs \"+proj=longlat +datum=WGS84 +ellps=WGS84\""
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}

	// remove temp GeoTIFF
	cmd = "rm " + filePath + "/tmp.tiff"
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}*/

	// return save result, (true if success, otherwise false)
	return true
}

// get the UTM zone from the EPSG code
func getUTMZone(EPSG int) int {
	if EPSG == 32632 {
		return 32
	} else if EPSG == 32635 {
		return 35
	}
	return 32
}

// save a map's GeoTIFFs to disk (calculate the geographic distance between points to determine the cluster size's length)
func saveGeoTIFFs(conf map[string]string, colorMap map[string]map[int]map[string]interface{}, filePath, mapName, metricName, date string, xLength, yLength float64) bool {
	var save bool

	// create folder if it does not exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, 0700)
	}

	db, err := sql.Open("mysql", conf["MySQL_username"]+":"+conf["MySQL_password"]+"@tcp("+conf["MySQL_hostname"]+":"+conf["MySQL_port"]+")/"+conf["MySQL_database"])

	// if there is an error opening the connection, handle it
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	var latitude float64
	var longitude float64
	var value float64

	// get the map's data
	results, err := db.Query("SELECT latitude, longitude, value FROM heatmap.data WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND clustered = " + conf["clustered"] + " AND date = '" + date + "'")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// save GeoTIFF for this row
	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&latitude, &longitude, &value)
		if err != nil {
			return false
			//log.Fatal(err)
			//panic(err.Error())
		}
		save = saveGeoTIFF(conf, colorMap, filePath, mapName, metricName, latitude, longitude, value, date, xLength, yLength, false)
	}
	// return save result, (true if success, otherwise false)
	return save
}

func saveGeoTIFFsNetCDFFile(conf map[string]string, colorMap map[string]map[int]map[string]interface{}, filePath, NetCDFFile, mapName, metricName, date string) bool {
	var latitude float32 = -100  // not set, define outside the range -90, 90
	var longitude float32 = -200 // not set, define outside the range -180, 180

	// create folder if it does not exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Mkdir(filePath, 0700)
	}

	// format the date
	dateString := strings.Replace(date, ":", "", -1)
	dateString = strings.Replace(dateString, "-", "", -1)
	dateString = strings.Replace(dateString, " ", "T", -1) + "Z"

	// format the filename
	fileName := filePath + "/" + dateString + ".tiff"

	// Open NetCDF file in read-only mode. The dataset is returned.
	ds, err := netcdf.OpenFile(NetCDFFile, netcdf.NOWRITE)
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	defer ds.Close()
	// get index's name
	var indexName string
	nVars, _ := ds.NVars()
	for i := 0; i < nVars; i++ {
		varn := ds.VarN(i)
		varName, _ := varn.Name()
		if varName != "crs" && varName != "time" && varName != "lat" && varName != "lon" {
			indexName = varName
		}
	}
	lat_v, err := ds.Var("lat")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	lon_v, err := ds.Var("lon")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	index_v, err := ds.Var(indexName)
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	// Read data from variables
	lat, err := netcdf.GetFloat32s(lat_v)
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	lon, err := netcdf.GetFloat32s(lon_v)
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	index, err := netcdf.GetFloat32s(index_v)
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	// Get the length of the dimensions of the data.
	index_dims, err := index_v.LenDims()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// calculate the bounding box
	var ymax float32 = -1 // max latitude
	var ymin float32 = -1 // min latitude
	var xmax float32 = -1 // max longitude
	var xmin float32 = -1 // min longitude
	// read the data
	for lat_index := 0; lat_index < int(index_dims[1]); lat_index++ {
		for lon_index := 0; lon_index < int(index_dims[2]); lon_index++ {
			if lat[lat_index] > ymax || ymax == -1 {
				ymax = lat[lat_index]
			}
			if lat[lat_index] < ymin || ymin == -1 {
				ymin = lat[lat_index]
			}
			if lon[lon_index] > xmax || xmax == -1 {
				xmax = lon[lon_index]
			}
			if lon[lon_index] < xmin || xmin == -1 {
				xmin = lon[lon_index]
			}
		}
	}

	db, err := sql.Open("mysql", conf["MySQL_username"]+":"+conf["MySQL_password"]+"@tcp("+conf["MySQL_hostname"]+":"+conf["MySQL_port"]+")/"+conf["MySQL_database"])

	// if there is an error opening the connection, handle it
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	// calculate the map's projection
	results, err := db.Query("SELECT projection FROM heatmap.metadata WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND clustered = " + conf["clustered"] + " AND date = '" + date + "'")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	var projection int

	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&projection)
	}
	if projection == 0 {
		projection = 4326
	}

	// make copies of latitudes and longitudes arrays and sort them
	lat_sorted := make([]float32, len(lat))
	lon_sorted := make([]float32, len(lon))
	copy(lat_sorted, lat)
	copy(lon_sorted, lon)

	// calculate the minimum delta latitude for this data set
	// minimum delta latitude between locations
	var yLength float32 = -1
	for lat_index := 0; lat_index < int(index_dims[1]); lat_index++ {
		for lon_index := 0; lon_index < int(index_dims[2]); lon_index++ {
			if latitude != -100 && latitude != lat_sorted[lat_index] {
				lat_delta := float32(math.Abs(float64(latitude) - float64(lat_sorted[lat_index])))
				// if the delta latitude is the minimum
				if yLength == -1 || lat_delta < yLength {
					yLength = lat_delta
				}
			}
			latitude = lat_sorted[lat_index]
		}
	}

	// calculate the minimum delta longitude for this data set
	// minimum delta longitude between locations
	var xLength float32 = -1
	for lat_index := 0; lat_index < int(index_dims[1]); lat_index++ {
		for lon_index := 0; lon_index < int(index_dims[2]); lon_index++ {
			if longitude != -200 && longitude != lon_sorted[lon_index] {
				lon_delta := float32(math.Abs(float64(longitude) - float64(lon_sorted[lon_index])))
				// if the delta longitude is the minimum
				if xLength == -1 || lon_delta < xLength {
					xLength = lon_delta
				}
			}
			longitude = lon_sorted[lon_index]
		}
	}

	// get map's pixel size
	ny := int((ymax-ymin)/yLength) + 1
	nx := int((xmax-xmin)/xLength) + 1

	// create each image's channel
	r_pixels := make([]uint8, ny*nx)
	g_pixels := make([]uint8, ny*nx)
	b_pixels := make([]uint8, ny*nx)
	alpha := make([]uint8, ny*nx)

	// get map's data from GRAL binary file
	// read NetCDF file
	//for time_index := 0; time_index < int(index_dims[0]); time_index++ {
	for time_index := 1; time_index < 2; time_index++ {
		for lat_index := 0; lat_index < int(index_dims[1]); lat_index++ {
			for lon_index := 0; lon_index < int(index_dims[2]); lon_index++ {
				// get metric name's color
				r, g, b := getColor(colorMap, metricName, float64(index[time_index*lat_index*lon_index]))
				// if there is no color map for this metric name, then return false
				if r == -1 {
					return false
					//log.Fatal(err)
					//panic(err.Error())
				}
				// set the pixel data
				loc := int32((lon[lon_index]-xmin)/xLength) + int32((ymax-lat[lat_index])/yLength)*int32(nx)
				r_pixels[loc] = uint8(r)
				g_pixels[loc] = uint8(g)
				b_pixels[loc] = uint8(b)
				alpha[loc] = 255
			}
		}
	}

	// set geotransform
	//xres := float64(xmax-xmin) / float64(nx)
	//yres := float64(ymax-ymin) / float64(ny)
	xres := float64(xmax-xmin+xLength) / float64(nx)
	yres := float64(ymax-ymin+yLength) / float64(ny)
	//geotransform := [6]float64{float64(xmin), xres, 0, float64(ymax), 0, -yres}
	geotransform := [6]float64{float64(xmin) - xres/2, xres, 0, float64(ymax) + yres/2, 0, -yres}

	driver, err := gdal.GetDriverByName("GTiff")
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	// create the 4-band raster file with transparency (r, g, b, alpha)
	dst_ds := driver.Create(filePath+"/tmp.tiff", nx, ny, 4, gdal.Byte, nil)
	//defer dst_ds.Close()

	spatialRef := gdal.CreateSpatialReference("") // establish encoding
	spatialRef.FromEPSG(projection)               // EPSG projection
	srString, err := spatialRef.ToWKT()
	dst_ds.SetProjection(srString)       // export coords to file
	dst_ds.SetGeoTransform(geotransform) // specify coords

	r_band := dst_ds.RasterBand(1)
	r_band.IO(gdal.Write, 0, 0, nx, ny, r_pixels, nx, ny, 0, 0)
	g_band := dst_ds.RasterBand(2)
	g_band.IO(gdal.Write, 0, 0, nx, ny, g_pixels, nx, ny, 0, 0)
	b_band := dst_ds.RasterBand(3)
	b_band.IO(gdal.Write, 0, 0, nx, ny, b_pixels, nx, ny, 0, 0)
	alpha_band := dst_ds.RasterBand(4)
	alpha_band.IO(gdal.Write, 0, 0, nx, ny, alpha, nx, ny, 0, 0) // set alpha channel as opaque (nodata = transparent)
	dst_ds.Close()

	// convert GeoTIFF from projection to WGS84
	cmd := conf["gdalwarp_path"] + " " + filePath + "/tmp.tiff" + " " + filePath + "/uncompressed.tiff -t_srs \"+proj=longlat +datum=WGS84 +ellps=WGS84\""
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// remove temp GeoTIFF
	cmd = "rm " + filePath + "/tmp.tiff"
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// compress GeoTIFF
	cmd = conf["gdal_translate_path"] + " -co compress=deflate -co zlevel=9 -co tiled=yes -co NUM_THREADS=ALL_CPUS --config GDAL_CACHEMAX 512 -of GTiff " + filePath + "/uncompressed.tiff" + " " + fileName
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}

	// remove uncompressed GeoTIFF
	cmd = "rm " + filePath + "/uncompressed.tiff"
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return false
		//log.Fatal(err)
		//panic(err.Error())
	}
	return true
}

// index this map into GeoServer
func indexMap(conf map[string]string, colorMap map[string]map[int]map[string]interface{}, archive, mapName, metricName string, filePath, date, geoTIFF string, xLength, yLength float64, projection, file int, fileType string) bool {
	var save bool
	// save map's GeoTIFFs (*.tiff) to disk

	/*if file == 1 && fileType == "NetCDF" { // if it is a NetCDF file
		save = saveGeoTIFFsNetCDFFile(conf, colorMap, filePath, NetCDFFile, mapName, metricName, date)
	} else */if projection != 4326 && projection != 0 && file == 1 { // if EPSG projection is not WGS84 and the map is on file, it is UTM and then create the whole image directly
		save = saveGeoTIFFDatasetFile(conf, colorMap, filePath, mapName, metricName, date)
	} else if projection != 4326 && projection != 0 { // if EPSG projection is not WGS84, it is UTM and then create the whole image directly
		save = saveGeoTIFFDataset(conf, colorMap, filePath, mapName, metricName, date)
	} else if xLength == 0 && yLength == 0 { // if xLength and yLength are = 0, save GeoTIFF estimating the cluster size
		save = saveGeoTIFFsDeltaLatLon(conf, colorMap, filePath, mapName, metricName, date)
	} else { // if xLength and yLength are != 0, save GeoTIFFs using the provided cluster size
		save = saveGeoTIFFs(conf, colorMap, filePath, mapName, metricName, date, xLength, yLength)
	}
	// if this map GeoTIFFs were not generated, then return false
	if !save {
		return false
	}

	// create the folder if it does not exist
	if _, err := os.Stat(filePath + "/merged"); os.IsNotExist(err) {
		os.Mkdir(filePath+"/merged", 0700)
	}

	// merge GeoTIFFs into one GeoTIFF (*.tiff) with gdal_merge.py
	files_to_mosaic, err := FilterDirsGlob(filePath, "/*.tiff")
	if err != nil {
		log.Fatal(err)
	}

	// if this map has > 1 GeoTIFF then merge them
	if len(files_to_mosaic) > 1 {
		// alternative method for merging GeoTIFFs: gdal_merge.py
		//command = "/home/debian/anaconda3/bin/python /home/debian/anaconda3/bin/gdal_merge.py -o " + filePath + "/merged/" + geoTIFF + " -of gtiff " + files_string
		/*cmd := conf["Python_path"] + " " + conf["Python_gdal_merge_path"] + " -o " + filePath + "/merged/" + geoTIFF + " -of gtiff " + filePath + "/*.tiff"
		_, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}*/
		//fmt.Println(out)

		// generate file list to be used with gdalbuildvrt
		f, err := os.OpenFile(filePath+"/input.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		for _, file_to_mosaic := range files_to_mosaic {
			if _, err = f.WriteString(file_to_mosaic + "\n"); err != nil {
				panic(err)
			}
		}
		/*cmd := "ls | grep \".*\\.tiff\" > " + filePath + "/input.txt"
		_, err := exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}*/

		// method for merging GeoTIFFs (no Python dependencies)
		// create a virtual mosaic of all GeoTIFFs
		//cmd := conf["gdalbuildvrt_path"] + " " + filePath + "/" + geoTIFF + ".vrt " + filePath + "/*.tiff"
		cmd := conf["gdalbuildvrt_path"] + " -input_file_list " + filePath + "/input.txt " + filePath + "/" + geoTIFF + ".vrt"
		_, err = exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(out)

		// convert the virtual mosaic to GeoTIFF
		//-co compress=deflate -co zlevel=9 -co tiled=yes
		//-co compress=lzw -co predictor=2 -co tiled=yes
		// https://gis.stackexchange.com/questions/241806/does-gdal-translate-support-multi-thread
		cmd = conf["gdal_translate_path"] + " -co compress=deflate -co zlevel=9 -co tiled=yes -co NUM_THREADS=ALL_CPUS --config GDAL_CACHEMAX 512 -of GTiff " + filePath + "/" + geoTIFF + ".vrt" + " " + filePath + "/merged/" + geoTIFF
		_, err = exec.Command("sh", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Println(out)
	} else { // else if this map has == 1 GeoTIFF then don't merge and copy the file to the /merged subfolder
		CopyFile(files_to_mosaic[0], filePath+"/merged/"+geoTIFF)
	}

	// check if this layer name is already present in GeoServer
	layer := conf["GeoServer_workspace"] + ":" + mapName
	isLayer := checkLayer(conf, layer)

	// compress merged GeoTIFF (and, if this is a new layer, .properties files) into archive
	compressFiles(conf, mapName, isLayer, filePath, archive)

	// send GeoTIFF archive to GeoServer
	sendGeoTIFF(conf, mapName, isLayer, filePath, archive)

	// temp, copy GeoTIFF archive into /home/debian/GeoTIFFs
	// create folder if does not exists
	/*dst = "/home/debian/GeoTIFFs"
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		os.Mkdir(dst, 0700)
	}
	CopyFile(filePath+"/merged/"+archive, dst+"/"+archive)*/

	// remove the folder
	os.RemoveAll(filePath)

	return true
}

// index maps into GeoServer
func indexMaps(conf map[string]string) {
	filePath := conf["data_folder"]

	db, err := sql.Open("mysql", conf["MySQL_username"]+":"+conf["MySQL_password"]+"@tcp("+conf["MySQL_hostname"]+":"+conf["MySQL_port"]+")/"+conf["MySQL_database"])

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	// get colors' maps
	colorMap := getColorsMaps(conf)

	// get data
	results, err := db.Query("SELECT a.map_name, a.metric_name, b.x_length, b.y_length, b.projection, b.file, b.fileType, a.date FROM heatmap.maps_completed a LEFT JOIN heatmap.metadata b ON a.map_name = b.map_name AND a.date = b.date WHERE a.completed = '1' AND a.indexed = '0' ORDER BY a.id DESC")
	//results, err := db.Query("SELECT a.map_name, a.metric_name, b.x_length, b.y_length, b.projection, b.file, b.fileType, a.date FROM heatmap.maps_completed a LEFT JOIN heatmap.metadata b ON a.map_name = b.map_name AND a.date = b.date WHERE a.map_name = 'GRALheatmapHelsinki6mPM' AND date(a.date) = '2019-04-26'")
	//results, err := db.Query("SELECT a.map_name, a.metric_name, b.x_length, b.y_length, b.projection, b.file, b.fileType, a.date FROM heatmap.maps_completed a LEFT JOIN heatmap.metadata b ON a.map_name = b.map_name AND a.date = b.date WHERE a.map_name = 'AirQualityPM10Average2HourHelsinkiJ'")
	if err != nil {
		panic(err.Error())
	}

	var mapName string
	var metricName string
	var xLength float64
	var yLength float64
	var date string
	var geoTIFFArchive string
	var geoTIFF string
	var projection int
	var file int
	var fileType []byte
	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&mapName, &metricName, &xLength, &yLength, &projection, &file, &fileType, &date)
		if err != nil {
			panic(err.Error())
		}
		// format date to be used in the file name
		dateString := strings.Replace(date, ":", "", -1)
		dateString = strings.Replace(dateString, "-", "", -1)
		dateString = strings.Replace(dateString, " ", "T", -1) + "Z"
		geoTIFF = dateString + ".tiff"
		geoTIFFArchive = mapName + "-" + dateString + ".zip"

		// index this map into GeoServer
		fmt.Println("Indexing map: " + mapName + " metric: " + metricName + " date: " + date)

		// get timestamp to concatenate to folder's name
		tmpPath := fmt.Sprintf("%d", time.Now().UnixNano())

		// index the map
		i := indexMap(conf, colorMap, geoTIFFArchive, mapName, metricName, filePath+"/"+mapName+"_"+tmpPath, date, geoTIFF, xLength, yLength, projection, file, string(fileType))

		// if this map was successfully indexed, then set this it as indexed into MySQL
		if i == true {
			setIndexedMap(conf, mapName, metricName, date, "1")
		} else { // else set is as not indexed due to an error
			setIndexedMap(conf, mapName, metricName, date, "-1")
		}
	}
}

// set a map as indexed into MySQL
func setIndexedMap(conf map[string]string, mapName, metricName, date, indexed string) {
	db, err := sql.Open("mysql", conf["MySQL_username"]+":"+conf["MySQL_password"]+"@tcp("+conf["MySQL_hostname"]+":"+conf["MySQL_port"]+")/"+conf["Mysql_database"])

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	stmt, err := db.Prepare(
		"UPDATE heatmap.maps_completed SET indexed = ? WHERE map_name = ? AND metric_name = ? AND date = ?",
	)
	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(indexed, mapName, metricName, date)
	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
}

// get the map's number of points NOT USED
func getMapPointsNumber(conf map[string]string, mapName, metricName, date string) int {
	db, err := sql.Open("mysql", conf["MySQL_username"]+":"+conf["MySQL_password"]+"@tcp("+conf["MySQL_hostname"]+":"+conf["MySQL_port"]+")/"+conf["MySQL_database"])

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	results, err := db.Query("SELECT COUNT(*) AS num FROM heatmap.data WHERE map_name = '" + mapName + "' AND metric_name = '" + metricName + "' AND date = '" + date + "'")

	if err != nil {
		return 0
	}

	var num int

	for results.Next() {
		// for each row, scan the result into variables
		err = results.Scan(&num)
		if err != nil {
			return 0
		}
	}

	return num
}

// read GRAL binary file
func readGRALFile(fileName string) {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	m := payload{}
	fileinfo, err := file.Stat()
	if err != nil {
		return
	}
	filesize := fileinfo.Size()
	// skip the first 4 bytes
	readNextBytes(file, 4)
	var i int64
	// read 12 bytes into struct [UTM Easting, UTM Northing, Value]
	for i = 4; i < filesize; i++ {
		data := readNextBytes(file, 12)
		buffer := bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.LittleEndian, &m)
		if err != nil {
			log.Fatal("binary.Read failed", err)
		}
	}
}

func readNextBytes(file *os.File, number int) []byte {
	bytes := make([]byte, number)
	_, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	return bytes
}

// this type represents a GRAL record [UTM Easting, UTM Northing, Value]
type payload struct {
	UTME  int32 // longitude field in MySQL database
	UTMN  int32 // latitude field in MySQL database
	Value float32
}

// find the minimum difference >0 between any two elements in an array
func findMinDiff(arr []int32) int32 {
	// Initialize difference as -1
	var diff int32 = -1

	// Find the min diff by comparing difference
	// of all possible pairs in given array
	for i := 0; i < len(arr)-1; i++ {
		for j := i + 1; j < len(arr); j++ {
			d := int32(math.Abs(float64(arr[i]) - float64(arr[j])))
			if d != 0 && (diff == -1 || d < diff) {
				diff = d
			}
		}
	}
	// Return min diff
	return diff
}

// find the minimum difference >0 between any two elements in an array
func findMinUTMDistance(arr1, arr2 []int32) int32 {
	// Initialize difference as -1
	var diff float64 = -1

	// Find the min diff by comparing difference
	// of all possible pairs in given array
	for i := 0; i < len(arr1)-1; i++ {
		for j := i + 1; j < len(arr1); j++ {
			d1 := math.Pow(float64(arr1[i])-float64(arr1[j]), 2)
			d2 := math.Pow(float64(arr2[i])-float64(arr2[j]), 2)
			d := math.Sqrt(d1 - d2)
			if d != 0 && (diff == -1 || d < diff) {
				diff = d
			}
		}
	}
	// Return min diff
	return int32(math.Round(diff))
}

// split a string array into chunks limited by lim
func split(buf []string, lim int) [][]string {
	var chunk []string
	chunks := make([][]string, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}

// write a file
func writeFile(filePath, filename, text string) {
	t := []byte(text)
	err := ioutil.WriteFile(filePath+"/"+filename, t, 0700)
	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
}

// get all file paths
func get_all_file_paths(directory string) []string {
	files, err := ioutil.ReadDir(directory)
	f := []string{}
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		f = append(f, directory+"/"+file.Name())
	}
	return f
}

func recursiveZip(pathToZip, destinationPath string) error {
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return err
	}
	myZip := zip.NewWriter(destinationFile)
	err = filepath.Walk(pathToZip, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, filepath.Dir(pathToZip))
		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	err = myZip.Close()
	if err != nil {
		return err
	}
	return nil
}

// compress files with zip
func compressFiles(conf map[string]string, mapName string, isLayer bool, filePath, zipFile string) {
	// write .properties files only if this map's layer does not exist in GeoServer
	if !isLayer {
		writeFile(filePath+"/merged", "datastore.properties",
			"SPI=org.geotools.data.postgis.PostgisNGJNDIDataStoreFactory\n"+
				"jndiReferenceName=jdbc/postgres"+"\n"+
				"Loose\\ bbox=true\n"+
				"preparedStatements=false")

		// write indexer.properties (contains PostgreSQL indexing properties)
		writeFile(filePath+"/merged", "indexer.properties",
			"TimeAttribute=ingestion\n"+
				"Schema=*the_geom:Polygon,location:String,ingestion:java.util.Date\n"+
				"PropertyCollectors=TimestampFileNameExtractorSPI[timeregex](ingestion)")

		// write timeregex.properties (contains the regular expression to parse the date from the TIFF's filename)
		writeFile(filePath+"/merged", "timeregex.properties", "regex=[0-9]{8}T[0-9]{6}Z")
	}
	// calling function to get all file paths in the directory
	//file_paths := get_all_file_paths(filePath)

	// writing files to a zipfile
	recursiveZip(filePath+"/merged/", filePath+"/"+zipFile)
	//fmt.Println("All files zipped successfully!")
}

// send GeoTIFF to GeoServer
func sendGeoTIFF(conf map[string]string, layer string, isLayer bool, filePath, filename string) {
	// create the layer
	// example with CURL: curl -v -u 'admin:password' -XPUT -H "Content-type:application/zip" --data-binary @init.zip http://localhost:8080/geoserver/rest/workspaces/sf/coveragestores/test_immos/file.imagemosaic
	// if this layer is not present in GeoServer, then make a PUT request for inserting the granule
	if !isLayer {
		f, err := os.Open(filePath + "/" + filename)
		if err != nil {
			// handle err
		}
		defer f.Close()
		req, err := http.NewRequest("PUT", conf["GeoServer_url"]+"/workspaces/"+conf["GeoServer_workspace"]+"/coveragestores/"+layer+"/file.imagemosaic?recalculate=nativebbox,latlonbbox", f)
		if err != nil {
			// handle err
		}
		req.SetBasicAuth(conf["GeoServer_username"], conf["GeoServer_password"])
		req.Header.Set("Content-Type", "application/zip")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	} else { // if this layer is present in GeoServer, then make a POST request for inserting the granule
		// example with CURL: curl -v -u 'admin:password' -XPOST -H "Content-type: application/zip" --data-binary @tiff.zip http://localhost:8080/geoserver/rest/workspaces/sf/coveragestores/test_immos/file.imagemosaic?recalculate=nativebbox,latlonbbox
		f, err := os.Open(filePath + "/" + filename)
		if err != nil {
			// handle err
		}
		defer f.Close()
		req, err := http.NewRequest("POST", conf["GeoServer_url"]+"/workspaces/"+conf["GeoServer_workspace"]+"/coveragestores/"+layer+"/file.imagemosaic?recalculate=nativebbox,latlonbbox", f)
		if err != nil {
			// handle err
		}
		req.SetBasicAuth(conf["GeoServer_username"], conf["GeoServer_password"])
		req.Header.Set("Content-Type", "application/zip")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}

	// if this layer is not present in GeoServer, then enable time dimension
	if !isLayer {
		// enable the time dimension
		// example with CURL: curl -v -u 'admin:password' -XPUT -H "Content-type:application/xml; charset=UTF-8" -d '<coverage><enabled>true</enabled><metadata><entry key="time"><dimensionInfo><enabled>true</enabled><presentation>LIST</presentation><units>ISO8601</units><defaultValue/></dimensionInfo></entry></metadata></coverage>' http://localhost:8080/geoserver/rest/workspaces/sf/coveragestores/test_immos/coverages/test_immos
		body := strings.NewReader(`<coverage><enabled>true</enabled><metadata><entry key="time"><dimensionInfo><enabled>true</enabled><presentation>LIST</presentation><units>ISO8601</units><defaultValue/></dimensionInfo></entry></metadata></coverage>`)
		req, err := http.NewRequest("PUT", conf["GeoServer_url"]+"/workspaces/"+conf["GeoServer_workspace"]+"/coveragestores/"+layer+"/coverages/"+layer, body)
		if err != nil {
			// handle err
		}
		req.SetBasicAuth(conf["GeoServer_username"], conf["GeoServer_password"])
		req.Header.Set("Content-Type", "application/xml; charset=UTF-8")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
		}
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
	}
}

// check if this layer name is already present in GeoServer
func checkLayer(conf map[string]string, layer string) bool {
	isLayer := false
	req, err := http.NewRequest("GET", conf["GeoServer_url"]+"/layers.json", nil)
	if err != nil {
		// handle err
	}
	req.SetBasicAuth(conf["GeoServer_username"], conf["GeoServer_password"])
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// get JSON result
	result := gjson.Get(string(body), "layers")
	result.ForEach(func(key, value gjson.Result) bool {
		//result1 := gjson.Get(value.String(), "layer").String()
		value.ForEach(func(key, value gjson.Result) bool {
			layerName := gjson.Get(value.String(), "name").String()
			if layerName == layer {
				isLayer = true
			}
			if isLayer {
				return false // stop iterating
			} else {
				return true // keep iterating
			}
		})
		if isLayer {
			return false // stop iterating
		} else {
			return true // keep iterating
		}
	})
	return isLayer
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func FilterDirs(dir, suffix string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), suffix) {
			res = append(res, filepath.Join(dir, f.Name()))
		}
	}
	return res, nil
}

func FilterDirsGlob(dir, suffix string) ([]string, error) {
	return filepath.Glob(filepath.Join(dir, suffix))
}

func appendToTextFile(filename, text string) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}

// make first character of a string upper case
func UcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// make first character of a string lower case
func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func main() {
	// Settings map
	conf := map[string]string{}
	// Default settings
	// GeoServer
	conf["GeoServer_username"] = "user"
	conf["GeoServer_password"] = "..."
	conf["GeoServer_url"] = "http://localhost:8080/geoserver/rest"
	conf["GeoServer_workspace"] = "Snap4City"
	// PostgreSQL
	conf["PostgreSQL_hostname"] = "localhost"
	conf["PostgreSQL_port"] = "5432"
	conf["PostgreSQL_database"] = "gisdata"
	conf["PostgreSQL_schema"] = "public"
	conf["PostgreSQL_username"] = "user"
	conf["PostgreSQL_password"] = "..."
	//MySQL
	conf["MySQL_hostname"] = "localhost"
	conf["MySQL_username"] = "user"
	conf["MySQL_password"] = "..."
	conf["MySQL_port"] = "3306"
	conf["MySQL_database"] = "heatmap"
	// Python paths
	conf["Python_path"] = "/usr/bin/python3"
	conf["Python_gdal_merge_path"] = "/usr/local/bin/gdal_merge.py"
	// GDAL paths
	conf["gdalbuildvrt_path"] = "/usr/local/bin/gdalbuildvrt"
	conf["gdal_translate_path"] = "/usr/local/bin/gdal_translate"
	conf["gdalwarp_path"] = "/usr/local/bin/gdalwarp"
	// if data is clustered
	conf["clustered"] = "0"
	// GRAL data folder
	conf["gral_data"] = "/home/debian/GRAL"
	// data folder where to write files
	// get the working directory
	filePath, err := os.Getwd()
	if err == nil {
		conf["data_folder"] = filePath + "/data"
	}

	// Custom settings
	// get c flag command line parameter
	c := flag.String("conf", "", "Configuration file path (JSON)")
	// parse flag
	flag.Parse()
	// don't use lowercase letter in struct members' initial letter, otherwise it does not work
	// https://stackoverflow.com/questions/24837432/golang-capitals-in-struct-fields
	type Configuration struct {
		GeoServerUsername   string
		GeoServerPassword   string
		GeoServerUrl        string
		GeoServerWorkspace  string
		PostgreSQLHostname  string
		PostgreSQLPort      string
		PostgreSQLDatabase  string
		PostgreSQLSchema    string
		PostgreSQLUsername  string
		PostgreSQLPassword  string
		MySQLHostname       string
		MySQLUsername       string
		MySQLPassword       string
		MySQLPort           string
		MySQLDatabase       string
		PythonPath          string
		PythonGdalMergePath string
		GdalBuildVrtPath    string
		GdalTranslatePath   string
		GdalWarpPath        string
		Clustered           string
		GRALData            string
		DataFolder          string
	}
	// if a configuration file (JSON) is specified as a command line parameter (-conf), then attempt to read it
	if *c != "" {
		configuration := Configuration{}
		file, err := os.Open(*c)
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&configuration)
		// if configuration file reading is ok, update the settings map
		if err == nil {
			// GeoServer
			conf["GeoServer_username"] = configuration.GeoServerUsername
			conf["GeoServer_password"] = configuration.GeoServerPassword
			conf["GeoServer_url"] = configuration.GeoServerUrl
			conf["GeoServer_workspace"] = configuration.GeoServerWorkspace
			// PostgreSQL
			conf["PostgreSQL_hostname"] = configuration.PostgreSQLHostname
			conf["PostgreSQL_port"] = configuration.PostgreSQLPort
			conf["PostgreSQL_database"] = configuration.PostgreSQLDatabase
			conf["PostgreSQL_schema"] = configuration.PostgreSQLSchema
			conf["PostgreSQL_username"] = configuration.PostgreSQLUsername
			conf["PostgreSQL_password"] = configuration.PostgreSQLPassword
			// MySQL
			conf["MySQL_hostname"] = configuration.MySQLHostname
			conf["MySQL_username"] = configuration.MySQLUsername
			conf["MySQL_password"] = configuration.MySQLPassword
			conf["MySQL_port"] = configuration.MySQLPort
			conf["MySQL_database"] = configuration.MySQLDatabase
			// Python
			conf["Python_path"] = configuration.PythonPath
			conf["Python_gdal_merge_path"] = configuration.PythonGdalMergePath
			// GDAL
			conf["gdalbuildvrt_path"] = configuration.GdalBuildVrtPath
			conf["gdal_translate_path"] = configuration.GdalTranslatePath
			conf["gdalwarp_path"] = configuration.GdalWarpPath
			// if data is clustered
			conf["clustered"] = configuration.Clustered
			// GRAL data folder
			conf["gral_data"] = configuration.GRALData
			// data folder where to write files
			conf["data_folder"] = configuration.DataFolder
		}
	}

	// create the data folder if it does not exist
	if _, err := os.Stat(conf["data_folder"]); os.IsNotExist(err) {
		os.Mkdir(conf["data_folder"], 0755)
	}
	// delete all data folder contents before start indexing (e.g. for ram disk cleanup)
	cmd := "rm -rf " + conf["data_folder"] + "/*"
	_, err = exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		log.Fatal(err)
	}

	if err == nil {
		indexMaps(conf)
	}
}
