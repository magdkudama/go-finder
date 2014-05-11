package finder

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func isHidden(directoryName string) bool {
	return (directoryName[0:1] == ".")
}

func isVCS(directoryName string) bool {
	vcs := []string{".svn", "_svn", "CVS", "_darcs", ".arch-params", ".monotone", ".bzr", ".git", ".hg"}
	for _, vcsElem := range vcs {
		if vcsElem == directoryName {
			return true
		}
	}
	return false
}

func sizeParser(value string) (n int64, err error) {
	elements := strings.Split(value, " ")

	if len(elements) != 1 && len(elements) != 2 {
		return 0, ErrInvalidArgument
	}

	result, convErr := strconv.ParseInt(elements[0], 10, 0)
	if convErr != nil {
		return 0, convErr
	}

	if len(elements) == 1 {
		return result, convErr
	}

	var multiplyBy int64
	switch {
	case elements[1] == "K":
		multiplyBy = 1000
	case elements[1] == "Ki":
		multiplyBy = 1024
	case elements[1] == "M":
		multiplyBy = 1000000
	case elements[1] == "Mi":
		multiplyBy = 1024 * 1024
	case elements[1] == "G":
		multiplyBy = 1000000000
	case elements[1] == "Gi":
		multiplyBy = 1024 * 1024 * 1024
	default:
		return 0, ErrInvalidArgument
	}

	return (result * multiplyBy), nil
}

func checkName(names []*regexp.Regexp, elementName string, add *bool) {
	if *add == true {
		if len(names) > 0 {
			someoneMatches := false
			for _, r := range names {
				if r.MatchString(elementName) {
					someoneMatches = true
					break
				}
			}
			if !someoneMatches {
				*add = false
			}
		}
	}
}

func checkNotName(names []*regexp.Regexp, elementName string, add *bool) {
	if *add == true {
		for _, r := range names {
			if r.MatchString(elementName) {
				*add = false
				break
			}
		}
	}
}

func checkSize(minSize int64, maxSize int64, size int64, add *bool) {
	if *add == true {
		if maxSize != -1 && minSize != -1 {
			if size < minSize || size > maxSize {
				*add = false
			}
		} else if maxSize != -1 {
			if size > maxSize {
				*add = false
			}
		} else if minSize != -1 {
			if size < minSize {
				*add = false
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
	absolutePath, _ := filepath.Abs(path)

	return len(strings.Split(absolutePath, separator))
}

func readDirectory(path string, depth int, baseDepth int, f finder) []os.FileInfo {
	var items []os.FileInfo

	var elements, err = ioutil.ReadDir(path)
	if err == nil {
		for _, element := range elements {
			if element.IsDir() {
				newPath := path + string(os.PathSeparator) + element.Name()
				if isValidDepth(newPath, f.depth, baseDepth) {
					if !f.excludeHidden || !isHidden(element.Name()) {
						if f.excludeVCS && !isVCS(element.Name()) || !f.excludeVCS {
							recElements := readDirectory(newPath, depth+1, baseDepth, f)
							for _, recElement := range recElements {
								items = append(items, recElement)
							}
						}
					}
				}
			} else {
				add := true
				checkName(f.namesLike, element.Name(), &add)
				checkNotName(f.namesNotLike, element.Name(), &add)
				checkSize(f.minSize, f.maxSize, element.Size(), &add)
				if add {
					items = append(items, element)
				}
			}
		}
	}

	return items
}
