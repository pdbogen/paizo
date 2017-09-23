package dom

import (
	"golang.org/x/net/html"
	"strings"
)

func Attribute(n *html.Node, attr string) *string {
	attr = strings.ToLower(attr)
	for _, nA := range n.Attr {
		if nA.Key == attr {
			s := new(string)
			*s = nA.Val
			return s
		}
	}
	return nil
}

type Matcher func(*html.Node) bool

func (a Matcher) And(b ...Matcher) Matcher {
	return func(n *html.Node) bool {
		for _, m := range append([]Matcher{a}, b...) {
			if !m(n) {
				return false
			}
		}
		return true
	}
}

func (a Matcher) Or(b ...Matcher) Matcher {
	return func(n *html.Node) bool {
		for _, m := range append([]Matcher{a}, b...) {
			if m(n) {
				return true
			}
		}
		return false
	}
}

// Find performs a depth-first search to locate and return the first node (include the given node; its children; or
// siblings) for which the matcher returns true.
func Find(n *html.Node, m Matcher) *html.Node {
	if n == nil {
		return nil
	}
	if m(n) {
		return n
	}
	if node := Find(n.FirstChild, m); node != nil {
		return node
	}
	if node := Find(n.NextSibling, m); node != nil {
		return node
	}
	return nil
}

func WithTag(tag string) Matcher {
	return func(n *html.Node) bool {
		return n.Data == tag
	}
}

func WithAttribute(attr, value string) Matcher {
	return func(n *html.Node) bool {
		attrValue := Attribute(n, attr)
		return attrValue != nil && *attrValue == value
	}
}

func FindAll(in *html.Node, matcher Matcher) []*html.Node {
	if in == nil {
		return []*html.Node{}
	}
	accum := []*html.Node{}
	if matcher(in) {
		accum = append(accum, in)
	}
	accum = append(accum, FindAll(in.FirstChild, matcher)...)
	accum = append(accum, FindAll(in.NextSibling, matcher)...)
	return accum
}

// FindParent finds the first parent node (or the node itself) directly in the given node's lineage for which the given
// matcher returns true.
func FindParent(in *html.Node, matcher Matcher) *html.Node {
	if in == nil {
		return nil
	}
	if matcher(in) {
		return in
	}
	return FindParent(in.Parent, matcher)
}