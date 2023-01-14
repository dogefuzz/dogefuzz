package listener

import (
	"context"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
)

type manager struct {
	cfg    *config.Config
	env    *env
	cancel context.CancelFunc
}

func NewManager(cfg *config.Config) *manager {
	return &manager{cfg: cfg, env: NewEnv(cfg)}
}

func (m *manager) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel

	listeners := m.getAvailableListeners()
	for _, id := range m.cfg.EventConfig.EnabledListeners {
		if listener, ok := listeners[id]; ok {
			go listener.StartListening(ctx)
			m.env.logger.Sugar().Infof("starting listener %s", id)
		} else {
			m.env.logger.Sugar().Warnf("ignore listener %s because it's not implemented", id)
		}
	}
}

func (m *manager) Shutdown() {
	m.cancel()
}

func (m *manager) getAvailableListeners() map[string]interfaces.Listener {
	return map[string]interfaces.Listener{
		m.env.ContractDeployerListener().Name():   m.env.ContractDeployerListener(),
		m.env.ExecutionAnalyticsListener().Name(): m.env.ExecutionAnalyticsListener(),
		m.env.FuzzerListener().Name():             m.env.FuzzerListener(),
		m.env.ReporterListener().Name():           m.env.ReporterListener(),
	}
}
