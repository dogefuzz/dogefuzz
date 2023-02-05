package entities

type Contract struct {
	Id            string `gorm:"primaryKey"`
	TaskId        string `gorm:"not null"`
	Address       string `gorm:"index"`
	Source        string `gorm:"not null"`
	CompiledCode  string `gorm:"not null"`
	AbiDefinition string `gorm:"not null"`
	Name          string `gorm:"not null"`
	CFG           string
	DistanceMap   string
	Functions     []Function
}
