package client

import (
	"bitbucket.org/lyndus/backend/internal/api"
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"bitbucket.org/lyndus/backend/domain/client"
	"bitbucket.org/lyndus/backend/domain/client/entity"
	"bitbucket.org/lyndus/backend/domain/client/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"bitbucket.org/lyndus/backend/infra/types"
	"github.com/go-chi/chi"
)

type RequestDomain struct {
	http.Request
}

func ImportEmployees(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	clientID, err := services.AuthService.CheckPermission(id, client.EMPLOYEE, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024) // 20 Mb
	err = r.ParseMultipartForm(20 * 1024 * 1024)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDMAXSIZE)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			//log.Println(err)
		}
	}(file)

	csvReader := csv.NewReader(file)

	line := 0
	var wrongLines []string
	var wrongCPFs []string
	var emptyNames []string
	var duplicateCPF []string
	var employees []entity.Employee
	//var errs []config.Err

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		line++

		if len(record) != 2 {
			wrongLines = append(wrongLines, strconv.Itoa(line))
		}

		if strings.ToLower(record[0]) == "nome" {
			continue
		}
		if strings.ToLower(record[0]) == "" {
			emptyNames = append(emptyNames, strconv.Itoa(line))
			continue
		}

		var cpf types.CPF
		err = cpf.ParseCPF(record[1])
		if err != nil {
			wrongCPFs = append(wrongCPFs, record[1])
			continue
		}

		for _, e := range employees {
			if e.CPF == cpf {
				duplicateCPF = append(duplicateCPF, strconv.Itoa(line))
			}
		}

		employee := entity.Employee{
			Name: record[0],
			CPF:  cpf,
		}
		employees = append(employees, employee)
	}
	var errors string = ""
	if len(wrongLines) > 1 || len(wrongCPFs) > 1 || len(emptyNames) > 1 || len(duplicateCPF) > 1 {

		if len(wrongLines) == 1 {
			//errs = append(errs, config.Err{Message: "Linha " + wrongLines[0] + " inválida."})
			errors += "Linha " + wrongLines[0] + " inválida.\n"
		} else if len(wrongLines) > 1 {
			//errs = append(errs, config.Err{Message: "Linhas " + strings.Join(wrongLines, ", ") + " inválidas."})
			errors += "Linhas " + strings.Join(wrongLines, ", ") + " inválidas." + "\n"
		}

		if len(wrongCPFs) == 1 {
			//errs = append(errs, config.Err{Message: "Linha " + wrongCPFs[0] + " com CPF inválido."})
			errors += "Linha " + wrongCPFs[0] + " com CPF inválido." + "\n"
		} else if len(wrongCPFs) > 1 {
			//errs = append(errs, config.Err{Message: "Linhas " + strings.Join(wrongCPFs, ", ") + " com CPFs, inválidos."})
			errors += "Linhas " + strings.Join(wrongCPFs, ", ") + " com CPFs, inválidos." + " inválidas." + "\n"
		}

		if len(emptyNames) == 1 {
			//errs = append(errs, config.Err{Message: "Linha " + emptyNames[0] + " com nome em branco."})
			errors += "Linha " + emptyNames[0] + " com nome em branco." + "\n"
		} else if len(emptyNames) > 1 {
			//errs = append(errs, config.Err{Message: "Linhas " + strings.Join(emptyNames, ", ") + " com nome em branco."})
			errors += "Linhas " + strings.Join(emptyNames, ", ") + " com nome em branco." + "\n"
		}

		if len(duplicateCPF) == 1 {
			//errs = append(errs, config.Err{Message: "CPF " + duplicateCPF[0] + " duplicado."})
			errors += "CPF " + duplicateCPF[0] + " duplicado." + "\n"
		} else if len(duplicateCPF) > 1 {
			//errs = append(errs, config.Err{Message: "CPFs " + strings.Join(duplicateCPF, ", ") + " duplicados."})
			errors += "CPFs " + strings.Join(duplicateCPF, ", ") + " duplicados." + "\n"
		}
	}

	if len(errors) > 0 {
		config.ResponsePerErr(w, err, config.CPFINVALID)
		return
	}

	for _, emp := range employees {
		emp.ClientID = clientID
		err = services.EmployeeService.GetEmployeeIDByCPF(&emp)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		emp.Active = true
		if emp.ID > 0 {
			err = services.EmployeeService.UpdateEmployeeActive(emp)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
			continue
		}

		per := entity.PersonBasic{
			CPF: emp.CPF,
		}

		err = services.PersonService.GetPersonIDByCPF(&per)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

		if per.ID == 0 {
			err = services.PersonService.CreatePersonBasic(&per)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
		}

		appUser := entity.AppUser{
			CPF:      emp.CPF,
			PersonID: per.ID,
		}

		err = services.AppUserService.GetAppUserIDByCPF(&appUser)

		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
		if appUser.ID == 0 {
			err = services.AppUserService.CreateAppUserBasic(&appUser)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
		}

		emp.PersonID = per.ID

		err = services.EmployeeService.CreateEmployee(&emp)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}

		category := entity.Category{
			ClientID: clientID,
		}

		err = services.CategoryService.GetDefaultCategory(&category)
		if err != nil {
			//log.Println(err)
		}
		err = services.EmployeeService.AddEmployeeCategory(emp, category.ID)
		if err != nil {
			//log.Println(err)
		}
	}

	w.WriteHeader(http.StatusOK)
	return

}

func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	clientID, err := services.AuthService.CheckPermission(id, client.EMPLOYEE, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var employeeCreate entity.EmployeeCreate
	err = json.NewDecoder(r.Body).Decode(&employeeCreate)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	employee := entity.Employee{
		ClientID: clientID,
		Name:     employeeCreate.Name,
		CPF:      employeeCreate.CPF,
	}

	err = services.EmployeeService.GetEmployeeIDByCPF(&employee)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if employee.ID != 0 {
		if !employee.Active {
			employee.Active = true
			err = services.EmployeeService.UpdateEmployeeActive(employee)
			if err != nil {
				config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
				return
			}
		} else {
			config.ResponsePerErr(w, err, config.CPFALREADYREGISTERED)
			return
		}

		resp := response.ID{
			ID: employee.ID,
		}

		config.JSONResponse(resp, http.StatusOK, w)
		return
	}

	person := entity.PersonBasic{
		CPF: employeeCreate.CPF,
	}

	err = services.PersonService.GetPersonIDByCPF(&person)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if person.ID == 0 {
		err = services.PersonService.CreatePersonBasic(&person)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
	}

	if person.ID == 0 {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	appUser := entity.AppUser{
		CPF:      employeeCreate.CPF,
		PersonID: person.ID,
	}

	err = services.AppUserService.GetAppUserIDByCPF(&appUser)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if appUser.ID == 0 {
		err = services.AppUserService.CreateAppUserBasic(&appUser)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
	}

	if appUser.ID == 0 {

		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return

	}
	employee.Active = true
	employee.PersonID = person.ID

	err = services.EmployeeService.CreateEmployee(&employee)

	if err != nil {
		//log.Println(err)
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}

	if len(employeeCreate.Categories) > 0 {
		for _, v := range employeeCreate.Categories {
			err := services.EmployeeService.AddEmployeeCategory(employee, v)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		category := entity.Category{
			ClientID: clientID,
		}

		err = services.CategoryService.GetDefaultCategory(&category)
		if err != nil {
			//log.Println(err)
		}
		err = services.EmployeeService.AddEmployeeCategory(employee, category.ID)
		if err != nil {
			//log.Println(err)
		}
	}

	resp := response.ID{
		ID: employee.ID,
	}

	config.JSONResponse(resp, http.StatusOK, w)

}
func UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.EMPLOYEE, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}
	var employee entity.Employee
	err = json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	employee.ClientID = clientID
	log.Println(employee)
	err = services.EmployeeService.UpdateEmployee(employee)
	if err != nil {
		//log.Println(err)
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	// Add or Remove categories
	// pegar todas as categorias
	// comparar com as que estão na requisição
	// se esta na requisição ignora , se não esta remove

	// usar um hash map para validar
	var categoriasMap = make(map[int]int)
	// fazer um get para categorias de um funcionario
	var tempEmployee = entity.Employee{ID: employee.ID}
	tempEmployee.Categories, err = services.EmployeeService.GetEmployeeCategories(tempEmployee.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	var employeesCategories = tempEmployee.Categories
	var requisiteCategories = employee.Categories
	const (
		REMOVER   = 0
		ADICIONAR = 1
		IGNORE    = 3
	)
	// Adiciono todas os id da requisição como para adicionar
	for i := range requisiteCategories {
		categoriasMap[requisiteCategories[i].ID] = ADICIONAR
	}

	// Valido se o id da categoria ja não esta adicionado então IGNORE
	// Caso o id não esteja na requisição setamos como REMOVER

	for i := range employeesCategories {
		// Estou testando o valor default
		if categoriasMap[employeesCategories[i].ID] == REMOVER {
			// Agora seto para que realmente existe um [:id] = REMOVER ,pois antes era undefined
			categoriasMap[employeesCategories[i].ID] = REMOVER
		} else {
			categoriasMap[employeesCategories[i].ID] = IGNORE
		}
	}

	for idCategory := range categoriasMap {
		if categoriasMap[idCategory] == ADICIONAR {
			err = services.EmployeeService.AddCategory(employee.ID, idCategory)
		} else if categoriasMap[idCategory] == REMOVER {
			err = services.EmployeeService.RemoveCategory(employee.ID, idCategory)
		}
	}

	w.WriteHeader(http.StatusOK)
	return
}

func EmployeeAddCategory(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.EMPLOYEE, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	employeeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	var categoryID response.ID
	err = json.NewDecoder(r.Body).Decode(&categoryID)

	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	employee := entity.Employee{
		ID:       employeeID,
		ClientID: clientID,
		Active:   true,
	}

	empID, err := services.EmployeeService.GetEmployeeID(employee)
	//empID, err := employee.GetEmployeeID()
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	if empID.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	categoryID, err = services.CategoryService.GetCategoryIDByID(clientID, categoryID.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	if categoryID.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	err = services.EmployeeService.AddEmployeeCategory(employee, categoryID.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	w.WriteHeader(http.StatusOK)
	return
}

func EmployeeDelCategory(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.EMPLOYEE, client.ALL)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	employeeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var categoryID response.ID
	err = json.NewDecoder(r.Body).Decode(&categoryID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	employee := entity.Employee{
		ID:       employeeID,
		ClientID: clientID,
		Active:   true,
	}

	empID, err := services.EmployeeService.GetEmployeeID(employee)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	if empID.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	categoryID, err = services.CategoryService.GetCategoryIDByID(clientID, categoryID.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	if categoryID.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	err = services.EmployeeService.RemoveEmployeeCategory(employee, categoryID.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	w.WriteHeader(http.StatusOK)
	return
}

func GetEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.EMPLOYEE, client.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	employeeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var employee entity.Employee
	employee.ID = employeeID
	employee.ClientID = clientID
	//employee.StatusID = client.EMPLOYEESTATUSACTIVE

	err = services.EmployeeService.GetEmployee(&employee)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if employee.CPF == "" {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	employee.Categories, err = services.EmployeeService.GetEmployeeCategories(employee.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if employee.Categories == nil {
		employee.Categories = make([]entity.Category, 0)
	}

	config.JSONResponse(employee, http.StatusOK, w)
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.EMPLOYEE, client.ALL)
	if err != nil {

		config.ResponsePerErr(w, err, config.UNAUTHORIZED)

		return
	}

	employeeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var employee entity.Employee
	employee.ID = employeeID
	employee.ClientID = clientID

	empID, err := services.EmployeeService.GetEmployeeID(employee)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	if empID.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	employee.Active = false
	err = services.EmployeeService.DeleteEmployee(clientID, employeeID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	w.WriteHeader(http.StatusOK)
	return
}

func ActiveEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.EMPLOYEE, client.ALL)
	if err != nil {

		config.ResponsePerErr(w, err, config.UNAUTHORIZED)

		return
	}

	employeeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var employee entity.Employee
	employee.ID = employeeID
	employee.ClientID = clientID

	empID, err := services.EmployeeService.GetEmployeeID(employee)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	if empID.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	err = services.EmployeeService.ActiveInactiveEmployee(clientID, employeeID, true)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	w.WriteHeader(http.StatusOK)
	return
}

func InactiveEmployee(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.EMPLOYEE, client.ALL)
	if err != nil {

		config.ResponsePerErr(w, err, config.UNAUTHORIZED)

		return
	}

	employeeID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	var employee entity.Employee
	employee.ID = employeeID
	employee.ClientID = clientID

	empID, err := services.EmployeeService.GetEmployeeID(employee)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	if empID.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	err = services.EmployeeService.ActiveInactiveEmployee(clientID, employeeID, false)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}
	w.WriteHeader(http.StatusOK)
	return
}

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	clientID, err := services.AuthService.CheckPermission(id, client.EMPLOYEE, client.VIEW)
	if err != nil {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var employeesResponse api.ListResponse
	employeesResponse.Active = true

	employeesResponse.ActivePage, _ = strconv.Atoi(r.URL.Query().Get("page"))
	itemsPerPage, _ := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	search := r.URL.Query().Get("search")
	orderBy := r.URL.Query().Get("orderBy")
	sortDesc, _ := strconv.ParseBool(r.URL.Query().Get("sortDesc"))
	activeSTR := r.URL.Query().Get("active")
	if activeSTR != "" {
		employeesResponse.Active, _ = strconv.ParseBool(activeSTR)
	}

	employeesResponse.Items, employeesResponse.TotalItems, employeesResponse.TotalPages, err =
		services.EmployeeService.GetEmployeesV2(
			employeesResponse.ActivePage,
			itemsPerPage,
			clientID,
			search,
			orderBy,
			sortDesc,
			employeesResponse.Active)

	if err != nil {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	config.JSONResponse(employeesResponse, http.StatusOK, w)
}
