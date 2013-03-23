/*
	(c) 2013 Carlos Reventlov, carlos@reventlov.com
*/

package main

import (
	"fmt"
	"github.com/gosexy/db"
	_ "github.com/gosexy/db/mongo"
	_ "github.com/gosexy/db/mysql"
	"github.com/gosexy/to"
	"github.com/gosexy/yaml"
)

// MySQL client
var myc db.Database

// MongoDB client
var mgc db.Database

func init() {
	yf, err := yaml.Open("settings.yaml")

	if err != nil {
		panic(err.Error())
	}

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

	mgc, err = db.Open("mongo",
		db.DataSource{
			Host:     to.String(yf.Get("database", "mongodb", "host")),
			Database: to.String(yf.Get("database", "mongodb", "name")),
		},
	)

	if err != nil {
		panic(err.Error())
	}
}

func main() {
	fmt.Printf("hello!\n")
}
