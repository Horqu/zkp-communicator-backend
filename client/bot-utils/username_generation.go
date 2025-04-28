package botutils

import (
	"client/internal"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gofrs/flock"
)

type BotCredential struct {
	Username   string `json:"username"`
	PrivateKey string `json:"private_key"`
}

// Generates username with pattern bot-<random_string>
// The random string is 16 characters long and consists of lowercase letters and digits.
func GenerateBotUsername() string {

	// Define the length of the random string
	const length = 16

	// Define the characters to choose from
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	// Create a byte slice to hold the random string
	b := make([]byte, length)

	// Generate a random string
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	// Return the generated username
	return "bot-" + string(b)
}

func GenerateBotMessage() string {
	// Define the length of the random string
	const length = 32

	// Define the characters to choose from
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"

	// Create a byte slice to hold the random string
	b := make([]byte, length)

	// Generate a random string
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}

	// Return the generated username
	return "bot-message-" + string(b)
}

func SaveUsernameAndPrivateKey(username string, privateKey string) {
	// Utwórz obiekt blokady dla pliku
	fileLock := flock.New("bot_credentials.json.lock")

	// Czekaj na odblokowanie pliku
	locked, err := fileLock.TryLock()
	if err != nil {
		log.Fatalf("Failed to acquire file lock: %v", err)
	}
	defer fileLock.Unlock()

	// Jeśli plik jest zablokowany, czekaj
	for !locked {
		log.Println("File is locked by another process. Waiting...")
		time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)
		locked, err = fileLock.TryLock()
		if err != nil {
			log.Fatalf("Failed to acquire file lock: %v", err)
		}
	}
	// Otwórz plik w trybie "append" (lub utwórz, jeśli nie istnieje)
	file, err := os.OpenFile("bot_credentials.json", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Odczytaj istniejącą zawartość pliku
	var existingData []map[string]string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&existingData); err != nil && err.Error() != "EOF" {
		log.Printf("Failed to decode existing JSON data: %v", err)
	}

	// Dodaj nowe dane do istniejącej zawartości
	newEntry := map[string]string{
		"username":    username,
		"private_key": privateKey,
	}
	existingData = append(existingData, newEntry)

	// Przenieś wskaźnik pliku na początek i zapisz zaktualizowane dane
	file.Truncate(0)
	file.Seek(0, 0)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Formatowanie JSON-a
	if err := encoder.Encode(existingData); err != nil {
		log.Fatalf("Failed to write JSON to file: %v", err)
	}

	log.Println("Username and private key appended to bot_credentials.json")
}

func LoadRandomBotCredential(filePath string) (BotCredential, error) {
	// Utwórz obiekt blokady dla pliku
	fileLock := flock.New(filePath + ".lock")

	// Czekaj na odblokowanie pliku
	locked, err := fileLock.TryLock()
	if err != nil {
		return BotCredential{}, err
	}
	defer fileLock.Unlock()

	// Jeśli plik jest zablokowany, czekaj losowo od 50 do 150 ms
	for !locked {
		log.Println("File is locked by another process. Waiting...")
		time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)
		locked, err = fileLock.TryLock()
		if err != nil {
			return BotCredential{}, err
		}
	}

	// Otwórz plik JSON
	file, err := os.Open(filePath)
	if err != nil {
		return BotCredential{}, err
	}
	defer file.Close()

	// Wczytaj dane z pliku
	var credentials []BotCredential
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&credentials); err != nil {
		return BotCredential{}, err
	}

	// Sprawdź, czy lista nie jest pusta
	if len(credentials) == 0 {
		return BotCredential{}, fmt.Errorf("no credentials found in file")
	}

	// Wybierz losowy rekord
	randomIndex := rand.Intn(len(credentials))
	return credentials[randomIndex], nil
}

func GetLoginMethod(method string) string {
	switch method {
	case "Schnorr":
		return "Schnorr"
	case "FeigeFiatShamir":
		return "FeigeFiatShamir"
	case "Sigma":
		return "Sigma"
	default:
		methods := []string{"Schnorr", "FeigeFiatShamir", "Sigma"}

		randomIndex := rand.Intn(len(methods))

		return methods[randomIndex]
	}
}

func IsFriendAlreadyAdded(friendList []internal.SimplifiedContact, username string) bool {
	for _, friend := range friendList {
		if friend.Username == username {
			return true
		}
	}
	return false
}
