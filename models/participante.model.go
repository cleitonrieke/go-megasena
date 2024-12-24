package models

import "time"

type Participante struct {
	Id                  int       `json:"id" gorm:"primaryKey"`
	NomeParticipante    string    `json:"nome_participante"`
	IdConvite           int       `json:"id_convite"`
	NumerosSelecionados string    `json:"numeros_selecionados"`
	IdPago              bool      `json:"id_pagamento"`
	DataPagamento       time.Time `json:"data_pagamento"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

func (p *Participante) TableName() string {
	// custom table name, this is default
	return "megavirada24.participantes"
}
