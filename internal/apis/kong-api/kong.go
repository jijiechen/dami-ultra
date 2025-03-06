package kong_api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

/*
curl -X POST http://ec2-54-166-250-69.compute-1.amazonaws.com:8001/routes -H 'application/json'
--data '{"name": "test_route"}'
*/

func ApplyKongConfig(url string, msg string) error {
	resp, err := http.Post(url, "application/json", bytes.NewReader([]byte(msg)))
	if err != nil {
		return fmt.Errorf("failed to send Post request, error: %v", err)
	}

	kongRespMsg, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse kong response, error: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		fmt.Println("==================", string(kongRespMsg))
		return nil

	case http.StatusBadRequest:
		return fmt.Errorf("failed to apply kong configuration, error: %s", string(kongRespMsg))

	default:
		return fmt.Errorf("receive kong http status code: %d, msg: %s", resp.StatusCode, string(kongRespMsg))
	}
}

/*
curl -X GET http://ec2-54-166-250-69.compute-1.amazonaws.com:8001/routes
*/

func GetKongConfig(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to send Get request, error: %v", err)
	}

	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse kong response, error: %v", err)
	}

	return string(msg), nil
}

func DeleteKongConfig(url string, routeName string) error {
	request, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", url, routeName), nil)
	if err != nil {
		return fmt.Errorf("failed to create Delete request, error: %v", err)
	}

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("failed to send Delete request, error: %v", err)
	}

	msg, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to parse kong response, error: %v", err)
	}

	fmt.Printf("=========== HttpStatusCode: %d, messages: %s\n", resp.StatusCode, msg)

	return nil
}
