package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		baseURL: fmt.Sprintf("https://api.telegram.org/bot%s", token),
		client: &http.Client{
			Timeout: time.Second * 30, // long pooling ожидание
		},
	}
}

func (c *Client) GetUpdates(offset int, timeout int) ([]Update, error) {
	url := fmt.Sprintf("%s/getUpdates?offset=%d&timeout=%d", c.baseURL, offset, timeout)
	resp, err := c.client.Get(url) // делаем запрос
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) // читаем поток body
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	var result struct { // обертка для структуры ответа
		OK     bool     `json:"ok"`
		Result []Update `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil { // кодируем body[]byte в структуру result
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	if !result.OK {
		return nil, fmt.Errorf("telegram API not OK")
	}
	return result.Result, nil
}

func (c *Client) SendMessage(chatID int, text string) error {
	url := fmt.Sprintf("%s/sendMessage", c.baseURL)

	payload := map[string]string{
		"chatID": strconv.FormatInt(int64(chatID), 10),
		"text":   text,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	resp, err := c.client.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("http post: %w", err)
	}
	defer resp.Body.Close()

	return nil
}
