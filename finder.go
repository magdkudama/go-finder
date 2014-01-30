package finder

import (
  "os"
  "regexp"
  "strings"
  "errors"
)

var (
  ErrInvalidArgument = errors.New("finder: invalid argument")
  ErrNotDir = errors.New("finder: not dir")
  ErrLogic = errors.New("finder: logic exception")
)

type finder struct {
  namesLike, namesNotLike []*regexp.Regexp
  in string
  depth int
  minSize, maxSize int64
}

type Finder struct {
  f finder
  err error
}

func Create(path string) *Finder {
  path = strings.Trim(path, " ")

  if path[(len(path) - 1):] == "/" && len(path) > 1 {
    path = path[:(len(path) - 1)]
  }

  baseFinder := Finder {
    f: finder { in: path, minSize: -1, maxSize: -1, depth: -1 },
    err: nil,
  }

  fi,err := os.Lstat(path)
  if err != nil {
    baseFinder.err = err
  } else if !fi.IsDir() {
    baseFinder.err = ErrNotDir
  }

  return &baseFinder
}

func (finder *Finder) Depth(depth int) *Finder {
  if depth < 0 {
    finder.err = ErrInvalidArgument
  } else {
    finder.f.depth = depth
  }

  return finder
}

func (finder *Finder) NotName(pattern string) *Finder {
  regexp, e := regexp.Compile(pattern)

  if e != nil {
    finder.err = e
  } else {
    finder.f.namesNotLike = append(finder.f.namesNotLike, regexp)
  }

  return finder
}

func (finder *Finder) Name(pattern string) *Finder {
  regexp, e := regexp.Compile(pattern)

  if e != nil {
    finder.err = e
  } else {
    finder.f.namesLike = append(finder.f.namesLike, regexp)
  }

  return finder
}

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

func (finder *Finder) Get() []os.FileInfo {
  return readDirectory(finder.f.in, 0, depth(finder.f.in), finder.f)
}
