package entities

type Contract struct {
	Id                 string `gorm:"primaryKey"`
	TaskId             string `gorm:"not null"`
	Address            string `gorm:"index"`
	Source             string `gorm:"not null"`
	DeploymentBytecode string `gorm:"not null"`
	RuntimeBytecode    string `gorm:"not null"`
	AbiDefinition      string `gorm:"not null"`
	Name               string `gorm:"not null"`
	CFG                string
	DistanceMap        string
	Functions          []Function
}
