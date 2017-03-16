package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type UserType uint8

const (
	Admin       = 1
	Manager     = 2
	Supervisor  = 3
	Operator    = 4
	Office      = 5
	Maintenance = 6
)

type Company struct {
	ID           bson.ObjectId `bson:"_id,omitempty"`
	Name         string        `bson:"name"`
	Address1     string        `bson:"address1"`
	Address2     string        `bson:"address2"`
	City         string        `bson:"city"`
	State        string        `bson:"state"`
	Postal       string        `bson:"postal"`
	PhoneNumber  string        `bson:"phoneNumber"`
	ContactEmail string        `bson:"contactEmail"`
	Created      time.Time     `bson:"created"`
	IsSuspended  bool          `bson:"isSuspended"`
	IsActive     bool          `bson:"isActive"`
}

type Machine struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	IsActive    bool          `bson:"isActive"`
	CompanyID   bson.ObjectId `bson:"companyId"`
	Name        string        `bson:"name"`
	Description string        `bson:"description"`
	Location    string        `bson:"location"`
	Created     time.Time     `bson:"created"`
}

type Part struct {
	ID          bson.ObjectId    `bson:"_id,omitempty"`
	Barcode     string           `bson:"barcode"`
	Name        string           `bson:"name"`
	Description string           `bson:"description"`
	ReferenceID string           `bson:"referenceId"`
	BatchID     bson.ObjectId    `bson:"batchId"`
	CompanyID   bson.ObjectId    `bson:"companyId"`
	Files       []PartFile       `bson:"files"`
	Attributes  []PartAttributes `bson:"attributes"`
	Created     time.Time        `bson:"created"`
}

type PartFile struct {
	MachineID     bson.ObjectId `bson:"machineId"`
	SystemID      string        `bson:"systemId"`
	FileExtension string        `bson:"extension"`
	FileName      string        `bson:"fileName"`
	Created       time.Time     `bson:"created"`
}

type PartAttributes struct {
	Name     string `bson:"name"`
	Value    string `bson:"value"`
	MappedTo string `bson:"mappedTo"`
}

type User struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	FirstName   string        `bson:"firstName"`
	LastName    string        `bson:"lastName"`
	Email       string        `bson:"email"`
	PhoneNumber string        `bson:"phoneNumber"`
	Username    string        `bson:"username"`
	Userlevel   UserType      `bson:"userlevel"`
	IsActive    bool          `bson:"isActive"`
	SMS         string        `bson:"SMS"`
	Created     time.Time     `bson:"created"`
	Roles       []Role        `bson:"roles"`
	CompanyID   bson.ObjectId `bson:"companyId"`
}

// Role - Defines a role used in the system
type Role struct {
	ID          bson.ObjectId `bson:"_id,omitempty"`
	Description string        `bson:"description"`
	IsActive    bool          `bson:"isActive"`
	Name        string        `bson:"name"`
}

// Token - represents the structure of a user token
type Token struct {
	UserID    string  `bson:"userID"`
	CompanyID string  `bson:"companyID"`
	Token     string  `bson:"token"`
	IssuedAt  float64 `bson:"issuedAt"`
	ExpiresAt float64 `bson:"expiresAt"`
}

// Login - represents login credentials
type Login struct {
	User      `bson:",inline"`
	Password  string    `bson:"password"`
	LastLogin time.Time `bson:"lastLogin"`
}
