package finder

import "testing"

var createProvider = []struct {
  in string
  shouldFail bool
  result string
} {
  {"fixture", false, "fixture"},
  {"fixture/", false, "fixture"},
  {"fixture/f1", false, "fixture/f1"},
  {"fixture/f1/", false, "fixture/f1"},
  {" fixture/f1/    ", false, "fixture/f1"},
  {"fixture/fake", true, ""},
}

func TestCreate(t *testing.T) {
  for i,elem := range createProvider {
    finder := Create(elem.in)
    if elem.shouldFail && finder.err == nil {
      t.Errorf("Create test, dataset %d. Expected error, result no error", i)
    } else if !elem.shouldFail && finder.err != nil {
      t.Errorf("Create test, dataset %d. Expected no error, result error", i)
    } else if !elem.shouldFail && finder.finder.in != elem.result {
      t.Errorf("Create test, dataset %d. Expected %q, result %q", i, elem.result, finder.finder.in)
    }
  }
}

var depthProvider = []struct {
  depth int
  shouldFail bool
} {
  {1, false},
  {0, false},
  {10, false},
  {-1, true},
  {-10, true},
}

func TestDepth(t *testing.T) {
  for i,elem := range depthProvider {
    finder := Create("fixture").Depth(elem.depth)
    if elem.shouldFail && finder.err == nil {
      t.Errorf("Depth test, dataset %d. Expected error, result no error", i)
    } else if !elem.shouldFail && finder.err != nil {
      t.Errorf("Depth test, dataset %d. Expected no error, result error", i)
    } else if !elem.shouldFail && finder.finder.depth != elem.depth {
      t.Errorf("Depth test, dataset %d. Expected %d, result %d", i, elem.depth, finder.finder.depth)
    }
  }
}

var nameProvider = []struct {
  pattern string
  negative bool
  shouldFail bool
} {
  {"test", false, false},
  {"(*", false, true},
  {"test", true, false},
  {"(*", true, true},
}

func TestNames(t *testing.T) {
  for i,elem := range nameProvider {
    finder := Create("fixture")
    if elem.negative {
      finder = finder.NotName(elem.pattern)
    } else {
      finder = finder.Name(elem.pattern)
    }
    if elem.shouldFail && finder.err == nil {
      t.Errorf("Names test, dataset %d. Expected error, result no error", i)
    } else if !elem.shouldFail && finder.err != nil {
      t.Errorf("Names test, dataset %d. Expected no error, result error", i)
    } else if !elem.shouldFail {
      if elem.negative {
        if len(finder.finder.namesNotLike) != 1 {
          t.Errorf("Names test (not), dataset %d. Expected slice length 1, result %d", i, len(finder.finder.namesNotLike))
        }
      } else {
        if len(finder.finder.namesLike) != 1 {
          t.Errorf("Names test, dataset %d. Expected slice length 1, result %d", i, len(finder.finder.namesNotLike))
        }
      }
    }
  }
}

var getProvider = []struct {
  finder FinderResult
  quantity int
} {
  {Create("fixture").Depth(0), 1},
  {Create("fixture").Depth(5), 10},
  {Create("fixture").Depth(5).NotName(".xml"), 9},
  {Create("fixture").Depth(4).Name(".xml").Name(".yml"), 2},
  {Create("fixture/f1").Depth(0), 0},
  {Create("fixture/f1").Depth(1), 1},
  {Create("fixture/f1").Depth(2), 3},
  {Create("fixture/f1").Depth(2).Name(".xml"), 1},
  {Create("fixture/f2").Depth(0), 1},
  {Create("fixture/f2").Depth(20), 1},
  {Create("fixture/f3").Depth(1), 5},
  {Create("fixture/f3").Depth(0), 4},
}

func TestGet(t *testing.T) {
  for i,elem := range getProvider {
    results := elem.finder.Get()
    if len(results) != elem.quantity {
      t.Errorf("Get test, dataset %d. Expected quantity %d, result %d", i, elem.quantity, len(results))      
    }
  }
}