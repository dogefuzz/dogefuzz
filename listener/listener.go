package listener

import "context"

type Listener interface {
	Name() string
	StartListening(ctx context.Context)
}
