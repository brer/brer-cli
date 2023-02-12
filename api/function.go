package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Function struct {
	Image      string        `json:"image"`
	SecretName string        `json:"secretName"`
	Env        []FunctionEnv `json:"env"`
}

type FunctionEnv struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	SecretKey string `json:"secretKey"`
}

func UpdateFunction(url string, token string, data *Function) error {
	var err error

	body, err := json.Marshal(&data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	// initialize http client
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println(res.StatusCode)

	var blob bytes.Buffer
	_, err = io.Copy(&blob, res.Body)
	if err != nil {
		return err
	}

	res.Body.Close()

	fmt.Println(bytes.NewBuffer(blob.Bytes()).String())

	return err
}
