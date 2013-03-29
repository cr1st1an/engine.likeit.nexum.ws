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
	"log"
	"time"
)

// Constants for the dummy ranker

// How much an Instagram comment is worth
const InstagramCommentFactor = 2.0

// How much an Instagram like is worth
const InstagramLikeFactor = 1.0

// How much a LikeIt like is worth
const LikeitLikeFactor = 100.0

// How much the age difference worths
const InstagramAgeFactor = -0.001

// Update time
const sleepTime = time.Second * 300

// MongoDB client
var mgc db.Database

// MySQL client
var myc db.Database

// Featured collection (mongo).
var featured db.Collection

// Photos collection (mongo).
var photos db.Collection

// Likes collection (mysql)
var likes db.Collection

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

	// Featured photos
	featured, _ = mgc.Collection("featured")

	// All photos
	photos, _ = mgc.Collection("photos")

	// This MySQL view holds media IDs and likes.
	likes = myc.ExistentCollection("ig_media_likes")
}

/*
	Does (future) magic and returns a rank value.

	map[
		created_time:1361805295
		ig_likes_count:80708
		_id:ObjectIdHex("51558cfbb79e40b00067f75f")
		likeit_count:6
		hand_picked:0
		rank:0
		id:399174211645560261_787132
		ig_comments_count:574
	]

*/

func getRank(photo map[string]interface{}) int64 {
	var rank float64

	age := time.Now().Unix() - dig.Int64(&photo, "created_time")

	rank = 0.0

	rank += dig.Float64(&photo, "ig_likes_count") * InstagramLikeFactor
	rank += dig.Float64(&photo, "ig_comments_count") * InstagramCommentFactor
	rank += dig.Float64(&photo, "likeit_count") * LikeitLikeFactor
	rank += float64(age) * InstagramAgeFactor

	return int64(rank)
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
			// Get likes
			likeData, _ := likes.Find(
				db.Cond{"id_ig_media": photo["id"]},
			)

			_, err = featured.Append(db.Item{
				"id":                photo["id"],
				"rank":              0, // Have to talk to Cris
				"hand_picked":       0, // Have to talk to Cris
				"likeit_count":      dig.Int64(&likeData, "c"),
				"ig_likes_count":    dig.Int64(&photo, "likes", "count"),
				"ig_comments_count": dig.Int64(&photo, "comments", "count"),
				"created_time":      dig.Int64(&photo, "created_time"),
				"modified":          time.Now().Unix(),
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

func updateRanks() error {
	// Destination collection
	res, err := featured.Query()
	if err != nil {
		return err
	}

	for {

		photo := map[string]interface{}{}
		err = res.Next(&photo)

		if err != nil {
			break
		}

		rank := getRank(photo)

		featured.Update(
			db.Cond{"_id": photo["_id"]},
			db.Set{
				"rank":     rank,
				"modified": time.Now().Unix(),
			},
		)

	}
	return nil
}

func main() {
	var err error
	for true {

		log.Printf("Attempt to add new photos...\n")
		err = addNewPhotos()
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
		}

		log.Printf("Attempt to update photo ranks...\n")
		err = updateRanks()
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
		}

		log.Printf("OK, back to sleep.\n")
		time.Sleep(sleepTime)
	}
}
