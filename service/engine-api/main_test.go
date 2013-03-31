/*
	(c) 2013 Carlos Reventlov, carlos@reventlov.com
*/

package main

import (
	"testing"
	"github.com/gosexy/db"
)

func TestList(t *testing.T) {
	f := &Featured{}
	data, err := f.List()
	if err != nil {
		t.Fatalf("Failed: %v.", err)
	}
	if data["data"] == nil {
		t.Fatalf("Data should not be nil.")
	}
}

func TestPhotosList(t *testing.T) {
	f := &Photos{}

	data, err := f.List()

	if err != nil {
		t.Fatalf("Failed: %v.", err)
	}
	if data["data"] == nil {
		t.Fatalf("Data should not be nil.")
	}

	f.Page = 2

	data, err = f.List()

	if err != nil {
		t.Fatalf("Failed: %v.", err)
	}
	if data["data"] == nil {
		t.Fatalf("Data should not be nil.")
	}

	f.Page = 200000

	data, err = f.List()

	if err != nil {
		t.Fatalf("Failed: %v.", err)
	}
	if data["data"] == nil {
		t.Fatalf("Data should not be nil.")
	}
	if len(data["data"].([]db.Item)) != 0 {
		t.Fatalf("Should not return items.")
	}
}
