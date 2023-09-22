package service

import "Patronus/model"

type AuthService interface {
	SignUpUser(signUpModel *model.UserSignUpModel) (*model.UserDBResponseModel, error)
	SignInUser(signInModel *model.UserSignInModel) (*model.UserDBResponseModel, error)
}
