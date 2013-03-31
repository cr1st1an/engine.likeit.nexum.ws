/*
	(c) 2013 Carlos Reventlov, carlos@reventlov.com
*/

package main

import (
	"github.com/gosexy/db"
	"github.com/gosexy/to"
	"testing"
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

func TestPhotosManualOrder(t *testing.T) {

	f := &Photos{}

	data, err := f.List()

	if err != nil {
		t.Fatalf("Failed: %v.", err)
	}
	if data["data"] == nil {
		t.Fatalf("Data should not be nil.")
	}

	for _, item := range data["data"].([]db.Item) {
		p := &Photo{}
		p.Id = to.String(item["id"])
		p.HandpickedRank = 10
		data2, err2 := p.SetHandpickedRank()
		if err2 != nil {
			t.Fatalf("Failed: %v.", err2)
		}
		if data2["success"] == nil {
			t.Fatalf("Expecting success")
		}
	}

}
