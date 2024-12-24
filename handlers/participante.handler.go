package handlers

import (
	"backlotofacil/models"
	"encoding/json"
	"net/http"
)

type participanteJson struct {
	NomeParticipante   string `json:"nome_participante"`
	Convite            string `json:"convite"`
	NumeroSelecionados string `json:"numeros_selecionados"`
}

func (h Handler) Registrar(w http.ResponseWriter, r *http.Request) {

	var participanteJson participanteJson
	var participante models.Participante

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&participanteJson)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var convite models.Convite
	res := h.DB.Where("chave_convite = ?", participanteJson.Convite).First(&convite)

	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusBadRequest)
		return
	}

	convite.IdUtilizado = true
	h.DB.Save(&convite)

	participante.IdConvite = convite.Id
	participante.NomeParticipante = participanteJson.NomeParticipante
	participante.NumerosSelecionados = participanteJson.NumeroSelecionados
	h.DB.Create(&participante)

	w.WriteHeader(http.StatusCreated)
}

func (h Handler) GetAllParticipantes(w http.ResponseWriter, r *http.Request) {
	var participantes []models.Participante
	h.DB.Find(&participantes)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participantes)
}
