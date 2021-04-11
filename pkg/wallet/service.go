package wallet

import (
	"errors"
	"github.com/AzizRahimov/wallet/pkg/types"
	"github.com/google/uuid"
)



var ErrPhoneRegistred  = errors.New("phone already registred")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrNotEnoughBalance = errors.New("Not enough balance")




type Service struct{
	nextAccountID int64
	accounts []*types.Account
	payments []*types.Payment
}




func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error)  {
		
	for _, account := range s.accounts{
		if account.Phone == phone{
			return nil, ErrPhoneRegistred
		}
		
	}

	 s.nextAccountID++
	acc := &types.Account{
		ID: s.nextAccountID,
		Phone: phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, acc)

	return acc, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) (error) {
	if amount <= 0{
		return ErrAmountMustBePositive
	}
	var account *types.Account

	for _, acc := range s.accounts{
		if acc.ID == accountID{
			account = acc
			break
		}
	}
	if account == nil{
		return ErrAccountNotFound
	}
	account.Balance += amount

	return nil
	
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory, ) (*types.Payment, error)  {
	if amount <= 0{
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account

	for _, acc := range s.accounts{
		if acc.ID == accountID{
			account = acc
			break
		}
	}
	if account == nil{
		return nil, ErrAccountNotFound
	}
	if account.Balance < amount{
		return nil, ErrNotEnoughBalance
	}
	account.Balance -= amount
	id := uuid.New().String()
	//! что хранит в себя Payment ?
	//! он хранит инфо о платеже, то есть с этого счета было совершно что то на такую-та сумму
	payment := &types.Payment{
		ID: id,
		AccountID: accountID,
		Amount: amount,
		Category: category,
		Status: types.PaymentStatusInProgress,

	}
	s.payments = append(s.payments, payment)




		return payment, nil
}


func (s *Service) FindAccountByID(accountID int64) (*types.Account, error)  {

	for _,  account := range s.accounts{
		if account.ID == accountID{
			return account, nil
		}
		
	}
	
	return nil, ErrAccountNotFound
}

func (s Service) FindPaymentByID(paymentID string) (*types.Payment, error)  {

	for _, payment :=  range s.payments{
		if payment.ID == paymentID{
			return payment, nil
		}
	}
	
	return nil, ErrPaymentNotFound
}


// Отменят платеж, то есть обратно возвращаем баланс
func (s *Service) Reject(paymentID string) error{
	// сперва находим 
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil{
		return err
	}
	//? тут нашли платеж
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil{
		return err
	}
	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}
	

//


func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	// мы должны ID по которому будет транзакция 
	pay, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}


	payment, err := s.Pay(pay.AccountID, pay.Amount, pay.Category)
	if err != nil {
		return nil, err
	}
	
	return payment, nil
}
