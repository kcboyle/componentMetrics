package handlers

import (
	"fmt"
	"html/template"
	"metricRevealer/metricnode"
	"net/http"
)

type metricsListingHandler struct {
	nodes *metricnode.NodeContainer
}

func NewListing(nodes *metricnode.NodeContainer) http.Handler {
	return metricsListingHandler{
		nodes: nodes,
	}
}

func (m metricsListingHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "html")

	m.render(res)
}

func (m metricsListingHandler) render(res http.ResponseWriter) {
	//TODO: test this using a bosh deployment instead of lattice (a.b.c.d... present)
	printNodeAsHTML := func(node metricnode.MetricNode) template.HTML {
		td := printNode(node, "")

		return template.HTML(td)
	}

	//TODO: put in html file and use http/template correctly
	//TODO: list all metricTypes for a particular origin(Need another row span *sigh*
	htmlTemplate := `
	<!DOCTYPE html>
	<html>
		<head>
			<style>
				table, th, td {
					border: 2px solid black;
					border-collapse: collapse;
					width: 75%;
					font-size:15pt;
				}
				th, td {
					padding: 10px;
					text-align: center;
				}
				tr:nth-child(even) {
					background-color: #CCCCCC;
				}
			</style>
		</head>
		<body>
			<h1>Loggregator Metrics</h1>
			<table>
				{{range $index, $node := .}}
				{{printNodeAsHTML $node}}
				{{end}}
			</table></body>
	</html>`

	t := template.New("t").Funcs(template.FuncMap{"printNodeAsHTML": printNodeAsHTML})

	t, err := t.Parse(htmlTemplate)
	if err != nil {
		fmt.Println("ERROR: could not parse template: ", err.Error())
		return
	}

	err = t.Execute(res, m.nodes.Nodes)
	if err != nil {
		panic(err)
	}
}

func printNode(node metricnode.MetricNode, td string) string {
	maxNodes := metricnode.MaxNodes(node.ChildNodes)
	td += "<tr>"

	if maxNodes > 1 {
		td += fmt.Sprintf("<td rowspan=%d>%s</td>", maxNodes, node.Value)
	} else {
		td += fmt.Sprintf("<td>%s</td>", node.Value)
	}

	if len(node.ChildNodes) > 1 {
		for _, n := range node.ChildNodes {
			td = printSingleNode(n, td)
		}
		return td + "</tr>"
	}

	if len(node.ChildNodes) == 1 {
		return printSingleNode(node.ChildNodes[0], td) + "</tr>"
	}
	return td + fmt.Sprintf("<td>%s</td></tr>", node.Value)
}

func printSingleNode(node *metricnode.MetricNode, td string) string {
	maxNodes := metricnode.MaxNodes(node.ChildNodes)
	if maxNodes > 1 {
		td = td + fmt.Sprintf("<td rowspan=%d>%s</td>", maxNodes, node.Value)
	} else {
		td = td + fmt.Sprintf("<td>%s</td>", node.Value)
	}

	if len(node.ChildNodes) > 1 {
		for _, node := range node.ChildNodes {
			td = printSingleNode(node, td)
		}
		return td
	}

	if len(node.ChildNodes) == 1 {
		return printSingleNode(node.ChildNodes[0], td)
	}

	return td
}