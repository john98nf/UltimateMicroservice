package model

type companyType int

const (
	corporations companyType = iota
	nonProfit
	cooperative
	soleProprietorship
)

type Company struct {
	id                 uint        `required`
	name               string      `required max:"15"`
	description        string      `max:"3000"`
	employees          int         `required`
	registrationStatus bool        `required`
	legalType          companyType `required`
}

func newCompany(newId uint,
	newName string,
	newDescription string,
	newEmployees int,
	newRegistrationStatus bool,
	newLegalType companyType) *Company {

	return &Company{
		id:                 newId,
		name:               newName,
		description:        newDescription,
		employees:          newEmployees,
		registrationStatus: newRegistrationStatus,
		legalType:          newLegalType,
	}
}
