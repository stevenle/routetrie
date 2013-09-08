# routetrie

routetrie is an implementation of a trie data structure that can be used for
efficient URL route lookups. It supports :param and \*wildcard placeholders.

## Usage

```go
package main

import (
  "github.com/stevenle/routetrie"
)

func main() {
  trie := routetrie.NewRouteTrie()
  trie.Add("example.com/foo", 1)
  trie.Add("example.com/foo/:bar", 2)
  trie.Add("example.com/foo/:bar/*baz", 3)

  trie.Get("example.com/foo")
  // => 1, {}
  trie.Get("example.com/foo/bar")
  // => 2, {"bar": "bar"}
  trie.Get("example.com/foo/bar/baz/qux")
  // => 3, {"bar": "bar", "baz": "baz/qux"}
}
```
