package wallet

import (
	"os"
	"fmt"
	"log"
	"strconv"
	"strings"
	"io/ioutil"
	"errors"
	"github.com/AzizRahimov/wallet/pkg/types"
	"github.com/google/uuid"
)



var ErrPhoneRegistred  = errors.New("phone already registred")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrNotEnoughBalance = errors.New("Not enough balance")
var ErrFavoriteNotFound = errors.New("favorite not found")
var ErrFileNotFound = errors.New("File Not found")





type Service struct{
	nextAccountID int64
	accounts  []*types.Account
	payments  []*types.Payment
	favorites []*types.Favorite
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

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)

	if err != nil {
		return nil, err
	}

	favoriteID := uuid.New().String()
	newFavorite := &types.Favorite{
		ID:        favoriteID,
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}

	s.favorites = append(s.favorites, newFavorite)
	return newFavorite, nil
}


func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}



//PayFromFavorite для совершения платежа в Избранное
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

//ExportToFile - для импорта данных
func (s *Service) ExportToFile(path string) error  {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}

		
	}()
	str := ""
	for _, account := range s.accounts {
		str += strconv.Itoa(int(account.ID)) + ";"
		str += string(account.Phone) + ";"
		str += strconv.Itoa(int(account.Balance)) + "|"
	}

	_, err = file.Write([]byte(str))
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}


func (s *Service) ImportFromFile(path string) error {
	
	byteData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return err
	}
	
	// обязательно ли преобразовать слайс байтов преобразовывать в стринг?
	data := string(byteData)
	fmt.Println(data)


	splitSlice := strings.Split(data, "|")
	// здесь он должен показать  весь список, тажкже убрать |
	fmt.Println(splitSlice)
	// тут он покажет 2 номер
	fmt.Println(splitSlice[1])
	// тут 1
	fmt.Println(splitSlice[0])
	//splitSlice = strings.Split(data, ";")
	for _, split := range splitSlice {
		if split != "" {
			datas := strings.Split(split, ";")

			id, err := strconv.Atoi(datas[0])
			if err != nil {
				log.Println(err)
				return err
			}
			fmt.Println(datas[0])
			balance, err := strconv.Atoi(datas[2])
			fmt.Println(datas[2])
			if err != nil {
				log.Println(err)
				return err
			}

			newAccount := &types.Account{
				ID:      int64(id),
				Phone:   types.Phone(datas[1]),
				Balance: types.Money(balance),
			}

			s.accounts = append(s.accounts, newAccount)
		}
	}

	return nil
}

// у нас же данные есть ведь? - например RegisterAc - который в памяти хранится
// если она пустая то дай ошибку
// нужна отдельная функция для создания файлов, чтобы 3 раза не писать одно и тоже

func (s *Service) Export(dir string) error {
		// внутри него данные
		// будет тру
	if s.accounts != nil{
		acc := ""
		for _, account := range s.accounts{
			acc += strconv.Itoa(int(account.ID)) + ";"
			acc +=string(account.Phone) + ";"
			acc += strconv.Itoa(int(account.Balance)) +";"
			acc += string('\n')
		}
		err := WriteToFile(dir +"/accounts.dump", acc)
		if err != nil {
			log.Print(err)
			return err
		}
	}
	if  s.payments != nil{
		pay := ""

		for _, payment := range s.payments {

			pay += 	payment.ID + ";"
			pay +=  strconv.Itoa(int(payment.AccountID)) + ";"
			pay +=  strconv.Itoa(int(payment.Amount)) + ";"
			pay +=  string(payment.Category) + ";"
			pay += string(payment.Status) + ";"
			pay += "\n"
			}
		err := WriteToFile(dir + "/payments.dump", pay)
		if err != nil {
			log.Print(err)
			return err
		}

	}
	if s.favorites != nil{
		fav := ""
		for _, favorite := range s.favorites{
			fav += favorite.ID + ";"
			fav += strconv.Itoa(int(favorite.AccountID)) + ";"
			fav += favorite.Name + ";"
			fav += strconv.Itoa(int(favorite.Amount)) + ";"
			fav += string(favorite.Category) + ";"
			fav += "\n"
		}
		err := WriteToFile(dir + "/favorites.dump", fav)
		if err != nil {
			log.Print(err)
			return err
		}

	}







	return  nil

}




func WriteToFile(path string, data string)error  {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		err =  file.Close()
		if err != nil {
			log.Print(err)
			return
		}
	}()
	// он возвращает кол-во байтов

	_, err = file.WriteString(data)	
	if err != nil {
		log.Print(err)
		return  err
	}
	return nil
}
func (s *Service) Import(dir string) error {
	err := s.actionByAccounts(dir + "/accounts.dump")
	if err != nil {
		log.Println("err from actionByAccount")
		return err
	}

	err = s.actionByPayments(dir + "/payments.dump")
	if err != nil {
		log.Println("err from actionByPayments")
		return err
	}

	err = s.actionByFavorites(dir + "/favorites.dump")
	if err != nil {
		log.Println("err from actionByFavorites")
		return err
	}

	return nil
}

func (s *Service) actionByAccounts(path string) error {
	byteData, err := ioutil.ReadFile(path)
	if err == nil {
		datas := string(byteData)
		splits := strings.Split(datas, "\n")

		for _, split := range splits {
			if len(split) == 0 {
				break
			}

			data := strings.Split(split, ";")

			id, err := strconv.Atoi(data[0])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			phone := types.Phone(data[1])

			balance, err := strconv.Atoi(data[2])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}
			// вот это зачем?
			account, err := s.FindAccountByID(int64(id))
			if err != nil {
				acc, err := s.RegisterAccount(phone)
				if err != nil {
					log.Println("err from register account")
					return err
				}

				acc.Balance = types.Money(balance)
			} else {
				account.Phone = phone
				account.Balance = types.Money(balance)
			}
		}
	} else {
		log.Println(ErrFileNotFound.Error())
	}

	return nil
}

func (s *Service) actionByPayments(path string) error {
	byteData, err := ioutil.ReadFile(path)
	if err == nil {
		datas := string(byteData)
		splits := strings.Split(datas, "\n")

		for _, split := range splits {
			if len(split) == 0 {
				break
			}

			data := strings.Split(split, ";")
			id := data[0]

			accountID, err := strconv.Atoi(data[1])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			amount, err := strconv.Atoi(data[2])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			category := types.PaymentCategory(data[3])

			status := types.PaymentStatus(data[4])

			payment, err := s.FindPaymentByID(id)
			if err != nil {
				newPayment := &types.Payment{
					ID:        id,
					AccountID: int64(accountID),
					Amount:    types.Money(amount),
					Category:  types.PaymentCategory(category),
					Status:    types.PaymentStatus(status),
				}

				s.payments = append(s.payments, newPayment)
			} else {
				payment.AccountID = int64(accountID)
				payment.Amount = types.Money(amount)
				payment.Category = category
				payment.Status = status
			}
		}
	} else {
		log.Println(ErrFileNotFound.Error())
	}

	return nil
}

func (s *Service) actionByFavorites(path string) error {
	byteData, err := ioutil.ReadFile(path)
	if err == nil {
		datas := string(byteData)
		splits := strings.Split(datas, "\n")

		for _, split := range splits {
			if len(split) == 0 {
				break
			}

			data := strings.Split(split, ";")
			id := data[0]

			accountID, err := strconv.Atoi(data[1])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			name := data[2]

			amount, err := strconv.Atoi(data[3])
			if err != nil {
				log.Println("can't parse str to int")
				return err
			}

			category := types.PaymentCategory(data[4])

			favorite, err := s.FindFavoriteByID(id)
			if err != nil {
				newFavorite := &types.Favorite{
					ID:        id,
					AccountID: int64(accountID),
					Name:      name,
					Amount:    types.Money(amount),
					Category:  types.PaymentCategory(category),
				}

				s.favorites = append(s.favorites, newFavorite)
			} else {
				favorite.AccountID = int64(accountID)
				favorite.Name = name
				favorite.Amount = types.Money(amount)
				favorite.Category = category
			}
		}
	} else {
		log.Println(ErrFileNotFound.Error())
	}

	return nil
}
