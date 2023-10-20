package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Transaction struct {
	Amount float64
	Time   time.Time
}

const MaxSingleTransactionAmount = 10000 // Threshold for single transaction
const SpikeThresholdMultiplier = 1.2     // Threshold for transaction rate spike
// SpikeWindow is the duration of time used to detect spikes in the system.
const SpikeWindow = time.Hour

// SpikeCompareWindow is the time window used to compare the current value with the previous value to detect spikes.
const SpikeCompareWindow = time.Minute * 10

// SpikeReportingWindow is the time window used to report spikes in the system.
const SpikeReportingWindow = time.Minute * 10

// transactionChannel is a channel used to send Transaction objects.
var transactionChannel = make(chan Transaction)

// Maintain a record of transaction times for spike detection.
var transactionTimes = []time.Time{}

func main() {
	// Create a channel to receive termination signals.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Use a new source of randomness with the current time as a seed.
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Start a goroutine to populate transactions into the channel indefinitely.
	go generateTransactions(random)

	// Start another goroutine to process transactions in parallel.
	go processTransactions()

	// Start another goroutine to report transaction spikes.
	go reportTransactionSpike()

	// Wait for termination signal.
	<-stopChan
	println("Waiting for all transactions to be processed...")
	println("Terminating...")
	// Gracefully shutdown the application.
	close(transactionChannel)
}

// transactionRateInSpikeWindow calculates the transaction rate in the spike window.
// It returns a float64 value representing the transaction rate.
func transactionRateInSpikeWindow() float64 {
	// Calculate the transaction rate in the spike window.
	transactionRate := float64(len(transactionTimes)) / SpikeWindow.Seconds()
	return transactionRate
}

// transactionRateInSpikeCompareWindow calculates the transaction rate in the spike compare window.
// It filters transaction times in the spike compare window and calculates the transaction rate.
// Returns the transaction rate as a float64.
func transactionRateInSpikeCompareWindow() float64 {
	// filter transaction times in the spike compare window.
	var filteredTransactionTimes []time.Time
	for i, t := range transactionTimes {
		if time.Since(t) > SpikeCompareWindow {
			filteredTransactionTimes = transactionTimes[i+1:]
		} else {
			break
		}
	}
	// Calculate the transaction rate in the spike compare window.
	transactionRate := float64(len(filteredTransactionTimes)) / SpikeCompareWindow.Seconds()
	return transactionRate
}

// reportTransactionSpike checks for a sudden spike in transaction rate and reports an alert if the spike threshold is exceeded.
// It compares the transaction rate in the spike compare window to the transaction rate in the spike window.
// If the transaction rate in the spike compare window is greater than the spike threshold multiplier times the transaction rate in the spike window,
// an alert is printed with the transaction rates in both windows.
// It sleeps for the spike reporting window duration before checking again.
func reportTransactionSpike() {
	for {
		// Check for a sudden spike in transaction rate.
		_transactionRateInSpikeCompareWindow := transactionRateInSpikeCompareWindow()
		_transactionRateInSpikeWindow := transactionRateInSpikeWindow()
		if _transactionRateInSpikeCompareWindow > SpikeThresholdMultiplier*(_transactionRateInSpikeWindow) {
			println("Alert : Sudden spike in transaction rate")
			fmt.Printf("Transaction rate in spike window: %.2f\n", _transactionRateInSpikeWindow)
			fmt.Printf("Transaction rate in spike compare window: %.2f\n", _transactionRateInSpikeCompareWindow)
			fmt.Println("-------------------------------------------------------")
		}
		time.Sleep(SpikeReportingWindow)
	}
}

// generateTransactions generates random transactions and sends them to the transactionChannel.
// It introduces a random time sleep between transactions (between 100ms and 5000ms).
// The function exits the goroutine if the transactionChannel is closed.
func generateTransactions(random *rand.Rand) {
	for {
		select {
		case <-transactionChannel:
			return // Exit goroutine if the channel is closed.
		default:
			amount := random.Float64() * 11000 // Random transaction amount
			now := time.Now()
			transaction := Transaction{Amount: amount, Time: now}
			transactionChannel <- transaction

			// Introduce random time sleep between transactions (between 100ms and 5000ms).
			sleepDuration := time.Millisecond * time.Duration(random.Intn(4900)+100)
			time.Sleep(sleepDuration)
		}
	}
}

// processTransactions processes transactions from a channel and checks for high single transaction amounts and spikes in transaction times.
// If a single transaction exceeds the threshold, it calls notifyHighAmount with the message "High single transaction amount" and the transaction details.
// It also tracks transaction times for spike detection and removes transactions older than SpikeWindow.
func processTransactions() {

	for transaction := range transactionChannel {
		// Check if a single transaction exceeds the threshold.
		if transaction.Amount > MaxSingleTransactionAmount {
			notifyHighAmount("High single transaction amount", transaction)
		}

		// Track transaction times for spike detection.
		transactionTimes = append(transactionTimes, transaction.Time)

		// Remove transactions older than SpikeWindow.
		for i, t := range transactionTimes {
			if time.Since(t) > SpikeWindow {
				transactionTimes = transactionTimes[i+1:]
			} else {
				break
			}
		}
	}
}

// notifyHighAmount sends a notification for a high amount transaction.
// It takes an alertType string and a Transaction struct as parameters.
// The function prints the alert type, transaction amount, and transaction time to the console.
func notifyHighAmount(alertType string, transaction Transaction) {
	// You can implement the notification logic here, e.g., sending an email or SMS.
	fmt.Printf("Alert: %s\n", alertType)
	fmt.Printf("Transaction Amount: %.2f\n", transaction.Amount)
	fmt.Printf("Transaction Time: %s\n", transaction.Time)
	fmt.Println("-------------------------------------------------------")
}
