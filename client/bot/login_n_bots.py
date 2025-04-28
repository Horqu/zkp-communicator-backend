import os
import subprocess
import sys
from datetime import datetime
from logging.handlers import RotatingFileHandler
import logging

# Konfiguracja loggera
log_dir = "bot-login-logs"
os.makedirs(log_dir, exist_ok=True)
main_log_file = os.path.join(log_dir, f"main-log-{datetime.now().strftime('%Y-%m-%d_%H-%M-%S')}.log")

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[
        logging.StreamHandler(),  # Logi na konsolę
        RotatingFileHandler(main_log_file, maxBytes=5 * 1024 * 1024, backupCount=1)  # Logi do pliku
    ]
)

# Funkcja do logowania pojedynczego bota
def login_bot(bot_id, login_method):
    log_file = os.path.join(log_dir, f"login-bot-{bot_id}-{datetime.now().strftime('%Y-%m-%d_%H-%M-%S')}.log")
    try:
        logging.info(f"Uruchamianie bota-{bot_id} z metodą logowania: {login_method}...")
        # Uruchomienie bota za pomocą subprocess
        result = subprocess.run(
            ["go", "run", "bot.go", "login", login_method, "1", "true"],  # "1" jako disconnectTime, "true" dla LOGIN_ONLY
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
        
        logging.info(f"Bot-{bot_id} zakończył działanie. Wynik zapisany do logu.")
    except Exception as e:
        with open(log_file, "a") as log:
            log.write(f"[{datetime.now()}] Bot-{bot_id} - Error: {str(e)}\n")
        logging.error(f"Bot-{bot_id} napotkał błąd: {str(e)}")

# Główna funkcja do logowania N botów
def main():
    # Sprawdź, czy podano argumenty
    if len(sys.argv) != 3:
        logging.error("Użycie: python3 login_n_bots.py <liczba_botów> <login_method>")
        return

    try:
        # Pobierz argumenty
        N = int(sys.argv[1])
        login_method = sys.argv[2]
    except ValueError:
        logging.error("Podano nieprawidłowe argumenty.")
        return

    # Zapisz czas rozpoczęcia
    start_time = datetime.now()
    logging.info(f"Rozpoczęcie logowania botów: {start_time}")

    # Logowanie N botów
    for i in range(N):
        login_bot(i + 1, login_method)

    # Zapisz czas zakończenia
    end_time = datetime.now()
    logging.info(f"Zakończenie logowania botów: {end_time}")

    # Oblicz czas trwania
    duration = end_time - start_time
    logging.info(f"Czas trwania programu: {duration}")

    logging.info(f"Wszystkie {N} boty zostały zalogowane.")

if __name__ == "__main__":
    main()