package blocker

import (
	"log"
	"testing"
)

func TestBlocker(t *testing.T) {
	var list RipIPList
	err := list.LoadFromFile("allowedlist")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	log.Printf("length of entries: %v", len(list))
	list.Dump()
	DefaultAllowing.Dump()
	log.Print("***split***")
}
