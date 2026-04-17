package models

type Product struct {
	Codigo    string `json:"codigo" db:"codigo"`
	Descricao string `json:"descricao" db:"descricao"`
	Saldo     int    `json:"saldo" db:"saldo"`
}
