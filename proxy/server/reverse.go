package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// ReverseProxy структура для хранения целевого сервера
type ReverseProxy struct {
	host string
	port string
}

// NewReverseProxy создаёт новый инстанс прокси-сервера
func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

// ReverseProxy мидлварь для проксирования запросов
func (rp *ReverseProxy) ReverseProxy() http.Handler {
	targetURL, err := url.Parse("http://" + rp.host + ":" + rp.port)
	if err != nil {
		panic(err) // Обработка ошибки, если URL недопустим
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Установка заголовков CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Если это предварительный запрос, завершите его здесь
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		r.URL.Host = targetURL.Host
		r.URL.Scheme = targetURL.Scheme
		r.Host = targetURL.Host

		// Проксирование запроса
		proxy.ServeHTTP(w, r)

	})
}
