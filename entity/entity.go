package entity

import "time"

type Entity interface {
	SetID(id string)
	GetID() string
	IsSame(entity Entity) bool
}

func (c *Codecast) SetID(id string) {
	c.ID = id
}

func (c *Codecast) GetID() string {
	return c.ID
}

func (c *Codecast) IsSame(entity Entity) bool {
	return c.ID != "" && c.ID == entity.GetID()
}

type Codecast struct {
	ID              string
	Title           string
	PublicationDate time.Time
}

type User struct {
	Username string
	ID       string
}

func (u *User) SetID(id string) {
	u.ID = id
}

func (u *User) GetID() string {
	return u.ID
}

func (u *User) IsSame(entity Entity) bool {
	return u.ID != "" && u.ID == entity.GetID()
}

type License struct {
	User     User
	Codecast Codecast
	Type     int8
}
