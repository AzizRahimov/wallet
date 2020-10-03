package wallet

import (
	"github.com/AzizRahimov/wallet/pkg/types"
	"testing"
)

func TestService_RegisterAccount(t *testing.T) {
	svc := &Service{}
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


func TestService_Reject_success(t *testing.T) {
	// создаем сервис
	svc := &Service{}
	
	// регистрация пользователя
	phone := types.Phone("+992938151007")
	account, err := svc.RegisterAccount(phone)
	if err != nil{
		t.Errorf("can't register accoount, error = %v", err)
	}
	// пополняем счет, даем, тот ак, который только что создали
	err = svc.Deposit(account.ID, 1000)
	if err != nil {
		t.Errorf("Reject(): can't deposit account, error = %v", err)
	}
	//Todo: отнять сумму

	payment, err := svc.Pay(account.ID, 500, "auto")
	if err != nil {
		t.Errorf("Reject(): не может создать платеж, ошибка = %v", err)
	}
	
	err = svc.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
	}

}




func TestService_Reject_fail(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+992938151007")

	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = svc.Deposit(account.ID, 1000_00)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	payment, err := svc.Pay(account.ID, 100_00, "auto")
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	pay, err := svc.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	editPayID := pay.ID + "31231:)"
	err = svc.Reject(editPayID)
	if err == nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}
}