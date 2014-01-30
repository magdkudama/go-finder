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

type Finder struct {
  namesLike, namesNotLike []*regexp.Regexp
  in string
  depth int
  minSize, maxSize int64
}

type FinderResult struct {
  finder Finder
  err error
}

func Create(path string) FinderResult {
  path = strings.Trim(path, " ")

  if path[(len(path) - 1):] == "/" && len(path) > 1 {
    path = path[:(len(path) - 1)]
  }

  baseFinder := FinderResult {
    finder: Finder { in: path, minSize: -1, maxSize: -1, depth: -1 },
    err: nil,
  }

  fi,err := os.Lstat(path)
  if err != nil {
    baseFinder.err = err
  } else if !fi.IsDir() {
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
    f.finder.namesLike = append(f.finder.namesLike, regexp)
  }

  return f
}

func (f FinderResult) MinSize(size string) FinderResult {
  result, err := sizeParser(size)

  if err == nil {
    if f.finder.maxSize != -1 {
      if result > f.finder.maxSize {
        f.err = ErrLogic
      } else {
        f.finder.minSize = result
      }
    } else {
      f.finder.minSize = result
    }
  } else {
    f.finder.minSize = result
    f.err = err
  }

  return f
}

func (f FinderResult) MaxSize(size string) FinderResult {
  result, err := sizeParser(size)

  if err == nil {
    if f.finder.minSize != -1 {
      if result < f.finder.minSize {
        f.err = ErrLogic
      } else {
        f.finder.maxSize = result
      }
    } else {
      f.finder.maxSize = result
    }
  } else {
    f.finder.maxSize = result
    f.err = err
  }

  return f
}

func (f FinderResult) Get() []os.FileInfo {
  return readDirectory(f.finder.in, 0, depth(f.finder.in), f.finder)
}
