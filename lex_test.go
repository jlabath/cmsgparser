package cmsgparser

import "testing"

func TestLex1(t *testing.T) {
	var test1 string = "Added filtering by dossier codes\nhttps://trello.com/c/4x225sdf/182-departure-departure-service-doc-req-field-extension"
	lexptr := lex("test l1", test1)
	item := lexptr.nextItem()
	equals(t, itemText, item.typ)
	equals(t, "Added filtering by dossier codes\n", item.val)
	item = lexptr.nextItem()
	equals(t, itemLink, item.typ)
	equals(t, "https://trello.com/c/4x225sdf", item.val)
	item = lexptr.nextItem()
	equals(t, itemEOF, item.typ)
}

func TestLex2(t *testing.T) {
	var test1 string = "Added filtering by dossier codes\nhttps://trello.com/c/4x225sdf/182-departure-departure-service-doc-req-field-extension move to Done"
	lexptr := lex("test l2", test1)
	item := lexptr.nextItem()
	equals(t, itemText, item.typ)
	equals(t, "Added filtering by dossier codes\n", item.val)
	item = lexptr.nextItem()
	equals(t, itemLink, item.typ)
	equals(t, "https://trello.com/c/4x225sdf", item.val)
	item = lexptr.nextItem()
	equals(t, itemMoveDestination, item.typ)
	equals(t, "Done", item.val)
}

func TestLex3(t *testing.T) {
	var test1 string = "Added filtering by dossier codes\nhttps://trello.com/c/4x225sdf/182-departure-departure-service-doc-req-field-extension\nOK\n"
	lexptr := lex("test l3", test1)
	item := lexptr.nextItem()
	equals(t, itemText, item.typ)
	equals(t, "Added filtering by dossier codes\n", item.val)
	item = lexptr.nextItem()
	equals(t, itemLink, item.typ)
	equals(t, "https://trello.com/c/4x225sdf", item.val)
	item = lexptr.nextItem()
	equals(t, itemText, item.typ)
	equals(t, "OK\n", item.val)
}

func TestLex4(t *testing.T) {
	var test1 string = "Added filtering by dossier codes\nhttps://trello.com/c/4x225sdf/182-departure-departure-service-doc-req-field-extension\nhttps://trello.com/c/skumba move Super List\nOK\n"
	lexptr := lex("test l4", test1)
	item := lexptr.nextItem()
	equals(t, itemText, item.typ)
	equals(t, "Added filtering by dossier codes\n", item.val)
	item = lexptr.nextItem()
	equals(t, itemLink, item.typ)
	equals(t, "https://trello.com/c/4x225sdf", item.val)
	item = lexptr.nextItem()
	equals(t, itemLink, item.typ)
	equals(t, "https://trello.com/c/skumba", item.val)
	item = lexptr.nextItem()
	equals(t, itemMoveDestination, item.typ)
	equals(t, "Super List", item.val)
	item = lexptr.nextItem()
	equals(t, itemText, item.typ)
	equals(t, "OK\n", item.val)
}

func TestLex5(t *testing.T) {
	var test1 string = "Added filtering by dossier codes\nhttps://trello.com/c/4x225sdf/182-departure-departure-service-doc-req-field-extension\nhttps://trello.com/c/skumba Move Super List\nOK\n"
	lexptr := lex("test l5", test1)
	item := lexptr.nextItem()
	equals(t, itemText, item.typ)
	equals(t, "Added filtering by dossier codes\n", item.val)
	item = lexptr.nextItem()
	equals(t, itemLink, item.typ)
	equals(t, "https://trello.com/c/4x225sdf", item.val)
	item = lexptr.nextItem()
	equals(t, itemLink, item.typ)
	equals(t, "https://trello.com/c/skumba", item.val)
	item = lexptr.nextItem()
	equals(t, itemMoveDestination, item.typ)
	equals(t, "Super List", item.val)
	item = lexptr.nextItem()
	equals(t, itemText, item.typ)
	equals(t, "OK\n", item.val)
	item = lexptr.nextItem()
	equals(t, itemEOF, item.typ)
}

func TestLex6(t *testing.T) {
	var test1 string = "Added filtering by dossier codes\nhttps://trello.com/c/"
	lexptr := lex("test l6", test1)
	item := lexptr.nextItem()
	equals(t, itemText, item.typ)
	equals(t, "Added filtering by dossier codes\n", item.val)
	item = lexptr.nextItem()
	equals(t, itemError, item.typ)
	equals(t, "Invalid trello card link", item.val)
	//everything after this should be EOF since channel is closed
	item = lexptr.nextItem()
	equals(t, itemEOF, item.typ)
	item = lexptr.nextItem()
	equals(t, itemEOF, item.typ)
}

func TestLex7(t *testing.T) {
	var test1 string = ""
	lexptr := lex("test l7", test1)
	item := lexptr.nextItem()
	equals(t, itemEOF, item.typ)
	item = lexptr.nextItem()
	equals(t, itemEOF, item.typ)
	lexptr = lex("test l7a", "    	")
	item = lexptr.nextItem()
	equals(t, itemText, item.typ)
	equals(t, "    	", item.val)
	item = lexptr.nextItem()
	equals(t, itemEOF, item.typ)
}
