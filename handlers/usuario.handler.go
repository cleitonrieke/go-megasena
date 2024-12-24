package handlers

import (
	"backlotofacil/models"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type retornoToken struct {
	Token string `json:"token"`
}

func (h Handler) CriarUsuario(w http.ResponseWriter, r *http.Request) {

	var usuario models.Usuario

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if (usuario.UserId == "") || (usuario.Password == "") {
		http.Error(w, "UserId e Password devem ser informados", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(usuario.Password), 10)
	usuario.Password = string(hash)

	res := h.DB.Create(&usuario)

	if res.Error != nil {
		http.Error(w, res.Error.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h Handler) Autenticar(w http.ResponseWriter, r *http.Request) {
	println("dentro do autenticador...")
	var usuario models.Usuario

	err := json.NewDecoder(r.Body).Decode(&usuario)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	println("recebido o usuario: ", usuario.UserId)

	if (usuario.UserId == "") || (usuario.Password == "") {
		http.Error(w, "UserId e Password devem ser informados", http.StatusBadRequest)
		return
	}

	var usuarioBanco models.Usuario
	h.DB.First(&usuarioBanco, "user_id = ?", usuario.UserId)
	if usuarioBanco.Id == 0 {
		http.Error(w, "Usuário ou Senha inválidos", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(usuarioBanco.Password), []byte(usuario.Password))
	if err != nil {
		http.Error(w, "Usuário ou Senha inválidos", http.StatusBadRequest)
		return
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usuarioBanco.Id,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(retornoToken{
		Token: tokenString})
}
