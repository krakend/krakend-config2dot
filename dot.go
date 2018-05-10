package dot

import (
	"bytes"
	"io"
	"text/template"

	"github.com/devopsfaith/krakend/config"
)

type ServiceConfig config.ServiceConfig

func (s ServiceConfig) WriteTo(w io.Writer) (n int64, err error) {
	return WriteDot(w, config.ServiceConfig(s))
}

func WriteDot(w io.Writer, cfg config.ServiceConfig) (n int64, err error) {
	t := template.New("dot")
	var buf bytes.Buffer
	if err = template.Must(t.Parse(tmplGraph)).Execute(&buf, cfg); err != nil {
		return 0, err
	}
	return buf.WriteTo(w)
}

const tmplGraph = `digraph krakend { {{ $port := .Port }}
    label="KrakenD Gateway";
    labeljust="l";
    fontname="Ubuntu";
    fontsize="13";
    rankdir="LR";
    bgcolor="aliceblue";
    style="solid";
    penwidth="0.5";
    pad="0.0";
    nodesep="0.35";

    node [shape="ellipse" style="filled" fillcolor="honeydew" fontname="Ubuntu" penwidth="1.0" margin="0.05,0.0"];

	{{ range $i, $endpoint := .Endpoints }}
    {{printf "subgraph \"cluster_%s\" {" .Endpoint }}
    	label="{{ .Endpoint }}";
    	bgcolor="lightgray";
    	shape="box";
    	style="solid";

        "{{ .Endpoint }}" [ shape=record, label="{ { Timeout | {{.Timeout.String}} } | { CacheTTL | {{.CacheTTL.String}} } | { Output | {{.OutputEncoding}} } | { QueryString | {{.QueryString}} } }" ]
        {{ if .ExtraConfig }}"extra_{{$i}}" [ shape=record, label="{ {ExtraConfig} {{ range $key, $value := .ExtraConfig }} | { {{ $key }} {{ range $k, $v := $value }}| { {{$k}} | {{$v}} } {{ end }} }{{ end }} }" ]{{ end }}
	    {{ range $j, $backend := .Backend }}
	    {{printf "subgraph \"cluster_%s\" {" .URLPattern }}
	    	label="{{ .URLPattern }}";
	    	bgcolor="beige";
	    	shape="box";
	    	style="solid";
        	"in_{{$i}}_{{$j}}" [ shape=record, label="{ {sd|{{ if .SD }}{{ .SD }}{{ else }}static{{ end }} } | { Hosts | {{.Host}} } | { Encoding | {{ if .Encoding }}{{ .Encoding }}{{ else }}JSON{{ end }} } }" ]
        {{ if .ExtraConfig }}"extra_{{$i}}_{{$j}}" [ shape=record, label="{ { ExtraConfig {{ range $key, $v := .ExtraConfig }} | {{ $key }} {{ end }} } }" ]{{ end }}
	    {{println "}" }}
	    "{{ $endpoint.Endpoint }}" -> in_{{$i}}_{{$j}} [ label="x{{ .ConcurrentCalls }}"]{{ end }}
    {{ println "}" }}{{ end }}
    {{ range .Endpoints }}
    ":{{ $port }}" -> "{{ .Endpoint }}" [ label="{{ .Method }}"]{{ end }}
}
`
