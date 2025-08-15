package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ToJSON(data interface{}) string {
	bytes, _ := json.MarshalIndent(data, "", "  ")
	return string(bytes)
}

func ParseFilterParams(c *gin.Context) (int, int, string, string, string) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	search := c.Query("search")

	return page, size, startDate, endDate, search
}

func FormatRupiah(amount float64) string {
	rp := fmt.Sprintf("Rp %s", formatWithThousandSeparator(amount))
	return rp
}

func formatWithThousandSeparator(n float64) string {
	s := fmt.Sprintf("%.0f", n) // no decimal places
	var result string
	count := 0

	for i := len(s) - 1; i >= 0; i-- {
		result = string(s[i]) + result
		count++
		if count%3 == 0 && i != 0 {
			result = "." + result
		}
	}
	return result
}
