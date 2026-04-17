package estoque

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

type Product struct {
	Codigo    string `json:"codigo"`
	Descricao string `json:"descricao"`
	Saldo     int    `json:"saldo"`
}

type responseEnvelope[T any] struct {
	Data T `json:"data"`
}

type decrementStockRequest struct {
	Quantidade int `json:"quantidade"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) GetProductByCode(ctx context.Context, codigo string) (*Product, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/products/%s", c.baseURL, codigo), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var envelope responseEnvelope[Product]
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			return nil, err
		}
		return &envelope.Data, nil
	case http.StatusNotFound:
		return nil, ErrProductNotFound
	default:
		return nil, fmt.Errorf("estoque get product returned status %d", resp.StatusCode)
	}
}

func (c *Client) DecrementStock(ctx context.Context, codigo string, quantidade int) (*Product, error) {
	body, err := json.Marshal(decrementStockRequest{Quantidade: quantidade})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, fmt.Sprintf("%s/api/products/%s/decrement", c.baseURL, codigo), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var envelope responseEnvelope[Product]
		if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
			return nil, err
		}
		return &envelope.Data, nil
	case http.StatusNotFound:
		return nil, ErrProductNotFound
	case http.StatusConflict:
		return nil, ErrInsufficientStock
	default:
		return nil, fmt.Errorf("estoque decrement stock returned status %d", resp.StatusCode)
	}
}
