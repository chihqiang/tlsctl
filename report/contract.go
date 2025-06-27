package report

type IReport interface {
	WithEnvConfig() error
	SendText(msg string) error
}
