package accrual

import (
	"context"
	"errors"
	"time"

	"github.com/71g3pf4c3/gophermart/internal/entity"
	"github.com/71g3pf4c3/gophermart/pkg/logger"
)

// OrderRepo is the minimal repository interface needed by the processor.
type OrderRepo interface {
	GetPendingOrders(ctx context.Context) ([]entity.Order, error)
	UpdateOrderStatus(ctx context.Context, number string, status entity.OrderStatus, accrual *float64) error
}

// Processor polls the accrual service and updates order statuses.
type Processor struct {
	baseURL      string
	pollInterval time.Duration
	logger       logger.Interface
	repo         OrderRepo
	client       *Client
}

// New creates a background accrual processor.
func New(baseURL string, pollInterval time.Duration, l logger.Interface, repo OrderRepo) *Processor {
	return &Processor{
		baseURL:      baseURL,
		pollInterval: pollInterval,
		logger:       l,
		repo:         repo,
		client:       newClient(baseURL),
	}
}

// Run starts polling loop and exits on context cancel.
func (p *Processor) Run(ctx context.Context) {
	if p.baseURL == "" {
		p.logger.Info("accrual processor - disabled (no base URL)")
		<-ctx.Done()
		return
	}

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

	p.logger.Info("accrual processor - started")
	defer p.logger.Info("accrual processor - stopped")

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.processOnce(ctx)
		}
	}
}

func (p *Processor) processOnce(ctx context.Context) {
	orders, err := p.repo.GetPendingOrders(ctx)
	if err != nil {
		p.logger.Error(err)
		return
	}

	for _, order := range orders {
		if ctx.Err() != nil {
			return
		}

		resp, err := p.client.GetOrderAccrual(ctx, order.Number)
		if err != nil {
			var tooMany ErrTooManyRequests
			if errors.As(err, &tooMany) {
				p.logger.Warn("accrual processor - rate limited, sleeping %s", tooMany.RetryAfter)
				select {
				case <-ctx.Done():
					return
				case <-time.After(tooMany.RetryAfter):
				}
				return
			}
			p.logger.Error(err)
			continue
		}

		if resp == nil {
			continue
		}

		status := mapAccrualStatus(resp.Status)
		if status == "" {
			continue
		}

		if err = p.repo.UpdateOrderStatus(ctx, order.Number, status, resp.Accrual); err != nil {
			p.logger.Error(err)
		}
	}
}

func mapAccrualStatus(s string) entity.OrderStatus {
	switch s {
	case "REGISTERED":
		return entity.OrderStatusNew
	case "PROCESSING":
		return entity.OrderStatusProcessing
	case "INVALID":
		return entity.OrderStatusInvalid
	case "PROCESSED":
		return entity.OrderStatusProcessed
	default:
		return ""
	}
}
