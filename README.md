# Dynamic DNS Updater (DDNS)

Этот код работает только в системах на базе Linux, так как использует netlink — механизм, специфичный для Linux.

EN: This code only works on Linux because it uses netlink, a Linux-specific mechanism.

## Описание проекта (Русский)

**Dynamic DNS Updater (DDNS)** — это утилита на языке Go, которая автоматически обновляет A-запись домена через API DNS-провайдера Beget.com при изменении IPv4-адреса. Утилита работает в фоновом режиме, отслеживая изменения IP-адресов, и записывает события в лог-файл.

### Основной функционал
- Отслеживание изменений IPv4-адресов в реальном времени.
- Автоматическое обновление A-записи через API DNS-провайдера.
- Логирование событий в файл `/var/log/ddns/history.log`.
- Работа в режиме демона с автоматическим запуском при старте системы.

---

### Настройка `.env` файла

Пример `.env` файла:

```ini
API_URL=https://api.beget.com/api/dns/changeRecords
USER_LOGIN=userlogin          # Логин для API
USER_PASSWORD=password        # Пароль для API
DOMAIN=example.com            # Домен для обновления записи
LOG_DIRECTORY=/var/log/ddns   # Каталог для логов
```

**Описание параметров**:
- `API_URL` — URL для доступа к API DNS-провайдера.
- `USER_LOGIN` — логин для авторизации в API.
- `USER_PASSWORD` — пароль для авторизации в API.
- `DOMAIN` — имя домена для обновления записи.
- `LOG_DIRECTORY` — путь к каталогу для хранения логов.

---

### Установка и запуск

1. Добавьте файл `.env` в директорию с приложением.
2. Сделайте установочный скрипт исполняемым:
   ```bash
   chmod +x install_ddns.sh
   ```
3. Запустите установку:
   ```bash
   sudo ./install_ddns.sh
   ```
4. Проверьте статус сервиса:
   ```bash
   systemctl status ddns.service
   ```

---

### Управление сервисом

Сервис DDNS работает в режиме демона, его можно контролировать с помощью `systemd`:

- Запуск:
  ```bash
  sudo systemctl start ddns.service
  ```
- Остановка:
  ```bash
  sudo systemctl stop ddns.service
  ```
- Перезапуск:
  ```bash
  sudo systemctl restart ddns.service
  ```
- Проверка статуса:
  ```bash
  sudo systemctl status ddns.service
  ```

---

## Project Description (English)

**Dynamic DNS Updater (DDNS)** is a Go utility that automatically updates the A-record of a domain via the Beget.com DNS provider API when the IPv4 address changes. The utility runs in the background, monitoring IP address changes, and records events in a log file.

### Main Functionality
- Monitoring IPv4 address changes in real time.
- Automatic A-record update via the DNS provider API.
- Logging events to the file `/var/log/ddns/history.log`.
- Working in daemon mode with automatic startup at system startup.

---

### Setting up `.env` file

Example of `.env` file:

```ini
API_URL=https://api.beget.com/api/dns/changeRecords
USER_LOGIN=userlogin # API login
USER_PASSWORD=password # API password
DOMAIN=example.com # Domain for updating records
LOG_DIRECTORY=/var/log/ddns # Directory for logs
```

**Parameter description**:
- `API_URL` — URL for accessing the DNS provider API.
- `USER_LOGIN` — login for authorization in the API.
- `USER_PASSWORD` — password for authorization in the API.
- `DOMAIN` — domain name for updating records.
- `LOG_DIRECTORY` — path to the directory for storing logs.

---

### Installation and launch

1. Add the `.env` file to the application directory.
2. Make the installation script executable:
```bash
chmod +x install_ddns.sh
```
3. Run the installation:
```bash
sudo ./install_ddns.sh
```
4. Check the service status:
```bash
systemctl status ddns.service
```

---

### Service management

The DDNS service runs in daemon mode, it can be controlled with `systemd`:

- Start:
```bash
sudo systemctl start ddns.service
```
- Stop:
```bash
sudo systemctl stop ddns.service
```
- Restart:
```bash
sudo systemctl restart ddns.service
```
- Check the status:
```bash
sudo systemctl status ddns.service
```

---

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
