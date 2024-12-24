package handlers

import (
	"backlotofacil/models"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (h Handler) GetAllConvites(w http.ResponseWriter, r *http.Request) {

	var convites []models.Convite

	h.DB.Preload("Participante").Find(&convites)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(convites)
}

func (h Handler) GerarConvites(w http.ResponseWriter, r *http.Request) {
	log.Println("Entrando no gera convites")
	parmQtd := r.URL.Query().Get("qtd")
	var qtdConvites int
	qtdConvites, err := strconv.Atoi(parmQtd)
	if err == nil {
		log.Println("Qtd de convites recebidos:", qtdConvites)
		for i := 0; i < qtdConvites; i++ {
			idConvite := uuid.New()
			convite := models.Convite{ChaveConvite: idConvite.String()}
			h.DB.Create(&convite)
		}
	} else {
		log.Println("Erro na conversão da Qtd de convites recebidos:", err, "-->", parmQtd)
	}

	w.WriteHeader(http.StatusCreated)
}

func (h Handler) ValidarConvite(w http.ResponseWriter, r *http.Request) {
	parmCodigoChave := r.URL.Query().Get("id")

	log.Println("validando o convite ", parmCodigoChave)

	var convite models.Convite

	res := h.DB.Debug().Where("chave_convite = ? and id_utilizado = false", parmCodigoChave).First(&convite)
	if res.Error != nil {
		log.Println("Erro ao buscar convite:", res.Error.Error())
		w.WriteHeader(http.StatusNotFound)
	} else {
		log.Println("Convite: ", convite)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(convite)
	}
}
