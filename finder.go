package finder

import (
	"errors"
	"os"
	"regexp"
	"strings"
)

var (
	ErrInvalidArgument = errors.New("finder: invalid argument")
	ErrNotDir          = errors.New("finder: not dir")
	ErrLogic           = errors.New("finder: logic exception")
)

type finder struct {
	namesLike, namesNotLike   []*regexp.Regexp
	in                        string
	depth                     int
	excludeHidden, excludeVCS bool
	minSize, maxSize          int64
}

type Finder struct {
	f   finder
	err error
}

// This method initializes the Finder component.
//
// Example:
//    finder := Create("/home/myuser/mydir")
func Create(path string) *Finder {
	path = strings.Trim(path, " ")

	if path[(len(path)-1):] == "/" && len(path) > 1 {
		path = path[:(len(path) - 1)]
	}

	baseFinder := Finder{
		f:   finder{in: path, minSize: -1, maxSize: -1, depth: -1},
		err: nil,
	}

	fi, err := os.Lstat(path)
	if err != nil {
		baseFinder.err = err
	} else if !fi.IsDir() {
		baseFinder.err = ErrNotDir
	}

	return &baseFinder
}

// Specify the depth to search (0 is the initial directory).
// If not specified, it will search over all directory structure recursively
//
// Example:
//    finder := Create("/home/myuser/mydir").Depth(1)
func (finder *Finder) Depth(depth int) *Finder {
	if depth < 0 {
		finder.err = ErrInvalidArgument
	} else {
		finder.f.depth = depth
	}

	return finder
}

// Add a pattern that excludes files by name
//
// Example:
//    finder := Create("/home/myuser/mydir").NotName(".py")
func (finder *Finder) NotName(pattern string) *Finder {
	regexp, e := regexp.Compile(pattern)

	if e != nil {
		finder.err = e
	} else {
		finder.f.namesNotLike = append(finder.f.namesNotLike, regexp)
	}

	return finder
}

// Add patterns that excludes files by name
//
// Example:
//    finder := Create("/home/myuser/mydir").NotNames([]string{".py", ".php"})
func (finder *Finder) NotNames(patterns []string) *Finder {
	var err error = nil
	for _, pattern := range patterns {
		finder := finder.NotName(pattern)
		if finder.err != nil {
			err = finder.err
		}
	}

	finder.err = err
	return finder
}

// Add pattern that includes only files that match the pattern (by name)
//
// Example:
//    finder := Create("/home/myuser/mydir").Name(".py")
func (finder *Finder) Names(patterns []string) *Finder {
	var err error = nil
	for _, pattern := range patterns {
		finder := finder.Name(pattern)
		if finder.err != nil {
			err = finder.err
		}
	}

	finder.err = err
	return finder
}

// Add patterns that includes only files that match the patterns (by name)
//
// Example:
//    finder := Create("/home/myuser/mydir").Names([]string{".py", ".php"})
func (finder *Finder) Name(pattern string) *Finder {
	regexp, e := regexp.Compile(pattern)

	if e != nil {
		finder.err = e
	} else {
		finder.f.namesLike = append(finder.f.namesLike, regexp)
	}

	return finder
}

// Filters files by size (size must be greater)
//
// Examples:
//    finder := Create("/home/myuser/mydir").MinSize("1 K")
//    finder := Create("/home/myuser/mydir").MinSize("230")
//    finder := Create("/home/myuser/mydir").MinSize("2 Gi")
func (finder *Finder) MinSize(size string) *Finder {
	result, err := sizeParser(size)

	if err == nil {
		if finder.f.maxSize != -1 {
			if result > finder.f.maxSize {
				finder.err = ErrLogic
			} else {
				finder.f.minSize = result
			}
		} else {
			finder.f.minSize = result
		}
	} else {
		finder.f.minSize = result
		finder.err = err
	}

	return finder
}

// Filters files by size (size must be less than the specified)
//
// Examples:
//    finder := Create("/home/myuser/mydir").MaxSize("1 K")
//    finder := Create("/home/myuser/mydir").MaxSize("230")
//    finder := Create("/home/myuser/mydir").MaxSize("2 Gi")
func (finder *Finder) MaxSize(size string) *Finder {
	result, err := sizeParser(size)

	if err == nil {
		if finder.f.minSize != -1 {
			if result < finder.f.minSize {
				finder.err = ErrLogic
			} else {
				finder.f.maxSize = result
			}
		} else {
			finder.f.maxSize = result
		}
	} else {
		finder.f.maxSize = result
		finder.err = err
	}

	return finder
}

// Excludes hidden directories (does not parse its contents)
//
// Example:
//    finder := Create("/home/myuser/mydir").ExcludeHidden()
func (finder *Finder) ExcludeHidden() *Finder {
	finder.f.excludeHidden = true
	return finder
}

// Excludes VCS directories (does not parse its contents)
//
// Example:
//    finder := Create("/home/myuser/mydir").ExcludeVCS()
func (finder *Finder) ExcludeVCS() *Finder {
	finder.f.excludeVCS = true
	return finder
}

// Gets a slide of os.FileInfo elements with the matched files
//
// Example:
//    finder := Create("/home/myuser/mydir").Get()
func (finder *Finder) Get() []os.FileInfo {
	return readDirectory(finder.f.in, 0, depth(finder.f.in), finder.f)
}
