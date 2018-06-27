package main

import "fmt"

type Iterator struct {
	last  int
	vals  []interface{}
	child *Iterator
}

func NewIterator(vals []interface{}) *Iterator {
	return &Iterator{
		last: 0,
		vals: vals,
	}
}

func (i *Iterator) Next() int {
	if i.child != nil {
		return i.child.Next()
	}

	defer func() {
		i.last++
	}()

	if i.last < len(i.vals) {
		val := i.vals[i.last]
		if v, ok := val.(int); ok {
			return v
		} else if v, ok := val.([]interface{}); ok && len(v) > 0 {
			i.child = NewIterator(v)
			return i.child.Next()
		}

		i.last++
		return i.Next()
	}

	return 0
}

func (i *Iterator) HasNext() bool {
	if i.child != nil {
		if i.child.HasNext() {
			return true
		}
		i.child = nil
	}

	if i.last < len(i.vals) {
		return true
	}

	return false
}

func main() {
	list := []interface{}{1, []interface{}{2, []interface{}{3, 4}}, []interface{}{}, []interface{}{5}, 6}
	//list := []interface{}{1, 2, 3, 4}
	it := NewIterator(list)

	for it.HasNext() {
		fmt.Println(it.Next())
	}
}
