package bank_test

import (
	"fmt"
	"testing"

	bank "."
)

func TestBank(t *testing.T) {
	done := make(chan struct{})

	if ok := bank.Withdraw(300); ok {
		t.Errorf("negative withdraw operation: %d", bank.Balance())
	}

	// Alice
	go func() {
		bank.Deposit(200)
		fmt.Println("=", bank.Balance())
		done <- struct{}{}
	}()

	// Bob
	go func() {
		bank.Deposit(100)
		done <- struct{}{}
	}()

	<-done
	<-done

	if got, want := bank.Balance(), 300; got != want {
		t.Errorf("Balance = %d, expected %d", got, want)
	}

	if ok := bank.Withdraw(50); ok {
		fmt.Println("=", bank.Balance())
		if got, want := bank.Balance(), 250; got != want {
			t.Errorf("Balance = %d, expected %d", got, want)
		}
	} else {
		t.Errorf("Balance = %d, withdraw %d failed", bank.Balance(), 50)
	}

}
