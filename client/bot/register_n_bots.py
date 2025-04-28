import os
import subprocess
import sys
from datetime import datetime

# Funkcja do uruchamiania pojedynczego bota
def run_bot(bot_id, log_file):
    try:
        # Uruchomienie bota za pomocą subprocess
        result = subprocess.run(
            ["go", "run", "bot.go", "register"],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )

        # Zapisz wynik do logu
        with open(log_file, "a") as log:
            log.write(f"[{datetime.now()}] Bot-{bot_id}:\n")
            log.write(result.stdout)
            log.write(result.stderr)
            log.write("\n")
        
        print(f"Bot-{bot_id} zakończył działanie. Wynik zapisany do logu.")
    except Exception as e:
        with open(log_file, "a") as log:
            log.write(f"[{datetime.now()}] Bot-{bot_id} - Error: {str(e)}\n")
        print(f"Bot-{bot_id} napotkał błąd: {str(e)}")

# Główna funkcja do uruchamiania N botów
def main():
    # Sprawdź, czy podano argument N
    if len(sys.argv) != 2:
        print("Użycie: python3 register_n_bots.py <liczba_botów>")
        return

    try:
        # Pobierz liczbę botów z argumentu
        N = int(sys.argv[1])
    except ValueError:
        print("Podano nieprawidłową liczbę.")
        return

    # Upewnij się, że folder logów istnieje
    log_dir = "bot-logs"
    os.makedirs(log_dir, exist_ok=True)
    log_file = os.path.join(log_dir, "register.log")

    # Uruchom N botów sekwencyjnie
    for i in range(N):
        run_bot(i + 1, log_file)

    print(f"Wszystkie {N} boty zakończyły działanie. Logi zapisano w {log_file}.")

if __name__ == "__main__":
    main()