package hashring

// StringNodeLocator derived from NodeLocator, but explicit specialization with string
type StringNodeLocator struct {
	nl *NodeLocator
}

func NewStringNodeLocator(opts ...NodeLocatorOption) *StringNodeLocator {
	return &StringNodeLocator{nl: New(opts...)}
}

// GetAllNodes returns all available nodes
func (c *StringNodeLocator) GetAllNodes() []string {
	return nodesToStrings(c.nl.GetAllNodes()...)
}

// GetPrimaryNode returns the first available node for a name, such as “127.0.0.1:11311-0” for "Alice"
func (c *StringNodeLocator) GetPrimaryNode(name string) (string, bool) {
	node, has := c.nl.GetPrimaryNode(name)
	return nodeToString(node), has
}

// GetMaxHashKey returns the last available node's HashKey
// that is, Maximum HashKey in the Hash Cycle
func (c *StringNodeLocator) GetMaxHashKey() (uint32, error) {
	return c.nl.GetMaxHashKey()
}

// SetNodes setups the NodeLocator with the list of nodes it should use.
// If there are existing nodes not present in nodes, they will be removed.
// @param nodes a List of Nodes for this NodeLocator to use in
// its continuum
func (c *StringNodeLocator) SetNodes(nodes ...string) {
	c.nl.SetNodes(stringsToNodes(nodes...)...)
}

// RemoveAllNodes removes all nodes in the continuum....
func (c *StringNodeLocator) RemoveAllNodes() {
	c.nl.RemoveAllNodes()
}

// AddNodes inserts nodes into the consistent hash cycle.
func (c *StringNodeLocator) AddNodes(nodes ...string) {
	c.nl.AddNodes(stringsToNodes(nodes...)...)
}

// Remove removes nodes from the consistent hash cycle...
func (c *StringNodeLocator) RemoveNodes(nodes ...string) {
	c.nl.RemoveNodes(stringsToNodes(nodes...)...)
}

// Get returns an element close to where name hashes to in the nodes.
func (c *StringNodeLocator) Get(name string) (string, bool) {
	node, has := c.nl.GetPrimaryNode(name)
	if node == nil {
		return "", has
	}
	return nodeToString(node), has
}

// GetTwo returns the two closest distinct elements to the name input in the nodes.
func (c *StringNodeLocator) GetTwo(name string) (string, string, bool) {
	firstNode, secondNode, has := c.nl.GetTwo(name)
	return nodeToString(firstNode), nodeToString(secondNode), has
}

// GetN returns the N closest distinct elements to the name input in the nodes.
func (c *StringNodeLocator) GetN(name string, n int) ([]string, bool) {
	nodes, has := c.nl.GetN(name, n)
	return nodesToStrings(nodes...), has
}

func stringsToNodes(nodes ...string) []Node {
	var _nodes []Node
	for _, node := range nodes {
		_nodes = append(_nodes, StringNode(node))
	}
	return _nodes
}

func nodesToStrings(nodes ...Node) []string {
	var _nodes []string
	for _, node := range nodes {
		_nodes = append(_nodes, nodeToString(node))
	}
	return _nodes
}

func nodeToString(node Node) string {
	if node == nil {
		return ""
	}
	return string(node.(StringNode))
}
