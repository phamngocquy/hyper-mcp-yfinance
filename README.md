# Extism Go PDK Plugin

A [hyper-mcp](https://github.com/tuananh/hyper-mcp) plugin retrieves stock pricing info from [Yahoo Finance](https://finance.yahoo.com/).

<p align="center">
  <img src="./assets/yfinance.png">
</p>

## Usage

The tool accepts a single parameter:

- `symbol`: The stock symbol
- `interval`: Time to aggregate data: 1d,5d,1mo,3mo,6mo,1y,2y,5y,10y,ytd. The default is: 1d

- Add the plugin to your hyper-mcp configuration:

```json
{
  "plugins": [
    {
      "name": "yfinance",
      "path": "oci://ghcr.io/phamngocquy/hyper-mcp-yfinance:latest",
      "runtime_config": {
        "allowed_host": "query2.finance.yahoo.com"
      }
    }
  ]
}
```

See more documentation at https://github.com/extism/go-pdk and
[join us on Discord](https://extism.org/discord) for more help.
