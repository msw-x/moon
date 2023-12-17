package urest

type AuthRequest[Account any, Session any, T any] struct {
	Account Account
	Session Session
	Request[T]
}
