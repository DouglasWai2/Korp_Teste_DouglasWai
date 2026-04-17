package models

type NotaFiscal struct {
	Numero int64  `json:"numero" db:"numero"`
	Status string `json:"status" db:"status"`
}

type NotaFiscalItem struct {
	CodigoProduto string `json:"codigo_produto" db:"codigo_produto"`
	Quantidade    int    `json:"quantidade" db:"quantidade"`
}

type NotaFiscalDetail struct {
	Numero int64            `json:"numero"`
	Status string           `json:"status"`
	Itens  []NotaFiscalItem `json:"itens"`
}
