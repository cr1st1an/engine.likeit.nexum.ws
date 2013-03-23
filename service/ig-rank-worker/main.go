/*
	(c) 2013 Carlos Reventlov, carlos@reventlov.com
*/

package main

import (
	"github.com/gosexy/db"
	_ "github.com/gosexy/db/mongo"
	_ "github.com/gosexy/db/mysql"
	"github.com/gosexy/dig"
	"github.com/gosexy/to"
	"github.com/gosexy/yaml"
	// "log"
)

// MongoDB client
var mgc db.Database

// Featured collection.
var featured db.Collection

var photos db.Collection

/*
	Connects to MySQL and MongoDB using settings.yaml.

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

/*
	Adds new photos with a null rank.
*/
func addNewPhotos() error {
	var err error

	// Querying photos that have not been imported yet.
	res, err := photos.Query(db.Cond{"imported": false})

	for {

		photo := map[string]interface{}{}
		err = res.Next(&photo)

		if err != nil {
			break
		}

		exists, _ := featured.Count(db.Cond{"id": photo["id"]})

		if exists == 0 {
			_, err = featured.Append(db.Item{
				"id":                photo["id"],
				"rank":              0,
				"hand_picked":       0,
				"likeit_count":      0,
				"ig_likes_count":    dig.Int64(&photo, "likes", "count"),
				"ig_comments_count": dig.Int64(&photo, "likes", "count"),
				"created_time":      dig.Int64(&photo, "created_time"),
			})
		}

		if err == nil {
			photos.Update(
				db.Cond{"id": photo["id"]},
				db.Set{"imported": true},
			)
		}

	}

	return nil
}

func updateRanks() {
	// Destination collection
}

func main() {
	addNewPhotos()
	updateRanks()
}
