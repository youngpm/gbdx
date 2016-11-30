package cmd

import (
	"bytes"
	"strings"
	"sort"
	"strconv"
)

func getSearchAreaFromArgs(inVal []string) (string) {
	subStr := "POLYGON(("
	suffixStr := "))"
	for i := 0; i < len(inVal); i += 1 {
		if strings.Contains(inVal[i], subStr) {
			var wkt bytes.Buffer
			for j := i; j < len(inVal); j += 1 {
				wkt.WriteString(inVal[j])
				if strings.HasSuffix(inVal[j], suffixStr) {
					return wkt.String()
				} else {
					wkt.WriteString(" ")
				}
			}
		}
	}
	return ""
}

func getFilterFromArgs(inVal []string) (string) {
	if len(inVal) == 5 {
		return inVal[4]
	} else if len(inVal) > 5 {
        	suffixStr := "))"
		for i := 5; i < len(inVal); i += 1 {
			if strings.HasSuffix(inVal[i],suffixStr) {
				var filter bytes.Buffer
				for j := i + 1; j < len(inVal); j += 1 {
					filter.WriteString(inVal[j])
                        		if j != len(inVal) - 1 {
						filter.WriteString(" ")
					}
				}
				return filter.String()
       	                 }
		}
	}
	return ""
}

func polygon_from_bounds(X0 float64, Y0 float64, X1 float64, Y1 float64) (string) {
	prefixStr := "POLYGON(("
        suffixStr := "))"

	x0 := strconv.FormatFloat(X0, 'f', 8, 64)
	y0 := strconv.FormatFloat(Y0, 'f', 8, 64)
	x1 := strconv.FormatFloat(X1, 'f', 8, 64)
	y1 := strconv.FormatFloat(Y1, 'f', 8, 64)

	var wkt bytes.Buffer
	wkt.WriteString(prefixStr)
	wkt.WriteString(x0)
	wkt.WriteString(" ")
	wkt.WriteString(y1)
	wkt.WriteString(",")
	wkt.WriteString(x1)
        wkt.WriteString(" ")
        wkt.WriteString(y1)
        wkt.WriteString(",")
	wkt.WriteString(x1)
        wkt.WriteString(" ")
        wkt.WriteString(y0)
        wkt.WriteString(",")
	wkt.WriteString(x0)
        wkt.WriteString(" ")
        wkt.WriteString(y0)
        wkt.WriteString(",")
	wkt.WriteString(x0)
        wkt.WriteString(" ")
        wkt.WriteString(y1)
	wkt.WriteString(suffixStr)
	return wkt.String()
}

func getSmallerSearches(searchAreaWkt string) ([]string) {
	if len(searchAreaWkt) == 0 {
		return []string{}
	}

	// the size in degrees of the side that we will search
	D := 1.4

	str_points := strings.Trim(searchAreaWkt, "POLYGON((")
	str_points = strings.Trim(str_points, "))")
	arr_points := strings.Split(str_points, ",")

	Xs := make(map[string]float64)
	Ys := make(map[string]float64)
	for i := 0; i < len(arr_points); i += 1 {
		s := strings.Split(arr_points[i], " ")
		_, exist := Xs[s[0]]
		if exist == false {
			fnum, err := strconv.ParseFloat(s[0], 64)
			if err == nil {
				Xs[s[0]] = fnum
			}
		}
		_, exist = Ys[s[1]]
                if exist == false {
			fnum, err := strconv.ParseFloat(s[1], 64)
			if err == nil {
				Ys[s[1]] = fnum
			}
                }
	}

	x_val := make([]float64, 0, len(Xs))
	y_val := make([]float64, 0, len(Ys))
	if len(Xs) == 2 && len(Ys) == 2 {
		for _, value := range Xs {
			x_val = append(x_val, value)
		}
		for _, value := range Ys {
                        y_val = append(y_val, value)
                }
		sort.Float64s(x_val)
		sort.Float64s(y_val)
		if (x_val[1] - x_val[0]) > D {
			index := int((x_val[1] - x_val[0]) / D)
			multiplier := (x_val[1] - x_val[0]) / float64(index + 1)
			for i := 1; i <= index; i += 1 {
				x_val = append(x_val, x_val[0] + float64(i) * multiplier)
			}
			sort.Float64s(x_val)
		}
		if (y_val[1] - y_val[0]) > D {
			index := int((y_val[1] - y_val[0]) / D)
                        multiplier := (y_val[1] - y_val[0]) / float64(index + 1)
                        for i := 1; i <= index; i += 1 {
                                y_val = append(y_val, y_val[0] + float64(i) * multiplier)
                        }
			sort.Float64s(y_val)
                }
	} else {
		// geometry package needed to split parallelogram or multi-sided polygon so return
		return []string{searchAreaWkt}
	}

	// return search polygons
	polygons := make([]string, 0, 1)
	for y := 0; y < len(y_val) - 1; y += 1 {
		for x := 0; x < len(x_val) - 1; x += 1 {
			subsearchpoly := polygon_from_bounds(x_val[x], y_val[y], x_val[x+1], y_val[y+1])
			polygons = append(polygons, subsearchpoly)
		}
	}

	return polygons
}
