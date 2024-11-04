package observable

type Observer interface {
	Update(message string)
}