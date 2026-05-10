package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ErrTooManyRequests is returned when accrual service returns 429.
type ErrTooManyRequests struct {
	RetryAfter time.Duration
}

func (e ErrTooManyRequests) Error() string {
	return fmt.Sprintf("accrual: too many requests, retry after %s", e.RetryAfter)
}

// AccrualResponse is the response from the accrual service.
type AccrualResponse struct {
	Order   string   `json:"order"`
	Status  string   `json:"status"`
	Accrual *float64 `json:"accrual,omitempty"`
}

// Client calls the external accrual service.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// newClient creates an accrual HTTP client.
func newClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetOrderAccrual calls GET /api/orders/{number}.
// Returns (nil, nil) when the order is not registered in accrual (204).
// Returns ErrTooManyRequests when rate limited (429).
func (c *Client) GetOrderAccrual(ctx context.Context, number string) (*AccrualResponse, error) {
	url := c.baseURL + "/api/orders/" + number

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("accrual client - NewRequest: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, fmt.Errorf("accrual client - Do: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var result AccrualResponse
		if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("accrual client - Decode: %w", err)
		}
		return &result, nil

	case http.StatusNoContent:
		return nil, nil

	case http.StatusTooManyRequests:
		retryAfter := 60 * time.Second
		if v := resp.Header.Get("Retry-After"); v != "" {
			if secs, err := strconv.Atoi(v); err == nil {
				retryAfter = time.Duration(secs) * time.Second
			}
		}
		return nil, ErrTooManyRequests{RetryAfter: retryAfter}

	default:
		return nil, fmt.Errorf("accrual client - unexpected status: %d", resp.StatusCode)
	}
}
