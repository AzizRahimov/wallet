package main

import (
	"github.com/AzizRahimov/wallet/pkg/wallet"
	"log"
)

func main() {
	svc := &wallet.Service{}
	_, err := svc.RegisterAccount("+992938151007")
	if err != nil {
		log.Print(err)
		return
	}
	err = svc.ExportToFile("data/export.txt")
	if err != nil {
		log.Print(err)
		return
	}
}