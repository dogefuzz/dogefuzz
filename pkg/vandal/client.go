package vandal

import (
	"bufio"
	"fmt"
	"net/http"
	"regexp"
)

const DELIMITER = "================================"

type VandalClient interface {
	Decompile() ([]Block, []Function, error)
}

type vandalClient struct {
}

func NewVandalClient() *vandalClient {
	return &vandalClient{}
}

func (c *vandalClient) Decompile() ([]Block, []Function, error) {
	resp, err := http.Get("http://localhost:5005")
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("bad status :%s", resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
	scanner.Split(bufio.ScanLines)

	blocks := make([]Block, 0)
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

	functions := make([]Function, 0)
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

	return blocks, functions, nil
}
