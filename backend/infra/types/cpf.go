package types

import (
	"bitbucket.org/lyndus/backend/infra/utils"
	"errors"
)

type CPF string

var CPFError = errors.New("cpf invalido")

func (cpf *CPF) ParseCPF(cpfStr string) (err error) {
	limpo := utils.Limpar([]byte(cpfStr))
	if utils.ValidateCPF(limpo) {
		*cpf = CPF(limpo)
		return nil
	}
	return CPFError
}

//MarshalJSON convert cpf to string json
func (cpf CPF) MarshalJSON() ([]byte, error) {
	return []byte("\"" + cpf[:3] + "." + cpf[3:6] + "." + cpf[6:9] + "-" + cpf[9:] + "\""), nil
}

//UnmarshalJSON convert string json to cpf type
func (cpf *CPF) UnmarshalJSON(b []byte) error {
	limpo := utils.Limpar(b)

	if utils.ValidateCPF(limpo) {
		*cpf = CPF(limpo)
		return nil
	}
	return CPFError
}
