package teaagent

type EventInterface interface {
	AsJSON() (data []byte, err error)
}
