import os
import subprocess
import sys
from datetime import datetime
from logging.handlers import RotatingFileHandler
import logging

# Konfiguracja loggera
log_dir = "bot-login-unauthorized-logs"
os.makedirs(log_dir, exist_ok=True)
main_log_file = os.path.join(log_dir, f"main-log-{datetime.now().strftime('%Y-%m-%d_%H-%M-%S')}.log")

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(message)s",
    handlers=[
        logging.StreamHandler(),  # Logi na konsolę
        RotatingFileHandler(main_log_file, maxBytes=5 * 1024 * 1024 * 1024, backupCount=1)  # Logi do jednego pliku
    ]
)

# Funkcja do logowania pojedynczego bota
def login_bot(bot_id, login_method):
    try:
        logging.info(f"Uruchamianie próby nr {bot_id} z metodą logowania: {login_method}...")
        # Uruchomienie bota za pomocą subprocess
        result = subprocess.run(
            ["go", "run", "bot.go", "login", login_method, "1", "true", "true"],  # "1" jako disconnectTime, "true" dla LOGIN_ONLY, "true" dla UNAUTHORIZED_ATTEMPT
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )

        # Zapisz wynik do głównego logu
        logging.info(f"[{datetime.now()}] Bot-{bot_id} STDOUT:\n{result.stdout}")
        logging.error(f"[{datetime.now()}] Bot-{bot_id} STDERR:\n{result.stderr}")
        
        logging.info(f"Bot-{bot_id} zakończył działanie.")
    except Exception as e:
        logging.error(f"Bot-{bot_id} napotkał błąd: {str(e)}")

# Główna funkcja do logowania N botów
def main():
    # Sprawdź, czy podano argumenty
    if len(sys.argv) != 3:
        logging.error("Użycie: python3 login-unauthorized-n-tries.py <liczba_prob> <login_method>")
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

    logging.info(f"Wszystkie {N} próby zostały zrealizowane.")

if __name__ == "__main__":
    main()