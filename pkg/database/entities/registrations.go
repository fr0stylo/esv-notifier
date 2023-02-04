package entities

import (
	"bytes"
	"encoding/gob"

	"github.com/samber/lo"
)

type Registrations []int64

func (r *Registrations) Update(value any, item []byte) ([]byte, error) {
	i := append(*r, value.(int64))

	i = lo.Uniq(i)

	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode(i); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Registrations) Insert(value any) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode([]int64{value.(int64)}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Registrations) Marshal() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode(r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *Registrations) Unmarshal(value []byte) (any, error) {
	if err := gob.NewDecoder(bytes.NewReader(value)).Decode(r); err != nil {
		return nil, err
	}

	return r, nil
}
