package dot_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	dot "github.com/krakend/krakend-config2dot/v3"
	"github.com/luraproject/lura/v3/config"
)

func printDot(s string) {
	lines := strings.Split(s, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimRight(line, " \t")
		if line == "" && len(out) > 0 && out[len(out)-1] == "" {
			continue
		}
		out = append(out, line)
	}
	fmt.Println(strings.Join(out, "\n"))
}

func render(cfg config.ServiceConfig) {
	var buf bytes.Buffer
	if _, err := dot.WriteDot(&buf, cfg); err != nil {
		fmt.Println("error:", err)
		return
	}
	printDot(buf.String())
}

func ExampleWriteDot() {
	render(config.ServiceConfig{
		Port: 8080,
		Endpoints: []*config.EndpointConfig{
			{
				Endpoint:       "/users",
				Method:         "GET",
				Timeout:        2 * time.Second,
				CacheTTL:       300 * time.Second,
				OutputEncoding: "json",
				QueryString:    []string{"page"},
				Backend: []*config.Backend{
					{
						URLPattern:      "/v1/users",
						Host:            []string{"http://users.example.com"},
						ConcurrentCalls: 1,
					},
				},
			},
		},
	})
	// Output:
	// digraph krakend {
	//     label="KrakenD Gateway";
	//     labeljust="l";
	//     fontname="Ubuntu";
	//     fontsize="13";
	//     rankdir="LR";
	//     bgcolor="aliceblue";
	//     style="solid";
	//     penwidth="0.5";
	//     pad="0.0";
	//     nodesep="0.35";
	//
	//     node [shape="ellipse" style="filled" fillcolor="honeydew" fontname="Ubuntu" penwidth="1.0" margin="0.05,0.0"];
	//
	//     subgraph "cluster_/users" {
	//     	label="/users";
	//     	bgcolor="lightgray";
	//     	shape="box";
	//     	style="solid";
	//
	//         "/users" [ shape=record, label="{ { Timeout | 2s } | { CacheTTL | 5m0s } | { Output | json } | { QueryString | [page] } }" ]
	//
	// 	    subgraph "cluster_/v1/users" {
	// 	    	label="/v1/users";
	// 	    	bgcolor="beige";
	// 	    	shape="box";
	// 	    	style="solid";
	//         	"in_0_0" [ shape=record, label="{ {sd|static } | { Hosts | [http://users.example.com] } | { Encoding | JSON } }" ]
	//
	// 	    }
	//
	// 	    "/users" -> in_0_0 [ label="x1"]
	//     }
	//
	//     ":8080" -> "/users" [ label="GET"]
	// }
}

func ExampleWriteDot_multipleBackends() {
	render(config.ServiceConfig{
		Port: 8080,
		Endpoints: []*config.EndpointConfig{
			{
				Endpoint:       "/profile",
				Method:         "GET",
				Timeout:        time.Second,
				OutputEncoding: "json",
				Backend: []*config.Backend{
					{
						URLPattern:      "/v1/users",
						Host:            []string{"http://users.example.com"},
						ConcurrentCalls: 2,
					},
					{
						URLPattern:      "/v1/orders",
						Host:            []string{"http://orders.example.com"},
						Encoding:        "json",
						ConcurrentCalls: 1,
					},
				},
			},
		},
	})
	// Output:
	// digraph krakend {
	//     label="KrakenD Gateway";
	//     labeljust="l";
	//     fontname="Ubuntu";
	//     fontsize="13";
	//     rankdir="LR";
	//     bgcolor="aliceblue";
	//     style="solid";
	//     penwidth="0.5";
	//     pad="0.0";
	//     nodesep="0.35";
	//
	//     node [shape="ellipse" style="filled" fillcolor="honeydew" fontname="Ubuntu" penwidth="1.0" margin="0.05,0.0"];
	//
	//     subgraph "cluster_/profile" {
	//     	label="/profile";
	//     	bgcolor="lightgray";
	//     	shape="box";
	//     	style="solid";
	//
	//         "/profile" [ shape=record, label="{ { Timeout | 1s } | { CacheTTL | 0s } | { Output | json } | { QueryString | [] } }" ]
	//
	// 	    subgraph "cluster_/v1/users" {
	// 	    	label="/v1/users";
	// 	    	bgcolor="beige";
	// 	    	shape="box";
	// 	    	style="solid";
	//         	"in_0_0" [ shape=record, label="{ {sd|static } | { Hosts | [http://users.example.com] } | { Encoding | JSON } }" ]
	//
	// 	    }
	//
	// 	    "/profile" -> in_0_0 [ label="x2"]
	// 	    subgraph "cluster_/v1/orders" {
	// 	    	label="/v1/orders";
	// 	    	bgcolor="beige";
	// 	    	shape="box";
	// 	    	style="solid";
	//         	"in_0_1" [ shape=record, label="{ {sd|static } | { Hosts | [http://orders.example.com] } | { Encoding | json } }" ]
	//
	// 	    }
	//
	// 	    "/profile" -> in_0_1 [ label="x1"]
	//     }
	//
	//     ":8080" -> "/profile" [ label="GET"]
	// }
}

func ExampleWriteDot_extraConfig() {
	render(config.ServiceConfig{
		Port: 8080,
		Endpoints: []*config.EndpointConfig{
			{
				Endpoint:       "/protected",
				Method:         "GET",
				Timeout:        time.Second,
				OutputEncoding: "json",
				ExtraConfig: config.ExtraConfig{
					"auth/validator": map[string]any{
						"alg":     "RS256",
						"jwk_url": "https://issuer.example.com/jwks",
					},
					"qos/ratelimit/router": map[string]any{
						"max_rate": 100,
					},
				},
				Backend: []*config.Backend{
					{
						URLPattern:      "/v1/protected",
						Host:            []string{"http://backend.example.com"},
						ConcurrentCalls: 1,
						ExtraConfig: config.ExtraConfig{
							"qos/circuit-breaker": map[string]any{
								"interval": 60,
							},
						},
					},
				},
			},
		},
	})
	// Output:
	// digraph krakend {
	//     label="KrakenD Gateway";
	//     labeljust="l";
	//     fontname="Ubuntu";
	//     fontsize="13";
	//     rankdir="LR";
	//     bgcolor="aliceblue";
	//     style="solid";
	//     penwidth="0.5";
	//     pad="0.0";
	//     nodesep="0.35";
	//
	//     node [shape="ellipse" style="filled" fillcolor="honeydew" fontname="Ubuntu" penwidth="1.0" margin="0.05,0.0"];
	//
	//     subgraph "cluster_/protected" {
	//     	label="/protected";
	//     	bgcolor="lightgray";
	//     	shape="box";
	//     	style="solid";
	//
	//         "/protected" [ shape=record, label="{ { Timeout | 1s } | { CacheTTL | 0s } | { Output | json } | { QueryString | [] } }" ]
	//         "extra_0" [ shape=record, label="{ {ExtraConfig}  | { auth/validator | { alg | RS256 } | { jwk_url | https://issuer.example.com/jwks }  } | { qos/ratelimit/router | { max_rate | 100 }  } }" ]
	//
	// 	    subgraph "cluster_/v1/protected" {
	// 	    	label="/v1/protected";
	// 	    	bgcolor="beige";
	// 	    	shape="box";
	// 	    	style="solid";
	//         	"in_0_0" [ shape=record, label="{ {sd|static } | { Hosts | [http://backend.example.com] } | { Encoding | JSON } }" ]
	//         "extra_0_0" [ shape=record, label="{ { ExtraConfig  | qos/circuit-breaker  } }" ]
	// 	    }
	//
	// 	    "/protected" -> in_0_0 [ label="x1"]
	//     }
	//
	//     ":8080" -> "/protected" [ label="GET"]
	// }
}

func ExampleWriteDot_serviceDiscovery() {
	render(config.ServiceConfig{
		Port: 8080,
		Endpoints: []*config.EndpointConfig{
			{
				Endpoint:       "/catalog",
				Method:         "GET",
				Timeout:        time.Second,
				OutputEncoding: "json",
				Backend: []*config.Backend{
					{
						URLPattern:      "/catalog",
						Host:            []string{"catalog.service.consul"},
						SD:              "dns",
						Encoding:        "xml",
						ConcurrentCalls: 1,
					},
				},
			},
		},
	})
	// Output:
	// digraph krakend {
	//     label="KrakenD Gateway";
	//     labeljust="l";
	//     fontname="Ubuntu";
	//     fontsize="13";
	//     rankdir="LR";
	//     bgcolor="aliceblue";
	//     style="solid";
	//     penwidth="0.5";
	//     pad="0.0";
	//     nodesep="0.35";
	//
	//     node [shape="ellipse" style="filled" fillcolor="honeydew" fontname="Ubuntu" penwidth="1.0" margin="0.05,0.0"];
	//
	//     subgraph "cluster_/catalog" {
	//     	label="/catalog";
	//     	bgcolor="lightgray";
	//     	shape="box";
	//     	style="solid";
	//
	//         "/catalog" [ shape=record, label="{ { Timeout | 1s } | { CacheTTL | 0s } | { Output | json } | { QueryString | [] } }" ]
	//
	// 	    subgraph "cluster_/catalog" {
	// 	    	label="/catalog";
	// 	    	bgcolor="beige";
	// 	    	shape="box";
	// 	    	style="solid";
	//         	"in_0_0" [ shape=record, label="{ {sd|dns } | { Hosts | [catalog.service.consul] } | { Encoding | xml } }" ]
	//
	// 	    }
	//
	// 	    "/catalog" -> in_0_0 [ label="x1"]
	//     }
	//
	//     ":8080" -> "/catalog" [ label="GET"]
	// }
}

func ExampleWriteDot_multipleEndpoints() {
	render(config.ServiceConfig{
		Port: 8080,
		Endpoints: []*config.EndpointConfig{
			{
				Endpoint:       "/users",
				Method:         "GET",
				Timeout:        time.Second,
				OutputEncoding: "json",
				Backend: []*config.Backend{
					{
						URLPattern:      "/v1/users",
						Host:            []string{"http://users.example.com"},
						ConcurrentCalls: 1,
					},
				},
			},
			{
				Endpoint:       "/users",
				Method:         "POST",
				Timeout:        3 * time.Second,
				OutputEncoding: "no-op",
				Backend: []*config.Backend{
					{
						URLPattern:      "/v1/users",
						Host:            []string{"http://users.example.com"},
						ConcurrentCalls: 1,
					},
				},
			},
		},
	})
	// Output:
	// digraph krakend {
	//     label="KrakenD Gateway";
	//     labeljust="l";
	//     fontname="Ubuntu";
	//     fontsize="13";
	//     rankdir="LR";
	//     bgcolor="aliceblue";
	//     style="solid";
	//     penwidth="0.5";
	//     pad="0.0";
	//     nodesep="0.35";
	//
	//     node [shape="ellipse" style="filled" fillcolor="honeydew" fontname="Ubuntu" penwidth="1.0" margin="0.05,0.0"];
	//
	//     subgraph "cluster_/users" {
	//     	label="/users";
	//     	bgcolor="lightgray";
	//     	shape="box";
	//     	style="solid";
	//
	//         "/users" [ shape=record, label="{ { Timeout | 1s } | { CacheTTL | 0s } | { Output | json } | { QueryString | [] } }" ]
	//
	// 	    subgraph "cluster_/v1/users" {
	// 	    	label="/v1/users";
	// 	    	bgcolor="beige";
	// 	    	shape="box";
	// 	    	style="solid";
	//         	"in_0_0" [ shape=record, label="{ {sd|static } | { Hosts | [http://users.example.com] } | { Encoding | JSON } }" ]
	//
	// 	    }
	//
	// 	    "/users" -> in_0_0 [ label="x1"]
	//     }
	//
	//     subgraph "cluster_/users" {
	//     	label="/users";
	//     	bgcolor="lightgray";
	//     	shape="box";
	//     	style="solid";
	//
	//         "/users" [ shape=record, label="{ { Timeout | 3s } | { CacheTTL | 0s } | { Output | no-op } | { QueryString | [] } }" ]
	//
	// 	    subgraph "cluster_/v1/users" {
	// 	    	label="/v1/users";
	// 	    	bgcolor="beige";
	// 	    	shape="box";
	// 	    	style="solid";
	//         	"in_1_0" [ shape=record, label="{ {sd|static } | { Hosts | [http://users.example.com] } | { Encoding | JSON } }" ]
	//
	// 	    }
	//
	// 	    "/users" -> in_1_0 [ label="x1"]
	//     }
	//
	//     ":8080" -> "/users" [ label="GET"]
	//     ":8080" -> "/users" [ label="POST"]
	// }
}

func ExampleWriteDot_noEndpoints() {
	render(config.ServiceConfig{Port: 8080})
	// Output:
	// digraph krakend {
	//     label="KrakenD Gateway";
	//     labeljust="l";
	//     fontname="Ubuntu";
	//     fontsize="13";
	//     rankdir="LR";
	//     bgcolor="aliceblue";
	//     style="solid";
	//     penwidth="0.5";
	//     pad="0.0";
	//     nodesep="0.35";
	//
	//     node [shape="ellipse" style="filled" fillcolor="honeydew" fontname="Ubuntu" penwidth="1.0" margin="0.05,0.0"];
	//
	// }
}

func ExampleServiceConfig_WriteTo() {
	cfg := config.ServiceConfig{
		Port: 8080,
		Endpoints: []*config.EndpointConfig{
			{
				Endpoint:       "/users",
				Method:         "GET",
				Timeout:        2 * time.Second,
				CacheTTL:       300 * time.Second,
				OutputEncoding: "json",
				QueryString:    []string{"page"},
				Backend: []*config.Backend{
					{
						URLPattern:      "/v1/users",
						Host:            []string{"http://users.example.com"},
						ConcurrentCalls: 1,
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if _, err := dot.ServiceConfig(cfg).WriteTo(&buf); err != nil {
		fmt.Println("error:", err)
		return
	}
	printDot(buf.String())
	// Output:
	// digraph krakend {
	//     label="KrakenD Gateway";
	//     labeljust="l";
	//     fontname="Ubuntu";
	//     fontsize="13";
	//     rankdir="LR";
	//     bgcolor="aliceblue";
	//     style="solid";
	//     penwidth="0.5";
	//     pad="0.0";
	//     nodesep="0.35";
	//
	//     node [shape="ellipse" style="filled" fillcolor="honeydew" fontname="Ubuntu" penwidth="1.0" margin="0.05,0.0"];
	//
	//     subgraph "cluster_/users" {
	//     	label="/users";
	//     	bgcolor="lightgray";
	//     	shape="box";
	//     	style="solid";
	//
	//         "/users" [ shape=record, label="{ { Timeout | 2s } | { CacheTTL | 5m0s } | { Output | json } | { QueryString | [page] } }" ]
	//
	// 	    subgraph "cluster_/v1/users" {
	// 	    	label="/v1/users";
	// 	    	bgcolor="beige";
	// 	    	shape="box";
	// 	    	style="solid";
	//         	"in_0_0" [ shape=record, label="{ {sd|static } | { Hosts | [http://users.example.com] } | { Encoding | JSON } }" ]
	//
	// 	    }
	//
	// 	    "/users" -> in_0_0 [ label="x1"]
	//     }
	//
	//     ":8080" -> "/users" [ label="GET"]
	// }
}

func TestWriteDot_templateError(t *testing.T) {
	cfg := config.ServiceConfig{
		Port: 8080,
		Endpoints: []*config.EndpointConfig{
			{
				Endpoint:    "/users",
				Method:      "GET",
				ExtraConfig: config.ExtraConfig{"broken": "not-an-object"},
			},
		},
	}

	n, err := dot.WriteDot(io.Discard, cfg)
	if err == nil {
		t.Fatal("expected an error, got none")
	}
	if n != 0 {
		t.Errorf("expected no bytes written, got %d", n)
	}
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) { return 0, errors.New("write failed") }

func TestWriteDot_writeCount(t *testing.T) {
	cfg := config.ServiceConfig{
		Port: 8080,
		Endpoints: []*config.EndpointConfig{
			{
				Endpoint: "/users",
				Method:   "GET",
				Backend: []*config.Backend{
					{URLPattern: "/v1/users", Host: []string{"http://users.example.com"}},
				},
			},
		},
	}

	var buf bytes.Buffer
	n, err := dot.WriteDot(&buf, cfg)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if n != int64(buf.Len()) {
		t.Errorf("expected %d bytes written, got %d", buf.Len(), n)
	}

	if _, err := dot.WriteDot(errWriter{}, cfg); err == nil {
		t.Error("expected the writer error to be reported, got none")
	}
}
