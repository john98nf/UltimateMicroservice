package dataManagement

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/john98nf/UltimateMicroservice/cmd/app/model"
)

var ResourceNotFoundError error = errors.New("Resource not found!")
var DuplicateResource error = errors.New("Resource already exists!")

type dbController struct {
	db  *sql.DB
	cfg *mysql.Config
}

type MiddlewareController struct {
	dbCtrl *dbController
}

func newDBController(env map[string]string) *dbController {

	return &dbController{
		cfg: &mysql.Config{
			User:   env["DBUSER"],
			Passwd: env["DBPASSWORD"],
			Net:    "tcp",
			Addr:   env["DBENDPOINT"],
			DBName: env["DBSCHEMA"],
		},
	}
}

func InitiallizeNewMiddlewareController(env map[string]string) *MiddlewareController {
	var ctrl *MiddlewareController = &MiddlewareController{}

	ctrl.dbCtrl = newDBController(env)

	if err := ctrl.dbCtrl.establishConnection(); err != nil {
		log.Fatal(err)
	}

	if err := ctrl.dbCtrl.pingConnection(); err != nil {
		log.Fatal(err)
	}

	log.Println("Successful database connection!")

	return ctrl
}

func (ctrl *MiddlewareController) TestConnection() error {
	if err := ctrl.dbCtrl.pingConnection(); err != nil {
		return err
	}
	return nil
}

func (ctrl *MiddlewareController) GetCompany(id uint) (*model.Company, error) {

	var comp *model.Company = model.NullCompany()
	row := ctrl.dbCtrl.db.QueryRow("SELECT ID, "+
		"NAME, "+
		"IFNULL(DESCRIPTION,''), "+
		"EMPLOYEES, "+
		"REGISTRATION_STATUS, "+
		"LEGAL_TYPE "+
		"FROM COMPANIES WHERE ID = ?", id)
	if err := row.Scan(&comp.Id,
		&comp.Name,
		&comp.Description,
		&comp.Employees,
		&comp.RegistrationStatus,
		&comp.LegalType); err != nil {
		if err == sql.ErrNoRows {
			return comp, ResourceNotFoundError
		}
		return comp, fmt.Errorf("company %d: %v", id, err)
	}
	return comp, nil
}

func (ctrl *MiddlewareController) CreateCompany(cmp *model.Company) error {
	sqlStatement := "INSERT INTO COMPANIES (ID, NAME, DESCRIPTION, EMPLOYEES, REGISTRATION_STATUS, LEGAL_TYPE) VALUES (?, ?, ?, ?, ?, ?)"
	if _, err := ctrl.dbCtrl.db.Exec(sqlStatement,
		cmp.Id,
		cmp.Name,
		cmp.Description,
		cmp.Employees,
		cmp.RegistrationStatus,
		cmp.LegalType); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return DuplicateResource
		}
		return err
	}
	return nil
}

func (ctrl *MiddlewareController) ModifyCompany(cmp *model.Company) error {
	sqlStatement := "UPDATE COMPANIES SET NAME = ?, DESCRIPTION = ?, EMPLOYEES = ?, REGISTRATION_STATUS = ?, LEGAL_TYPE = ? WHERE ID = ?"
	res, err := ctrl.dbCtrl.db.Exec(sqlStatement,
		cmp.Name,
		cmp.Description,
		cmp.Employees,
		cmp.RegistrationStatus,
		cmp.LegalType,
		cmp.Id)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return DuplicateResource
		}
		return err
	}
	if rows, err := res.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return ResourceNotFoundError
	}
	return nil
}

func (ctrl *MiddlewareController) DeleteCompany(cmpID uint) error {
	sqlStatement := "DELETE FROM COMPANIES WHERE ID = ?"
	res, err := ctrl.dbCtrl.db.Exec(sqlStatement, cmpID)
	if err != nil {
		return err
	}
	if rows, err := res.RowsAffected(); err != nil {
		return err
	} else if rows != 1 {
		return ResourceNotFoundError
	}
	return nil
}

func (ctrl *dbController) establishConnection() error {
	var err error
	ctrl.db, err = sql.Open("mysql", ctrl.cfg.FormatDSN())
	if err != nil {
		return err
	}
	return nil
}

func (ctrl *dbController) pingConnection() error {
	if err := ctrl.db.Ping(); err != nil {
		return err
	}
	return nil
}
