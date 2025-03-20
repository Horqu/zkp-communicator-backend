package wsresponses

import (
	"encoding/json"
	"log"

	db "github.com/Horqu/zkp-communicator-backend/cmd/database"
	"github.com/Horqu/zkp-communicator-backend/cmd/database/models"
	"github.com/Horqu/zkp-communicator-backend/cmd/internal"
)

type SimplifiedContact struct {
	Username string `json:"username"`
	Status   string `json:"status"`
}

func ResponseSolveSuccess(publicKey string, friendList []models.Contact) internal.Response {
	// Tworzymy nową listę uproszczonych kontaktów
	var simplifiedContacts []SimplifiedContact

	for _, contact := range friendList {
		// Pobieramy username kontaktu na podstawie jego ID
		var user models.User
		if err := db.GetDBInstance().First(&user, contact.ContactID).Error; err != nil {
			log.Printf("Failed to fetch user for contact ID %d: %v", contact.ContactID, err)
			continue
		}

		// Dodajemy uproszczony kontakt do listy
		simplifiedContacts = append(simplifiedContacts, SimplifiedContact{
			Username: user.Username,
			Status:   contact.Status,
		})
	}

	// Tworzymy strukturę odpowiedzi
	responseData := struct {
		PublicKey  string              `json:"publicKey"`
		FriendList []SimplifiedContact `json:"friendList"`
	}{
		PublicKey:  publicKey,
		FriendList: simplifiedContacts,
	}

	// Serializujemy odpowiedź do JSON
	data, err := json.Marshal(responseData)
	if err != nil {
		log.Printf("Failed to marshal response data: %v", err)
		return internal.Response{
			Command: internal.ResponseSolveSuccess,
			Data:    "Error",
		}
	}

	return internal.Response{
		Command: internal.ResponseSolveSuccess,
		Data:    string(data),
	}
}
