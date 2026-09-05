package types

import (
	"bitbucket.org/lyndus/backend/infra/utils"
	"errors"
)

type DOC string

var DOCError = errors.New("documento invalido")

//ParseDOC parse string to doc
func (doc *DOC) ParseDOC(docStr string) (err error) {
	limpo := utils.Limpar([]byte(docStr))
	if len(limpo) == 14 {
		if utils.ValidateCNPJ(limpo) {
			*doc = DOC(limpo)
			return nil
		}
	} else if len(limpo) == 11 {
		if utils.ValidateCPF(limpo) {
			*doc = DOC(limpo)
			return nil
		}
	}

	return DOCError
}

//MarshalJSON convert doc to string json
func (doc DOC) MarshalJSON() ([]byte, error) {
	if len(doc) == 11 {
		doc = "\"" + doc[:3] + "." + doc[3:6] + "." + doc[6:9] + "-" + doc[9:] + "\""
		return []byte(doc), nil
	} else if len(doc) == 14 {
		doc = "\"" + doc[:2] + "." + doc[2:5] + "." + doc[5:8] + "/" + doc[8:12] + "-" + doc[12:] + "\""
	}
	return []byte(doc), nil
}

//UnmarshalJSON convert string json to doc type
func (doc *DOC) UnmarshalJSON(b []byte) error {
	limpo := utils.Limpar(b)
	if len(limpo) == 14 {
		if utils.ValidateCNPJ(limpo) {
			*doc = DOC(limpo)
			return nil
		}
	} else if len(limpo) == 11 {
		if utils.ValidateCPF(limpo) {
			*doc = DOC(limpo)
			return nil
		}
	}
	return DOCError
}
