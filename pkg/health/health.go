package health

import (
	"github.com/go-chi/render"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HealthCheckType string

const (
	HEALTH_TYPE_STARTUP   HealthCheckType = "startup"
	HEALTH_TYPE_READINESS HealthCheckType = "readiness"
	HEALTH_TYPE_LIVENESS  HealthCheckType = "liveness"
)

type HealthService struct {
	router chi.Router
	checks map[HealthCheckType]map[string]CheckFunc
}

func NewHealth(opts ...Option) (*HealthService, error) {
	h := &HealthService{
		router: chi.NewRouter(),
		checks: make(map[HealthCheckType]map[string]CheckFunc),
	}

	handleOptions(h, opts)

	h.router.Get("/startup", func(w http.ResponseWriter, r *http.Request) {
		err := h.runChecks(HEALTH_TYPE_STARTUP, w, r)
		if err != nil {
			slog.Error(
				"health check process failed",
				slog.String("type", "startup"),
				slog.String("error", err.Error()),
			)
		}
	})

	h.router.Get("/readiness", func(w http.ResponseWriter, r *http.Request) {
		err := h.runChecks(HEALTH_TYPE_READINESS, w, r)
		if err != nil {
			slog.Error(
				"health check process failed",
				slog.String("type", "readiness"),
				slog.String("error", err.Error()),
			)
		}
	})

	h.router.Get("/liveness", func(w http.ResponseWriter, r *http.Request) {
		err := h.runChecks(HEALTH_TYPE_LIVENESS, w, r)
		if err != nil {
			slog.Error(
				"health check process failed",
				slog.String("type", "liveness"),
				slog.String("error", err.Error()),
			)
		}
	})

	return h, nil
}

func (h *HealthService) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.router.ServeHTTP(writer, request)
}

func (h *HealthService) AddCheck(checkType HealthCheckType, name string, check CheckFunc) {
	if h.checks[checkType] == nil {
		h.checks[checkType] = make(map[string]CheckFunc)
	}

	h.checks[checkType][name] = check
}

func (h *HealthService) AddStartupCheck(name string, check CheckFunc) {
	h.AddCheck(HEALTH_TYPE_STARTUP, name, check)
}

func (h *HealthService) AddReadinessCheck(name string, check CheckFunc) {
	h.AddCheck(HEALTH_TYPE_READINESS, name, check)
}

func (h *HealthService) AddLivenessCheck(name string, check CheckFunc) {
	h.AddCheck(HEALTH_TYPE_LIVENESS, name, check)
}

func (h *HealthService) runChecks(checkType HealthCheckType, w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	resp := &Response{
		Success: true,
		Checks:  make(map[string]CheckResult),
	}

	for name, check := range h.checks[checkType] {
		if err := check(ctx); err != nil {
			resp.Checks[name] = CheckResult{
				Success: false,
				Message: err.Error(),
			}

			resp.Success = false
		} else {
			resp.Checks[name] = CheckResult{
				Success: true,
				Message: "OK",
			}
		}
	}

	if resp.Success {
		render.Status(r, http.StatusOK)
	} else {
		render.Status(r, http.StatusBadRequest)
	}

	return render.Render(w, r, resp)
}
