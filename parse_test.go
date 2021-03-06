package cmsgparser

import "testing"

func TestParser1(t *testing.T) {
	var test1 string = "Added filtering by dossier codes\nhttps://trello.com/c/4x225sdf/182-departure-departure-service-doc-req-field-extension\nhttps://trello.com/c/skumba Move Super List\nOK\n"
	root, err := Parse(test1)
	ok(t, err)
	equals(t, len(root.Children()), 4)
	node := root
	equals(t, RootNode, node.Type)
	node = root.Children()[0]
	equals(t, TextNode, node.Type)
	equals(t, "Added filtering by dossier codes\n", node.Value)
	node = root.Children()[1]
	equals(t, LinkNode, node.Type)
	equals(t, node.Value, "https://trello.com/c/4x225sdf")
	node = root.Children()[2]
	equals(t, LinkNode, node.Type)
	equals(t, node.Value, "https://trello.com/c/skumba")
	equals(t, len(node.Children()), 1)
	equals(t, node.Children()[0].Type, MoveActionNode)
	equals(t, node.Children()[0].Value, "Super List")
	node = root.Children()[3]
	equals(t, TextNode, node.Type)
	equals(t, "OK\n", node.Value)
}

func TestParser2(t *testing.T) {
	var test1 string = "Added filtering by dossier codes\nhttps://trello.com/c/"
	_, err := Parse(test1)
	assert(t, err != nil, "Error should not be nil")
	equals(t, err.Error(), "Parse error [Invalid trello card link] near \"Added filtering by dossier codes\nhttps://trello.com/c/\"")
}

func TestParser3(t *testing.T) {
	var test3 string = "Added filtering by dossier codes\r\nhttps://trello.com/c/4x225sdf/182-departure-departure-service-doc-req-field-extension\r\nhttps://trello.com/c/skumba Move Super List\r\nOK\r\n"
	root, err := Parse(test3)
	ok(t, err)
	equals(t, 4, len(root.Children()))
	node := root
	equals(t, RootNode, node.Type)
	node = root.Children()[0]
	equals(t, TextNode, node.Type)
	equals(t, "Added filtering by dossier codes\r\n", node.Value)
	node = root.Children()[1]
	equals(t, LinkNode, node.Type)
	equals(t, node.Value, "https://trello.com/c/4x225sdf")
	node = root.Children()[2]
	equals(t, LinkNode, node.Type)
	equals(t, node.Value, "https://trello.com/c/skumba")
	equals(t, len(node.Children()), 1)
	equals(t, node.Children()[0].Type, MoveActionNode)
	equals(t, node.Children()[0].Value, "Super List")
	node = root.Children()[3]
	equals(t, TextNode, node.Type)
	equals(t, "OK\r\n", node.Value)
}

func TestParser4(t *testing.T) {
	var test1 string = ""
	root, err := Parse(test1)
	ok(t, err)
	equals(t, len(root.Children()), 0)
}

func TestParser5(t *testing.T) {
	var s string
	root, err := Parse(s)
	ok(t, err)
	equals(t, len(root.Children()), 0)
}

func TestParser6(t *testing.T) {
	var test1 string = "Foo bar https://trello.com/c/foo boom"
	root, err := Parse(test1)
	ok(t, err)
	equals(t, 3, len(root.Children()))
}

func TestParserOnlyUnreservedChars(t *testing.T) {
	var test1 string = "Foo bar (https://trello.com/c/foo) boom"
	root, err := Parse(test1)
	ok(t, err)
	equals(t, 3, len(root.Children()))
	node := root.Children()[1]
	equals(t, LinkNode, node.Type)
	equals(t, "https://trello.com/c/foo", node.Value)
	//this is correct since there is a space here expect
	//the normal text parsing to kick in since it's not an instruction
	node = root.Children()[2]
	equals(t, TextNode, node.Type)
	equals(t, "boom", node.Value)
}
