package reporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

type callbackReporter struct {
	callbackEndpoint string
}

func NewCallbackReporter(callbackEndpoint string) *callbackReporter {
	return &callbackReporter{callbackEndpoint: callbackEndpoint}
}

func (r *callbackReporter) SendOutput(ctx context.Context, report common.TaskReport) error {
	url := r.callbackEndpoint + "/dogefuzz/callback"
	obj, err := json.Marshal(report)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(obj))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("the request returned with status code: %s", resp.Status)
	}
	return nil
}
