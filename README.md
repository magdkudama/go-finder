go-finder
=========

[![Build Status](https://travis-ci.org/magdkudama/go-finder.png?branch=master)](https://travis-ci.org/magdkudama/go-finder)

"Finder" is a simple library to find files in your filesystem (as it's name suggests), using pure Go (no external dependencies).

It's heavily inspired (I mean... it's a bad copy of...) on the awesome [Symfony Finder Component](https://github.com/symfony/Finder)

~~Please, don't use it until I write some tests~~ (this library is just my playground). But feel free to help me improve the library, as it's my very first Go code.

[Click to view documentation (auto-generated)](https://godoc.org/github.com/magdkudama/go-finder)

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
		Names([]string{".py",".json"}).
		NotName(".html").
		NotNames([]string{".xml",".yml"}).
		MinSize("2 K").
		MaxSize("4 Mi").
		ExcludeHidden().
		ExcludeVCS().
		Get()

	for _, element := range results {
		fmt.Println(element.Name())
	}
}
```

## To Do

* ~~Support filtering by size~~ Done! And supports defining min size and max size in a "fancy" format
* ~~Exclude hidden directories~~
* ~~Exclude VCS directories~~
* Add documentation on methods
* Exclude directories even before parsing
* Allow minDepth and maxDepth
* Clean-up API and internal methods

## Contributors

- Magd Kudama [magdkudama]