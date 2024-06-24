package dataManagement

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

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
