package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	AeroHost = "aerodatabox.p.rapidapi.com"
	AeroKey  = "c83526ca45msheb18036b37ecbd2p16357ejsn2938be1e64c0"
)

type FlightNumber struct {
	Default string
}

func limitFlights(flights []AeroFlight, max int) []AeroFlight {
	if len(flights) <= max {
		return flights
	}
	return flights[:max]
}

func (fn *FlightNumber) UnmarshalJSON(b []byte) error {

	if len(b) > 0 && b[0] == '"' {
		return json.Unmarshal(b, &fn.Default)
	}

	var temp struct {
		Default string `json:"default"`
	}
	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	fn.Default = temp.Default
	return nil
}

type AeroFlight struct {
	Number FlightNumber `json:"number"`
	Status string       `json:"status"`

	Airline struct {
		Name string `json:"name"`
	} `json:"airline"`

	Departure struct {
		ScheduledTime struct {
			UTC   string `json:"utc"`
			Local string `json:"local"`
		} `json:"scheduledTime"`

		RevisedTime struct {
			UTC   string `json:"utc"`
			Local string `json:"local"`
		} `json:"revisedTime"`

		Gate    string `json:"gate"`
		Airport struct {
			IATA string `json:"iata"`
		} `json:"airport"`
	} `json:"departure"`

	Arrival struct {
		ScheduledTime struct {
			UTC   string `json:"utc"`
			Local string `json:"local"`
		} `json:"scheduledTime"`

		Gate    string `json:"gate"`
		Airport struct {
			IATA string `json:"iata"`
		} `json:"airport"`
	} `json:"arrival"`
}

func mapAeroFlights(arr []AeroFlight, dep []AeroFlight) gin.H {

	var mappedArrivals []gin.H
	var mappedDepartures []gin.H

	for _, f := range arr {
		mappedArrivals = append(mappedArrivals, gin.H{
			"number": f.Number.Default,
			"airline": gin.H{
				"name": f.Airline.Name,
			},
			"departure": gin.H{
				"airport": gin.H{
					"iata": f.Departure.Airport.IATA,
				},
			},
			"status": f.Status,
			"gate":   f.Arrival.Gate,
			"arrival": gin.H{
				"scheduledTimeLocal": f.Arrival.ScheduledTime.Local,
			},
		})
	}

	for _, f := range dep {
		mappedDepartures = append(mappedDepartures, gin.H{
			"number": f.Number.Default,
			"airline": gin.H{
				"name": f.Airline.Name,
			},
			"arrival": gin.H{
				"airport": gin.H{
					"iata": f.Arrival.Airport.IATA,
				},
			},
			"status": f.Status,
			"gate":   f.Departure.Gate,
			"departure": gin.H{
				"scheduledTimeLocal": f.Departure.ScheduledTime.Local,
			},
		})
	}

	return gin.H{
		"arrivals":   mappedArrivals,
		"departures": mappedDepartures,
	}
}

func GetAeroFlights(c *gin.Context) {

	iata := strings.ToUpper(strings.TrimSpace(c.Query("iata")))
	if iata == "" {
		c.JSON(400, gin.H{"error": "Missing IATA code"})
		return
	}

	url := fmt.Sprintf(
		"https://%s/flights/airports/iata/%s?offsetMinutes=-120&durationMinutes=720&withLeg=true&direction=Both&withCancelled=true&withCodeshared=true&withCargo=true&withPrivate=true&withLocation=false",
		AeroHost, iata,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Key", AeroKey)
	req.Header.Add("X-RapidAPI-Host", AeroHost)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "AeroDataBox request failed"})
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{"error": "AeroDataBox returned an error", "body": bodyString})
		return
	}

	var data struct {
		Arrivals   []AeroFlight `json:"arrivals"`
		Departures []AeroFlight `json:"departures"`
	}

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		c.JSON(500, gin.H{"error": "JSON decode error", "body": bodyString})
		return
	}

	arrivalsLimited := limitFlights(data.Arrivals, 20)
	departuresLimited := limitFlights(data.Departures, 20)

	mapped := mapAeroFlights(arrivalsLimited, departuresLimited)

	c.JSON(200, mapped)
}

func GetAirlineName(c *gin.Context) {
	code := strings.ToUpper(strings.TrimSpace(c.Query("airline_code")))
	if code == "" {
		c.JSON(400, gin.H{"error": "Missing airline code"})
		return
	}

	url := fmt.Sprintf(
		"https://%s/flights/number/%s?withAircraftImage=false&withLocation=false",
		AeroHost, code,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Key", AeroKey)
	req.Header.Add("X-RapidAPI-Host", AeroHost)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "AeroDataBox request failed"})
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{
			"error": "AeroDataBox error",
			"body":  bodyString,
		})
		return
	}

	var flights []struct {
		Airline struct {
			Name string `json:"name"`
			IATA string `json:"iata"`
			ICAO string `json:"icao"`
		} `json:"airline"`
	}

	if err := json.Unmarshal(bodyBytes, &flights); err != nil {
		c.JSON(500, gin.H{"error": "JSON decode error", "body": bodyString})
		return
	}

	if len(flights) == 0 {
		c.JSON(404, gin.H{"error": "Flight not found"})
		return
	}

	c.JSON(200, gin.H{
		"airline_name": flights[0].Airline.Name,
		"iata":         flights[0].Airline.IATA,
		"icao":         flights[0].Airline.ICAO,
	})
}

func SearchFlights(c *gin.Context) {
	flightNumber := strings.ToUpper(strings.TrimSpace(c.Query("flight")))

	if flightNumber == "" {
		c.JSON(400, gin.H{"error": "Missing flight number or date"})
		return
	}

	url := fmt.Sprintf(
		"https://%s/flights/number/%s",
		AeroHost, flightNumber,
	)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Key", AeroKey)
	req.Header.Add("X-RapidAPI-Host", AeroHost)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": "AeroDataBox request failed"})
		return
	}
	defer resp.Body.Close()

	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == 204 {
		c.JSON(200, gin.H{"flights": []gin.H{}})
		return
	}

	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{
			"error": string(bodyBytes),
		})
		return
	}

	var flights []struct {
		Number  string `json:"number"`
		Status  string `json:"status"`
		Airline struct {
			Name string `json:"name"`
		} `json:"airline"`
		Aircraft struct {
			Model string `json:"model"`
			Reg   string `json:"reg"`
		} `json:"aircraft"`
		Departure struct {
			Airport struct {
				IATA string `json:"iata"`
				Name string `json:"name"`
			} `json:"airport"`
			ScheduledTime struct {
				Local string `json:"local"`
			} `json:"scheduledTime"`
		} `json:"departure"`
		Arrival struct {
			Airport struct {
				IATA string `json:"iata"`
				Name string `json:"name"`
			} `json:"airport"`
			ScheduledTime struct {
				Local string `json:"local"`
			} `json:"scheduledTime"`
		} `json:"arrival"`
	}

	json.Unmarshal(bodyBytes, &flights)

	var result []gin.H
	for _, f := range flights {
		result = append(result, gin.H{
			"number":   f.Number,
			"airline":  f.Airline.Name,
			"from":     f.Departure.Airport.IATA,
			"fromName": f.Departure.Airport.Name,
			"to":       f.Arrival.Airport.IATA,
			"toName":   f.Arrival.Airport.Name,
			"time":     f.Departure.ScheduledTime.Local,
			"arrival":  f.Arrival.ScheduledTime.Local,
			"aircraft": f.Aircraft.Model,
			"reg":      f.Aircraft.Reg,
			"status":   f.Status,
		})
	}

	c.JSON(200, gin.H{"flights": result})
}
