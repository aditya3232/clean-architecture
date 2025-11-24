package conv

import "strconv"

func LatLngToString(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}

func StringToInt64(s string) (int64, error) {
	newData, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}

	return newData, nil
}
