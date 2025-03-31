package main

import (
	"fmt"
	victor "victorgo/binding"
)

func main() {
	index, err := victor.AllocIndex(0, 1, 4)
	if err != nil {
		panic(err)
	}
	defer index.DestroyIndex()

	fmt.Println(index.Insert(1, []float32{1, 2, 3, 4}))

	fmt.Println(index.Insert(2, []float32{5, 6, 7, 8}))
	fmt.Println(index.Insert(3, []float32{9, 10, 11, 12}))

	r, e := index.Search([]float32{1, 2, 3, 4}, 4)
	if e != nil {
		panic(e)
	}
	fmt.Println(r.ID, r.Distance)

}
