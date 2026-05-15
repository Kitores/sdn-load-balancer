package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"text/template"
)

type AlertData struct {
	Annotations  Annotations `json:"annotations"`
	EndsAt       string      `json:"endsAt"`
	StartsAt     string      `json:"startsAt"`
	GeneratorURL string      `json:"generatorURL"`
	Labels       Labels      `json:"labels"`
}

type Annotations struct {
	Summary string `json:"summary"`
}

type Labels struct {
	AlertName string `json:"alertname"`
	ID        string `json:"id"`
	Instance  string `json:"instance"`
	Interface string `json:"interface"`
	Job       string `json:"job"`
	Severity  string `json:"severity"`
}

func extractContainerNumber(summary string) int {
	var n int
	re := regexp.MustCompile(`Switching traffic to (\S+)`)
	matches := re.FindStringSubmatch(summary)
	var name string
	if len(matches) > 1 {
		name = string(matches[1])
	}
	if name == "d1" {
		n = 1
	} else {
		n = 2
	}
	return n
}

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
			return
		}
		log.Println(string(body))
		var alert AlertData
		err = json.Unmarshal(body, &alert)
		if err != nil {
			log.Println("Error: ", err)
			return
		}
		var n int
		n = extractContainerNumber(alert.Annotations.Summary)

		log.Println("n = ", n)
		err = swithToHostN(n)
		if err != nil {
			log.Printf("Error to switch port %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	defer r.Body.Close()
	// чёто с body + вызов switchToHostN

	w.WriteHeader(http.StatusOK)
	return
}

func main() {
	// надо придумать как регистрировать момент алёрта(превышения нормы) и реализовтаь логику перенаправления
	log.Println("Starting listenning on port :9091")
	http.HandleFunc("/", alertHandler)
	log.Fatal(http.ListenAndServe("0.0.0.0:9091", nil))
}
