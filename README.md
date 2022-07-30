# Hazelcast Go Client Benchmark

**Requirements**

* Go 1.18

## Build

    $ git clone https://github.com/yuce/hazelcast-go-client-benchmark.git
    $ go build .

## Usage

Run `./hazelcast-go-client-benchmark` with a configuration file, e.g.,:

    $ ./hazelcast-go-client-benchmark configs/128x1024_500keys_nc.json

## Configuration

Here's a sample JSON configuration:

      "Client": {
        "Cluster": {
          "Network": {
            "Addresses": [
              "localhost:5701"
            ]
          }
        },
        "NearCaches": [
          {
            "Name": "mymap*",
            "InMemoryFormat": "binary",
            "Eviction": {
              "Policy": "lru",
              "Size": 500
            }
          }
        ]
      },
      "MapName": "mymap1",
      "KeyCount": 500,
      "Repeat": 100,
      "GoroutineCount": 1,
      "EntryGenerator": "sized128x1024",
      "Warmup": false
    }

Check out https://pkg.go.dev/github.com/hazelcast/hazelcast-go-client#hdr-Configuration for the `Client` configuration.
The program configuration is in `config.go` in this project.
The config fields can be mapped to JSON fields straightforwardly.
The following program configuration keys are supported:

* `MapName`: Map name to be used during the benchmark.
* `KeyCount`: Number of unique keys to be used.
* `Repeat`: Number of `Map.Get`s to do (after the warmup, if configured)
* `GoroutineCount`: Number of goroutines to do `Map.Get`s. 
* `Warmup`: If `true`, runs populates the Near Cache before starting the benchmark. 
* `Client`: Hazelcast client configuration. Check out https://pkg.go.dev/github.com/hazelcast/hazelcast-go-client#hdr-Configuration. 
* `EntryGenerator`: Entry generator to use for generating keys and values. See the section below.

## Entry Generator

Checkout `entry_generators.go` for sample entry generators and how to add yours.
Here are the builtin ones:

* `identity`: Both keys and values are the `int64` index. 
* `sized128x1024`: String key of size 128 bytes, value of size 1024 bytes.
* `sized128x4096`: String key of size 128 bytes, value of size 4096 bytes.

## License

See: [LICENSE](LICENSE)
