CREATE TABLE IF NOT EXISTS produtos (
    codigo CHAR(6) UNIQUE,
    descricao VARCHAR(255) NULL,
    saldo NUMERIC
);

CREATE TABLE IF NOT EXISTS notas_fiscais (
    numero BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    status VARCHAR(10) NOT NULL,
    CONSTRAINT chk_notas_fiscais_status
        CHECK (Status IN ('Aberta', 'Fechada'))
);

CREATE TABLE IF NOT EXISTS nf_produtos (
    codigo_produto CHAR(6) NOT NULL,
    numero_nf NUMERIC NOT NULL,
    PRIMARY KEY (codigo_produto, numero_nf),
    CONSTRAINT fk_nf_produtos_produto
        FOREIGN KEY (codigo_produto) REFERENCES Produto(Codigo),
    CONSTRAINT fk_nf_produtos_nota_fiscal
        FOREIGN KEY (numero_nf) REFERENCES Notas_fiscais(Numero)
);