package accrual

import (
	"context"
	"time"

	"github.com/71g3pf4c3/gophermart/pkg/logger"
)

// Processor polls accrual service in background.
type Processor struct {
	baseURL      string
	pollInterval time.Duration
	logger       logger.Interface
}

// New creates a background processor.
func New(baseURL string, pollInterval time.Duration, l logger.Interface) *Processor {
	return &Processor{
		baseURL:      baseURL,
		pollInterval: pollInterval,
		logger:       l,
	}
}

// Run starts polling loop and exits on context cancel.
func (p *Processor) Run(ctx context.Context) {
	if p.baseURL == "" {
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
			p.logger.Debug("accrual processor - poll tick")
		}
	}
}
