package users

import (
	"context"
	"errors"

	"gopkg.in/go-playground/validator.v9"
)

type Service interface {
	Create(ctx context.Context, user User) error
	Get(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, email string) error
	GetAll(context.Context) ([]User, error)
}

type Repository interface {
	Create(ctx context.Context, user User) error
	GetByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, email string) error
	GetAll(ctx context.Context) ([]User, error)
}

type UserService struct {
	repository Repository
}

func NewUserService(repo Repository) *UserService {
	return &UserService{
		repository: repo,
	}
}

//Create - add a new user into database
func (us *UserService) Create(ctx context.Context, user User) error {

	v := validator.New()

	if errVal := v.Struct(user); errVal != nil {
		errorMessage := errVal.(validator.ValidationErrors)[0]
		return errors.New(errorMessage.Field() + " is not valid")
	}

	dbUser, err := us.repository.GetByEmail(ctx, user.Email)

	if err != nil {
		return err
	}

	if dbUser.UserId > 0 {

		return errors.New("user with Email already exists")
	}

	errCreate := us.repository.Create(ctx, user)

	if errCreate != nil {
		return errCreate
	}

	return nil

}

//GetByEmail - retrieves the information of a user based on the email
func (us *UserService) GetByEmail(ctx context.Context, email string) (User, error) {

	dbUser, err := us.repository.GetByEmail(ctx, email)

	if err != nil {
		return User{}, err
	}

	if dbUser.UserId == 0 {
		return User{}, errors.New("user not found")
	}

	return dbUser, nil

}

//GetAll -  gets all the existing users
func (us *UserService) GetAll(ctx context.Context) ([]User, error) {

	users, err := us.repository.GetAll(ctx)

	if err != nil {
		return []User{}, err
	}

	return users, nil
}

//Update - validates the data and updates the user information
func (us *UserService) Update(ctx context.Context, user User) error {

	v := validator.New()

	if errVal := v.Struct(user); errVal != nil {
		errorMessage := errVal.(validator.ValidationErrors)[0]
		return errors.New(errorMessage.Field() + " is not valid")
	}

	usrToUpdate, errU := us.repository.GetByEmail(ctx, user.Email)

	if errU != nil {
		return errU
	}

	if usrToUpdate.UserId == 0 {
		return errors.New("user not found")
	}

	user.UserId = usrToUpdate.UserId

	if err := us.repository.Update(ctx, user); err != nil {
		return errors.New("cannot update the user")
	}

	return nil
}

//Delete - removes a user
func (us *UserService) Delete(ctx context.Context, email string) error {

	userToDelete, err := us.repository.GetByEmail(ctx, email)

	if err != nil {
		return err
	}

	if userToDelete.UserId == 0 {
		return errors.New("user not found")
	}

	if errD := us.repository.Delete(ctx, userToDelete.Email); errD != nil {
		return errD
	}

	return nil
}
