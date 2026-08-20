package alert

import (
	"context"
	"log"

	"github.com/jb843051627/slag-pool/internal/model"
	"github.com/jb843051627/slag-pool/internal/service"
)

type Engine struct {
	alerts *service.AlertService
}

func NewEngine(as *service.AlertService) *Engine {
	return &Engine{alerts: as}
}

func (e *Engine) Evaluate(ctx context.Context) error {
	active, err := e.alerts.ListByStatus(ctx, model.AlertStatusActive)
	if err != nil {
		return err
	}
	for _, a := range active {
		log.Printf("evaluating alert %d: %s", a.ID, a.Message)
	}
	return nil
}

func (e *Engine) Run(ctx context.Context) {
	if err := e.Evaluate(ctx); err != nil {
		log.Printf("alert engine evaluate: %v", err)
	}
}

type Notifier struct {
	channels []string
}

func NewNotifier() *Notifier {
	return &Notifier{channels: []string{"console", "log"}}
}

func (n *Notifier) Notify(level, message string) error {
	log.Printf("[%s] %s", level, message)
	return nil
}

func (n *Notifier) NotifyAlert(alert *model.Alert) error {
	return n.Notify(alert.Level, alert.Message)
}

type Router struct {
	alerts *service.AlertService
}

func NewRouter(as *service.AlertService) *Router {
	return &Router{alerts: as}
}

func (r *Router) Route(ctx context.Context, alertID int64) error {
	a, err := r.alerts.Get(ctx, alertID)
	if err != nil {
		return err
	}
	routedTo := "duty_operator"
	switch a.Level {
	case model.AlertLevelCritical, model.AlertLevelEmergency:
		routedTo = "emergency_team"
	case model.AlertLevelWarning:
		routedTo = "shift_supervisor"
	}
	return r.alerts.Acknowledge(ctx, alertID, routedTo)
}
