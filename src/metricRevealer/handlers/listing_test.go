package handlers_test

import (
	"metricRevealer/handlers"

	"metricRevealer/metricnode"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Listing", func() {
	var (
		req      *http.Request
		recorder *httptest.ResponseRecorder
		handler  http.Handler
		nodes    *metricnode.NodeContainer
	)

	BeforeEach(func() {
		recorder = httptest.NewRecorder()
		nodes = metricnode.NewNodeContainer()
		handler = handlers.NewListing(nodes)
		req = &http.Request{}
	})

	It("returns valid boilerplate html with no nodes", func() {
		handler.ServeHTTP(recorder, req)

		body := recorder.Body.String()

		Expect(body).ToNot(ContainSubstring("<td>"))
	})

	It("returns html with populated node data", func() {
		node1 := metricnode.NewMetricNode("this.is.metricOne")
		node2 := metricnode.NewMetricNode("this.is.metricTwo")
		nodes.Store(node1)
		nodes.Store(node2)

		generated := "<td>this</td><td>is</td><td><table border=0><tr><td>metricOne</td></tr><tr><td>metricTwo</td></tr></table></td>"

		handler.ServeHTTP(recorder, req)

		Expect(recorder.Body.String()).To(ContainSubstring(generated))
	})

	FIt("returns html with interactive rowspans", func() {
		node1 := metricnode.NewMetricNode("this.is.metricOne")
		node3 := metricnode.NewMetricNode("this.is.metricOne.metricThree")
		node4 := metricnode.NewMetricNode("this.is.metricOne.metricFour")
		node2 := metricnode.NewMetricNode("this.is.metricTwo")
		node5 := metricnode.NewMetricNode("this.is.metricTwo.metricFive")
		node6 := metricnode.NewMetricNode("now.be.aNewMetric")

		nodes.Store(node1)
		nodes.Store(node2)
		nodes.Store(node3)
		nodes.Store(node4)
		nodes.Store(node5)
		nodes.Store(node6)

		handler.ServeHTTP(recorder, req)

		generated := `
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
				    <tr>
				        <td rowspan="3">this</td>
				        <td rowspan="3">is</td>
				        <td rowspan="2">metricone</td>
				        <td>metricthree</td>
				    </tr>
				    <tr>
				        <td>metricfour</td>
				    </tr>
				    <tr>
				        <td>metrictwo</td>
				        <td>metricfive</td>
				    </tr>
				    <tr>
				        <td>now</td>
				        <td>is</td>
				        <td>newmetric</td>
				    </tr>
			</table>
		</body>
	</html>`

		Expect(recorder.Body.String()).To(ContainSubstring(generated))
	})
})
