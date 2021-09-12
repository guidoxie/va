package spider

import (
	"log"
	"testing"
)

func TestGet(t *testing.T) {
	ep, err := Get10BondEP("2021-09-11")
	if err != nil {
		t.Error(err)
	}
	log.Println(ep)

}
