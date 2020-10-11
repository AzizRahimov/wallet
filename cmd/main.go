package main

import (
	"github.com/AzizRahimov/wallet/pkg/wallet"
	"log"
)

func main() {
	s := &wallet.Service{}

	_, err := s.RegisterAccount("+992938151007")
	if err != nil {
		log.Println(err)
		return
	}

	err = s.ExportToFile("data/export.txt")
	if err != nil {
		log.Println(err)
		return
	}
}