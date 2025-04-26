package sbmap

import (
	"fmt"
	"hash/maphash"
	"slices"
	"strings"
)

func Example_simple() {
	h := New[int, string]()
	h.Put(1, "one")
	h.Put(2, "two")
	h.Put(3, "three")

	if val, ok := h.Get(1); ok {
		fmt.Println(val)
	}
	if _, ok := h.Get(4); !ok {
		fmt.Println("4 was not in the map!")
	}

	h.Remove(1)
	fmt.Println(h.Len())

	fmt.Println("Keys:")
	keys := slices.Collect(h.Keys())
	slices.Sort(keys)
	fmt.Println(keys)

	//Output:
	// one
	// 4 was not in the map!
	// 2
	// Keys:
	// [2 3]
}

func Example_customEqAndHashFuncs() {
	seed := maphash.MakeSeed()

	h := NewCustom[string, string](
		4,
		func(l, r string) bool { return strings.ToLower(l) == strings.ToLower(r) },
		func(v string) uint64 { return maphash.String(seed, strings.ToLower(v)) },
	)
	h.Put("one", "one")
	h.Put("two", "two")
	h.Put("three", "three")
	h.Put("ThReE", "four")

	if val, ok := h.Get("three"); ok {
		fmt.Println(val)
	}

	h.Remove("tWO")
	fmt.Println(h.Len())

	fmt.Println("Keys:")
	keys := slices.Collect(h.Keys())
	slices.Sort(keys)
	fmt.Println(keys)

	//Output:
	// four
	// 2
	// Keys:
	// [one three]
}
