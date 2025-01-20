package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/vishvananda/netlink"
)

// Глобальные переменные, инициализируемые из .env
var (
	apiURL       string
	userLogin    string
	userPassword string
	domain       string
	logDirectory string
	logFilePath  string
)

// DNSRecord содержит структуру записи A
type DNSRecord struct {
	Priority int    `json:"priority"`
	Value    string `json:"value"`
}

// InputData структура данных для API
type InputData struct {
	FQDN    string                 `json:"fqdn"`
	Records map[string][]DNSRecord `json:"records"`
}

// Инициализация параметров из .env
func initEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %v", err)
	}

	apiURL = os.Getenv("API_URL")
	userLogin = os.Getenv("USER_LOGIN")
	userPassword = os.Getenv("USER_PASSWORD")
	domain = os.Getenv("DOMAIN")
	logDirectory = os.Getenv("LOG_DIRECTORY")
	logFilePath = filepath.Join(logDirectory, "history.log")

	if apiURL == "" || userLogin == "" || userPassword == "" || domain == "" || logDirectory == "" {
		log.Fatalf("Некоторые параметры отсутствуют в .env файле")
	}
}

// Инициализация логирования
func initLogger() {
	if err := os.MkdirAll(logDirectory, 0755); err != nil {
		log.Fatalf("Не удалось создать каталог для логов: %v", err)
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Не удалось открыть файл лога: %v", err)
	}

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

// Получение внешнего IP-адреса
func getExternalIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")
	if err != nil {
		return "", fmt.Errorf("ошибка при получении внешнего IP: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неожиданный статус ответа от API: %s", resp.Status)
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа: %v", err)
	}

	return string(ip), nil
}

// Отправка запроса к API для обновления A записи
func updateARecord(newIP string) error {
	data := InputData{
		FQDN: domain,
		Records: map[string][]DNSRecord{
			"A": {
				{
					Priority: 10,
					Value:    newIP,
				},
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("ошибка кодирования JSON: %v", err)
	}

	query := url.Values{}
	query.Set("login", userLogin)
	query.Set("passwd", userPassword)
	query.Set("input_format", "json")
	query.Set("output_format", "json")
	query.Set("input_data", string(jsonData))

	resp, err := http.PostForm(apiURL, query)
	if err != nil {
		return fmt.Errorf("ошибка при вызове API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус ответа: %s", resp.Status)
	}

	log.Printf("Запись A успешно обновлена на IP: %s", newIP)
	fmt.Printf("Запись A успешно обновлена на IP: %s\n", newIP)
	return nil
}

// Проверка, является ли адрес IPv4
func isIPv4(ip net.IP) bool {
	return ip != nil && ip.To4() != nil
}

// Отслеживание изменений IPv4-адресов
func watchIPv4Changes() {
	updates := make(chan netlink.AddrUpdate)
	done := make(chan struct{})

	if err := netlink.AddrSubscribe(updates, done); err != nil {
		log.Fatalf("Ошибка подписки на события: %v", err)
	}

	log.Println("Отслеживание изменений IPv4-адресов начато...")
	fmt.Println("Отслеживание изменений IPv4-адресов начато...")

	var previousIP string
	for update := range updates {
		ip := update.LinkAddress.IP
		if isIPv4(ip) {
			currentIP := ip.String()
			if currentIP != previousIP {
				currentIPGlobal, err := getExternalIP()
				if err != nil {
					log.Printf("Ошибка определения внешнего IP: %v", err)
					continue
				}
				log.Printf("Изменение IP: старый %s -> новый %s", previousIP, currentIP)
				fmt.Printf("Изменение IP: старый %s -> новый %s\n", previousIP, currentIP)
				log.Printf("Изменение внешнего IP -> новый %s", currentIPGlobal)
				fmt.Printf("Изменение внешнего IP -> новый %s\n", currentIPGlobal)
				if err := updateARecord(currentIPGlobal); err != nil {
					log.Printf("Ошибка обновления записи A: %v", err)
				}
				previousIP = currentIP
			}
		}
	}
}

func main() {
	initEnv()
	initLogger()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go watchIPv4Changes()

	<-signalChan
	log.Println("Завершение программы...")
	fmt.Println("Завершение программы...")
}
