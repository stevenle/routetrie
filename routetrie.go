// Copyright 2013 Steven Le. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package routetrie implements a trie data structure that can be used for
// efficient URL route lookups.
package routetrie

import (
	"errors"
	"strings"
)

// RouteTrie is a trie data structure that can be used for efficient URL route
// lookups.
type RouteTrie struct {
	children      map[string]*RouteTrie
	paramChild    *paramChild
	wildcardChild *wildcardChild
	value         interface{}
}

func NewRouteTrie() *RouteTrie {
	t := RouteTrie{
		children: make(map[string]*RouteTrie),
	}
	return &t
}

// Add sets a value value to the trie at a particular route. The route is
// typically a URL path with or without host, e.g. "example.com/foo" or just
// "/foo", and can contain :param placeholders and *wildcard placeholders.
func (t *RouteTrie) Add(route string, value interface{}) error {
	route = normalize(route)

	// If the end has been reached, save the value to the node.
	if route == "" {
		t.value = value
		return nil
	}

	// Split off the head path from the route.
	// For example: 'foo/bar/baz' => 'foo', 'bar/baz'.
	head, tail := split(route)

	var nextNode *RouteTrie
	if head[0] == ':' {
		paramChildName := head[1:]
		if t.paramChild == nil {
			t.paramChild = &paramChild{
				name: paramChildName,
				trie: NewRouteTrie(),
			}
		} else {
			if t.paramChild.name != paramChildName {
				return errors.New("Cannot set two different :param names at the same path")
			}
		}
		nextNode = t.paramChild.trie
	} else if head[0] == '*' {
		t.wildcardChild = &wildcardChild{
			name:  head[1:],
			value: value,
		}
		return nil
	} else {
		nextNode = t.children[head]
		if nextNode == nil {
			nextNode = NewRouteTrie()
			t.children[head] = nextNode
		}
	}

	return nextNode.Add(tail, value)
}

// Get finds the value of the trie at a particular route, preferring direct
// matches over :param matches, and :param matches over *wildcard matches.
// If a match is found, a map of placeholder values is also returned.
func (t *RouteTrie) Get(route string) (interface{}, map[string]string) {
	params := make(map[string]string)
	value := t.getValue(route, params)
	if value != nil {
		return value, params
	}
	return nil, nil
}

func (t *RouteTrie) getValue(route string, params map[string]string) interface{} {
	route = normalize(route)
	if route == "" {
		return t.value
	}
	head, tail := split(route)

	// Check for direct matches.
	childNode := t.children[head]
	if childNode != nil {
		value := childNode.getValue(tail, params)
		if value != nil {
			return value
		}
	}

	// Check for matches on :param paths.
	if t.paramChild != nil {
		value := t.paramChild.trie.getValue(tail, params)
		if value != nil {
			params[t.paramChild.name] = head
			return value
		}
	}

	// Check for matches on *wildcard paths.
	if t.wildcardChild != nil {
		if t.wildcardChild.name != "" {
			params[t.wildcardChild.name] = strings.TrimRight(route, "/")
		}
		return t.wildcardChild.value
	}

	return nil
}

type paramChild struct {
	name string
	trie *RouteTrie
}

type wildcardChild struct {
	name  string
	value interface{}
}

func split(route string) (string, string) {
	index := strings.Index(route, "/")
	if index == -1 {
		return route, ""
	}
	return route[:index], route[index+1:]
}

func normalize(route string) string {
	// Remove leading slashes.
	return strings.TrimLeft(route, "/")
}
