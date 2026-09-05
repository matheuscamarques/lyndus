package types

import (
	"errors"

	"bitbucket.org/lyndus/backend/infra/utils"
)

type CNPJ string

var CNPJError = errors.New("cnpj invalido")

//ParseCNPJ parse string to CNPJ
func (cnpj *CNPJ) ParseCNPJ(cnpjStr string) (err error) {
	limpo := utils.Limpar([]byte(cnpjStr))
	if utils.ValidateCNPJ(limpo) {
		*cnpj = CNPJ(limpo)
		return nil
	}
	return CNPJError
}

//MarshalJSON convert cnpj to string json
func (cnpj CNPJ) MarshalJSON() ([]byte, error) {
	cnpj = "\"" + cnpj[:2] + "." + cnpj[2:5] + "." + cnpj[5:8] + "/" + cnpj[8:12] + "-" + cnpj[12:] + "\""
	return []byte(cnpj), nil
}

//UnmarshalJSON convert string json to cnpj type
func (cnpj *CNPJ) UnmarshalJSON(b []byte) error {
	limpo := utils.Limpar(b)
	if utils.ValidateCNPJ(limpo) {
		*cnpj = CNPJ(limpo)
		return nil
	}
	return CNPJError
}
