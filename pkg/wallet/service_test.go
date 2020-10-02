package wallet

import (
	"testing"
)

func TestService_RegisterAccount(t *testing.T) {
	svc := Service{}
	account, err :=  svc.RegisterAccount("+992938151007")

	if err != nil {
		t.Errorf("не получилось создать ак, получили: %v", account)
	}
	

}

func TestService_FindbyAccountById_success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9929351007")
	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("не удалось найти аккаунт, получили: %v", account)
	}
	
}

func TestService_FindByAccountByID_notFound(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+992938151007")
	account, err := svc.FindAccountByID(2)
	if err == nil {
		t.Errorf("аккаунт не найден, аккаунт: %v", account)
	}
	
	
}