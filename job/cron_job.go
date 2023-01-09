package job

type CronJob interface {
	Id() string
	CronConfig() string
	Handler()
}
