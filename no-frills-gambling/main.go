package main

import (
	"fmt"
	"math/rand"

	log "github.com/sirupsen/logrus"
)

var (
	// Balance in pennies
	balance = 10000

	// Chance of winning
	winChance = 0.49

	// Stake amount
	stake = 1
)

func main() {
	log.Infof("Zero Frills Gambling!")
	log.Infof("Completely Automated!")
	log.Infof("Lose money without all that pointless clicking!")

	log.Infof("Starting Balance: %s", toCurrency(balance))

	for {
		// if you've run out of money...
		if balance < stake {
			break
		}

		balance = balance - stake + gamble(stake, winChance)

		log.Infof("Balance: %s", toCurrency(balance))
	}

	log.Fatalf("You don't have enough money to gamble anymore. Goodbye.")
}

func toCurrency(pennies int) string {
	str := fmt.Sprintf("£%.2f", float64(pennies)/100.0)
	return str
}

func gamble(stake int, winChance float64) int {

	log.Infof("Gambling %s...", toCurrency(stake))

	if rand.Float64() < winChance {
		log.Infof("WIN!")
		return stake * 2
	} else {
		log.Warnf("LOSS!")
		return 0
	}

}
