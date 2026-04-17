package models

type NotaFiscal struct {
	Numero int64  `json:"numero" db:"numero"`
	Status string `json:"status" db:"status"`
}
