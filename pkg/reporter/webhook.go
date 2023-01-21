package reporter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/dogefuzz/dogefuzz/pkg/common"
	"github.com/dogefuzz/dogefuzz/pkg/interfaces"
)

var ErrNonSuccessResponse = errors.New("the response was not success")

type webhookReporter struct {
	client          interfaces.HttpClient
	webhookEndpoint string
	timeout         time.Duration
}

func NewWebhookReporter(client interfaces.HttpClient, webhookEndpoint string, timeout time.Duration) *webhookReporter {
	return &webhookReporter{client: client, webhookEndpoint: webhookEndpoint, timeout: timeout}
}

func (r *webhookReporter) SendOutput(ctx context.Context, report common.TaskReport) error {
	url := r.webhookEndpoint + "/dogefuzz/webhook"
	obj, err := json.Marshal(report)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(obj))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return ErrNonSuccessResponse
	}
	return nil
}
