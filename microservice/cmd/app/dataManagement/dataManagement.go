package dataManagement

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/john98nf/UltimateMicroservice/cmd/app/model"
	"golang.org/x/crypto/bcrypt"
)

var ResourceNotFoundError error = errors.New("Resource not found")
var DuplicateResource error = errors.New("Resource already exists")
var NoResourceModification error = errors.New("No modification took place")
var UnavailableUUIDGeneration error = errors.New("UUID generation currently unavailable")
var UserAuthenticationFailed error = errors.New("User authentication failed")
var InvalidToken error = errors.New("Invalid Token")

type Claims struct {
	Username string
	jwt.StandardClaims
}

type dbController struct {
	db  *sql.DB
	cfg *mysql.Config
}

type MiddlewareController struct {
	dbCtrl *dbController
	JWTKey []byte
}

func newDBController() *dbController {

	return &dbController{
		cfg: &mysql.Config{
			User:   os.Getenv("DBUSER"),
			Passwd: os.Getenv("DBPASSWORD"),
			Net:    "tcp",
			Addr:   os.Getenv("DBENDPOINT"),
			DBName: os.Getenv("DBSCHEMA"),
		},
	}
}

func InitiallizeNewMiddlewareController() *MiddlewareController {
	var ctrl *MiddlewareController = &MiddlewareController{}

	ctrl.JWTKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	ctrl.dbCtrl = newDBController()

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

func (ctrl *MiddlewareController) GetCompany(id uuid.UUID) (*model.Company, error) {

	var comp *model.Company = &model.Company{}
	row := ctrl.dbCtrl.db.QueryRow("SELECT BIN_TO_UUID(ID), "+
		"NAME, "+
		"IFNULL(DESCRIPTION,''), "+
		"EMPLOYEES, "+
		"REGISTRATION_STATUS, "+
		"LEGAL_TYPE "+
		"FROM COMPANIES WHERE ID = UUID_TO_BIN(?)", id.String())
	if err := row.Scan(&comp.Id,
		&comp.Name,
		&comp.Description,
		&comp.Employees,
		&comp.RegistrationStatus,
		&comp.LegalType); err != nil {
		if err == sql.ErrNoRows {
			return comp, ResourceNotFoundError
		}
		return comp, err
	}
	return comp, nil
}

func (ctrl *MiddlewareController) CreateCompany(cmp *model.Company) error {
	id, err := ctrl.produceNewRandomUUID()
	if err != nil {
		return err
	}
	cmp.Id = id

	sqlStatement := "INSERT INTO COMPANIES (ID, NAME, DESCRIPTION, EMPLOYEES, REGISTRATION_STATUS, LEGAL_TYPE) VALUES (UUID_TO_BIN(?), ?, ?, ?, ?, ?)"
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

func (ctrl *MiddlewareController) ModifyCompany(id uuid.UUID,
	name *string,
	description *string,
	employees *int,
	registrationStatus *bool,
	legalType *string) (*model.Company, error) {

	cmp, err := ctrl.GetCompany(id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		cmp.Name = *name
	}
	if description != nil {
		cmp.Description = *description
	}
	if employees != nil {
		cmp.Employees = *employees
	}
	if registrationStatus != nil {
		cmp.RegistrationStatus = *registrationStatus
	}
	if legalType != nil {
		cmp.LegalType = *legalType
	}
	sqlStatement := "UPDATE COMPANIES SET NAME = ?, DESCRIPTION = ?, EMPLOYEES = ?, REGISTRATION_STATUS = ?, LEGAL_TYPE = ? WHERE ID = UUID_TO_BIN(?)"
	if _, err := ctrl.dbCtrl.db.Exec(sqlStatement,
		cmp.Name,
		cmp.Description,
		cmp.Employees,
		cmp.RegistrationStatus,
		cmp.LegalType,
		cmp.Id); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, DuplicateResource
		}
		return nil, err
	}
	return cmp, nil
}

func (ctrl *MiddlewareController) DeleteCompany(id uuid.UUID) error {
	sqlStatement := "DELETE FROM COMPANIES WHERE ID = UUID_TO_BIN(?)"
	res, err := ctrl.dbCtrl.db.Exec(sqlStatement, id)
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

func (ctrl *MiddlewareController) produceNewRandomUUID() (uuid.UUID, error) {
	id := uuid.New()
	var temp struct{}
	for i := 0; i < 100; i++ {
		row := ctrl.dbCtrl.db.QueryRow("SELECT ID FROM COMPANIES WHERE ID = ?", id.String())
		if err := row.Scan(&temp); err != nil {
			if err == sql.ErrNoRows {
				return id, nil
			}
			return uuid.Nil, UnavailableUUIDGeneration
		}
	}
	return uuid.Nil, UnavailableUUIDGeneration
}

func (ctrl *MiddlewareController) FetchPasswordHash(usr string) ([]byte, error) {
	var hashedPassword []byte
	row := ctrl.dbCtrl.db.QueryRow("SELECT PASSWORD_HASH FROM USERS WHERE USERNAME = ?", usr)
	if err := row.Scan(&hashedPassword); err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func (ctrl *MiddlewareController) ValidateToken(tokenStr string) error {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("bad signed method received")
		}
		return ctrl.JWTKey, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return InvalidToken
	}

	return nil
}

func (ctrl *MiddlewareController) GenerateTokenJWT(creds *model.Credentials) (string, error) {

	storedHash, err := ctrl.FetchPasswordHash(creds.Username)
	if err != nil {
		log.Println(err)
		return "", UserAuthenticationFailed
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(creds.Password)); err != nil {
		return "", UserAuthenticationFailed
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(ctrl.JWTKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (ctrl *MiddlewareController) CreateUser(creds model.Credentials) error {

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	sqlStatement := "INSERT INTO USERS (USERNAME, PASSWORD_HASH) VALUES (?, ?)"
	if _, err := ctrl.dbCtrl.db.Exec(sqlStatement,
		creds.Username,
		passwordHash); err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return DuplicateResource
		}
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
