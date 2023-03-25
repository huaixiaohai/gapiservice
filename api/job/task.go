package job

type IJob interface {
	Run()
	GetSpec() string
}
