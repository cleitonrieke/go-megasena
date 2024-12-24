package main

import (
	"backlotofacil/db"
	"backlotofacil/handlers"
	"backlotofacil/middleware"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	godotenv.Load()
	r := mux.NewRouter()
	db := db.Init()
	h := handlers.Handler{DB: db}
	middleware := middleware.HandlerMiddleware{DB: db}
	r.HandleFunc("/convites", middleware.ValidarJWT(h.GetAllConvites)).Methods("GET")
	r.HandleFunc("/convites/gerar", middleware.ValidarJWT(h.GerarConvites)).Methods("GET")
	r.HandleFunc("/convites/validar-convite", h.ValidarConvite).Methods("GET")
	r.HandleFunc("/participantes", middleware.ValidarJWT(h.GetAllParticipantes)).Methods("GET")
	r.HandleFunc("/participantes/registrar", h.Registrar).Methods("POST")
	r.HandleFunc("/usuarios/cadastrar", h.CriarUsuario).Methods("POST")
	r.HandleFunc("/usuarios/autenticar", h.Autenticar).Methods("POST")
	r.HandleFunc("/jogos/listar-numeros", middleware.ValidarJWT(h.ResumoNumerosLotoFacil)).Methods("GET")

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:9000", "http://localhost:9001", "https://quasar-mega-sena.onrender.com/"}, //you service is available and allowed for this base url
		AllowedMethods: []string{
			http.MethodGet, //http methods for your app
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},

		AllowedHeaders: []string{
			"*", "authentication", //or you can your header key values which you are using in your application

		},
	})

	http.ListenAndServe(":8080", corsOpts.Handler(r))
}
