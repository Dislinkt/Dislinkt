package validator

import (
	"fmt"
	"regexp"

	goValidator "github.com/go-playground/validator/v10"
)

func InitValidator() *goValidator.Validate {
	customValidator := goValidator.New()
	err := customValidator.RegisterValidation("username", UsernameValidation)
	err = customValidator.RegisterValidation("password", PasswordValidation)
	err = customValidator.RegisterValidation("alphaenum", AlphaEnumValidation)
	err = customValidator.RegisterValidation("precondition", PreconditionValidation)
	err = customValidator.RegisterValidation("title", TitleValidation)
	if err != nil {
		return nil
	}
	return customValidator
}

func UsernameValidation(fl goValidator.FieldLevel) bool {
	usernameRegExp := "[a-zA-Z0-9!?#$@.*+_]*[!?#$@.*+_][a-zA-Z0-9!?#$@.*+_]*"
	matched, err := regexp.MatchString(usernameRegExp, fl.Field().String())
	if err != nil {
		return false
	}
	return matched
}
func AlphaEnumValidation(fl goValidator.FieldLevel) bool {
	RegExp := "[a-zA-Z0-9 ,']+"
	matched, err := regexp.MatchString(RegExp, fl.Field().String())
	if err != nil {
		return false
	}
	return matched
}
func TitleValidation(fl goValidator.FieldLevel) bool {
	RegExp := "[a-zA-Z '!?.,\"]+"
	matched, err := regexp.MatchString(RegExp, fl.Field().String())
	if err != nil {
		return false
	}
	return matched
}
func PreconditionValidation(fl goValidator.FieldLevel) bool {
	RegExp := "[a-zA-Z ,']+"
	matched, err := regexp.MatchString(RegExp, fl.Field().String())
	if err != nil {
		return false
	}
	return matched
}

func PasswordValidation(fl goValidator.FieldLevel) bool {
	const uppercase = `[A-Z]{1}`
	const lowercase = `[a-z]{1}`
	const number = `[0-9]{1}`
	const specialCharacters = `[!?#$@.*+_]{1}`

	if matched, err := regexp.MatchString(uppercase, fl.Field().String()); !matched || err != nil {
		fmt.Println("Your password should contain at least one uppercase letter.")
		return false
	}

	if matched, err := regexp.MatchString(lowercase, fl.Field().String()); !matched || err != nil {
		fmt.Println("Your password should contain at least one lowercase letter.")
		return false
	}

	if matched, err := regexp.MatchString(number, fl.Field().String()); !matched || err != nil {
		fmt.Println("Your password should contain at least one number.")
		return false
	}

	if matched, err := regexp.MatchString(specialCharacters, fl.Field().String()); !matched || err != nil {
		fmt.Println("Your password should contain at least one special character.")
		return false
	}
	fmt.Println("Thanks, you have entered a password in a valid format!")
	return true
}

func UsernameValidationString(username string) bool {
	usernameRegExp := "[a-zA-Z0-9!?#$@.*+_]*[!?#$@.*+_][a-zA-Z0-9!?#$@.*+_]*"
	matched, err := regexp.MatchString(usernameRegExp, username)
	if err != nil {
		return false
	}
	return matched
}
func EmailValidationString(email string) bool {
	RegExp := `^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`
	matched, err := regexp.MatchString(RegExp, email)
	if err != nil {
		return false
	}
	return matched
}
func UUIDValidation(uuid string) bool {
	RegExp := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`
	matched, err := regexp.MatchString(RegExp, uuid)
	if err != nil {
		return false
	}
	return matched
}
func JobOfferSearchValidation(search string) bool {
	RegExp := "[a-zA-Z0-9 ,]*"
	matched, err := regexp.MatchString(RegExp, search)
	if err != nil {
		return false
	}
	return matched
}
func UserSearchValidation(search string) bool {
	RegExp := "[a-zA-Z0-9!?#$@.*+_]*"
	matched, err := regexp.MatchString(RegExp, search)
	if err != nil {
		return false
	}
	return matched
}
