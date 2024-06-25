package model

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

type Company struct {
	Id                 uint
	Name               string `max:"15"`
	Description        string `max:"3000"`
	Employees          int
	RegistrationStatus bool
	LegalType          string
}

func NewCompany(id uint,
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

func NullCompany() *Company {
	return &Company{
		Id:                 0,
		Name:               "",
		Description:        "",
		Employees:          0,
		RegistrationStatus: false,
		LegalType:          "NoValidType",
	}
}

func VerifyCompanyType(legalType string) bool {
	_, ok := legalTypes[companyType(legalType)]
	return ok
}
