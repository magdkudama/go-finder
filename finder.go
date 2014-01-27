package finder

import (
  "os"
  "regexp"
  "path/filepath"
  "strings"
  "errors"
  "io/ioutil"
)

var (
  ErrInvalidArgument = errors.New("finder: invalid argument")
  ErrNotDir = errors.New("finder: not dir")
)

type Finder struct {
  namesLike, namesNotLike []*regexp.Regexp
  in string
  depth int
}

type FinderResult struct {
  finder Finder
  err error
}

func checkName(names []*regexp.Regexp, elementName string, hasToMatch bool, add *bool) {
  if *add == true {
    for _,r := range names {
      if r.MatchString(elementName) == (!hasToMatch) {
        *add = false
        break
      }
    }
  }
}

func isValidDepth(path string, maxDepth int, baseDepth int) bool {
  pathDepth := depth(path)

  if maxDepth > -1 {
    maxDepth += baseDepth
    if pathDepth > maxDepth {
      return false
    }
  }

  return true
}

func depth(path string) int {
  separator := string(os.PathSeparator)
  absolutePath,_ := filepath.Abs(path)

  return len(strings.Split(absolutePath, separator))
}

func readDirectory(path string, depth int, baseDepth int, f Finder) []os.FileInfo {
  var items []os.FileInfo

  var elements,err = ioutil.ReadDir(path)
  if(err == nil) {
    for _,element := range elements {
      if element.IsDir() {
        newPath := path + string(os.PathSeparator) + element.Name()
        if isValidDepth(newPath, f.depth, baseDepth) {
          recElements := readDirectory(newPath, depth + 1, baseDepth, f)
          for _,recElement := range recElements {
            items = append(items, recElement)
          }
        }
      } else {
        add := true
        checkName(f.namesLike, element.Name(), true, &add)
        checkName(f.namesNotLike, element.Name(), true, &add)
        if add {
          items = append(items, element)
        }
      }
    }
  }

  return items
}

func Create(path string) FinderResult {
  path = strings.Trim(path, " ")

  baseFinder := FinderResult {
    finder: Finder { in: path },
    err: nil,
  }

  if path[(len(path) - 1):] == "/" {
    path = path[:(len(path) - 1)]
  }

  fi,err := os.Lstat(path)
  if err != nil {
    baseFinder.err = err
  }

  if !fi.IsDir() {
    baseFinder.err = ErrNotDir
  }

  return baseFinder
}

func (f FinderResult) Depth(depth int) FinderResult {
  if depth < 0 {
    f.err = ErrInvalidArgument
  } else {
    f.finder.depth = depth
  }

  return f
}

func (f FinderResult) NotName(pattern string) FinderResult {
  regexp, e := regexp.Compile(pattern)

  if e != nil {
    f.err = e
  } else {
    f.finder.namesNotLike = append(f.finder.namesNotLike, regexp)
  }

  return f
}

func (f FinderResult) Name(pattern string) FinderResult {
  regexp, e := regexp.Compile(pattern)

  if e != nil {
    f.err = e
  } else {
    f.finder.namesLike = append(f.finder.namesNotLike, regexp)
  }

  return f
}

func (f FinderResult) Get() []os.FileInfo {
  return readDirectory(f.finder.in, 0, depth(f.finder.in), f.finder)
}
