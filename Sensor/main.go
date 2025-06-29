package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Reading struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	SensorID    string    `gorm:"index;not null" json:"sensorId"`
	Date        time.Time `json:"date"`
	Temperature float64   `gorm:"type:decimal(5,1)" json:"temperature"`
}

func main() {
	url := os.Getenv("URL") + "/readings"

	var prevTemp float64 = 10

	for {

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

		payload := []Reading{}

		payload = append(payload, Reading{
			SensorID:    os.Getenv("SENSOR_ID"),
			Date:        time.Now().UTC(),
			Temperature: temp,
		})

		jsonData, err := json.Marshal(payload)
		if err != nil {
			log.Println("JSON marshal error:", err)
			continue
		}

		resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println("POST failed:", err)
		} else {
			log.Println("Status Code:", resp.StatusCode)
			resp.Body.Close()
		}

		time.Sleep(5 * time.Second)
	}
}
