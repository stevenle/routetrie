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

package routetrie

import (
	"testing"
)

func TestRouteTrie(t *testing.T) {
	trie := NewRouteTrie()

	// Verify a direct path.
	trie.Add("example.com/foo", 1)
	value, _ := trie.Get("example.com/foo")
	if value != 1 {
		t.Errorf("Expected value to be 1, but was: %v", value)
	}

	// Verify another direct path.
	trie.Add("example.com/foo/bar", 2)
	value, _ = trie.Get("example.com/foo/bar")
	if value != 2 {
		t.Errorf("Expected value to be 2, but was: %v", value)
	}

	// Verify a :param path.
	trie.Add("example.com/foo/:bar/baz", 3)
	value, params := trie.Get("example.com/foo/bar/baz")
	if value != 3 {
		t.Errorf("Expected value to be 3, but was: %v", value)
	}
	if params["bar"] != "bar" {
		t.Errorf("Expected value of :bar to be \"bar\", but was: %v", params["bar"])
	}

	// Verify an unnamed *wildcard path.
	trie.Add("example.com/*", 4)
	value, params = trie.Get("example.com/abc/def")
	if value != 4 {
		t.Errorf("Expected value to be 4, but was: %v", value)
	}
	if len(params) != 0 {
		t.Errorf("Expected empty map, but was: %v", params)
	}

	// Verify a named *wildcard path.
	trie.Add("example.com/foo/*wildcard", 5)
	value, params = trie.Get("example.com/foo/blahblahblah")
	if params["wildcard"] != "blahblahblah" {
		t.Errorf("Expected *wildcard to be \"blahblahblah\" but was %v", params["wildcard"])
	}

	// Verify multiple matching routes.
	trie.Add("example.com/one/:two", 6)
	trie.Add("example.com/:one/two", 7)
	trie.Add("example.com/one/*two", 8)
	value, params = trie.Get("example.com/one/two")
	if value != 6 {
		t.Errorf("Expected value to be 6, but was: %v", value)
	}
	if len(params) != 1 || params["two"] != "two" {
		t.Errorf("Unexpected :param values: %v", params)
	}

	// Verify routes that walk up and back down a trie.
	trie.Add("example.com/test/one/:two", 9)
	trie.Add("example.com/test/:one/two/three", 10)
	value, params = trie.Get("example.com/test/one/two/three")
	if value != 10 {
		t.Errorf("Expected value to be 10, but was: %v", value)
	}
	if len(params) != 1 || params["one"] != "one" {
		t.Errorf("Unexpected :param values: %v", params)
	}
}
