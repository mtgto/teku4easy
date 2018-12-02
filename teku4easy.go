package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"math"
	"os"
	"strconv"
)

// 大字の定義
type oaza struct {
	name string
	city string
	pos position
}

type position struct {
	latitude float64
	longitude float64
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "引数に大字のCSVファイルを渡してください\n")
		os.Exit(1)
	}
	oazas := loadCsv(flag.Arg(0))
	//fmt.Printf("%v\n", oazas)
	minLat, minLong := 10000.0, 10000.0
	maxLat, maxLong := 0.0, 0.0
	for _, oaza := range oazas {
		minLat = math.Min(minLat, oaza.pos.latitude)
		maxLat = math.Max(maxLat, oaza.pos.latitude)
		minLong = math.Min(minLong, oaza.pos.longitude)
		maxLong = math.Max(maxLong, oaza.pos.longitude)
	}
	pos, result := findMostCongested(&oazas, 0.01, 0.01, position{latitude: minLat, longitude: minLong}, position{latitude: maxLat, longitude: maxLong})
	fmt.Printf("緯度 %v, 経度 %v 大字の数%v\n", pos.latitude, pos.longitude, len(result))
	for _, oaza := range result {
		fmt.Printf("%v, %v (%v, %v)\n", oaza.city, oaza.name, oaza.pos.latitude, oaza.pos.longitude)
	}
}

// min - maxの矩形内で、width, height (単位は緯度経度) 以内に一番多くの大字を含む地点とその数を返す
func findMostCongested(oazas *[]oaza, width, height float64, minPos, maxPos position) (position, []oaza) {
	var bestPos position
	var bestResult []oaza = make([]oaza, 0, 0)
	for lat := minPos.latitude; lat < maxPos.latitude; lat += width {
		for long := minPos.longitude; long < maxPos.longitude; long += height {
			result := make([]oaza, 0)
			for _, oaza := range *oazas {
				if lat - width <= oaza.pos.latitude && oaza.pos.latitude <= lat + width && long - height <= oaza.pos.longitude && oaza.pos.longitude <= long + height {
					result = append(result, oaza)
				}
			}
			if len(bestResult) < len(result) {
				//fmt.Printf("%v, %v, %v\n", lat, long, result)
				bestResult = result
				bestPos = position{latitude: lat, longitude: long}
			}
		}
	}
	return bestPos, bestResult
}

// 大字CSVデータを読み込む。ShiftJISになっているのとヘッダ行があるので処理する
func loadCsv(in string) []oaza {
	// 今回は離島は計算しない
	count := 5283
	oazas := make([]oaza, 0, count)
	file, err := os.Open(in)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(transform.NewReader(file, japanese.ShiftJIS.NewDecoder()))
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	for i, record := range records {
		if i >= 1 && i < count {
			latitude, err := strconv.ParseFloat(record[6], 64)
			if err != nil {
				panic(err)
			}
			longitude, err := strconv.ParseFloat(record[7], 64)
			if err != nil {
				panic(err)
			}
			oazas = append(oazas, oaza{name: record[5], city: record[3], pos: position{latitude: latitude, longitude: longitude}})
		}
	}
	return oazas
}
