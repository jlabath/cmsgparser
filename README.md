### Commit message parser

Is the parsing part of a github/trello integration used at my workplace,
meant to enable team of developers to automatically include comments on trello cards.
This new parser also supports multiple cards in the same commit message.
And the ability to move the trello card to a different list.

Example. 

	Series of tasks completed.
	https://trello.com/c/foo
	https://trello.com/c/bar move to Done
	Few more items still remaining.

Parse of this commit message, would result in RootNode,
containing 4 children nodes TextNode, LinkNode, LinkNode, TextNode
with the second link containing additional child MoveActionNode.

Parser inspired by Lexical Scanning in GO talk by Rob Pike
http://cuddle.googlecode.com/hg/talk/lex.html

### Syntax

* Trello card references can be anywhere in the text.
* After a trello card reference a text "move" or "move to" gains special meaning and means that the trello card should be moved to the specified list.

  Example

	https://trello.com/c/bar move to Review and QA

  would move the card with id _bar_ to a list named _Review and QA_.
* The list name is terminated by a newline.
* Simlarly if you want to avoid chance of any special actions simply enter newline after a trello card reference.
* Each commit message can have multiple text and trello card references and order does not matter.
