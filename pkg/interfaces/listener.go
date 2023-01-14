package interfaces

import "context"

type Listener interface {
	Name() string
	StartListening(ctx context.Context)
}

type Manager interface {
	Start()
	Shutdown()
}
