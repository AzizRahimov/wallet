package wallet

import (
	"fmt"
	
	"reflect"
	"testing"

	"github.com/AzizRahimov/wallet/pkg/types"
	"github.com/google/uuid"
)

type testAccount struct{
	phone types.Phone
	balance types.Money
	payments []struct{
		amount types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone: "+9920001",
	balance: 100_000,
	payments: []struct{
		amount types.Money
		category types.PaymentCategory
	}{
		{amount: 10_00, category: "auto"},
	},
}





type testService struct{
	*Service
}

func newService() *testService  {
	return &testService{Service: &Service{}}
	
}

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error)  {
	// Рег пользователя 
	account, err := s.RegisterAccount(data.phone)
	if err != nil{
		return nil, nil, fmt.Errorf("can't register Account, error = %v", err)
	}
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit account, error = %v", err)
	}
	//? выполняем платежи
	//? можем создать слайс сразу нужной длинны, поскольку знаем размер
	//?
	payments := make([]*types.Payment, len(data.payments)) // хмм, так он же пуст

	for i, payment :=  range data.payments{
		// здесь работаем через индекс, а не через append
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}
	
	return account, payments, nil
}




func (s *testService) addAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error)  {
	account, err := s.RegisterAccount(phone)
	if err != nil {
		return nil, fmt.Errorf("cant't register account, error = %v", err)
	}
	err =s.Deposit(account.ID, balance)
	if err != nil {
		return nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	return account, nil
	
}

func TestService_FindPaymentByID_success(t *testing.T) {
	// создаем сервис
	svc := newService()
	_, payments, err := svc.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
	}

	//! пробуем найти платеж
	payment := payments[0]
	
	got, err := svc.FindPaymentByID(payment.ID)
		if err != nil {
			t.Errorf("FindPaymentByID(): error %v", err)
			return
		}
		if !reflect.DeepEqual(payment, got){
			t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
			return
		}

}












// func TestService_FindPaymentByID_success(t *testing.T) {
// 	s := newService()
// 	account, err := s.addAccountWithBalance("+0001", 1000)
// 	if err != nil{
// 		t.Error(err)
// 		return
// 	}
// 	//! осуществляем платеж
// 	payment, err := s.Pay(account.ID, 500, "auto")
// 	if err != nil {
// 		t.Errorf("FindPaymentByID(): can't create payment, error = %v", err)
// 		return
// 	}

// 	// try to find paymentID
// 	got, err := s.FindPaymentByID(payment.ID)
// 	if err != nil{
// 		t.Errorf("FindPaymentByID(): error = %v", err)
// 	}
// 	// переменная payment - Хранит в себя лишь 1 поле структуру Payment - которое получили через Pay
// 	// переменная got:-  тоже самое, хранит в себя указатель на структуру Payment - которое получили через
// 	//  FindPaymentByID
// 	if !reflect.DeepEqual(payment, got){
// 		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
// 		return
// 	}
	
// }


func TestService_FindPaymentByID_fail(t *testing.T) {
	s := newService()
	account, err := s.addAccountWithBalance("+0001", 1000)
	if err != nil{
		t.Error(err)
		return
	}
	//! осуществляем платеж
	_, err = s.Pay(account.ID, 500, "auto")
	if err != nil {
		t.Errorf("FindPaymentByID(): can't create payment, error = %v", err)
		return
	}

	// try to find paymentID - whicch not exist

	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil{
		t.Errorf("FindPaymentByID(): must return error = %v", err)
		return
	}
	// err == ErrorNotFound
	if err != ErrPaymentNotFound{
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
	}

	// переменная payment - Хранит в себя лишь 1 поле структуру Payment - которое получили через Pay
	// переменная got:-  тоже самое, хранит в себя указатель на структуру Payment - которое получили через
	//  FindPaymentByID
	// if !reflect.DeepEqual(payment, got){
	// 	t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
	// 	return
	// }
	
}








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



func TestService_Repeat_success(t *testing.T) {
	svc := newService()
	_, payments, err :=svc.addAccount(defaultTestAccount)
	if err != nil {
		t.Errorf("Repeat error = %v", err)
	}

	payment := payments[0]
	_, err =svc.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeat(), can't repeat payment, error =  %v", err)
	}
	
}



