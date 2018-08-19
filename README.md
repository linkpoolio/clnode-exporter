# Chainlink Node Stats Exporter
A lightweight Golang metrics/stats collector and exporter for [Chainlink](https://chain.link/) nodes. It exports stats via Prometheus, REST and WebSockets. Supports scraping multiple nodes by a config file or a REST API.

**NOTE:** Currently not unit/integration tested, although has been extensively tested in practise with multiple nodes.

**NOTE 2:** This exporter relies on the Chainlink node endpoint `/stats` that currently hasn't been merged to master.

### Setup Instructions
#### Local Install
Make sure [Golang](https://golang.org/pkg/) is installed.

Build:
```
make build
```

Then run the exporter:
```
./clnode-exporter
```

#### Docker
To run the container:
```
docker run -it -p 8080:8080 linkpoolio/clnode-exporter
```

Container also supports passing in CLI arguments.

### Configuration
#### `port`
The port for the exporter to listen on for the Prometheus, REST and WS endpoints. Defaults to 8080.

**Example Usage:** `-port=8080`
#### `clientApiPort`
The port to be used for the client management API. See `clientApi` for enabling. Defaults to 8082.

**Example Usage:** `-clientApiPort=8082`
#### `tickerInterval`
The interval to be used for scraping metrics from each node. Defaults to 15s. Argument needs to be specified in Golangs `time.Duration` format.

**Example Usage:** `-tickerInterval=1m`
#### `configFile`
The location of the configuration file that contains the node URLs and credentials. Needs to be specified in the following format:
```json
[
  {
    "url": "https://localhost:6688",
    "username": "chainlink",
    "password": "twochains"
  }
]
```
Multiple nodes are specified by simply adding more objects to the JSON array.

This configuration file is refreshed at intervals specified by `tickerInterval`. This allows nodes to be removed/added during service operation without restart.

**Example Usage:** `-configFile=/etc/exporter/nodes.json`
#### `prom`
Flag for whether to enable Prometheus exporting on the `/metrics` endpoint. Defaults to true.

**Example Usage:** `-prom=false`

#### `clientApi`
Flag for whether to enable the client management api. See the `Service Discovery` section for more information. Defaults to false.

When enabled, it completely disables refreshing configuration from file.

**IMPORTANT:** Do not expose this port publicly in any circumstance. Ideally it should only be used within scope of localhost communication, or opened up internally with SSL.

**Example Usage:** `-clientApi=true`
#### `stats`
Flag for whether to enable REST stats exporting on the `/stats` endpoint. Defaults to true.

**Example Usage:** `-stats=false`
#### `ws`
Flag for whether to enable WebSocket stats exporting on the `/ws` endpoint. Defaults to true.

**Example Usage:** `-ws=false`
#### `debug`
Flag for more verbose logging. Debugging being enabled will show the metric scraping cycles and every endpoint being called on each node.

### Client Management API
This exporter supports custom service discovery implementations with the client management API. The API, enabled with `-clientApi=true` allows the list of nodes to be updated during service operation without storing any credentials on disk.

For example, a POST to `localhost:8082/clients` with the body:
```json
[
  {
    "url": "https://localhost:6688",
    "username": "chainlink",
    "password": "twochains"
  }
]
```
Will inform the service to scrape a node located at `localhost` with the above credentials. If the url being sent already exists in the service, it is ignored and no change occurs. If a node which is present is not in the request, it will be removed from the service.

Upon posting, the service will check that it can authenticate against the node. If it fails, the node will be ignored.

#### Use Cases
Following examples where the Client Management API is needed:
- Discovery by SRV records.
- Discovery with services like Consul, Zookeeper, etcd.

To enable service discovery, a custom lightweight service should be built that posts to the `/clients` endpoint every so-often or when it recognises a change in nodes.

### Usage
#### Stats
By default, the stats endpoint will be on port 8080 at `/stats`. Example output:
```json
{
    "0x89f94f8D23187f59404737D96D500a1A5186544E": {
        "totalSpecs": 1,
        "totalBridges": 2,
        "address": "0x89f94f8D23187f59404737D96D500a1A5186544E",
        "ethBalance": "2894840440000000000",
        "linkBalance": "0",
        "job_spec_stats": {
            "id": "*",
            "run_count": 165,
                "adaptor_count": {
                "ethtx": 1,
                "ethuint256": 1,
                "httpget": 1,
                "jsonparse": 1,
                "multiply": 1
            },
            "status_count": {
                "completed": 158,
                "errored": 7
            },
            "param_count": {
            "url": [
                {
                "value": "https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD,EUR,JPY",
                "count": 165
                }
            ]}
        }
    }
}
```

#### Prometheus
By default Prometheus metrics will be exposed on port 8080 at `/metrics`, output snippet:
```
# HELP node_eth_balance The Ethereum balance of the nodes wallet.
# TYPE node_eth_balance gauge
node_eth_balance{address="0x89f94f8D23187f59404737D96D500a1A5186544E"} 2.89150742e+18
node_eth_balance{address="0xBe6EDAab06316ca6e3CB0E2F7830304D518A53ce"} 2.8944938e+18
# HELP node_link_balance The LINK balance of the nodes wallet.
# TYPE node_link_balance gauge
node_link_balance{address="0x89f94f8D23187f59404737D96D500a1A5186544E"} 0
node_link_balance{address="0xBe6EDAab06316ca6e3CB0E2F7830304D518A53ce"} 0
```

#### WebSocket
By default the WebSocket endpoint will be exposed on port 8080 at `/ws`. To subscribe to an address for real-time updates, send:
```json
{
  "addresses": [
    "0xBe6EDAab06316ca6e3CB0E2F7830304D518A53ce"
  ]
}
```
Then if the address exists, you should see output that is in the same schema as `/stats`.