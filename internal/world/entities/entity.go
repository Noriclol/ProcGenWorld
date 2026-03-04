package entities

import "github.com/eihigh/vec"

type EntityKind int

const (
	KindPerson  EntityKind = iota
	KindAnimal
	KindMonster
)

type Entity interface {
	GetUID() uint64
	GetID() string
	GetPosition() vec.Vec2
	GetKind() EntityKind
}

type BaseEntity struct {
	UID      uint64
	ID       string
	Position vec.Vec2
	Kind     EntityKind
}

func (e *BaseEntity) GetUID() uint64       { return e.UID }
func (e *BaseEntity) GetID() string        { return e.ID }
func (e *BaseEntity) GetPosition() vec.Vec2 { return e.Position }
func (e *BaseEntity) GetKind() EntityKind  { return e.Kind }

type Person struct{ BaseEntity }
type Animal struct{ BaseEntity }
type Monster struct{ BaseEntity }
