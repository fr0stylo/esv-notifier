package telegram

import "fmt"

type SpecialistItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (s SpecialistItem) String() string {
	return fmt.Sprintf("(%d) %s", s.ID, s.Name)
}
