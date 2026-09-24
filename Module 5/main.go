package main

import (
	"errors"
	"fmt"
)

type BankAccount struct {
	Owner string
	amount float64
}

func (bankAccount *BankAccount) Deposit(amount float64) error {
	if amount < 0 {
		return errors.New("Amount is negative")
	}

	bankAccount.amount += amount

	return nil
}

func (bankAccount *BankAccount) Withdraw(amount float64) error {
	if amount < 0 {
		return errors.New("Amount is negative")
	}

	if bankAccount.amount < amount {
		return errors.New("Insufficient Funds")
	}

	bankAccount.amount -= amount

	return nil
}

func (bankAccount BankAccount) GetAmount() float64 {
	return bankAccount.amount
}

func (bankAccount BankAccount) Display() {
	fmt.Printf("Compte de %s : $%.2f\n", bankAccount.Owner, bankAccount.amount)
}

func main() {
	bankAccount := BankAccount{
		Owner: "Alice",
		amount: 1000,
	}

	bankAccount.Display()

	if err := bankAccount.Deposit(500); err != nil {
		fmt.Printf("Impossible de déposer les fonds : %s\n", err.Error())
	}

	bankAccount.Display()

	if err := bankAccount.Withdraw(200); err != nil {
		fmt.Printf("Impossible de retirer les fonds : %s\n", err.Error())
	}

	bankAccount.Display()

	if err := bankAccount.Withdraw(2000); err != nil {
		fmt.Printf("Impossible de retirer les fonds : %s\n", err.Error())
	}

	bankAccount.Display()
}
