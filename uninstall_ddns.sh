#!/bin/bash

# Проверка прав root
if [ "$EUID" -ne 0 ]; then
  echo "Пожалуйста, запустите скрипт с правами root."
  exit 1
fi

echo "Остановка сервиса ddns..."
systemctl stop ddns.service >/dev/null 2>&1

echo "Отключение автозапуска сервиса..."
systemctl disable ddns.service >/dev/null 2>&1

echo "Удаление файла сервиса..."
rm -f /etc/systemd/system/ddns.service

echo "Перезагрузка systemd для удаления следов сервиса..."
systemctl daemon-reload
systemctl reset-failed

echo "Удаление исполняемого файла приложения..."
rm -f /usr/bin/ddns/ddns

echo "Удаление конфигурационного файла .env..."
rm -f /usr/bin/ddns/.env

echo "Удаление логов..."
rm -rf /var/log/ddns

echo "Приложение DDNS успешно удалено."
