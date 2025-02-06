package helpers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func RoundToDecimals(value float64, decimals int) float64 {
	if decimals < 0 {
		decimals = 0
	}
	multiplier := math.Pow(10, float64(decimals))
	return math.Round(value*multiplier) / multiplier
}

func SendRequest(params map[string]string, method string, endpoint string, service string) (string, error) {
	// Собираем параметры в строку запроса
	data := url.Values{}
	for key, value := range params {
		data.Set(key, value)
	}
	postData := data.Encode()

	// Подготовка HTTP-запроса
	var req *http.Request
	var err error

	if strings.ToUpper(method) == "POST" {
		req, err = http.NewRequest("POST", endpoint, bytes.NewBufferString(postData))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else if strings.ToUpper(method) == "GET" {
		req, err = http.NewRequest("GET", endpoint+"?"+postData, nil)
	} else {
		return "", fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	if service == "fastex" {
		hmacSigner := hmac.New(sha512.New, []byte(os.Getenv("FASTEX_PRIVATE_KEY")))
		hmacSigner.Write([]byte(postData))
		sign := hex.EncodeToString(hmacSigner.Sum(nil))

		headers := map[string]string{
			"Key":  os.Getenv("FASTEX_PUBLIC_KEY"),
			"Sign": sign,
		}
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	// Выполняем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute request: %v", err)
	}

	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			fmt.Printf("failed to close response body: %v\n", cerr)
		}
	}()

	// Чтение ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	return string(body), nil
}
