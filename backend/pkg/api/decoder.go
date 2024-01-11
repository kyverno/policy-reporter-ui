package api

import (
	"encoding/json"
	"io"
)

func DecodeList[T any](r io.Reader) ([]T, error) {
	list := make([]T, 0)
	err := json.NewDecoder(r).Decode(&list)

	return list, err
}

func DecodeMap[R comparable, T any](r io.Reader) (map[R]T, error) {
	mapping := make(map[R]T)
	err := json.NewDecoder(r).Decode(&mapping)

	return mapping, err
}

func Decode[T any](r io.Reader) (*T, error) {
	model := new(T)
	err := json.NewDecoder(r).Decode(model)

	return model, err
}
