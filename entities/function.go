package entities

type Function struct {
	Id            string `gorm:"primaryKey"`
	Name          string `gorm:"index;not null"`
	NumberOfArgs  int64  `gorm:"not null"`
	Payable       bool   `gorm:"not null"`
	IsConstructor bool   `gorm:"not null"`
	ContractId    string `gorm:"not null"`
	Transactions  []Transaction
}
