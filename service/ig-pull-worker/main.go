/*
	(c) 2013 Carlos Reventlov, carlos@reventlov.com
*/

package main

import (
	"errors"
	"github.com/gosexy/db"
	_ "github.com/gosexy/db/mongo"
	_ "github.com/gosexy/db/mysql"
	"github.com/gosexy/to"
	"github.com/gosexy/yaml"
	"github.com/xiam/instagram"
	"log"
	"time"
)

const queryLimit = 200

const sleepTime = time.Second * 300

// MySQL client
var myc db.Database

// MongoDB client
var mgc db.Database

// Instagram client
var igc *instagram.Client

// Media structure (on MySQL db)
type myMedia struct {
	IdIgMedia string
	C         int64
}

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

	// MySQL client.
	myc, err = db.Open("mysql",
		db.DataSource{
			Host:     to.String(yf.Get("database", "mysql", "host")),
			Database: to.String(yf.Get("database", "mysql", "name")),
			User:     to.String(yf.Get("database", "mysql", "user")),
			Password: to.String(yf.Get("database", "mysql", "password")),
		},
	)
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

	// Instagram client.
	igc = instagram.New()
	igc.SetAccessToken(to.String(yf.Get("providers", "instagram", "access_token")))
}

/*
	Receives a photo ID and queries the instagram API for data about that photo.

	If the photo does not exists appends it to the "photos" collection of the
	MongoDB instance.
*/
func pullPhoto(mediaId string) error {

	var err error

	// Destination collection
	photos, _ := mgc.Collection("photos")

	// Does the photo exists?
	exists, err := photos.Count(db.Cond{"id": mediaId})

	if err != nil {
		return err
	}

	if exists == 0 {
		res := map[string]interface{}{}

		// Query the instagram server.
		igc.Media(&res, mediaId)

		if res["data"] == nil {
			return errors.New("No data received.")
		}

		data := res["data"].(map[string]interface{})
		data["imported"] = false

		// Append photo to database.
		_, err = photos.Append(data)

		if err != nil {
			return err
		}
	}

	return nil
}

/*
	Gets new photos from the "id_media_likes" table and updates each one of them.

	TODO:
		- Queue requests or implement retry to avoid abusing API rate limit.
		- Mark photos that have been already queried to avoid querying them again.
		- Implement WaitGroup to make concurrent queries.
*/
func updateMedia() error {

	// This MySQL view holds media IDs and likes.
	view := myc.ExistentCollection("ig_media_likes")

	// A simple query.
	res, err := view.Query(
		db.Sort{"c": -1},
		//db.Limit(queryLimit),
	)

	if err != nil {
		return err
	}

	/*
		count, _ := view.Count(
			db.Sort{"c": -1},
			db.Limit(queryLimit),
		)

		log.Printf("Got %d new items.\n", count)
	*/

	var media myMedia

	for {

		// Iterating over results.
		err = res.Next(&media)

		if err != nil {
			break
		}

		// Updating photo.
		err = pullPhoto(media.IdIgMedia)

		if err != nil {
			// Log to stderr.
			log.Printf("Error while pulling photo id %v: %v\n", media.IdIgMedia, err)
		}

	}

	return nil
}

func main() {
	for true {
		log.Printf("Updating...\n")
		updateMedia()
		time.Sleep(sleepTime)
	}
}
