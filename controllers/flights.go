package controllers

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	selectedPlanes []string
	selectionOnce  sync.Once
	selectionMutex sync.Mutex
)

type AdsbResponse struct {
	AC []map[string]interface{} `json:"ac"`
}

func fetchAllPlanes() (*AdsbResponse, error) {
	url := "https://api.adsb.lol/v2/lat/0/lon/0/dist/20000"

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data AdsbResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func selectPlanesOnce(all *AdsbResponse) {
	selectionOnce.Do(func() {

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(all.AC), func(i, j int) {
			all.AC[i], all.AC[j] = all.AC[j], all.AC[i]
		})

		for _, p := range all.AC {
			if hex, ok := p["hex"].(string); ok {
				selectedPlanes = append(selectedPlanes, hex)
				if len(selectedPlanes) == 200 {
					break
				}
			}
		}
	})
}

func GetFlights(c *gin.Context) {

	allPlanes, err := fetchAllPlanes()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch ADSB planes"})
		return
	}

	selectPlanesOnce(allPlanes)

	filtered := []map[string]interface{}{}

	selectionMutex.Lock()
	for _, ac := range allPlanes.AC {
		if hex, ok := ac["hex"].(string); ok {
			for _, selected := range selectedPlanes {
				if hex == selected {
					filtered = append(filtered, ac)
					break
				}
			}
		}
	}
	selectionMutex.Unlock()

	c.JSON(200, gin.H{
		"ac": filtered,
	})
}

func GetPlaneImage(c *gin.Context) {
	reg := c.Param("reg")
	reg = strings.TrimSpace(strings.ToUpper(reg))

	apiURL := "https://api.planespotters.net/pub/photos/reg/" + reg

	resp, err := http.Get(apiURL)
	if err != nil {
		c.JSON(200, gin.H{
			"image": "https://via.placeholder.com/400x200?text=No+Image",
		})
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if photos, ok := result["photos"].([]interface{}); ok && len(photos) > 0 {
		if first, ok := photos[0].(map[string]interface{}); ok {

			if tl, ok := first["thumbnail"].(map[string]interface{}); ok {
				if src, ok := tl["src"].(string); ok {
					c.JSON(200, gin.H{"image": src})
					return
				}
			}
		}
	}

	c.JSON(200, gin.H{
		"image": "https://via.placeholder.com/400x200?text=No+Image",
	})
}
