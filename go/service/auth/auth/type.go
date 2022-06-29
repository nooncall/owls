package auth

const (
	DB = "db"
)

const (
	StatusCancelApply = "cancel_apply"
	StatusPass        = "paas"
	StatusReject      = "reject"
)

func Types() []string {
	return []string{DB}
}
