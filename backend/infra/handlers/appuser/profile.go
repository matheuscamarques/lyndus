package appuser

import (
	"bitbucket.org/lyndus/backend/domain/appuser/entity"
	"bitbucket.org/lyndus/backend/domain/appuser/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/logger"
	"bitbucket.org/lyndus/backend/infra/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
)

func GetProfile(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	appUser, err := services.UserUserService.GetAppUserFullByID(id)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	gender, err := services.UserGenderService.GetGenderByID(appUser.GenderID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	appUser.GenderName = gender.Desc
	config.JSONResponse(appUser, http.StatusOK, w)
	return
}

func PutProfile(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	var appUser entity.AppUserCreate
	err := json.NewDecoder(r.Body).Decode(&appUser)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	appUser.ID = id

	err = services.UserUserService.UpdateAppUserProfile(appUser)
	if err != nil {
		//todo melhorar/ajustar tratamento de erro.
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func PostPhoto(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	var err error
	r.Body = http.MaxBytesReader(w, r.Body, 20*1024*1024) // 20 Mb
	err = r.ParseMultipartForm(20 * 1024 * 1024)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDMAXSIZE)
		return
	}

	file, _, err := r.FormFile("photo")
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

	filePath := utils.GetPathFile(id, id, config.Config.BaseStaticPath, modulePath, "profile_", "png")

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

func GetPhoto(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	filePath := utils.GetPathFile(id, id, config.Config.BaseStaticPath, modulePath, "profile_", "png")
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
	w.Header().Set("Content-Disposition", "attachment; filename=profile.png")
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fi.Size()))

	//stream the body to the client without fully loading it into memory
	_, err = io.Copy(w, f)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
}

func GetPhoto64(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)

	filePath := utils.GetPathFile(id, id, config.Config.BaseStaticPath, modulePath, "profile_", "png")

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	image, err := utils.PathImageToBase64Str(filePath)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(image, http.StatusOK, w)
}
