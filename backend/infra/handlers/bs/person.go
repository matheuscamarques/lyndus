package bs

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"

	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/logger"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"bitbucket.org/lyndus/backend/infra/types"
	"bitbucket.org/lyndus/backend/infra/utils"
	"go.uber.org/zap"
)

//SavePerson create new bs_person
func SavePerson(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BANNER, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var person entity.Person
	err = json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		log.Println(err)
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	person.BsID = bsID

	personDB, err := services.PersonService.FindPerson(bsID, person.CPF, person.Phone)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	var appUserID int

	appUserID, err = services.AppUserService.GetAppUserIDByCPF(person.CPF)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if appUserID != 0 {
		person.AppUserID = &appUserID
	}

	if personDB.ID != 0 {

		if person.Phone == personDB.Phone {
			config.ResponsePerErr(w, err, config.PHONEALREADYREGISTERED)
			return
		} else if *person.CPF == *personDB.CPF {
			config.ResponsePerErr(w, err, config.CPFALREADYREGISTERED)
			return
		}
		//
		//if bs_person.CPF != personDB.CPF && bs_person.Phone == personDB.Phone {
		//	bs_person.ID = personDB.ID
		//	if appUserID != 0 {
		//		//  atualizar cpf e app_user_id
		//		bs_person.AppUserID = &appUserID
		//		err = services.PersonService.UpdatePersonCPFAndAppUser(bsID, personDB.ID, bs_person.CPF, appUserID)
		//		if err != nil {
		//			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		//			return
		//		}
		//	} else {
		//		// atualizar cpf
		//		err = services.PersonService.UpdatePersonCPF(bsID, bs_person.ID, bs_person.CPF)
		//		if err != nil {
		//			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		//			return
		//		}
		//	}
		//} else {
		//	if bs_person.CPF == ""{
		//		config.ResponsePerErr(w, err, config.PHONEALREADYREGISTERED)
		//		return
		//	}
		//	config.ResponsePerErr(w, err, config.CPFALREADYREGISTERED)
		//	return
		//}
	}

	person.ID, err = services.PersonService.CreatePerson(person)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	resp := response.ID{
		ID: person.ID,
	}
	config.JSONResponse(resp, http.StatusOK, w)
}

//UpdatePerson update bs_person
func UpdatePerson(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PERSON, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var person entity.Person
	err = json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var personDB entity.Person
	//repo := repository.NewPersonRepository()

	personDB, err = services.PersonService.GetPersonID(bsID, person.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if personDB.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	if personDB.AppUserID != nil && *personDB.AppUserID == 0 {

		appUserID, err := services.AppUserService.GetAppUserIDByCPF(person.CPF)
		if err != nil {
			//log.Println(err)
		}
		if appUserID != 0 {
			person.AppUserID = &appUserID
		}
	}

	person.BsID = bsID
	err = services.PersonService.UpdatePerson(person)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func GetPersons(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PERSON, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	//todo pagination
	search := r.URL.Query().Get("search")
	//repo := repository.NewPersonRepository()

	var persons []entity.Person

	if search != "" {
		var cpf types.CPF

		phone := utils.OnlyDigits(search)

		err = cpf.ParseCPF(search)
		if err == nil {

			persons, err = services.PersonService.GetPersonsByCPF(bsID, cpf)
			if len(persons) == 0 {

				appUser, err := services.AppUserService.GetAppUserByCPF(cpf)
				if err != nil {
					logger.Error("Erro ao buscar app bs_user:", zap.String("error", err.Error()))
				}
				if appUser.ID != 0 {
					var person entity.Person
					person.AppUserID = &appUser.ID
					person.Name = appUser.Name
					person.Phone = appUser.Cellphone
					person.Birthdate = &appUser.Birthdate

					person.ID, err = services.PersonService.CreatePerson(person)
					if err != nil {
						logger.Error("Erro ao salvar dados bs_person from app bs_user:", zap.String("error", err.Error()))
					}

					persons = []entity.Person{
						{
							ID:        person.ID,
							Name:      person.Name,
							CPF:       person.CPF,
							Phone:     person.Phone,
							Birthdate: person.Birthdate,
						},
					}
				}
			}
		} else if len(phone) > 5 {
			persons, err = services.PersonService.GetPersonsByPhone(bsID, phone)
			if err != nil {
				logger.Error("Erro ao listar bs_person por telefone:", zap.String("error", err.Error()))
			}
		} else {

			persons, err = services.PersonService.GetPersonsByName(bsID, search)

			if err != nil {
				logger.Error("Erro ao bs_person por nome:", zap.String("error", err.Error()))
			}
		}
	} else {
		persons, err = services.PersonService.GetPersons(bsID)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

	}

	config.JSONResponse(persons, http.StatusOK, w)
}

func SearchPerson(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PERSON, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	name := chi.URLParam(r, "name")

	var persons []entity.Person
	persons, err = services.PersonService.GetPersonsByName(bsID, name)
	if err != nil {
		logger.Error("Erro ao bs_person por nome:", zap.String("error", err.Error()))
	}
	log.Println(persons, name)
	if persons == nil {
		persons = make([]entity.Person, 0)
	}

	config.JSONResponse(persons, http.StatusOK, w)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	personID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || personID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	bsID, err := services.AuthService.CheckPermission(id, bs.SERVICE, bs.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	person, err := services.PersonService.GetPerson(bsID, personID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	config.JSONResponse(person, http.StatusOK, w)
}
