package stream

type GroupType int

type Pipeline interface {
}

// TreePipeline is a collection of collector and processors which has a tree topology
type TreePipeline struct {
	Name       string
	Collector  *Collector
	Processors []*PipelineNode
}

// SetRoot set the collector of the pipeline
func (t *TreePipeline) SetRoot(collector *Collector) {
	t.Collector = collector
}

func (t *TreePipeline) addNode(node *PipelineNode, parallelism int) {
	t.Processors = append(t.Processors, node)
}

// PipelineNode represents a processor node
type PipelineNode struct {
	ProcessHandler
	children []*PipelineNode
}

func (n *PipelineNode) addNode(node *PipelineNode, parallelism int) {
	n.children = append(n.children, node)
}

// NewTreePipeline creates a TreePipeline
func NewTreePipeline(name string, collector *Collector) *TreePipeline {
	return &TreePipeline{Name: name, Collector: collector}
}
