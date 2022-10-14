package db

import (
	"github.com/gongbell/contractfuzzer/db/domain"
)

var Oracles = map[string]domain.Oracle{
	"delegate":             {Id: "714ae490-edd9-4f66-ae39-edd9ce2eb8bd", Name: "delegate"},
	"exception_disorder":   {Id: "62aeff49-6103-408c-ac49-c123ecf18b94", Name: "exception_disorder"},
	"gasless_send":         {Id: "14d6f912-3b8a-4a99-858d-0466142220b2", Name: "gasless_send"},
	"number_dependency":    {Id: "cb0d2424-16c7-4475-b139-a85f3acb0cee", Name: "number_dependency"},
	"reentrancy":           {Id: "358508e7-5a6b-4025-a712-f69890696587", Name: "reentrancy"},
	"timestamp_dependency": {Id: "0d199a66-c483-4002-8a0c-2f73af7f989c", Name: "timestamp_dependency"},
}

func (m SQLiteManager) Seed() error {
	tx, _ := m.db.Begin()

	// Adding oracles
	for _, oracle := range Oracles {
		_, err := m.db.Exec(`INSERT INTO oracles(id, name) (?, ?)`, oracle.Id, oracle.Name)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}
