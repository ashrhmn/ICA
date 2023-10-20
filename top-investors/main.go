package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
)

// Transaction represents an individual transaction.
type Transaction struct {
	InvestorID  string
	SyndicateID string
	Amount      float64
}

// Investor represents information about an investor, including unique syndicates and total investment.
type Investor struct {
	ID              string
	Syndicates      map[string]struct{}
	TotalInvestment float64
}

func main() {
	// Open and read data from the CSV file
	csvFile, err := os.Open("transactions.csv")
	if err != nil {
		fmt.Println("Error opening the CSV file:", err)
		return
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	lines, err := csvReader.ReadAll()
	if err != nil {
		fmt.Println("Error reading the CSV file:", err)
		return
	}

	// Create a slice to store transaction data
	transactions := make([]Transaction, 0, len(lines))

	// Parse CSV data and store it in the transactions slice
	for _, line := range lines {
		if len(line) < 3 {
			continue
		}
		investorID := line[0]
		syndicateID := line[1]
		amount, err := strconv.ParseFloat(line[2], 64)
		if err != nil {
			continue
		}
		transactions = append(transactions, Transaction{investorID, syndicateID, amount})
	}

	// Create a map to store investor data
	investors := make(map[string]*Investor)

	// Mutex to protect concurrent map access
	var mu sync.Mutex

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Process transactions using goroutines
	for _, transaction := range transactions {
		wg.Add(1)
		go func(t Transaction) {
			defer wg.Done()

			investorID := t.InvestorID
			syndicateID := t.SyndicateID
			amount := t.Amount

			// Lock the mutex to protect concurrent map access
			mu.Lock()

			// If the investor is not in the map, create a new entry
			if _, ok := investors[investorID]; !ok {
				investors[investorID] = &Investor{
					ID:         investorID,
					Syndicates: make(map[string]struct{}),
				}
			}

			// Track unique syndicates using a map
			investors[investorID].Syndicates[syndicateID] = struct{}{}

			// Calculate total investment
			investors[investorID].TotalInvestment += amount

			// Unlock the mutex
			mu.Unlock()
		}(transaction)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Convert investor data from the map to a slice
	investorList := make([]Investor, 0, len(investors))
	for _, investor := range investors {
		investorList = append(investorList, *investor)
	}

	// Sort investors by unique syndicates and total investment
	sort.Slice(investorList, func(i, j int) bool {
		syndicateCountI := len(investorList[i].Syndicates)
		syndicateCountJ := len(investorList[j].Syndicates)
		if syndicateCountI == syndicateCountJ {
			return investorList[i].TotalInvestment > investorList[j].TotalInvestment
		}
		return syndicateCountI > syndicateCountJ
	})

	// Get the top 5 investors
	top5Investors := investorList
	if len(investorList) > 5 {
		top5Investors = investorList[:5]
	}

	// Print the top 5 investors
	fmt.Println("Top 5 Investors:")
	for i, investor := range top5Investors {
		fmt.Printf(
			"%d. Investor ID: %s, Unique Syndicates: %d, Total Investment: %.2f\n",
			i+1, investor.ID, len(investor.Syndicates), investor.TotalInvestment,
		)
	}
}
