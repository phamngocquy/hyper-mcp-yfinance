package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	pdk "github.com/extism/go-pdk"
)

var base_url = "https://query2.finance.yahoo.com/v8/finance/chart/%s?interval=%s"
var user_agent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36"

func Call(input CallToolRequest) (CallToolResult, error) {
	args := input.Params.Arguments
	if args == nil {
		return CallToolResult{}, errors.New("Arguments must be provided")
	}
	params := args.(map[string]interface{})
	data, err := getStockData(params)
	if err != nil {
		return createErrorResult(fmt.Sprintf("Error getting stock data: %v", err)), err
	}
	result := CallToolResult{
		Content: []Content{
			{
				Type: ContentTypeText,
				Text: &data,
			},
		},
	}

	return result, nil
}

func Describe() ([]byte, error) {
	desc := map[string]interface{}{
		"name":        "hyper-mcp-yfinance",
		"description": "Stock pricing from Yahoo Finance",
		"version":     "1.0.0",
		"schema": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"symbol": map[string]interface{}{
					"type":        "string",
					"description": "Stock ticker symbol, e.g., AAPL or NVDA.",
				},
				"interval": map[string]interface{}{
					"type":        "string",
					"description": "Time range interval for data aggregation; supports '1d' and '5d'. Default is '1d'.",
				},
			},
			"required": []string{"symbol"},
		},
	}
	return json.Marshal(desc)
}

func createErrorResult(message string) CallToolResult {
	isError := true
	return CallToolResult{
		Content: []Content{
			{
				Type: ContentTypeText,
				Text: &message,
			},
		},
		IsError: &isError,
	}
}

func getStockData(args map[string]interface{}) (string, error) {
	symbol, ok := args["symbol"].(string)
	if !ok {
		return "", errors.New("symbol must be provided")
	}

	interval, ok := args["interval"].(string)
	if !ok {
		interval = "1d"
	}

	url := fmt.Sprintf(base_url,
		url.QueryEscape(symbol),
		url.QueryEscape(interval),
	)

	req := pdk.NewHTTPRequest(pdk.MethodGet, url)
	req.SetHeader("user-agent", user_agent)

	resp := req.Send()

	pdk.Log(pdk.LogDebug, url)

	var response ChartData
	err := json.Unmarshal([]byte(resp.Body()), &response)
	if err != nil {
		return "", errors.New("failed to parse response")
	}

	var output strings.Builder
	for _, result := range response.Chart.Result {
		output.WriteString(fmt.Sprintf("Symbol: %s | ", result.Meta.Symbol))
		output.WriteString(fmt.Sprintf("Exchange: %s |", result.Meta.FullExchangeName))
		output.WriteString(fmt.Sprintf("Timezone: %s |", result.Meta.ExchangeTimezoneName))
		output.WriteString(fmt.Sprintf("Currency: %s \n\n", result.Meta.Currency))

		for idx, quote := range result.Indicators.Quote {
			output.WriteString(fmt.Sprintf("Open: %.6f,", quote.Open[idx]))
			output.WriteString(fmt.Sprintf("Close: %.6f,", quote.Close[idx]))
			output.WriteString(fmt.Sprintf("High: %.6v,", quote.High[idx]))
			output.WriteString(fmt.Sprintf("Low: %.6f,", quote.Low[idx]))
			output.WriteString(fmt.Sprintf("Volume: %v,", quote.Volume[idx]))
		}
	}
	text := output.String()
	return text, nil
}
