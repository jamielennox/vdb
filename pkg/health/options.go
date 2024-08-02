package health

type options struct {
	checks map[HealthCheckType]map[string]CheckFunc
}

type Option func(*options)

func WithCheck(checkType HealthCheckType, name string, check CheckFunc) Option {
	return func(o *options) {
		if o.checks[checkType] == nil {
			o.checks[checkType] = make(map[string]CheckFunc)
		}

		o.checks[checkType][name] = check
	}
}

func WithStartupCheck(name string, check CheckFunc) Option {
	return WithCheck(HEALTH_TYPE_STARTUP, name, check)
}

func WithReadinessCheck(name string, check CheckFunc) Option {
	return WithCheck(HEALTH_TYPE_READINESS, name, check)
}

func WithLivenessCheck(name string, check CheckFunc) Option {
	return WithCheck(HEALTH_TYPE_LIVENESS, name, check)
}

func handleOptions(h *HealthService, opts []Option) options {
	o := options{
		checks: make(map[HealthCheckType]map[string]CheckFunc),
	}

	for _, opt := range opts {
		opt(&o)
	}

	for checkType, checks := range o.checks {
		for name, check := range checks {
			h.AddCheck(checkType, name, check)
		}
	}

	return o
}
