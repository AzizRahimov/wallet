package wallet

import (
	"github.com/AzizRahimov/wallet/pkg/types"
	"reflect"
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
	// тут даст false, так как err (уже имеет что то внутри)
	if err == nil {
		t.Errorf("аккаунт не найден, аккаунт: %v", account)
	}
	
	
}

func TestFindPaymentByID_success(t *testing.T) {
	// cоздаем сервис
	svc := &Service{}
	// создаем регистрацию
	phone := types.Phone("+992938151007")
	account, err := svc.RegisterAccount(phone)
	if err != nil {
		t.Errorf("не удалось зарегестрироваться, Ошибка = %v", err)
		return
	}
	// пополняем счет
	err = svc.Deposit(account.ID, 1000)
	if err != nil {
		t.Errorf("ошибка при пополнении баланса, ошибка = %v", err)
		return
	}
	// осуществляем платеж на его счет
	pay, err := svc.Pay(account.ID,500, "auto")
	if err != nil {
		t.Errorf("ошибка payment error = %v", err)
		return
	}
	got, err := svc.FindPaymentByID(pay.ID)
	if err != nil{
		t.Errorf("FindPayment(): error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, pay){
		t.Errorf("FindPayment(): wrong payment returned = %v", err)
		return
	}
}




func TestService_Reject_success(t *testing.T) {
	// cоздаем сервис 
	svc := &Service{}
	// создаем регистрацию
	phone := types.Phone("+992938151007")
	account, err := svc.RegisterAccount(phone)
	if err != nil {
			t.Errorf("не удалось зарегестрироваться, Ошибка = %v", err)
		return
	}
	// пополняем счет
	err = svc.Deposit(account.ID, 1000)
	if err != nil {
		t.Errorf("ошибка при пополнении баланса, ошибка = %v", err)
		return
	}
	// осуществляем платеж на его счет
	pay, err := svc.Pay(account.ID,500, "auto")
	if err != nil {
		t.Errorf("ошибка payment error = %v", err)
		return
	}
	// делаем отмену платежа
	err = svc.Reject(pay.ID)
	if err != nil {
		t.Errorf("ошибка при отмене платежа, Ошибка = %v", err)
		return
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
		t.Errorf("\ngot > %v \nwant > nil",  pay)
	}

	editPayID := "231231"
	err = svc.Reject(editPayID)
	if err == nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}
}



func TestService_Repeat_success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+99938151007")

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

	pay, err = svc.Repeat(pay.ID)
	if err != nil {
		t.Errorf("Repeat(): Error(): can't pay for an account(%v): %v", pay.ID, err)
	}
}

func TestService_Favorite_success_user(t *testing.T) {
	svc := Service{}

	account, err := svc.RegisterAccount("+992938151007")
	if err != nil {
		t.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	payment, err := svc.Pay(account.ID, 10_00, "auto")
	if err != nil {
		t.Errorf("Pay() Error() can't pay for an account(%v): %v", account, err)
	}

	favorite, err := svc.FavoritePayment(payment.ID, "megafon")
	if err != nil {
		t.Errorf("FavoritePayment() Error() can't for an favorite(%v): %v", favorite, err)
	}

	paymentFavorite, err := svc.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("PayFromFavorite() Error() can't for an favorite(%v): %v", paymentFavorite, err)
	}
}