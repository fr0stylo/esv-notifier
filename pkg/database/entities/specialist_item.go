package entities

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type SpecialistItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (s SpecialistItem) String() string {
	return fmt.Sprintf("(%d) %s", s.ID, s.Name)
}

func (r *SpecialistItem) Update(value any, item []byte) ([]byte, error) {
	return r.Marshal()
}

func (r *SpecialistItem) Insert(value any) ([]byte, error) {
	return r.Marshal()
}

func (r *SpecialistItem) Marshal() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := gob.NewEncoder(buf).Encode(r); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (r *SpecialistItem) Unmarshal(value []byte) (any, error) {
	if err := gob.NewDecoder(bytes.NewReader(value)).Decode(r); err != nil {
		return nil, err
	}

	return r, nil
}
