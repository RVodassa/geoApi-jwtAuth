package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}

func HandleGeocode(w http.ResponseWriter, r *http.Request) {
	apiKey, secretKey := os.Getenv("ApiKey"), os.Getenv("SecretKey")
	geoService := NewGeoService(apiKey, secretKey)

	var data GeocodeRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("ERROR: Не удалось прочитать JSON данные: %v\n", err)
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	if data.Lat == "" || data.Lng == "" {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	log.Printf("INFO: Запрос на геокодирование: lat=%s, lng=%s\n", data.Lat, data.Lng)

	geoRes, err := geoService.GeoCode(data.Lat, data.Lng)
	if err != nil {
		//log.Println("ERROR: Не удалось выполнить геокодирование:", err)
		http.Error(w, "Ошибка геокодирования", http.StatusInternalServerError)
		return
	}

	// Подготовка ответа в формате, который ожидает клиент
	response := SearchResponse{
		Addresses: geoRes,
	}

	// Логирование ответа перед отправкой
	log.Printf("DEBUG: Response sent to client: %+v\n", response)

	res, err := json.Marshal(response)
	if err != nil {
		log.Println("ERROR: Не удалось маршализировать JSON:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	apiKey, secretKey := os.Getenv("ApiKey"), os.Getenv("SecretKey")
	geoService := NewGeoService(apiKey, secretKey)

	var data SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Printf("ERROR: Не удалось прочитать JSON данные: %v\n", err)
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	if data.Query == "" {
		http.Error(w, "Неверные данные", http.StatusBadRequest)
		return
	}

	log.Printf("INFO: Запрос на поиск: query=%s", data.Query)

	geoRes, err := geoService.AddressSearch(data.Query)
	if err != nil {
		//log.Println("ERROR: Не удалось найти адрес:", err)
		http.Error(w, "Ошибка поиска адреса", http.StatusInternalServerError)
		return
	}

	// Подготовка ответа в формате, который ожидает клиент
	response := SearchResponse{
		Addresses: geoRes,
	}

	// Логирование ответа перед отправкой
	log.Printf("DEBUG: Response sent to client: %+v\n", response)

	res, err := json.Marshal(response)
	if err != nil {
		log.Println("ERROR: Не удалось маршализировать JSON:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
