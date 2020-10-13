package spec

// Get retrieves a message specification by its name.
func Get(name string) *Node {
	return msgDir[name]
}

// Add registers a new message specification.
func Add(name string, root *Node) {
	msgDir[name] = root
}

var msgDir = map[string]*Node{}
