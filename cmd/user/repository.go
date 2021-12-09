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
	querySaveUser    = `INSERT INTO user(name,pwdhash,age,aditional_information,email) values (?,?,?,?,?);`
	querySaveParents = `INSERT INTO parent(parent, son) (?,?)`
	queryDeleteUser  = `DELETE FROM user WHERE email = ?`
	queryGetUser     = `SELECT Id, name, pwdhash, age, aditional_information FROM user WHERE email = ?`
	queryGetUsers    = `SELECT Id, email, name, pwdhash, age, aditional_information FROM user`
	queryUpdateUSer  = `UPDATE user SET name= ?, age =?, aditional_information =? WHERE email = ?`
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

func (r *MySQLUserRepository) Get(email string) (pb.UserResponse, error) {
	fmt.Printf("service.grpc.repository GET user email :%v \n", email)
	stmt, err := r.ConnectionClient.Prepare(queryGetUser)
	defer stmt.Close()
	if err != nil {
		return pb.UserResponse{}, fmt.Errorf(
			fmt.Sprintf("Connetion Error, Couldn't get User With ID: %s in database, Error: %v", email, err.Error()),
		)
	}
	var userResponse pb.UserResponse

	errExec := stmt.QueryRow(email).Scan(
		&userResponse.UserId,
		&userResponse.Name,
		&userResponse.PwdHash,
		&userResponse.Age,
		&userResponse.AdditionalInformation,
	)
	fmt.Printf("service.grpc.repository GET user row:%v \n", userResponse)
	if errExec != nil {
		return pb.UserResponse{}, fmt.Errorf(
			fmt.Sprintf("Database Exec Error, Couldn't get User With ID: %s in database, Error: %v", email, errExec.Error()),
		)
	}

	return userResponse, nil

}

func (r *MySQLUserRepository) Update(userRequest pb.UserRequest) error {
	user, err := r.Get(userRequest.Email)
	if err != nil {
		return fmt.Errorf(
			fmt.Sprintf("Database Exec Error, Couldn't Update User With ID: %s in database, Error: %v", user.Email, err.Error()),
		)
	}

	equal, userToUpdate := getUserToUpdate(user, userRequest)
	fmt.Printf("equal? : %+v \n userToUpdate %+v \n", equal, userToUpdate)
	if !equal {
		fmt.Println("updating")
		stmtUpdateUser, err := r.ConnectionClient.Prepare(queryUpdateUSer)
		fmt.Println("preparing")
		if err != nil {
			return fmt.Errorf(
				fmt.Sprintf("Connetion Error, Couldn't save User With ID: %s in database, Error: %v", user.Email, err.Error()),
			)
		}
		fmt.Println("executing")
		_, errExec := stmtUpdateUser.Exec(
			userToUpdate.Name,
			userToUpdate.Age,
			userToUpdate.AdditionalInformation,
			userRequest.Email,
		)
		if errExec != nil {
			return fmt.Errorf(
				fmt.Sprintf("Database Exec Error, Couldn't save User With ID: %s in database, Error: %v", user.Email, errExec.Error()),
			)
		}
		fmt.Println("non error updating")
	}

	return nil
}
func (r *MySQLUserRepository) Delete(email string) error {
	stmtSaveUser, err := r.ConnectionClient.Prepare(queryDeleteUser)

	if err != nil {
		return fmt.Errorf("Database Connection Error, Couldn't delete User With ID: %s in database", email)
	}

	_, err = stmtSaveUser.Exec(email)

	if err != nil {
		return fmt.Errorf("Database Exec Error, Couldn't delete User With ID: %s in database", email)
	}

	return nil
}

func (r *MySQLUserRepository) GetAll() (pb.UserColletionResponse, error) {
	fmt.Printf("service.grpc.repository getAll\n")
	stmt, err := r.ConnectionClient.Prepare(queryGetUsers)
	if err != nil {
		return pb.UserColletionResponse{}, fmt.Errorf(
			fmt.Sprintf("Database prepare stmt Error, Couldn't get Users in database,\n Error: %v", err.Error()),
		)
	}

	rows, errExec := stmt.Query()
	fmt.Sprintf("service.grpc.repository getAll execute \n")

	if errExec != nil {
		return pb.UserColletionResponse{}, fmt.Errorf(
			fmt.Sprintf("Database query Error, Couldn't get Users in database,\n Error: %v", errExec.Error()),
		)
	}
	var usersResponse pb.UserColletionResponse
	for rows.Next() {
		var user pb.UserResponse
		err := rows.Scan(
			&user.UserId,
			&user.Email,
			&user.Name,
			&user.PwdHash,
			&user.Age,
			&user.AdditionalInformation,
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

func getUserToUpdate(userDB pb.UserResponse, userRequest pb.UserRequest) (bool, pb.UserRequest) {
	equal := true
	var userToUpdate pb.UserRequest
	fmt.Printf("\n user request to update  \n  %+v \n  usuario response \n %+v", userRequest, userDB)

	if userDB.Name != userRequest.Name {
		fmt.Printf("\n name diferent \n")
		userToUpdate.Name = userRequest.Name
		equal = false
	}

	if userDB.AdditionalInformation != userRequest.AdditionalInformation {
		userToUpdate.AdditionalInformation = userRequest.AdditionalInformation

		fmt.Printf("\n information diferent \n")
		equal = false
	}

	if userDB.PwdHash != userRequest.PwdHash {
		userToUpdate.PwdHash = userRequest.PwdHash

		fmt.Printf("\n hash diferent \n")
		equal = false
	}

	if userDB.Age != userRequest.Age {
		userToUpdate.Age = userRequest.Age

		fmt.Printf("\n age diferent \n")
		equal = false
	}

	if equal {
		return equal, pb.UserRequest{}
	}
	return equal, userToUpdate
}
