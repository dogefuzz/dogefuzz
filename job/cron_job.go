package job

type CronJob interface {
	Id()
	CronConfig() string
	Handler()
}
