package utils

import "math/rand"

func Some[T interface{}](slice []T, has func(T) bool) bool {
	for _, item := range slice {
		_has := has(item)
		if _has {
			return true
		}
	}
	return false
}

func Has[T comparable](slice *[]T, has T) bool {
	for _, item := range *slice {

		if has == item {
			return true
		}
	}
	return false
}

func HasByFunc[T interface{}](slice *[]T, is func(T) bool) bool {
	for _, item := range *slice {
		_has := is(item)
		if _has {
			return true
		}
	}
	return false
}

// 使切片中的元素都唯一
func Unique[T comparable](slice []T) *[]T {
	var newSlice []T
	var _map = make(map[T]int, 0)
	for _, item := range slice {
		_, has := _map[item]
		if !has {
			newSlice = append(newSlice, item)
			_map[item] = 1
		}
	}
	return &newSlice
}

// 使切片中的元素都唯一
func Filter[T interface{}](slice *[]T, filter func(T) bool) *[]T {
	var newSlice []T

	for _, item := range *slice {
		need := filter(item)
		if need {
			newSlice = append(newSlice, item)
		}
	}
	return &newSlice
}

func Find[T interface{}](slice *[]T, find func(T) bool) *T {
	for _, item := range *slice {
		ok := find(item)
		if ok {
			return &item
		}
	}
	return nil
}

func Map[T interface{}, O interface{}](slice *[]T, mapFunc func(T) O) *[]O {
	var newSlice []O

	for _, item := range *slice {
		o := mapFunc(item)
		newSlice = append(newSlice, o)
	}
	return &newSlice
}

func Shuffle[T interface{}](slice *[]T) {
	for i := range *slice {
		j := rand.Intn(i + 1)
		(*slice)[i], (*slice)[j] = (*slice)[j], (*slice)[i]
	}
}
