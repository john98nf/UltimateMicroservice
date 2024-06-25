package model

import (
	"github.com/google/uuid"
)

type companyType string

const (
	NoValidType        companyType = ""
	Corporations       companyType = "Corporations"
	NonProfit          companyType = "NonProfit"
	Cooperative        companyType = "Cooperative"
	SoleProprietorship companyType = "SoleProprietorship"
)

var legalTypes map[companyType]struct{} = map[companyType]struct{}{
	NoValidType:        struct{}{},
	Corporations:       struct{}{},
	NonProfit:          struct{}{},
	Cooperative:        struct{}{},
	SoleProprietorship: struct{}{},
}

type Credentials struct {
	Username string
	Password string
}

type Company struct {
	Id                 uuid.UUID
	Name               string `max:"15"`
	Description        string `max:"3000"`
	Employees          int
	RegistrationStatus bool
	LegalType          string
}

func NewCompany(id uuid.UUID,
	name string,
	description string,
	employees int,
	registrationStatus bool,
	legalType string) *Company {

	return &Company{
		Id:                 id,
		Name:               name,
		Description:        description,
		Employees:          employees,
		RegistrationStatus: registrationStatus,
		LegalType:          legalType,
	}
}

func VerifyCompanyType(legalType string) bool {
	_, ok := legalTypes[companyType(legalType)]
	return ok
}
