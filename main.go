package main

import (
	"net/http"

	"newtest/internal/handlers"
	"newtest/internal/kpiclient"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

const apiURL = "https://development.kpi-drive.ru/_api/"

func main() {
	//создаем клиент для работы с апи
	kpiVlient := kpiclient.NewKPIClient(apiURL)

	//инициализирукм хендлер
	st := handlers.Init(kpiVlient)

	//создаем роутер
	r := mux.NewRouter()
	r.HandleFunc("/upload-facts", st.NewFacts).Methods(http.MethodPost)

	log.Info("server is running")

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
