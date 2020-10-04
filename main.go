package main

import (
	"fmt"
	"github.com/AzizRahimov/wallet/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	account, err := svc.RegisterAccount("+992938151007")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = svc.Deposit(1, 500)
	if err != nil {
		fmt.Println(err)
		return
	}


	pay, err := svc.Pay(account.ID, 500, "auto")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = svc.Reject(pay.ID)
	if err != nil {
		fmt.Println(err)
		return
	}


	fmt.Println(pay)
	fmt.Println(account.Balance)
}