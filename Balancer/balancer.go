package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"text/template"
)

func swithToHostN(number int) error {
	// Меняем запись в arp таблице на мак адрес сервиса N
	dNMac := fmt.Sprintf("00:00:00:00:00:0%s", strconv.Itoa(number))

	cmd := exec.Command("docker", "exec", "mn.client", "arp", "-s", "10.0.0.100", dNMac)
	err := cmd.Start()
	if err != nil {
		return err
	}

	// Переписываем конфиг faucet и заставляем docker его перечитать
	err = rewriteConfig(number)
	if err != nil {
		return err
	}
	return nil
}

func rewriteConfig(port int) error {
	// реализовать логику переписывания файла конфигурации
	// Текущая идея такова: этот скрипт работает как systemd юнит и по алёрту от prometheus О том что конкретный service захлёбывается, пере
	// <- перевыбирает свободный сервис переписывает конфиг faucet(надо скопировать faucet.yaml.tmpl) и отправляет докеру сигнал перечитать конфиг
	tmpl, err := template.ParseFiles("faucet.yaml.tmpl")
	if err != nil {
		return err
	}
	outFile, err := os.OpenFile("/etc/faucet/faucet.yaml", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	err = tmpl.Execute(outFile, port)
	if err != nil {
		return err
	}

	cmd := exec.Command("docker", "kill", "-s", "SIGHUP", "faucet")
	err = cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		log.Println("[!] Alert from Prometheus! Starting the switchover...")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error read request body %v\n", err)

		}
		log.Println(string(body))
		err = swithToHostN(1)
		if err != nil {
			log.Printf("Error to switch port %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	// чёто с body + вызов switchToHostN

	w.WriteHeader(http.StatusOK)

}

func main() {
	// надо придумать как регистрировать момент алёрта(превышения нормы) и реализовтаь логику перенаправления
	http.HandleFunc("/", alertHandler)
	log.Fatal(http.ListenAndServe(":9091", nil))
}
