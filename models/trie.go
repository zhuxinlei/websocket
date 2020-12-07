package models

import (
	"fmt"
	"strings"
)

const TrieSeparator = "."

// separated by '.', root val is empty string.
type Trie struct {
	Val      string
	Children []*Trie
	Valid    bool
}

func NewTrie(val string) *Trie {
	return &Trie{
		Val:      val,
		Children: make([]*Trie, 0),
	}
}

func (n *Trie) Add(val string) {
	values := strings.Split(val, TrieSeparator)
	next := n
	for _, val := range values {
		exists := false
		for _, sub := range next.Children {
			if sub.Val == val {
				exists = true
				next = sub
				break
			}
		}
		if !exists {
			child := NewTrie(val)
			next.Children = append(next.Children, child)
			next = child
		}
	}
	next.Valid = true
}

func (n *Trie) Exist(val string) bool {
	values := strings.Split(val, TrieSeparator)
	next := n
	for _, val := range values {
		exists := false
		for _, sub := range next.Children {
			if sub.Val == val {
				exists = true
				next = sub
				break
			}
		}
		if !exists {
			return false
		}
	}
	return next.Valid
}

func (n *Trie) IsDynamic(val string) (interface{}, bool) {
	tempSlash := strings.Split(val, "/")
	len := len(tempSlash)
	if len < 2 {
		return nil, false
	}
	index := len - 1
	temp := tempSlash[:index]
	tempString := strings.Replace(strings.Trim(fmt.Sprint(temp), "[]"), " ", ",", -1)
	values := strings.Split(tempString, TrieSeparator)

	next := n

	for _, val := range values {
		exists := false
		for _, sub := range next.Children {
			if sub.Val == val {
				exists = true
				next = sub
				break
			}
		}
		if !exists {
			return nil, false
		}
	}
	return tempSlash[len-1], next.Valid
}
