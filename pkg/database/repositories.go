package database

import "github.com/fr0stylo/esveikata-registracija/pkg/database/entities"

type Databases struct {
	specialists   *Database[*entities.SpecialistItem]
	registrations *Database[*entities.Registrations]
}

func NewDatabases(path string) (*Databases, error) {
	registrationRepository, err := NewDatabase[*entities.Registrations](path, "registrations")
	if err != nil {
		return nil, err
	}
	specialistRepository, err := NewDatabase[*entities.SpecialistItem](path, "specialists")
	if err != nil {
		return nil, err
	}

	return &Databases{
		specialists:   specialistRepository,
		registrations: registrationRepository,
	}, nil
}

func (d *Databases) Close() {
	d.Specialists().Close()
	d.Registrations().Close()
}

func (d *Databases) Specialists() *Database[*entities.SpecialistItem] {
	return d.specialists
}

func (d *Databases) Registrations() *Database[*entities.Registrations] {
	return d.registrations
}
