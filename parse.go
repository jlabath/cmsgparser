package cmsgparser

import "fmt"

//Parse parses the input commit message and returns the root node
//of the parse tree for the commit message
func Parse(input string) (*Node, error) {
	var err error
	var last *Node
	root := &Node{
		Type: RootNode}
	lexer := lex("lex", input)
	//since we are so simple just stick it in one loop
	for done := 0; done == 0; {
		item := lexer.nextItem()
		switch item.typ {
		case itemError:
			err = fmt.Errorf("Parse error [%s] near \"%s\"", item.val, input[0:lexer.pos])
			done = 1
		case itemEOF:
			done = 1 //ok end
		case itemText:
			newn := &Node{Type: TextNode, Value: item.val}
			root.AddChild(newn)
			last = newn
		case itemLink:
			newn := &Node{Type: LinkNode, Value: item.val}
			root.AddChild(newn)
			last = newn
		case itemMoveDestination:
			newn := &Node{Type: MoveActionNode, Value: item.val}
			//todo if last not a link error and end
			last.AddChild(newn)
			last = newn
		default:
			panic("missing a case to handle")
		}
	}
	return root, err
}
