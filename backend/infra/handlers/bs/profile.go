package bs

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"

	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/repository"
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/logger"
	"bitbucket.org/lyndus/backend/infra/utils"
	"go.uber.org/zap"
)

func Profile(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BANNER, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var bsCompany entity.Company
	repo := repository.NewBSCompanyRepository()
	bsCompany, err = repo.GetCompanyByBsID(bsID)
	if err == sql.ErrNoRows {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	bsCompany.WeekDays, err = repo.GetBSWeekDays(bsID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(bsCompany, http.StatusOK, w)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PROFILE, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var bsCompany entity.Company
	err = json.NewDecoder(r.Body).Decode(&bsCompany)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	//bs.ID id of company not BS.
	if bsCompany.ID != bsID {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	repo := repository.NewBSCompanyRepository()

	err = repo.UpdateCompanyByBS(bsCompany)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	if len(bsCompany.WeekDays) > 0 {
		err = repo.DeleteBSAllWeekDays(bsID)
		if err != nil {
			logger.Error("Erro ao deletar weekday do bs:", zap.String("error", err.Error()))
		}
		for _, weekDay := range bsCompany.WeekDays {
			err = repo.AddBSWeekDay(bsID, weekDay)
			if err != nil {
				logger.Error("Erro ao adicionar weekday do bs:", zap.String("error", err.Error()))
			}
		}
	}

	w.WriteHeader(http.StatusOK)
	return
}

//SaveProfileLogo func set profile Logo
func SaveProfileLogo(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PROFILE, bs.ALL)
	if err != nil || bsID == 0 {
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
		if err == http.ErrMissingFile {
			config.ResponsePerErr(w, err, config.FILEINVALID)
			return
		}
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	defer file.Close()

	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {

		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	fileType := http.DetectContentType(buff)
	if fileType != "image/jpeg" && fileType != "image/png" {
		config.ResponsePerErr(w, err, config.FILEINVALID)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	img, _, err := image.Decode(file)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	filePath := utils.GetPathFile(bsID, bsID, config.Config.BaseStaticPath, modulePath, "logo", "png")

	log.Println(filePath)
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		panic(err)
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	//_, err = io.Copy(f, file)
	//if err != nil {
	//	//log.Println(err)
	//	errs := []config.Err{
	//		{
	//			Message: "Erro interno no sistema, se persistir contatar suporte.",
	//			Code:    "",
	//		},
	//	}
	//	config.JSONResponse(config.ResponseErr{Errs: errs}, api.StatusInternalServerError, w)
	//	return
	//}

	err = r.MultipartForm.RemoveAll()
	if err != nil {
		//panic(err)
		logger.Error("Erro ao deletar Multi Part Form:", zap.String("error", err.Error()))
	}

	w.WriteHeader(http.StatusOK)
	return
}

func GetProfileLogo(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PROFILE, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	filePath := utils.GetPathFile(bsID, bsID, config.Config.BaseStaticPath, modulePath, "logo", "png")
	f, err := os.Open(filePath)
	if os.IsNotExist(err) {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	fi, err := f.Stat()
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
	w.Header().Set("Content-Disposition", "attachment; filename=logo.png")
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))

	//stream the body to the client without fully loading it into memory
	_, err = io.Copy(w, f)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
}
