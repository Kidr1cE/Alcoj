package worker

type Worker interface {
	Register()
	Run()
}
