package bs

import (
	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetBankInfo(w http.ResponseWriter, r *http.Request){
	id := r.Context().Value("id").(int)

	bsID, err := services.AuthService.CheckPermission(id, bs.BANK, bs.VIEW)

	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	bank,err := services.BS.GetBankDetails(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if bank.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	config.JSONResponse(bank,http.StatusOK,w)
	return
}

func SaveBankInfo(w http.ResponseWriter, r *http.Request){
	id := r.Context().Value("id").(int)

	bsID, err := services.AuthService.CheckPermission(id, bs.BANK, bs.ALL)
	fmt.Println(bsID)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}


	bank,err := services.BS.GetBankDetails(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&bank)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if bank.ID == 0 {
		// CREATE
		bank.BsID = bsID
		err := services.BS.CreateBankDetails(bank)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
	}else{
		err := services.BS.UpdateBankDetails(bank)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
	}



	config.JSONResponse(bank,http.StatusOK,w)
	return
}





