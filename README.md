routetrie
=========

Efficient URL route lookups for Golang

routetrie is an implementation of a trie data structure that can be used for
efficient URL route lookups.

Named :params and \*wildcard placeholders are supported.

```go
package main

import (
  "github.com/stevenle/routetrie"
)

func main() {
  // Initialize the trie.
  trie := routetrie.NewRouteTrie()

  // Add routes.
  trie.Add("example.com/foo", 1)
  trie.Add("example.com/foo/:bar", 2)
  trie.Add("example.com/foo/:bar/*baz", 3)

  // Find routes.
  trie.Get("example.com/foo")  // => 1, {}
  trie.Get("example.com/foo/bar")  // => 2, {"bar": "bar"}
  trie.Get("example.com/foo/bar/baz/qux")  // => 3, {"bar": "bar", "baz": "baz/qux"}
}
```
