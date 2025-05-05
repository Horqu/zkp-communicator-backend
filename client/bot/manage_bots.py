import os
import subprocess
import sys
import time
from datetime import datetime, timedelta
from logging.handlers import RotatingFileHandler
import logging
import threading
import signal
import random

# Flaga do kontrolowania działania programu
running = True

# Konfiguracja loggera
log_dir = "bot-logs"
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

# Funkcja do obsługi sygnału przerwania (Ctrl+C)
def handle_sigint(signum, frame):
    global running
    logging.info("Graceful shutdown initiated...")
    running = False

# Funkcja do uruchamiania pojedynczego bota
def run_bot(bot_id, log_file, login_method, disconnect_time):
    while running:
        try:
            logging.info(f"Uruchamianie bota-{bot_id}...")
            # Uruchomienie bota za pomocą subprocess
            result = subprocess.run(
                ["go", "run", "bot.go", "login", login_method, disconnect_time],
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

        # Jeśli bot się wyłączy, uruchom go ponownie
        if running:
            delay = random.randint(30, 60)  # Losowe opóźnienie od 30 do 60 sekund
            logging.info(f"Restarting Bot-{bot_id} in {delay} seconds...")
            time.sleep(delay)

# Główna funkcja do uruchamiania N botów
def main():
    global running

    # Sprawdź, czy podano argumenty
    if len(sys.argv) != 5:
        logging.error("Użycie: python3 manage_bots.py <liczba_botów> <login_method> <disconnect_time> <czas_działania_w_minutach>")
        return

    try:
        # Pobierz argumenty
        N = int(sys.argv[1])
        login_method = sys.argv[2]
        disconnect_time = sys.argv[3]
        run_duration = int(sys.argv[4])  # Czas działania w minutach
    except ValueError:
        logging.error("Podano nieprawidłowe argumenty.")
        return

    # Ustaw czas zakończenia działania
    end_time = datetime.now() + timedelta(minutes=run_duration)

    # Uruchom N botów w osobnych wątkach
    threads = []
    for i in range(N):
        log_file = os.path.join(log_dir, f"manage-bots-{datetime.now().strftime('%Y-%m-%d_%H-%M-%S')}-bot-{i+1}.log")
        
        thread = threading.Thread(target=run_bot, args=(i + 1, log_file, login_method, disconnect_time))
        threads.append(thread)
        thread.start()

    # Główna pętla kontrolująca czas działania
    try:
        while running and datetime.now() < end_time:
            time.sleep(1)
    except KeyboardInterrupt:
        handle_sigint(None, None)

    # Poczekaj na zakończenie wszystkich wątków
    running = False
    for thread in threads:
        thread.join()

    logging.info("Wszystkie boty zostały zakończone.")

if __name__ == "__main__":
    # Obsługa sygnału Ctrl+C
    signal.signal(signal.SIGINT, handle_sigint)
    main()