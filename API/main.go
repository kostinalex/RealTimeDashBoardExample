package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var maxNumberOfReadingsPerCall int = 50
var numberOfMonthsToSeed time.Duration = 3

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	formatted := time.Time(t).Format("Jan 02, 15:04")
	return []byte(`"` + formatted + `"`), nil
}

type ReadingDTO struct {
	Date        JSONTime `json:"date"`
	Temperature float64  `json:"temperature"`
}

type ReadingResponse struct {
	Name     string       `json:"name"`
	Readings []ReadingDTO `json:"readings"`
}

type Sensor struct {
	ID          string    `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Readings    []Reading `gorm:"foreignKey:SensorID" json:"readings,omitempty"`
	Temperature float64   `gorm:"type:decimal(5,1)" json:"temperature"`
}

type Reading struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SensorID    string    `gorm:"index;not null" json:"sensorId"`
	Date        time.Time `json:"date"`
	Temperature float64   `gorm:"type:decimal(5,1)" json:"temperature"`
	Sensor      Sensor    `gorm:"constraint:OnDelete:CASCADE;foreignKey:SensorID;references:ID" json:"-"`
}

func main() {

	dsn := "host=" + os.Getenv("HOST") + " user=alex password=KKzj5NP6CkvvuMqTClIL dbname=testdb port=" + os.Getenv("PORT") + " sslmode=disable TimeZone=UTC"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Println("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&Sensor{}, &Reading{}); err != nil {
		log.Fatal("Migration failed:", err)
	}

	var count int64
	db.Model(&Sensor{}).Count(&count)

	if count == 0 {
		sensors := []Sensor{
			{ID: "abc1", Name: "Sensor 1", Temperature: 40},
			{ID: "abc2", Name: "Sensor 2", Temperature: 15},
		}
		if err := db.Create(&sensors).Error; err != nil {
			log.Println("Failed to insert initial sensors:", err)
		}
		log.Println("Inserted default sensors")
	}

	seedReadings(db)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/delete", func(c echo.Context) error {

		if err := db.Exec("truncate table public.readings;").Error; err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "Failed to delete readings",
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	e.GET("/sensors", func(c echo.Context) error {

		var sensors []Sensor
		if err := db.Find(&sensors).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query sensors"})
		}

		return c.JSON(http.StatusOK, sensors)
	})

	e.GET("/readings/:sensorId/:start/:end", func(c echo.Context) error {
		sensorID := c.Param("sensorId")
		startStr := c.Param("start")
		endStr := c.Param("end")

		start, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid start time"})
		}

		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid end time"})
		}

		var sensor Sensor
		if err := db.Where("id = ?", sensorID).
			Find(&sensor).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query sensor"})
		}

		var allReadings []Reading
		if err := db.Where("sensor_id = ? and date >= ? and date <= ?", sensorID, start, end).
			Order("date").
			Find(&allReadings).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to query readings"})
		}

		var response []ReadingDTO

		if len(allReadings) <= maxNumberOfReadingsPerCall {
			for _, r := range allReadings {
				response = append(response, ReadingDTO{
					Date:        JSONTime(r.Date),
					Temperature: r.Temperature,
				})
			}
			return c.JSON(http.StatusOK, ReadingResponse{Name: sensor.Name, Readings: response})
		}

		//reduce the number of readings (bad idea, for testing only)
		step := float64(len(allReadings)) / float64(maxNumberOfReadingsPerCall)
		for i := 0; i < maxNumberOfReadingsPerCall; i++ {
			idx := int(float64(i) * step)
			if idx >= len(allReadings) {
				idx = len(allReadings) - 1
			}
			response = append(response, ReadingDTO{
				Date:        JSONTime(allReadings[idx].Date),
				Temperature: allReadings[idx].Temperature,
			})
		}

		return c.JSON(http.StatusOK, ReadingResponse{Name: sensor.Name, Readings: response})
	})

	e.POST("/readings", func(c echo.Context) error {
		var readings []Reading
		if err := c.Bind(&readings); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "JSON parse failed"})
		}

		if len(readings) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "No readings"})
		}

		if err := db.Create(&readings).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to insert readings"})
		}

		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"version": "1.0"})
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func seedReadings(db *gorm.DB) {
	var count int64
	db.Model(&Reading{}).Where("date > ?", time.Now().Add(-5*time.Hour)).Count(&count)

	if count > 0 {
		log.Println("Recent readings found. Skipping seeding.")
		return
	}

	log.Println("No recent readings found. Seeding last " + numberOfMonthsToSeed.String() + " months of readings...")

	sensors := []string{"abc1", "abc2"}
	start := time.Now().Add(-(numberOfMonthsToSeed * 30 * 24) * time.Hour)
	end := time.Now()

	interval := 5 * time.Second
	batchSize := 1000

	for _, sensorID := range sensors {
		log.Println("Seeding readings for sensor with ID", sensorID)
		readings := make([]Reading, 0, batchSize)
		current := start

		var prevTemp float64 = 10

		for current.Before(end) {

			delta := float64(rand.Intn(41)-20) / 10.0
			temp := prevTemp + delta

			if temp < 10 {
				temp = 10
			}
			if temp > 40 {
				temp = 40
			}

			temp = float64(int(temp*10)) / 10.0
			prevTemp = temp

			readings = append(readings, Reading{
				SensorID:    sensorID,
				Date:        current,
				Temperature: temp,
			})

			if len(readings) >= batchSize {
				if err := db.Create(&readings).Error; err != nil {
					log.Fatal("Failed inserting readings:", err)
				}
				readings = readings[:0]
			}

			current = current.Add(interval)
		}

		if len(readings) > 0 {
			if err := db.Create(&readings).Error; err != nil {
				log.Println("Failed inserting final batch:", err)
			}
		}
	}

	log.Println("Finished seeding readings.")
}
