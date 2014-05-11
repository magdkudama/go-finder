package finder

import "testing"

var createProvider = []struct {
	in         string
	shouldFail bool
	result     string
}{
	{"fixture", false, "fixture"},
	{"fixture/", false, "fixture"},
	{"fixture/f1", false, "fixture/f1"},
	{"fixture/f1/", false, "fixture/f1"},
	{" fixture/f1/    ", false, "fixture/f1"},
	{"fixture/fake", true, ""},
}

func TestCreate(t *testing.T) {
	for i, elem := range createProvider {
		finder := Create(elem.in)
		if elem.shouldFail && finder.err == nil {
			t.Errorf("Create test, dataset %d. Expected error, result no error", i)
		} else if !elem.shouldFail && finder.err != nil {
			t.Errorf("Create test, dataset %d. Expected no error, result error", i)
		} else if !elem.shouldFail && finder.f.in != elem.result {
			t.Errorf("Create test, dataset %d. Expected %q, result %q", i, elem.result, finder.f.in)
		}
	}
}

var depthProvider = []struct {
	depth      int
	shouldFail bool
}{
	{1, false},
	{0, false},
	{10, false},
	{-1, true},
	{-10, true},
}

func TestDepth(t *testing.T) {
	for i, elem := range depthProvider {
		finder := Create("fixture").Depth(elem.depth)
		if elem.shouldFail && finder.err == nil {
			t.Errorf("Depth test, dataset %d. Expected error, result no error", i)
		} else if !elem.shouldFail && finder.err != nil {
			t.Errorf("Depth test, dataset %d. Expected no error, result error", i)
		} else if !elem.shouldFail && finder.f.depth != elem.depth {
			t.Errorf("Depth test, dataset %d. Expected %d, result %d", i, elem.depth, finder.f.depth)
		}
	}
}

var nameProvider = []struct {
	pattern    string
	negative   bool
	shouldFail bool
}{
	{"test", false, false},
	{"(*", false, true},
	{"test", true, false},
	{"(*", true, true},
}

func TestNames(t *testing.T) {
	for i, elem := range nameProvider {
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
				if len(finder.f.namesNotLike) != 1 {
					t.Errorf("Names test (not), dataset %d. Expected slice length 1, result %d", i, len(finder.f.namesNotLike))
				}
			} else {
				if len(finder.f.namesLike) != 1 {
					t.Errorf("Names test, dataset %d. Expected slice length 1, result %d", i, len(finder.f.namesNotLike))
				}
			}
		}
	}
}

var namesProvider = []struct {
	patterns   []string
	negative   bool
	shouldFail bool
}{
	{[]string{"test1", "test2"}, false, false},
	{[]string{"test1", "(*"}, false, true},
	{[]string{"test1", "test2"}, true, false},
	{[]string{"test1", "(*"}, true, true},
}

func TestSliceNames(t *testing.T) {
	for i, elem := range namesProvider {
		finder := Create("fixture")
		if elem.negative {
			finder = finder.NotNames(elem.patterns)
		} else {
			finder = finder.Names(elem.patterns)
		}
		if elem.shouldFail && finder.err == nil {
			t.Errorf("Names test (slice), dataset %d. Expected error, result no error", i)
		} else if !elem.shouldFail && finder.err != nil {
			t.Errorf("Names test (slice), dataset %d. Expected no error, result error", i)
		} else if !elem.shouldFail {
			if elem.negative {
				if len(finder.f.namesNotLike) != 2 {
					t.Errorf("Names test (not) (slice), dataset %d. Expected slice length 1, result %d", i, len(finder.f.namesNotLike))
				}
			} else {
				if len(finder.f.namesLike) != 2 {
					t.Errorf("Names test (slice), dataset %d. Expected slice length 1, result %d", i, len(finder.f.namesNotLike))
				}
			}
		}
	}
}

var sizeProvider = []struct {
	value      string
	shouldFail bool
	result     int64
}{
	{"1", false, 1},
	{"1 K", false, 1000},
	{"1 Ki", false, 1024},
	{"1 M", false, 1000000},
	{"100 Mi", false, 104857600},
	{"23 G", false, 23000000000},
	{"12 Gi", false, 12884901888},
	{"12G", true, 0},
	{"12 T", true, 0},
	{"1221X", true, 0},
	{"12  Gi", true, 0},
}

func TestSizes(t *testing.T) {
	for i, elem := range sizeProvider {
		finder := Create("fixture")
		finder = finder.MinSize(elem.value)
		if elem.shouldFail && finder.err == nil {
			t.Errorf("Sizes test, dataset %d. Expected error, result no error", i)
		} else if !elem.shouldFail && finder.err != nil {
			t.Errorf("Sizes test, dataset %d. Expected no error, result error", i)
		} else if !elem.shouldFail {
			if finder.f.minSize != elem.result {
				t.Errorf("Sizes test, dataset %d. Expected value %d, result %d", i, elem.result, finder.f.minSize)
			}
		}
	}
}

var getProvider = []struct {
	finder   *Finder
	quantity int
}{
	{Create("fixture").Depth(0), 1},
	{Create("fixture").Depth(5), 12},
	{Create("fixture").Depth(5).NotName(".xml"), 11},
	{Create("fixture").Depth(4).Name(".xml").Name(".yml"), 2},
	{Create("fixture").Depth(5).NotNames([]string{".xml"}), 11},
	{Create("fixture").Depth(4).Names([]string{".xml", ".yml"}), 2},
	{Create("fixture").ExcludeHidden().Depth(1), 6},
	{Create("fixture").ExcludeVCS().Depth(1), 7},
	{Create("fixture").ExcludeVCS().ExcludeHidden().Depth(1), 6},
	{Create("fixture/f1").Depth(0), 0},
	{Create("fixture/f1").Depth(1), 1},
	{Create("fixture/f1").Depth(2), 3},
	{Create("fixture/f1").Depth(2).Name(".xml"), 1},
	{Create("fixture/f2").Depth(0), 1},
	{Create("fixture/f2").Depth(20), 1},
	{Create("fixture/f3").Depth(1), 5},
	{Create("fixture/f3").Depth(0), 4},
	{Create("fixture").MinSize("1 K"), 3},
	{Create("fixture").MinSize("3 Ki"), 0},
	{Create("fixture").MinSize("1").MaxSize("1 Ki"), 1},
	{Create("fixture/f3").Depth(0).MinSize("1 K"), 1},
}

func TestGet(t *testing.T) {
	for i, elem := range getProvider {
		results := elem.finder.Get()
		if len(results) != elem.quantity {
			t.Errorf("Get test, dataset %d. Expected quantity %d, result %d", i, elem.quantity, len(results))
		}
	}
}
