/*
	(c) 2013 Carlos Reventlov, carlos@reventlov.com
*/

package main

import (
	"github.com/gosexy/db"
	_ "github.com/gosexy/db/mongo"
	"github.com/gosexy/dig"
	"github.com/gosexy/to"
	"github.com/gosexy/yaml"
	"github.com/xiam/bridge"
	"log"
)

// Server address
var serverType = "tcp"
var serverAddr = "0.0.0.0:9192"

// MongoDB client
var mgc db.Database

// Featured collection (mongo).
var featured db.Collection

// Photos collection (mongo).
var photos db.Collection

/*
	Connects to MongoDB using settings.yaml.

	Performs initialization tasks.
*/
func init() {
	// Attempt to open settings file.
	yf, err := yaml.Open("../settings.yaml")

	if err != nil {
		panic(err.Error())
	}

	// MongoDB client.
	mgc, err = db.Open("mongo",
		db.DataSource{
			Host:     to.String(yf.Get("database", "mongo", "host")),
			Database: to.String(yf.Get("database", "mongo", "name")),
		},
	)

	if err != nil {
		panic(err.Error())
	}

	// Featured photos
	featured, _ = mgc.Collection("featured")

	// All photos
	photos, _ = mgc.Collection("photos")
}

// Featured photos
type Featured struct {
	Limit int
}

// List featured photos.
func (self *Featured) List() (map[string]interface{}, error) {
	response := []map[string]interface{}{}

	if self.Limit < 1 {
		self.Limit = 20
	}

	// Sort
	info, err := featured.FindAll(
		db.Sort{"rank": -1},
		db.Limit(self.Limit),
	)

	if err != nil {
		return nil, err
	}

	for i, _ := range info {
		photo_id := dig.String(&info[i], "id")
		photo, err := photos.Find(db.Cond{"id": photo_id})
		if err == nil {
			response = append(response, map[string]interface{}{
				"info":  info[i],
				"photo": photo,
			})
		}
	}

	data := map[string]interface{}{
		"data": response,
	}

	return data, nil
}

func main() {

	server := bridge.New(serverType, serverAddr)
	server.AddRoute("/api/featured", &Featured{})
	err := server.Start()

	if err != nil {
		log.Printf("Failed to start server: %s\n", err.Error())
	}

}
