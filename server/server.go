package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Cotacao struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Valor string `gorm:"valor"`
}

type APIResponse struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("cotacao.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Cotacao{})

	http.HandleFunc("/cotacao", handlerCotacao)
	http.ListenAndServe(":8080", nil)
}

func handlerCotacao(w http.ResponseWriter, r *http.Request) {
	c, err := buscaCotacao()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	dbctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()
	result := db.WithContext(dbctx).Create(&Cotacao{Valor: c.USDBRL.Bid})
	if result.Error != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(c.USDBRL.Bid))
}

func buscaCotacao() (*APIResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, error := client.Do(req)
	if error != nil {
		return nil, error
	}
	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)
	if error != nil {
		return nil, error
	}
	var apiResp APIResponse
	error = json.Unmarshal(body, &apiResp)
	if error != nil {
		return nil, error
	}
	return &apiResp, nil
}
