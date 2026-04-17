CREATE TABLE IF NOT EXISTS product (
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
    numero_nf BIGINT NOT NULL,
    quantidade NUMERIC NOT NULL DEFAULT 1,
    PRIMARY KEY (codigo_produto, numero_nf),
    CONSTRAINT fk_nf_produtos_produto
        FOREIGN KEY (codigo_produto) REFERENCES product(codigo),
    CONSTRAINT fk_nf_produtos_nota_fiscal
        FOREIGN KEY (numero_nf) REFERENCES notas_fiscais(numero)
);
