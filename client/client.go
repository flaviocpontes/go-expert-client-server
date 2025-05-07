package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	url := "http://localhost:8080/cotacao"

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatal("Erro acessando servidor:", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Erro ao ler a resposta:", err)
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s", resBody)
	if err != nil {
		log.Fatal("Error writing to file:", err)
	}
}
