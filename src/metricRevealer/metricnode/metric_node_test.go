package metricnode_test

import (
	"metricRevealer/metricnode"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MetricNode", func() {
	It("parses a metric into child nodes ", func() {
		str := metricnode.NewMetricNode("period.separated.string")
		Expect(str.Value).To(Equal("period"))
		Expect(str.ChildNodes[0].Value).To(Equal("separated"))
		Expect(str.ChildNodes[0].ChildNodes[0].Value).To(Equal("string"))
	})

	It("stores duplicates as child nodes", func() {
		nodeContainer := metricnode.NewNodeContainer()
		node1 := metricnode.NewMetricNode("this.is.metricOne")
		node2 := metricnode.NewMetricNode("this.is.metricTwo")
		nodeContainer.Store(node1)
		nodeContainer.Store(node2)

		Expect(nodeContainer.Nodes).To(HaveLen(1))                             //this
		Expect(nodeContainer.Nodes[0].ChildNodes).To(HaveLen(1))               //is
		Expect(nodeContainer.Nodes[0].ChildNodes[0].ChildNodes).To(HaveLen(2)) //metricOne, metricTwo
	})

	It("does not repeat the same metric", func() {
		nodeContainer := metricnode.NewNodeContainer()
		node1 := metricnode.NewMetricNode("this.is.metricOne")
		node2 := metricnode.NewMetricNode("this.is.metricOne")
		nodeContainer.Store(node1)
		nodeContainer.Store(node2)

		Expect(nodeContainer.Nodes).To(HaveLen(1))                             //this
		Expect(nodeContainer.Nodes[0].ChildNodes).To(HaveLen(1))               //is
		Expect(nodeContainer.Nodes[0].ChildNodes[0].ChildNodes).To(HaveLen(1)) //metricOne, metricTwo
	})

	It("returns the total number of child nodes for rowspan", func() {
		nodeContainer := metricnode.NewNodeContainer()
		node1 := metricnode.NewMetricNode("this.is.metricOne")
		node3 := metricnode.NewMetricNode("this.is.metricOne.metricThree")
		node4 := metricnode.NewMetricNode("this.is.metricOne.metricFour")
		node2 := metricnode.NewMetricNode("this.is.metricTwo")
		node5 := metricnode.NewMetricNode("this.is.metricTwo.metricFive")
		node6 := metricnode.NewMetricNode("this.is.aNewMetric")

		nodeContainer.Store(node1)
		nodeContainer.Store(node2)
		nodeContainer.Store(node3)
		nodeContainer.Store(node4)
		nodeContainer.Store(node5)
		nodeContainer.Store(node6)

		Expect(nodeContainer.MaxChildNodes()).To(Equal(4))
	})

	It("stores unique nodes and childNodes in the appropriate structures", func() {
		nodeContainer := metricnode.NewNodeContainer()
		node1 := metricnode.NewMetricNode("this.is.metricOne")
		node2 := metricnode.NewMetricNode("this.is.metricTwo")
		nodeContainer.Store(node1)
		nodeContainer.Store(node2)

		node3 := metricnode.NewMetricNode("this.now.metricThree")
		node4 := metricnode.NewMetricNode("this.isnot.metricThree")
		nodeContainer.Store(node3)
		nodeContainer.Store(node4)

		Expect(nodeContainer.Nodes).To(HaveLen(1))                             //this
		Expect(nodeContainer.Nodes[0].ChildNodes).To(HaveLen(3))               //is, now, isnot
		Expect(nodeContainer.Nodes[0].ChildNodes[0].ChildNodes).To(HaveLen(2)) //metricOne, metricTwo
		Expect(nodeContainer.Nodes[0].ChildNodes[1].ChildNodes).To(HaveLen(1)) //metricThree
		Expect(nodeContainer.Nodes[0].ChildNodes[2].ChildNodes).To(HaveLen(1)) //metricThree
	})
})
