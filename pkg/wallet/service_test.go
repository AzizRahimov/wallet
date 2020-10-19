package wallet

import (
	"fmt"
	"github.com/AzizRahimov/wallet/pkg/types"
	"reflect"
	"testing"
)

type testService struct {
	*Service
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone:   "+992938151007",
	balance: 10_000_00,
	payments: []struct {
		amount   types.Money
		category types.PaymentCategory
	}{{
		amount:   1000_00,
		category: "auto",
	}},
}




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

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error) {
	account, err := s.RegisterAccount(phone)

	if err != nil {
		return nil, fmt.Errorf("cant register account, error = %v", err)
	}

	err = s.Deposit(account.ID, balance)

	if err != nil {
		return nil, fmt.Errorf("cant deposit account, error = %v", err)
	}

	return account, nil
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("cant register account %v = ", err)
	}

	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("cant deposit account %v = ", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("cant make payment %v = ", err)
		}
	}

	return account, payments, nil
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


func BenchmarkSumPayments(b *testing.B) {
	s := newTestService()

	_, _, err := s.addAccount(testAccount{
		phone:   "+992935444994",
		balance: 1000_000_00,
		payments: []struct {
			amount   types.Money
			category types.PaymentCategory
		}{
			{
				amount:   1000_00,
				category: "auto",
			},
			{
				amount:   2000_00,
				category: "auto",
			},
			{
				amount:   3000_00,
				category: "auto",
			},
			{
				amount:   4000_00,
				category: "auto",
			},
			{
				amount:   5000_00,
				category: "auto",
			},
			{
				amount:   6000_00,
				category: "auto",
			},
			{
				amount:   1250_00,
				category: "auto",
			},
			{
				amount:   1870_00,
				category: "auto",
			},
			{
				amount:   9877_00,
				category: "auto",
			},
		},
	})

	if err != nil {
		b.Error(err)
		return
	}

	want := types.Money(3399700)

	for i := 0; i < b.N; i++ {
		result := s.SumPayments(6)
		if result != want {
			b.Fatalf("invalid result, got = %v want = %v", result, want)
		}
	}
}
