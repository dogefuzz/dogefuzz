package db

const migrationQuery = `
CREATE TABLE IF NOT EXISTS tasks(
	id TEXT PRIMARY KEY,
	duration INT NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions(
	id TEXT PRIMARY KEY,
	blockchain_hash TEXT NOT NULL UNIQUE,
	task_id TEXT NOT NULL,
	CONSTRAINT fk_transactions_task_id_tasks
		FOREIGN KEY (task_id)
			REFERENCES tasks (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE,
	contract_id TEXT NOT NULL,
	CONSTRAINT fk_transactions_contract_id_contracts
		FOREIGN KEY (contract_id)
			REFERENCES contracts (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE,
	detected_weaknesses TEXT
);

CREATE TABLE IF NOT EXISTS oracles(
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tasks_oracles(
	task_id TEXT NOT NULL,
	CONSTRAINT fk_tasks_oracles_task_id_tasks
		FOREIGN KEY (task_id)
			REFERENCES tasks (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE,
	oracle_id TEXT NOT NULL,
	CONSTRAINT fk_tasks_oracles_oracle_id_oracles
		FOREIGN KEY (oracle_id)
			REFERENCES oracles (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS contract(
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	address TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tasks_contracts(
	task_id TEXT NOT NULL,
	CONSTRAINT fk_tasks_contracts_task_id_tasks
		FOREIGN KEY (task_id)
			REFERENCES tasks (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE,
	contract_id TEXT NOT NULL,
	CONSTRAINT fk_tasks_contracts_contract_id_contracts
		FOREIGN KEY (contract_id)
			REFERENCES contracts (id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
);
`

func (m SQLiteManager) Migrate() error {
	_, err := m.db.Exec(migrationQuery)
	return err
}
