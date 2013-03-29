/*
	(c) 2013 Carlos Reventlov, carlos@reventlov.com
*/

package main

import (
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
