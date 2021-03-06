package namegenerator

import (
	petname "github.com/dustinkirkland/golang-petname"
	"math/rand"
	"time"
)

func Generate() string {
	// generate some entropy
	rand.Seed(time.Now().UTC().UnixNano())
	return petname.Generate(3, "-")
}