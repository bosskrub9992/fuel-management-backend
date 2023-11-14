package main

import (
	"log"
	"time"

	"github.com/bosskrub9992/fuel-management/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shopspring/decimal"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Static("/static", "./internal/templates/static")

	api := e.Group("/api")
	api.Use(middleware.CORS())

	api.GET("/fuel/usages", func(c echo.Context) error {
		return c.JSON(200, models.GetCarFuelUsageData{
			TodayDate:               time.Now().Format("2006-01-02"),
			LatestKilometerAfterUse: 0,
			LatestFuelPrice:         decimal.NewFromFloat(1.41),
			AllUsers: []models.User{
				{
					ID:       1,
					Nickname: "Boss",
				},
				{
					ID:       2,
					Nickname: "Best",
				},
				{
					ID:       3,
					Nickname: "Nut",
				},
				{
					ID:       4,
					Nickname: "Pat",
				},
			},
			Data: []models.CarFuelUsageDatum{
				{
					ID:                 1,
					FuelUseDate:        time.Now().Format("2006-01-02"),
					FuelPrice:          decimal.NewFromFloat(1.41),
					KilometerBeforeUse: 900,
					KilometerAfterUse:  800,
					Description:        "dinner",
					TotalMoney:         decimal.NewFromFloat(100),
					FuelUsers:          "Boss, Best",
				},
			},
			AllCars: []models.Car{
				{
					ID:   1,
					Name: "Mazda 2",
				},
				{
					ID:   2,
					Name: "Ford Fiesta",
				},
			},
			TotalRecord: 1,
			CurrentCar: models.Car{
				ID:   1,
				Name: "Mazda 2",
			},
			CurrentUser: models.UserWithImageURL{
				User: models.User{
					ID:       1,
					Nickname: "Boss",
				},
				UserImageURL: "http://localhost:8080/static/profile_image/BOSS.PNG",
			},
		})
	})

	api.GET("/users", func(c echo.Context) error {
		return c.JSON(200, models.GetUserData{
			Data: []models.GetUserDatum{
				{
					ID:              1,
					DefaultCarID:    1,
					Nickname:        "Boss",
					ProfileImageURL: "http://localhost:8080/static/profile_image/BOSS.PNG",
				},
				{
					ID:              2,
					DefaultCarID:    1,
					Nickname:        "Best",
					ProfileImageURL: "http://localhost:8080/static/profile_image/BEST.PNG",
				},
				{
					ID:              3,
					DefaultCarID:    2,
					Nickname:        "Nut",
					ProfileImageURL: "http://localhost:8080/static/profile_image/NUT.PNG",
				},
				{
					ID:              4,
					DefaultCarID:    1,
					Nickname:        "Pat",
					ProfileImageURL: "http://localhost:8080/static/profile_image/PAT.PNG",
				},
			},
		})
	})

	if err := e.Start(":8080"); err != nil {
		log.Println(err)
		return
	}
}
