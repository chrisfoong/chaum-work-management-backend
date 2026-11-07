package work

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

type ProfitConfirmation struct {
	TorID string `json:"tor_id"`
	Start string `json:"period_start"`
	End   string `json:"period_end"`
}

func payrollList(s *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, offset, e := paging(c)
		if e != nil {
			respond(c, nil, e)
			return
		}
		if month := c.Query("period_month"); month != "" {
			data, e := s.PayrollMonth(c.Request.Context(), actor(c), month, limit, offset)
			respond(c, gin.H{"data": data, "limit": limit, "offset": offset}, e)
			return
		}
		data, e := s.Payrolls(c.Request.Context(), actor(c), limit, offset)
		respond(c, gin.H{"data": data, "limit": limit, "offset": offset}, e)
	}
}
func (s *Service) profitRequest(c *gin.Context) (json.RawMessage, error) {
	start, end := c.Query("period_start"), c.Query("period_end")
	if start == "" && end == "" {
		month := c.Query("month")
		if e := Month(month); e != nil {
			return nil, e
		}
		d, _ := time.Parse("2006-01", month)
		start = d.Format("2006-01-02")
		end = d.AddDate(0, 1, -1).Format("2006-01-02")
	}
	return s.ProfitRange(c.Request.Context(), c.Query("tor_id"), start, end)
}
func (s *Service) ExportProfitPDF(c *gin.Context) {
	data, e := s.profitRequest(c)
	if e != nil {
		respond(c, nil, e)
		return
	}
	pdf, e := profitPDF(data)
	if e != nil {
		respond(c, nil, e)
		return
	}
	c.Header("Content-Disposition", `attachment; filename="chaum-profit.pdf"`)
	c.Data(200, "application/pdf", pdf)
}
