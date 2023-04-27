package main

import (
	"log"
	"testing"
)

func TestName(t *testing.T) {
	for i := 0; i < 20; i++ {
		log.Println(GenerateRandMAC())
	}
}
