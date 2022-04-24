package utils

import (
	"math"
	"strconv"
)

var (
	suffixes [5]string
)

func round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func GetFileSize(sizeString string) string {
	if sizeString == "" {
		return ""
	}
	size, err := strconv.ParseFloat(sizeString, 64)
	if err != nil {
		return ""
	} // This is in bytes
	suffixes[0] = "B"
	suffixes[1] = "KB"
	suffixes[2] = "MB"
	suffixes[3] = "GB"
	suffixes[4] = "TB"

	base := math.Log(size) / math.Log(1024)
	getSize := round(math.Pow(1024, base-math.Floor(base)), .5, 2)
	baseV := int(math.Floor(base))
	if baseV < 0 || baseV > 5 {
		return ""
	}
	getSuffix := suffixes[baseV]
	return strconv.FormatFloat(getSize, 'f', -1, 64) + " " + string(getSuffix)
}
