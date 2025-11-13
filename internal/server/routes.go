package server

import (
	"encoding/json"
	"net/http"

	repo "github.com/Hajdudev/invoice-flow/internal/adapters/postgresql/sqlc"
	"github.com/Hajdudev/invoice-flow/internal/users"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", s.healthHandler)
	userService := users.NewService(repo.New(s.db.Pool()), s.conn)
	userHandler := users.NewHandler(userService)
	r.Get("/register", userHandler.RegisterUser)

	// auth
	// r.Route("/auth", func(r chi.Router) {
	// 	// Routes for OAuth providers
	// 	r.Get("/{provider}", s.authProviderCallback)
	// 	r.Get("/{provider}/callback", s.getAuthCallback)
	// 	// Login and Signup
	// 	r.With(
	// 		routesmiddleware.ValidateRequest[request_types.UserLogin],
	// 	).Post("/login", s.login)
	// 	r.With(
	// 		routesmiddleware.ValidateRequest[request_types.UserSignup],
	// 	).Post("/login", s.signup)
	// })

	return r
}

type ValidationErrors struct {
	Errors map[string][]string `json:"errors"`
}

// func (s *Server) signup(w http.ResponseWriter, r *http.Request) {
// 	req, ok := routesmiddleware.GetDecodedRequest[request_types.UserSignup](r)
// 	if !ok {
// 		utils.Error(w, http.StatusInternalServerError, "Could not retrieve request from context")
// 		return
// 	}
// 	fmt.Println(req)
// }
//
// func (s *Server) login(w http.ResponseWriter, r *http.Request) {
// 	req, ok := routesmiddleware.GetDecodedRequest[request_types.UserLogin](r)
// 	if !ok {
// 		utils.Error(w, http.StatusInternalServerError, "Could not retrieve request from context")
// 		return
// 	}
// 	fmt.Println(req)
// }

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

// func (s *Server) getAuthCallback(w http.ResponseWriter, r *http.Request) {
// 	provider := chi.URLParam(r, "provider")
// 	r = r.WithContext(context.WithValue(r.Context(), contextkeys.Provider, provider))
//
// 	user, err := gothic.CompleteUserAuth(w, r)
// 	if err != nil {
// 		log.Printf("gothic.CompleteUserAuth failed: %v", err)
// 		utils.Error(w, http.StatusInternalServerError, "Failed to complete authentication")
// 	}
// 	fmt.Println(user)
// }
//
// func (s *Server) authProviderCallback(w http.ResponseWriter, r *http.Request) {
// 	provider := chi.URLParam(r, "provider")
// 	ctx := context.WithValue(r.Context(), contextkeys.Provider, provider)
// 	r = r.WithContext(ctx)
// 	gothic.BeginAuthHandler(w, r)
// }
