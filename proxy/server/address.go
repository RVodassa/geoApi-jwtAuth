package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ekomobile/dadata/v2/api/suggest"
	"github.com/ekomobile/dadata/v2/client"
)

// Структуры для запросов и ответов
type SearchRequest struct {
	Query string `json:"query"`
}

type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type Address struct {
	City   string `json:"city"`
	Street string `json:"street"`
	House  string `json:"house"`
	Lat    string `json:"lat"`
	Lon    string `json:"lon"`
}

type ResponseAddress struct {
	Addresses []*Address `json:"addresses"`
}

type GeoService struct {
	api       *suggest.Api
	apiKey    string
	secretKey string
}

// Интерфейс для провайдера геосервиса
type GeoProvider interface {
	AddressSearch(input string) ([]*Address, error)
	GeoCode(lat, lng string) ([]*Address, error)
}

func NewGeoService(apiKey, secretKey string) *GeoService {
	endpointUrl, err := url.Parse("https://suggestions.dadata.ru/suggestions/api/4_1/rs/")
	if err != nil {
		return nil
	}

	creds := client.Credentials{
		ApiKeyValue:    apiKey,
		SecretKeyValue: secretKey,
	}

	api := suggest.Api{
		Client: client.NewClient(endpointUrl, client.WithCredentialProvider(&creds)),
	}

	return &GeoService{
		api:       &api,
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

// Метод поиска адресов
func (g *GeoService) AddressSearch(input string) ([]*Address, error) {
	var res []*Address
	rawRes, err := g.api.Address(context.Background(), &suggest.RequestParams{Query: input})
	if err != nil {
		return nil, err
	}

	// Преобразуем результаты API DaData в массив объектов Address
	for _, r := range rawRes {
		if r.Data.City == "" || r.Data.Street == "" {
			continue
		}
		res = append(res, &Address{
			City:   r.Data.City,
			Street: r.Data.Street,
			House:  r.Data.House,
			Lat:    r.Data.GeoLat,
			Lon:    r.Data.GeoLon,
		})
	}

	return res, nil
}

// Метод геокодирования (поиск адресов по координатам)
func (g *GeoService) GeoCode(lat, lng string) ([]*Address, error) {
	httpClient := &http.Client{}
	data := strings.NewReader(fmt.Sprintf(`{"lat": %s, "lon": %s}`, lat, lng))

	req, err := http.NewRequest("POST", "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address", data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", g.apiKey))

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var geoCode geoCode
	if err := json.NewDecoder(resp.Body).Decode(&geoCode); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	var res []*Address
	for _, r := range geoCode.Suggestions {
		address := &Address{
			City:   string(r.Data.City),
			Street: string(r.Data.Street),
			House:  r.Data.House,
			Lat:    r.Data.GeoLat,
			Lon:    r.Data.GeoLon,
		}
		res = append(res, address)
	}

	return res, nil
}

// Структура для обработки ответа от API DaData для геолокации
type geoCode struct {
	Suggestions []struct {
		Data struct {
			City   string `json:"city"`
			Street string `json:"street"`
			House  string `json:"house"`
			GeoLat string `json:"geo_lat"`
			GeoLon string `json:"geo_lon"`
		} `json:"data"`
	} `json:"suggestions"`
}
