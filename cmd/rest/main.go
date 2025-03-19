package main

import (
	"time"

	"github.com/adiubaidah/syafiiyah-main/initiator"
)

func main() {

	time.Local = time.UTC
	initiator.Init()
}
