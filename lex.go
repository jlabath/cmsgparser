package cmsgparser

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// ideas taken from Lexical Scanning in GO talk by Rob Pike
// http://cuddle.googlecode.com/hg/talk/lex.html

// item represents a token returned from the scanner.
type item struct {
	typ itemType // Type, such as itemNumber.
	val string   // Value, such as "23.2".
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	if len(i.val) > 10 {
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

// itemType identifies the type of lex items.
type itemType int

const (
	itemEOF itemType = iota // eof good to have EOF first as that is what will be sent on closed channel
	// value is text of error
	itemError           //error occurred
	itemText            // the plain text of commit message
	itemLink            // the trello's card link
	itemMoveDestination // the destination tab to move trello's card to
)

const (
	trelloLink string = "https://trello.com/c/"
	eof        rune   = -1
	conMoveTo  string = "move to "
	conMove    string = "move "
)

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*lexer) stateFn

// lexer holds the state of the scanner.
type lexer struct {
	name  string    // used only for error reports.
	input string    // the string being scanned.
	start int       // start position of this item.
	pos   int       // current position in the input.
	width int       // width of last rune read from input.
	items chan item // channel of scanned items.
}

// run lexes the input by executing state functions until
// the state is nil.
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// next returns the next rune in the input.
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width =
		utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// peek returns the next rune without adjusting any offsets
func (l *lexer) peek() (r rune) {
	if l.pos >= len(l.input) {
		return eof
	}
	r, _ = utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

//skip advances pos by count and discards any input
func (l *lexer) skip(count int) {
	for i := 0; i < count; i++ {
		l.next()
	}
	l.ignore()
}

// backup steps back one rune.
// Can be called only once per call of next,
// due to the fact that we have stored only the width of latest
// rune
func (l *lexer) backup() {
	l.pos -= l.width
}

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

// nextItem returns the next item from the input.
func (l *lexer) nextItem() item {
	item := <-l.items
	return item
}

//lexText is for lexing plain text in the commit message
func lexText(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], trelloLink) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexLink // Next state.
		}
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF) // Useful to make EOF a token.
	return nil      // Stop the run loop.
}

//lexLink is for lexing once we ancountered trello link
func lexLink(l *lexer) stateFn {
	l.pos += len(trelloLink)
	cardIdStart := l.pos
	for {
		switch r := l.next(); {
		case r == eof || r == '\n' || unicode.IsSpace(r) || r == '/':
			if cardIdStart < l.pos-1 {
				l.backup()
				l.emit(itemLink)
				return lexAfterLink
			} else {
				return l.errorf("Invalid trello card link")
			}
		}
	}

}

//lexAfterLink is for lexing after we passed the card id in trello link
func lexAfterLink(l *lexer) stateFn {
	//here we expect bunch of characters we stop at eof newline or space
	for {
		switch r := l.next(); {
		case r == eof:
			l.emit(itemEOF)
			return nil
		case r == '\n':
			l.ignore()
			return lexText
		case unicode.IsSpace(r):
			l.ignore()
			return lexAfterLinkInstructions
		}
	}
}

//lexAfterLinkInstructions is for lexing after link instructions such as moving the card somewhere
func lexAfterLinkInstructions(l *lexer) stateFn {
	//move or move to is the only one we support right now
	if strings.ToLower(l.input[l.pos:l.pos+len(conMoveTo)]) == conMoveTo {
		l.skip(len(conMoveTo))
		return lexMoveDst
	} else if strings.ToLower(l.input[l.pos:l.pos+len(conMove)]) == conMove {
		l.skip(len(conMove))
		return lexMoveDst
	} else {
		//just ignore anything till end or newline
		for {
			switch r := l.next(); {
			case r == eof:
				l.emit(itemEOF)
				return nil
			case r == '\n':
				l.ignore()
				return lexText

			}
		}
	}
}

//lexMoveDst lexes destination for this card to be moved
//it ends at newline or eof
func lexMoveDst(l *lexer) stateFn {
	for {
		switch r := l.next(); {
		case r == eof:
			if l.pos > l.start+1 {
				l.emit(itemMoveDestination)
			}
			l.emit(itemEOF)
			return nil
		case r == '\r':
			if l.peek() == '\n' {
				l.backup()
				l.emit(itemMoveDestination)
				l.skip(2)
				return lexText
			}

		case r == '\n':
			l.backup()
			l.emit(itemMoveDestination)
			l.skip(1)
			return lexText

		}
	}
}

// lex is the actual function that starts the scanner
// it returns the scanner, and channel on which to
// receive the tokens
func lex(name, input string) *lexer {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run() //start the scanner
	return l
}