package utils

func IsStatusPositive(status int) bool {
	return status >= 200 && status < 300
}
