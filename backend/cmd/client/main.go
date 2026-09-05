package main

import (
	"log"
	"net/http"

	"bitbucket.org/lyndus/backend/domain/client/services"
	"bitbucket.org/lyndus/backend/infra/auth"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/db/postgres"
	"bitbucket.org/lyndus/backend/infra/handlers/client"
	"bitbucket.org/lyndus/backend/infra/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

var (
	err error
)

func main() {
	logger.New("Client")
	logger.Info("System UP")

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(logger.APILogger())
	//r.Use(render.SetContentType(render.ContentTypeJSON))

	err = config.ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	postgres.DB, err = postgres.InitDB(config.Config.PostgresHost, config.Config.PostgresUser, config.Config.PostgresPass, config.Config.PostgresDatabase, config.Config.PostgresPort)
	if err != nil {
		log.Fatal(err)
	}

	err = auth.InitKeysPem()
	if err != nil {
		log.Fatal(err)
	}

	auth.Service = auth.NewAuthenticationService()
	services.AuthService = services.NewClientAuthService()

	services.UserService = services.NewClientUserService()
	services.PersonService = services.NewClientPersonService()
	services.EmployeeService = services.NewClientEmployeeService()
	services.CompanyService = services.NewClientCompanyService()
	services.AppUserService = services.NewClientAppUserService()
	services.BenefitService = services.NewClientBenefitService()
	services.CategoryService = services.NewClientCategoryService()

	services.Poll = services.NewClientPollService()
	services.PollQuestion = services.NewClientPollQuestionService()
	services.PollCategory = services.NewClientPollCategoryService()
	services.PollChoice = services.NewClientPollChoiceService()

	cors := cors.New(cors.Options{
		//AllowOriginFunc:  AllowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	//workDir, _ := os.Getwd()
	//filesDir := filepath.Join(workDir, "dist")
	//config.FileServer(r, "/", api.Dir(filesDir))

	// Routes
	//r.Get("/", hello)
	r.Post("/login", client.LoginHandler)
	r.Post("/request-recovery", client.RequestRecovery)
	r.Post("/confirm-recovery", client.ConfirmRecovery)
	r.Post("/recovery", client.Recovery)
	r.Route("/", func(r chi.Router) {
		r.Use(auth.ValidateTokenMiddleware)

		r.Get("/permissions", client.GetPermissions)

		//TODO SÓ MASTER ALTERA MASTER
		//todo não cadastrar lista permissões vazia.
		r.Post("/user", client.CreateUser)
		r.Put("/user", client.UpdateUser)
		r.Get("/user/{id}", client.GetUser)
		r.Put("/user/{id}/active", client.ActiveUser)
		r.Put("/user/{id}/inactive", client.InactiveUser)
		r.Delete("/user/{id}", client.DeleteUser)

		r.Get("/user", client.GetAllUser)
		r.Get("/user/{id}/permissions", client.GetUserPermissions)

		r.Post("/password", client.ChangePassword)

		r.Get("/profile", client.Profile)
		r.Put("/profile", client.UpdateProfile)

		r.Post("/employee", client.CreateEmployee)
		r.Put("/employee", client.UpdateEmployee)
		r.Get("/employee/{id}", client.GetEmployee)
		r.Put("/employee/{id}/active", client.ActiveEmployee)
		r.Put("/employee/{id}/inactive", client.InactiveEmployee)
		r.Delete("/employee/{id}", client.DeleteEmployee)

		r.Get("/employee", client.GetEmployees)
		r.Put("/employee/{id}/category", client.EmployeeAddCategory)
		r.Delete("/employee/{id}/category", client.EmployeeDelCategory)
		r.Post("/employee/import", client.ImportEmployees)

		r.Post("/category", client.CreateCategory)
		r.Put("/category", client.UpdateCategory)
		r.Get("/category", client.GetCategories)
		r.Get("/category/list", client.GetCategoryList)
		r.Get("/category/{id}", client.GetCategory)
		r.Delete("/category/{id}", client.DeleteCategory)
		r.Put("/category/{id}/active", client.ActiveCategory)
		r.Put("/category/{id}/inactive", client.InactiveCategory)

		r.Post("/benefit", client.CreateBenefit)
		r.Put("/benefit", client.UpdateBenefit)
		r.Put("/benefit/cancel", client.BenefitCancel)
		r.Put("/benefit/requestTicket", client.BenefitRequestTicket)
		r.Put("/benefit/confirmPayment", client.BenefitConfirmPayment)

		// r.Put("/benefit/payment", client.BenefitPayment)
		r.Get("/benefit", client.GetBenefits)
		r.Get("/benefit/{id}", client.GetBenefit)
		r.Get("/benefit/{id}/downloadTicket", client.BenefitDownloadTicket)
		r.Put("/benefit/{id}/benefit-user", client.UpdateBenefits)

		r.Post("/poll", client.CreatePoll)
		r.Get("/poll", client.GetPollAll)
		r.Put("/poll", client.PutPoll)
		r.Get("/poll/{id}", client.GetPollById)

		r.Post("/poll/{id}/question", client.PullAddQuestion)
		r.Get("/poll/{id}/question", client.GetAllPullQuestions)
		r.Get("/poll/{id}/question/{questionID}", client.GetPullQuestion)
		r.Put("/poll/{id}/question", client.PutPullQuestion)

	})

	//log.Println("Running... port:", config.Config.Port)
	// Start server
	err = http.ListenAndServe(":"+config.Config.Port, r)
	if err != nil {
		log.Fatal(err)
	}

}
