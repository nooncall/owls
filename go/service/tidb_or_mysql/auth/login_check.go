package auth

type LoginChecker interface {
	Login(userName, pwd string) error
}

var loginService LoginChecker

func SetLoginService(impl LoginChecker) {
	loginService = impl
}
