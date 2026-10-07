package work

import (
	"encoding/csv"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strings"
)

func (s *Service) ExportProfit(c *gin.Context) {
	data, e := s.Profit(c.Request.Context(), c.Query("month"))
	if e != nil {
		respond(c, nil, e)
		return
	}
	var records []map[string]any
	if e = json.Unmarshal(data, &records); e != nil {
		respond(c, nil, e)
		return
	}
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", `attachment; filename="profit.csv"`)
	writer := csv.NewWriter(c.Writer)
	_ = writer.Write([]string{"tor_id", "project_name", "month", "revenue", "labor", "material", "profit", "estimated"})
	for _, record := range records {
		row := []string{}
		for _, key := range []string{"tor_id", "project_name", "month", "revenue", "labor", "material", "profit"} {
			value, _ := record[key].(string)
			if key == "project_name" && strings.ContainsAny(strings.TrimSpace(value)[:min(len(strings.TrimSpace(value)), 1)], "=+-@") {
				value = "'" + value
			}
			row = append(row, value)
		}
		row = append(row, "true")
		_ = writer.Write(row)
	}
	writer.Flush()
}
