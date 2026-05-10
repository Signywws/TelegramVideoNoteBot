package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	baseURL    = "https://api.telegram.org/bot"
	getFileURL = "https://api.telegram.org/file/bot"
)

type Client struct {
	Token   string
	baseURL string
	Client  *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		Token:   token,
		baseURL: baseURL + token,
		Client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// Получаем Updates
func (c *Client) GetUpdate(offset int, timeout int) ([]Update, error) {
	url := fmt.Sprintf("%s/getUpdates?offset=%d&timeout=%d", c.baseURL, offset, timeout)
	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %w", err)
	}

	var result Results
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Umnarshal body: %w", err)
	}

	return result.Results, nil

}

func (c *Client) SendMessage(chatId int, text string) error {
	url := fmt.Sprintf("%s/sendMessage", c.baseURL)
	reqBody := map[string]string{
		"chat_id": strconv.FormatInt(int64(chatId), 10),
		"text":    text,
	}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("Marshal body: %w", err)
	}

	resp, err := c.Client.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("SendMessage POST: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SendMessage responce: %w", err)
	}
	return nil
}

func (c *Client) GetFile(fileID string) (*File, error) {
	url := fmt.Sprintf("%s/getFile?file_id=%s", c.baseURL, fileID)

	resp, err := c.Client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("getFile request: %w", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var res struct {
		Ok          bool   `json:"ok"`
		Result      File   `json:"result"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(body, &res); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	if !res.Ok {
		return nil, fmt.Errorf("getFile not ok: ", res.Description)
	}

	return &res.Result, nil
}

func (c *Client) SendVideoNote(chatID int64, filePath string, duration int) error {
	url := fmt.Sprintf("%s/sendVideoNote", c.baseURL)
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	w.WriteField("chat_id", strconv.FormatInt(chatID, 10))
	w.WriteField("duration", strconv.Itoa(duration))

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	defer f.Close()

	part, err := w.CreateFormFile("video_note", filepath.Base(filePath))
	if err != nil {
		return fmt.Errorf("create form file: %w", err)
	}

	if _, err := io.Copy(part, f); err != nil {
		return fmt.Errorf("copy file: %w", err)
	}

	w.Close()

	req, err := http.NewRequest(http.MethodPost, url, &b)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("sendVideoNote request: %w", err)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	log.Printf("sendVideoNote response: %s", string(body))
	return nil
}

func (c *Client) DownloadFile(filePath string, writer io.Writer) error {
	url := fmt.Sprintf("%s%s/%s", getFileURL, c.Token, filePath)
	resp, err := c.Client.Get(url)
	if err != nil {
		return fmt.Errorf("download file: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download file: status %d", resp.StatusCode)
	}

	_, err = io.Copy(writer, resp.Body)
	return err

}
