package internal

import (
	"sync"
	"time"
)

var (
	sessionTimeLeft int = 120 // Czas sesji w sekundach (2 minuty)
	mu              sync.RWMutex
)

// StartSessionTimer uruchamia timer w osobnej gorutynie
func StartSessionTimer(resetChan <-chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mu.Lock()
			if sessionTimeLeft > 0 {
				sessionTimeLeft--
			}
			mu.Unlock()
		case reset := <-resetChan:
			if reset {
				mu.Lock()
				sessionTimeLeft = 120 // Resetujemy czas do 2 minut
				mu.Unlock()
			}
		}
	}
}

// GetSessionTimeLeft zwraca pozostały czas sesji
func GetSessionTimeLeft() int {
	mu.RLock()
	defer mu.RUnlock()
	return sessionTimeLeft
}
