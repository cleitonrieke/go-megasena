package handlers

import (
	"backlotofacil/models"
	"encoding/json"
	"math"
	"math/rand"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type KeyValue struct {
	Numero     int
	Ocorrencia int
}

type JsonSaida struct {
	Media               int
	Min                 int
	Max                 int
	DesvioPadrao        float64
	NumerosSelecionados []KeyValue
	NumerosFracos       []int
	NumerosMedios       []int
	NumerosFortes       []int
	Cartoes             []CartaoGerado
}

type CartaoGerado struct {
	Cartao []int
}

type JsonEntrada struct {
	QuantidadeCartoes int `json:"quantidade_cartoes"`
	QuantidadeNumeros int `json:"quantidade_numeros"`
}

func (h Handler) ResumoNumerosLotoFacil(w http.ResponseWriter, r *http.Request) {
	jsonEntrada := JsonEntrada{}
	err := json.NewDecoder(r.Body).Decode(&jsonEntrada)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var participantes []models.Participante
	h.DB.Find(&participantes)

	tabNumeros := make(map[int]int, 60)
	for i := 0; i < 60; i++ {
		tabNumeros[i+1] = 0
	}

	for _, participante := range participantes {
		numeros := strings.Split(participante.NumerosSelecionados, ",")
		for i := 0; i < len(numeros); i++ {
			numero, _ := strconv.Atoi(numeros[i])
			tabNumeros[numero]++
		}
	}

	var pairs []KeyValue
	for k, v := range tabNumeros {
		pairs = append(pairs, KeyValue{k, v})
	}

	// Sort slice by int value
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Ocorrencia < pairs[j].Ocorrencia
	})

	var qtdMin = pairs[0].Ocorrencia
	var qtdMax = pairs[len(pairs)-1].Ocorrencia

	var media = (qtdMin + qtdMax) / 2

	totalDesvio := 0

	for i := 0; i < len(pairs); i++ {
		totalDesvio += (pairs[i].Ocorrencia - media) * (pairs[i].Ocorrencia - media)
	}

	desvio := math.Sqrt(float64(totalDesvio / len(pairs)))
	limiteMinimo := media - int(math.Round(desvio))
	limiteMaximo := media + int(math.Round(desvio))

	var tabBaixos []int
	var tabMedios []int
	var tabFortes []int

	for i := 0; i < len(pairs); i++ {
		if pairs[i].Ocorrencia < limiteMinimo {
			tabBaixos = append(tabBaixos, pairs[i].Numero)
		} else if pairs[i].Ocorrencia > limiteMaximo {
			tabFortes = append(tabFortes, pairs[i].Numero)
		} else {
			tabMedios = append(tabMedios, pairs[i].Numero)
		}
	}

	var qtdFracas int
	var qtdMedias int
	var qtdFortes int
	var qtdCartoes int
	var qtdGerados int
	var qtdGeradosFracas int
	var qtdGeradosMedias int
	var qtdGeradosFortes int
	var randomNumber int
	var cartao []int
	var achouCartao bool
	var cartoesGerados []CartaoGerado

	// Cartoes de 17 Numeros
	if jsonEntrada.QuantidadeNumeros == 17 {
		qtdFracas = 1
		qtdMedias = 13
		qtdFortes = 3
	} else if jsonEntrada.QuantidadeNumeros == 16 {
		qtdFracas = 1
		qtdMedias = 12
		qtdFortes = 3
	} else if jsonEntrada.QuantidadeNumeros == 15 {
		qtdFracas = 1
		qtdMedias = 11
		qtdFortes = 3
	}
	qtdCartoes = jsonEntrada.QuantidadeCartoes
	qtdGerados = 0

	cartoesGerados = []CartaoGerado{}
	for qtdGerados < qtdCartoes {

		println("Inicio do Loop com qtdGerados: ", qtdGerados, " e qtdCartoes: ", qtdCartoes)
		cartao = []int{}
		qtdGeradosFracas = 0
		for qtdGeradosFracas < qtdFracas {
			randomNumber = rand.Intn(len(tabBaixos))
			println("Gerou o numero fraco ", tabBaixos[randomNumber])
			if !lo.Contains(cartao, tabBaixos[randomNumber]) {
				cartao = append(cartao, tabBaixos[randomNumber])
				qtdGeradosFracas++
				println("Adicionou o numero fraco ", tabBaixos[randomNumber], " qtdFracas ", qtdGeradosFracas)
			}
		}
		qtdGeradosMedias = 0
		for qtdGeradosMedias < qtdMedias {
			randomNumber = rand.Intn(len(tabMedios))
			println("Gerou o numero medio ", tabMedios[randomNumber])
			if !lo.Contains(cartao, tabMedios[randomNumber]) {
				cartao = append(cartao, tabMedios[randomNumber])
				qtdGeradosMedias++
				println("Adicionou o numero medio ", tabMedios[randomNumber], " qtdMedias ", qtdGeradosMedias)
			}
		}
		qtdGeradosFortes = 0
		for qtdGeradosFortes < qtdFortes {
			randomNumber = rand.Intn(len(tabFortes))
			println("Gerou o numero forte ", tabFortes[randomNumber])
			if !lo.Contains(cartao, tabFortes[randomNumber]) {
				cartao = append(cartao, tabFortes[randomNumber])
				qtdGeradosFortes++
				println("Adicionou o numero forte ", tabFortes[randomNumber], " qtdFortes ", qtdGeradosFortes)
			}
		}

		sort.Ints(cartao)
		achouCartao = false
		for i := 0; i < len(cartoesGerados); i++ {
			if reflect.DeepEqual(cartoesGerados[i].Cartao, cartao) {
				achouCartao = true
			}
		}
		if !achouCartao {
			cartoesGerados = append(cartoesGerados, CartaoGerado{Cartao: cartao})
			qtdGerados++
			println("Adicionou o cartao ", cartao, " qtdGerados ", qtdGerados)
		}
	}

	saida := JsonSaida{
		Media:               media,
		Min:                 limiteMinimo,
		Max:                 limiteMaximo,
		DesvioPadrao:        desvio,
		NumerosSelecionados: pairs,
		NumerosFracos:       tabBaixos,
		NumerosMedios:       tabMedios,
		NumerosFortes:       tabFortes,
		Cartoes:             cartoesGerados,
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(saida)
}
