// Exercise 9.1: Add a function Withdraw(amount int) bool to the
// gopl.io/ch9/bank1 program. The result should indicate whether the transaction
// succeeded or failed due to insufficient funds. The message sent to the
// monitor goroutine must contain both the amount to withdraw and a new channel
// over which the monitor goroutine can send the boolean result back to
// Withdraw.

// Package bank provides a concurrency-safe bank with one account.
package bank

type withdrawOperation struct {
	amount int
	ok     chan bool
}

var deposits = make(chan int) // send amount to deposit
var balances = make(chan int) // receive balance
var withdraws = make(chan withdrawOperation)

func Deposit(amount int) { deposits <- amount }
func Balance() int       { return <-balances }
func Withdraw(amount int) bool {
	ok := make(chan bool)
	withdraws <- withdrawOperation{amount, ok}
	return <-ok
}

func teller() {
	var balance int // balance is confined to teller goroutine
	for {
		select {
		case amount := <-deposits:
			balance += amount
		case op := <-withdraws:
			if balance >= op.amount {
				balance -= op.amount
				op.ok <- true
			} else {
				op.ok <- false
			}
		case balances <- balance:
		}
	}
}

func init() {
	go teller() // start the monitor goroutine
}

//!-
