go-finder
=========

"Finder" is a simple library to find files in your filesystem (as it's name suggests), using pure Go (no external dependencies).

Please, don't use it until I write some tests (this library is just my playgound). But feel free to help me improve the library, as it's my very first Go code.

## Quick Start

```go
package main

import (
	"github.com/magdkudama/finder"
	"fmt"
)

func main() {
	results := finder.
		Create("/my/path/").
		Depth(1).
		Name(".go").
		NotName(".html")
		Get()

	for _, element := range results {
		fmt.Println(element.Name())
	}
}
```
