package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

// MAKE A CLI FOR CREATE FILES TO DOMAIN
// Example go run command.go create client/client_poll
// go to caseUse.Folder domain/client
// go to caseUse.Folder domain/client/entity and create a file client_poll.go
// go to caseUse.Folder domain/client/contracts  create a file client_poll_repository_interface.go
// go to caseUse.Folder domain/client/contracts  create a file client_poll_service_interface.go
// go to caseUse.Folder domain/client/repository create a file client_poll_repository.go
// go to caseUse.Folder domain/client/bs_service create a file client_poll_service.go

const (
	DOMAIN_ROOT      = "../domain/"
	LAYER_CONTRACTS  = "contracts/"
	LAYER_SERVICE    = "services/"
	LAYER_REPOSITORY = "repository/"
	LAYER_ENTITY     = "entity/"
)

type types uint8

const (
	INTERFACES types = iota
	ENTITY
	REPOSITORY
	SERVICE
)

type Syntactic struct {
	Package string
	Name    string
	Type    string
	Methods []string
}

const templateT = `package {{.Package}}

type {{.Name}} {{.Type}} {
	{{if .Methods}}{{range .Methods}}{{.}}
	{{end}}{{end}}//TODO Validar	
}
`

type CaseUse struct {
	DDDLayer string
	Folder   string
	Methods  []string
}

//go run command.go create client/client_poller 15535
//go run command.go create client/salao "ID int, Nome string, Bruno string, ID_client int"

func main() {
	if len(os.Args) > 1 {
		// geral
		err := Command(os.Args)
		if err != nil {
			return
		}
	}
}

func Command(commands []string) error {
	mode := commands[1]
	command := commands[2]
	methods := []string{}

	//fmt.Println(commands)
	if len(commands) > 3 {
		methods = strings.Split(commands[3], ",")
		for i := range methods {
			methods[i] = strings.TrimSpace(methods[i])
		}
	}

	switch mode {
	case "create":
		CreateMode(command, methods)
	}

	return nil
}

func CreateMode(command string, methods []string) {

	// go run command.go create client/client_poller "ID int, Nome string"
	s := strings.Split(command, "/")
	var caseUse CaseUse
	domain := s[0]
	caseUse.DDDLayer = s[1]

	caseUse.Folder = DOMAIN_ROOT + domain + "/" + LAYER_CONTRACTS
	Make(INTERFACES, caseUse)

	caseUse.Folder = DOMAIN_ROOT + domain + "/" + LAYER_REPOSITORY
	Make(REPOSITORY, caseUse)

	caseUse.Folder = DOMAIN_ROOT + domain + "/" + LAYER_SERVICE
	Make(SERVICE, caseUse)

	caseUse.Folder = DOMAIN_ROOT + domain + "/" + LAYER_ENTITY
	caseUse.Methods = methods
	Make(ENTITY, caseUse)
}

func Make(option types, caseUse CaseUse) error {
	switch option {
	case INTERFACES:
		MakeInterfaceRepository(caseUse)
		MakeInterfaceService(caseUse)
		return nil
	case ENTITY:
		MakeEntity(caseUse)
		return nil
	case SERVICE:
		MakeService(caseUse)
		return nil
	case REPOSITORY:
		MakeRepository(caseUse)
		return nil
	}
	return fmt.Errorf("option not founded")
}

func createFile(fileName string) (*os.File, error) {
	f, err := os.Create(fileName)
	if err != nil {
		return f, err
	}

	err = f.Chmod(0644)
	if err != nil {
		return f, err
	}

	return f, err
}

func WriteFileSyntax(syntax Syntactic, f *os.File) error {
	tmpl, err := template.New("parse").Parse(templateT)
	if err != nil {
		return err
	}

	err = tmpl.Execute(f, syntax)
	if err != nil {
		return err
	}

	return nil
}

func CreateTemplate(repoFileName string, syntax Syntactic) error {
	f, err := createFile(repoFileName)
	if err != nil {
		return err
	}
	defer f.Close()

	err = WriteFileSyntax(syntax, f)
	if err != nil {
		return err
	}

	fmt.Println("file created: " + repoFileName)

	return nil
}

func MakeService(caseUse CaseUse) error {
	repoFileName := caseUse.DDDLayer + "_service" + ".go"
	repoFileName = caseUse.Folder + repoFileName

	if _, err := os.Stat(repoFileName); !os.IsNotExist(err) {
		return fmt.Errorf("file exists: " + repoFileName)
	}

	syntax := Syntactic{
		Package: "services",
		Name:    convertCommonCase(caseUse.DDDLayer + "_service"),
		Type:    "struct",
		Methods: caseUse.Methods,
	}

	err := CreateTemplate(repoFileName, syntax)
	if err != nil {
		return err
	}

	return nil
}

func MakeRepository(caseUse CaseUse) error {
	repoFileName := caseUse.DDDLayer + "_repository" + ".go"
	repoFileName = caseUse.Folder + repoFileName

	if _, err := os.Stat(repoFileName); !os.IsNotExist(err) {
		return fmt.Errorf("file exists: " + repoFileName)
	}

	syntax := Syntactic{
		Package: "repository",
		Name:    convertCommonCase(caseUse.DDDLayer + "_repository"),
		Type:    "struct",
		Methods: caseUse.Methods,
	}

	err := CreateTemplate(repoFileName, syntax)
	if err != nil {
		return err
	}

	return nil
}

func MakeEntity(caseUse CaseUse) error {
	repoFileName := caseUse.DDDLayer + ".go"
	repoFileName = caseUse.Folder + repoFileName

	if _, err := os.Stat(repoFileName); !os.IsNotExist(err) {
		return fmt.Errorf("file exists: " + repoFileName)
	}

	syntax := Syntactic{
		Package: "entity",
		Name:    convertCommonCase(caseUse.DDDLayer),
		Type:    "struct",
		Methods: caseUse.Methods,
	}

	err := CreateTemplate(repoFileName, syntax)
	if err != nil {
		return err
	}

	return nil
}

func MakeInterfaceRepository(caseUse CaseUse) error {
	repoFileName := caseUse.DDDLayer + "_" + "repository" + "_" + "interface" + ".go"
	repoFileName = caseUse.Folder + repoFileName

	if _, err := os.Stat(repoFileName); !os.IsNotExist(err) {
		return fmt.Errorf("file exists: " + repoFileName)
	}

	syntax := Syntactic{
		Package: "contracts",
		Name:    convertCommonCase(caseUse.DDDLayer + "_repository" + "_interface"),
		Type:    "interface",
		Methods: caseUse.Methods,
	}

	err := CreateTemplate(repoFileName, syntax)
	if err != nil {
		return err
	}

	return nil
}

func MakeInterfaceService(caseUse CaseUse) error {
	repoFileName := caseUse.DDDLayer + "_" + "bs_service" + "_" + "interface" + ".go"
	repoFileName = caseUse.Folder + repoFileName
	if _, err := os.Stat(repoFileName); !os.IsNotExist(err) {
		return fmt.Errorf("file exists: " + repoFileName)
	}

	syntax := Syntactic{
		Package: "contracts",
		Name:    convertCommonCase(caseUse.DDDLayer + "_service" + "_interface"),
		Type:    "interface",
		Methods: caseUse.Methods,
	}

	err := CreateTemplate(repoFileName, syntax)
	if err != nil {
		return err
	}

	return nil
}

func convertCommonCase(str string) string {
	s := strings.ReplaceAll(str, "_", " ")
	s = strings.Title(s)
	s = strings.ReplaceAll(s, " ", "")

	return s
}
