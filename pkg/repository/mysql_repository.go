package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/caarlos0/env/v6"
	"github.com/go-kit/kit/log"
	_ "github.com/go-sql-driver/mysql"

	"github.com/neidersalgado/go-grpc-api/pkg/users"
)

const (
	QUERYSAVEUSER     = `INSERT INTO user(name,pwdhash,age,aditional_information,email) values (?,?,?,?,?);`
	QUERYCREATEPARENT = `INSERT INTO parent(parent, son) (?,?)`
	QUERYDELETEUSER   = `DELETE FROM user WHERE email = ?`
	QUERYGETUSER      = `SELECT Id, name, pwdhash, age, aditional_information FROM user WHERE email = ?`
	QUERYGETUSERS     = `SELECT Id, email, name, pwdhash, age, aditional_information FROM user`
	QUERYUPDATEUSER   = `UPDATE user SET name= ?, age =?, aditional_information =? WHERE email = ?`
)

type config struct {
	User     string `env:"MYSQL_USER" envDefault:"root"`
	Password string `env:"MYSQL_PASSWORD" envDefault:"BulkD3v_mysql"`
	Port     string `env:"MYSQL_PORT" envDefault:":3306"`
	Host     string `env:"MYSQL_HOST" envDefault:"127.0.0.1"`
	//Host      string `env:"MYSQL_HOST" envDefault:"127.0.0.1"`
	DefaultDB string `env:"MYSQL_DEFAULTDB" envDefault:"users"`
}

func initMySQLRepository(log log.Logger) (*sql.DB, error) {

	cfg := config{}

	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s%s)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DefaultDB)
	db, err := sql.Open("mysql", connectionString)
	log.Log("Connection String", connectionString)
	if err != nil {
		log.Log("Connection", "error", err.Error())
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Log("Connection", "error", err.Error())
		return nil, err
	}
	return db, nil
}

type MySQLUserRepository struct {
	db     *sql.DB
	logger log.Logger
}

func NewMySQLUserRepository(logger log.Logger) (*MySQLUserRepository, error) {
	db, err := initMySQLRepository(logger)
	if err != nil {
		logger.Log("Connection", "Cant connect with bd")
		return nil, err
	}
	return &MySQLUserRepository{
		db:     db,
		logger: logger,
	}, nil
}

func (r *MySQLUserRepository) Create(ctx context.Context, user users.User) error {
	stmtSaveUser, err := r.db.Prepare(QUERYSAVEUSER)
	r.logger.Log("repository", fmt.Sprintf("creating user %+v", user))
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("Connetion Error, Couldn't save User With ID: %s in database, Error: %v", user.Email, err.Error()),
		)
	}
	_, errExec := stmtSaveUser.Exec(user.Name, user.PwdHash, user.Age, user.AdditionalInformation, user.Email)

	if errExec != nil {
		return fmt.Errorf(
			fmt.Sprintf("Database Exec Error, Couldn't save User With ID: %s in database, Error: %v", user.Email, errExec.Error()),
		)
	}

	return nil
}

func (r *MySQLUserRepository) GetByEmail(ctx context.Context, email string) (users.User, error) {
	stmt, err := r.db.Prepare(QUERYGETUSER)
	defer stmt.Close()
	if err != nil {
		return users.User{}, fmt.Errorf(
			fmt.Sprintf("Connetion Error, Couldn't get User With ID: %s in database, Error: %v", email, err.Error()),
		)
	}
	userDB := users.User{}

	row := stmt.QueryRow(email)
	err = row.Scan(
		&userDB.UserId,
		&userDB.Name,
		&userDB.PwdHash,
		&userDB.Age,
		&userDB.AdditionalInformation,
	)

	if err == sql.ErrNoRows {
		return userDB, nil
	}

	return userDB, err

}

func (r *MySQLUserRepository) Update(ctx context.Context, userToUpdate users.User) error {
	user, err := r.GetByEmail(ctx, userToUpdate.Email)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("Database Exec Error, Couldn't Update User With ID: %s in database, Error: %v", user.Email, err.Error()),
		)
	}

	equal, userToUpdate := getUserToUpdate(user, userToUpdate)
	if !equal {
		stmtUpdateUser, err := r.db.Prepare(QUERYUPDATEUSER)
		if err != nil {
			return fmt.Errorf(
				fmt.Sprintf("Connetion Error, Couldn't save User With ID: %s in database, Error: %v", user.Email, err.Error()),
			)
		}
		_, errExec := stmtUpdateUser.Exec(
			userToUpdate.Name,
			userToUpdate.Age,
			userToUpdate.AdditionalInformation,
			userToUpdate.Email,
		)
		if errExec != nil {
			return fmt.Errorf(
				fmt.Sprintf("Database Exec Error, Couldn't save User With ID: %s in database, Error: %v", user.Email, errExec.Error()),
			)
		}
	}

	return nil
}
func (r *MySQLUserRepository) Delete(ctx context.Context, email string) error {
	stmtSaveUser, err := r.db.Prepare(QUERYDELETEUSER)
	if err != nil {
		return fmt.Errorf("Database Connection Error, Couldn't delete User With ID: %s in database", email)
	}

	_, err = stmtSaveUser.Exec(email)

	if err != nil {
		return fmt.Errorf("Database Exec Error, Couldn't delete User With ID: %s in database", email)
	}

	return nil
}

func (r *MySQLUserRepository) GetAll(ctx context.Context) ([]users.User, error) {
	stmt, err := r.db.Prepare(QUERYGETUSERS)
	if err != nil {
		return []users.User{}, fmt.Errorf(
			fmt.Sprintf("Database prepare stmt Error, Couldn't get Users in database,\n Error: %v", err.Error()),
		)
	}

	rows, errExec := stmt.Query()

	if errExec != nil {
		return []users.User{}, fmt.Errorf(
			fmt.Sprintf("Database query Error, Couldn't get Users in database,\n Error: %v", errExec.Error()),
		)
	}

	var usersResponse []users.User
	for rows.Next() {
		var user users.User
		err := rows.Scan(
			&user.UserId,
			&user.Email,
			&user.Name,
			&user.PwdHash,
			&user.Age,
			&user.AdditionalInformation,
		)

		if err != nil {
			return []users.User{}, fmt.Errorf(
				fmt.Sprintf("Repository  mapping Error, Couldn't get Users in database,\n Error: %v", errExec.Error()),
			)
		}

		usersResponse = append(usersResponse, user)
	}

	return usersResponse, nil
}

func getUserToUpdate(userDB users.User, userToUpdate users.User) (bool, users.User) {
	equal := true
	var userUpdate users.User
	userUpdate.Email = userToUpdate.Email
	if userDB.Name != userToUpdate.Name {
		userUpdate.Name = userToUpdate.Name
		equal = false
	}

	if userDB.AdditionalInformation != userToUpdate.AdditionalInformation {
		userUpdate.AdditionalInformation = userToUpdate.AdditionalInformation
		equal = false
	}

	if userDB.PwdHash != userToUpdate.PwdHash {
		userUpdate.PwdHash = userToUpdate.PwdHash
		equal = false
	}

	if userDB.Age != userToUpdate.Age {
		userUpdate.Age = userToUpdate.Age
		equal = false
	}

	if equal {
		return equal, users.User{}
	}
	return equal, userUpdate
}
