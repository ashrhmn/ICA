#Top Investors
This script reads transaction data from a CSV file and calculates the top 5 investors based on the number of unique syndicates they invested in and the total amount of their investments.

###Usage

- Put your transaction data in a CSV file named transactions.csv in the same directory as the script. Each row should contain the following columns: investor ID, syndicate ID, and amount.
- Run the script: go run main.go
- The script will print the top 5 investors to the console, sorted by the number of unique syndicates they invested in (in descending order) and then by the total amount of their investments (in descending order). For each investor, the script will print their ID, the number of unique syndicates they invested in, and the total amount of their investments.

###Implementation

The script uses Go's built-in concurrency features to process the transaction data in parallel using goroutines. It reads the data from the CSV file, parses it into a slice of Transaction structs, and then processes each transaction using a goroutine. The goroutines update a shared map of Investor structs, which store information about each investor's unique syndicates and total investment. To protect concurrent access to the map, the script uses a mutex. Once all goroutines have finished, the script converts the map to a slice and sorts it based on the number of unique syndicates and total investment. Finally, the script prints the top 5 investors to the console.

#Alert System for High Transaction Rates
This script implements an alert system for detecting high transaction rates in a financial system. The system generates random transactions and processes them in parallel, checking for high single transaction amounts and sudden spikes in transaction rates. If a high transaction amount or a sudden spike is detected, an alert is printed to the console.

###How it Works
The script uses goroutines and channels to implement a concurrent system for generating and processing transactions. The generateTransactions function generates random transactions and sends them to the transactionChannel channel. The processTransactions function receives transactions from the channel and checks for high single transaction amounts and spikes in transaction rates. The reportTransactionSpike function checks for sudden spikes in transaction rates and reports an alert if the spike threshold is exceeded.

The script uses the following parameters to control the behavior of the system:

MaxSingleTransactionAmount: the maximum amount for a single transaction to trigger a high amount alert.
SpikeThresholdMultiplier: the multiplier for the transaction rate in the spike window to trigger a spike alert.
SpikeWindow: the time window for calculating the transaction rate in the spike window.
SpikeCompareWindow: the time window for calculating the transaction rate in the spike compare window.
SpikeReportingWindow: the time window for reporting alerts for transaction spikes.
How to Use
To use the script, simply run the main.go file with the Go compiler. The script will start generating transactions and processing them in parallel. If a high transaction amount or a sudden spike is detected, an alert will be printed to the console.

You can customize the behavior of the system by changing the parameters in the main.go file. For example, you can increase the MaxSingleTransactionAmount parameter to trigger alerts for higher transaction amounts, or decrease the SpikeThresholdMultiplier parameter to trigger alerts for smaller spikes in transaction rates.
