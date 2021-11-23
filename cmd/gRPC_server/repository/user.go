package repository

import (
	"database/sql"
	"fmt"

	"github.com/neidersalgado/go-camp-grpc/cmd/gRPC_server/pb"
)

type MySQLUserRepository struct {
	ConnectionClient *sql.DB
	//TODO slice the contention with permissions read an write
}

const (
	querySaveUser   = `INSERT INTO user(id,name, pwdhash,age,aditional_information) values (?,?,?,?,?);`
	queryDeleteUser = `DELETE FROM user WHERE Id = ?`
	queryGetUser    = `SELECT * FROM user WHERE Id = ?`
)

func NewMySQLUserRepository(connection *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{
		ConnectionClient: connection,
	}
}

func (r *MySQLUserRepository) Create(user pb.UserRequest) error {
	stmtSaveUser, err := r.ConnectionClient.Prepare(querySaveUser)

	if err != nil {
		return fmt.Errorf("Connetion Error, Couldn't save User With ID: %s in database, Error: %v", user.Id, err)
	}

	_, errExec := stmtSaveUser.Exec(user.Id, user.Name, user.PwdHash, user.Age, user.AdditionalInformation)

	if errExec != nil {
		return fmt.Errorf("Database Exec Error, Couldn't save User With ID: %s in database, Error: %v", user.Id, err)
	}

	return nil
}

func (r *MySQLUserRepository) Get(userID string) (pb.UserResponse, error) {
	var userResponse pb.UserResponse
	errExec := r.ConnectionClient.QueryRow(queryGetUser, userID).Scan(
		&userResponse.Id,
		&userResponse.PwdHash,
		&userResponse.Age,
		&userResponse.AdditionalInformation,
		&userResponse.Name)

	if errExec != nil {
		return pb.UserResponse{}, fmt.Errorf("Database Exec Error, Couldn't get User With ID: %s in database, Error: %v", userID, errExec)
	}

	return userResponse, nil
}

func (r *MySQLUserRepository) Update(pb.UserRequest) error {
	return nil
}

func (r *MySQLUserRepository) Delete(userID string) error {
	stmtSaveUser, err := r.ConnectionClient.Prepare(queryDeleteUser)

	if err != nil {
		return fmt.Errorf("Database Connection Error, Couldn't delete User With ID: %s in database", userID)
	}

	_, err = stmtSaveUser.Exec(userID)

	if err != nil {
		return fmt.Errorf("Database Exec Error, Couldn't delete User With ID: %s in database", userID)
	}

	return nil
}

func (r *MySQLUserRepository) GetAll() (pb.UserColletionResponse, error) {
	return pb.UserColletionResponse{}, nil
}
