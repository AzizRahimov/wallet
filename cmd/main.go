package main

import (
	"github.com/AzizRahimov/wallet/pkg/wallet"
	"log"
)

func main() {
	s := &wallet.Service{}

	account, err := s.RegisterAccount("+99293815999")
	if err != nil {
		log.Print(err)
		return
	}


	err = s.Deposit(account.ID, 1000)
	if err != nil {
		log.Print(err)
		}
	payment, err := s.Pay(account.ID, 100, "auto")


	fav, err := s.FavoritePayment(payment.ID, "shop")
	if err != nil {
		log.Print(err)
		return
	}
	_, err = s.PayFromFavorite(fav.ID)
	if err != nil {
		log.Print(err)
	}

	err = s.Export("data")
	if err != nil {
		log.Print(err)
	}





	//_, err := s.RegisterAccount("+992938151007")
	//_, err = s.RegisterAccount("+992938151008")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	//err = s.ExportToFile("data/export.txt")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//err := s.ImportFromFile("data/export.txt")
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//
	//log.Println(s.FindAccountByID(1))
	//log.Println(s.FindAccountByID(2))



}

