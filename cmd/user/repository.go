package user

import (
	"database/sql"
	"fmt"

	"github.com/neidersalgado/go-camp-grpc/cmd/user/pb"
)

type MySQLUserRepository struct {
	ConnectionClient *sql.DB
	//TODO slice the contention with permissions read an write
}

const (
	querySaveUser   = `INSERT INTO user(name,pwdhash,age,aditional_information,email) values (?,?,?,?,?);`
	queryDeleteUser = `DELETE FROM user WHERE Id = ?`
	queryGetUser    = `SELECT * FROM user WHERE Id = ?`
	queryGetUsers   = `SELECT * FROM user`
)

func NewMySQLUserRepository(connection *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{
		ConnectionClient: connection,
	}
}

func (r *MySQLUserRepository) Create(user pb.UserRequest) error {
	stmtSaveUser, err := r.ConnectionClient.Prepare(querySaveUser)

	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("Connetion Error, Couldn't save User With ID: %s in database, Error: %v", user.Email, err.Error()),
		)
	}
	fmt.Println(user.PwdHash)
	_, errExec := stmtSaveUser.Exec(user.Name, user.PwdHash, user.Age, user.AdditionalInformation, user.Email)

	if errExec != nil {
		return fmt.Errorf(
			fmt.Sprintf("Database Exec Error, Couldn't save User With ID: %s in database, Error: %v", user.Email, errExec.Error()),
		)
	}

	return nil
}

func (r *MySQLUserRepository) Get(userID string) (pb.UserResponse, error) {
	var userResponse pb.UserResponse
	errExec := r.ConnectionClient.QueryRow(queryGetUser, userID).Scan(
		&userResponse.UserId,
		&userResponse.PwdHash,
		&userResponse.Age,
		&userResponse.AdditionalInformation,
		&userResponse.Name)

	if errExec != nil {
		return pb.UserResponse{}, fmt.Errorf(
			fmt.Sprintf("Database Exec Error, Couldn't get User With ID: %s in database, Error: %v", userID, errExec.Error()),
		)
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
	var usersResponse pb.UserColletionResponse
	rows, errExec := r.ConnectionClient.Query(queryGetUsers)

	if errExec != nil {
		return pb.UserColletionResponse{}, fmt.Errorf(
			fmt.Sprintf("Database query Error, Couldn't get Users in database,\n Error: %v", errExec.Error()),
		)
	}

	for rows.Next() {
		var user pb.UserResponse
		err := rows.Scan(
			&user.UserId,
			&user.PwdHash,
			&user.Email,
			&user.Name,
			&user.Age,
			&user.AdditionalInformation,
			&user.Parents,
		)

		if err != nil {
			return pb.UserColletionResponse{}, fmt.Errorf(
				fmt.Sprintf("Repository  mapping Error, Couldn't get Users in database,\n Error: %v", errExec.Error()),
			)
		}

		usersResponse.Users = append(usersResponse.Users, &user)
	}

	return usersResponse, nil
}
