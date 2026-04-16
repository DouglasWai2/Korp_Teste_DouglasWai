package models

type Product struct {
	codigo    string `json:"codigo" db:"codigo"`
	descricao string `json:"descricao" db:"descricao"`
	saldo     int    `json:"saldo" db:"saldo"`
}
