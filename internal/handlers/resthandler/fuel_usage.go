package resthandler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type FuelUsageData struct {
	FuelUsageRecords  []FuelUsageRecord
	TodayDate         string
	LatestEndFuelLeft int
	LatestFuelPrice   float64
}

type FuelUsageRecord struct {
	FuelUsedDate  string
	PeopleUsed    string
	Description   string
	StartFuelLeft int
	EndFuelLeft   int
	Total         string
}

func (h *RESTHandler) FuelUsage(c echo.Context) error {
	fuelUsageData := FuelUsageData{
		FuelUsageRecords: []FuelUsageRecord{
			{
				FuelUsedDate:  time.Now().Format("2006-01-02"),
				PeopleUsed:    "นัท แพท บอส เบส",
				Description:   "กินข้าว",
				StartFuelLeft: 800,
				EndFuelLeft:   700,
				Total:         "100.00",
			},
			{
				FuelUsedDate:  time.Now().Format("2006-01-02"),
				PeopleUsed:    "นัท แพท บอส เบส",
				Description:   "กินข้าว",
				StartFuelLeft: 800,
				EndFuelLeft:   700,
				Total:         "100.00",
			},
			{
				FuelUsedDate:  time.Now().Format("2006-01-02"),
				PeopleUsed:    "นัท แพท บอส เบส",
				Description:   "กินข้าว",
				StartFuelLeft: 800,
				EndFuelLeft:   700,
				Total:         "100.00",
			},
		},
		TodayDate:         time.Now().Format("2006-01-02"),
		LatestEndFuelLeft: 900,
		LatestFuelPrice:   1.31,
	}
	return c.Render(http.StatusOK, "fuel-usage", fuelUsageData)
}
