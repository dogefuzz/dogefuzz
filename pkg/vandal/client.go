package vandal

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/dogefuzz/dogefuzz/pkg/common"
)

const DELIMITER = "================================"

type VandalDecompileRequest struct {
	Source string `json:"source"`
}

type vandalClient struct {
	endpoint string
}

func NewVandalClient(vandalEndpoint string) *vandalClient {
	return &vandalClient{endpoint: vandalEndpoint}
}

func (c *vandalClient) Decompile(ctx context.Context, source string) ([]common.Block, []common.Function, error) {
	body, err := json.Marshal(VandalDecompileRequest{Source: source})
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("bad status :%s", res.Status)
	}

	scanner := bufio.NewScanner(res.Body)
	scanner.Split(bufio.ScanLines)

	blocks := make([]common.Block, 0)
	currentBlock := make([]string, 0)
	for scanner.Scan() {
		value := scanner.Text()
		if value == "" {
			continue
		}

		if value == DELIMITER {
			blocks = append(blocks, NewBlockFromLines(currentBlock))
			currentBlock = make([]string, 0)
			continue
		}
		currentBlock = append(currentBlock, value)
	}

	functions := make([]common.Function, 0)
	currentFunction := make([]string, 0)
	for _, line := range currentBlock {
		regexFunction := regexp.MustCompile("Function (.*):")
		if regexFunction.MatchString(line) {
			if len(currentFunction) > 0 {
				functions = append(functions, NewFunctionFromLines(currentFunction))
				currentFunction = make([]string, 0)
			}
			continue
		}
		currentFunction = append(currentFunction, line)
	}
	if len(currentFunction) > 0 {
		functions = append(functions, NewFunctionFromLines(currentFunction))
	}

	return blocks, functions, nil
}
