package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var teamNames = []string{"John", "Matteo", "Julio", "Samuel"}

func generateNonce(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?/~`"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func addNonceToList(names []string, nonce string) string {
	items := strings.Join(names, "")
	items += " " + nonce
	return items
}

func hashWithNonce(stringList string) string {
	hash := sha256.Sum256([]byte(stringList))
	return hex.EncodeToString(hash[:])
}

func worker(teamNames []string, targetZeros int, results chan<- [2]string, wg *sync.WaitGroup) {
	defer wg.Done()
	targetPrefix := strings.Repeat("0", targetZeros)
	for {
		nonce := generateNonce(16)
		teamList := addNonceToList(teamNames, nonce)
		hashValue := hashWithNonce(teamList)
		if strings.HasPrefix(hashValue, targetPrefix) {
			results <- [2]string{hashValue, nonce}
			return
		}
	}
}

func findHashWithLeadingZeros(teamNames []string, targetZeros int, numWorkers int) (string, string) {
	results := make(chan [2]string)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(teamNames, targetZeros, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	result := <-results
	return result[0], result[1]
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	user_key := os.Getenv("USER_KEY")
	app_key := os.Getenv("APP_KEY")
	fmt.Println("Starting program...")
	// fmt.Println(user_key)
	// fmt.Println(app_key)
	startTime := time.Now()
	hashValue, nonce := findHashWithLeadingZeros(teamNames, 10, 4) // Using 4 workers
	endTime := time.Now()
	executionTime := endTime.Sub(startTime).Seconds()
	fmt.Printf("Hash: %s\nNonce: %s\n", hashValue, nonce)
	fmt.Printf("Execution time: %v seconds\n", executionTime)
	SendMessage(user_key, app_key, hashValue, nonce, executionTime)

}
