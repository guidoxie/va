package spider

import (
	"log"
	"testing"
)

func TestGet(t *testing.T) {
	ep, err := Get10BondEP("2021-11-16")
	if err != nil {
		t.Error(err)
	}
	log.Println(ep)

}
