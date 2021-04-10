package wallet

import (

"testing"
)


func TestService_RegisterAccount_success(t *testing.T) {
	svc := &Service{}
	account, err := svc.RegisterAccount("+0001")
	if err != nil{
		t.Error("не удалось создать акк", err, account)
	}
	
}


func TestService_RegisterAccount_alreadyRegistered(t *testing.T) {
	svc := &Service{}
	account, err := svc.RegisterAccount("+0001")
	if err != nil{
		t.Error("не удалось создать акк", err, account)
	}
	account, err = svc.RegisterAccount("+0001")
	if err == nil{
		t.Error(ErrPhoneRegistred, account)
	}
	
}


func TestService_Reject_success(t *testing.T) {
	svc := &Service{}
	account, err := svc.RegisterAccount("+0001")
	if err != nil{
		t.Error("не удалось создать акк", err, account)
	}
	err = svc.Deposit(account.ID, 2500)
	if err != nil{
		t.Errorf("error deposit: got: %v", err)
	}
	payment, err := svc.Pay(account.ID, 500, "auto")
	if err != nil{
		t.Errorf("in method Pay got error: %v", err)
	}
	
	err = svc.Reject(payment.ID)
	if err != nil{
		t.Errorf("in method Reject got error : %v", err)
	}



	
}
