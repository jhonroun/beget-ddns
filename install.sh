#!/bin/bash
export PATH=$PATH:/usr/local/go/bin
# Проверка прав root
if [ "$EUID" -ne 0 ]; then
  echo "Пожалуйста, запустите скрипт с правами root."
  exit 1
fi

# Компиляция приложения
echo "Компиляция приложения..."
mkdir -p /usr/bin/ddns
go build -o /usr/bin/ddns/ddns main.go
if [ $? -ne 0 ]; then
  echo "Ошибка компиляции. Убедитесь, что Go установлен."
  exit 1
fi

# Копирование .env файла
echo "Копирование .env файла..."
cp .env /usr/bin/ddns/.env

# Установка прав доступа для каталога логов
echo "Создание каталога для логов..."
mkdir -p /var/log/ddns
chmod 755 /var/log/ddns

# Установка сервиса
echo "Установка сервиса..."
cat <<EOF > /etc/systemd/system/ddns.service
[Unit]
Description=Dynamic DNS Updater Service
After=network.target

[Service]
EnvironmentFile=/usr/bin/ddns/.env
ExecStart=/usr/bin/ddns/ddns
Restart=always
User=root
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=ddns

[Install]
WantedBy=multi-user.target
EOF

# Перезагрузка systemd
systemctl daemon-reload

# Включение и запуск сервиса
systemctl enable ddns.service
systemctl start ddns.service

echo "Сервис установлен и запущен."
