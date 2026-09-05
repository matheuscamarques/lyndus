package bs

import (
	"bitbucket.org/lyndus/backend/domain/constants"
	"bitbucket.org/lyndus/backend/infra/types"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"net/http"
	"os"
	"strconv"

	"bitbucket.org/lyndus/backend/domain/bs"
	"bitbucket.org/lyndus/backend/domain/bs/entity"
	"bitbucket.org/lyndus/backend/domain/bs/services"
	"bitbucket.org/lyndus/backend/infra/config"
	"bitbucket.org/lyndus/backend/infra/logger"
	"bitbucket.org/lyndus/backend/infra/rest/response"
	"bitbucket.org/lyndus/backend/infra/utils"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

//SaveBM create new BM
func SaveBM(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BM, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var bmCreate entity.BM
	err = json.NewDecoder(r.Body).Decode(&bmCreate)
	if err == types.CNPJError {
		config.ResponsePerErr(w, err, config.CNPJINVALID)
		return
	} else if err == types.CPFError {
		config.ResponsePerErr(w, err, config.CPFINVALID)
		return
	} else if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	// uitls
	if !bmCreate.ValidateServiceStepTime() {
		config.ResponsePerErr(w, errors.New("serviços duration não coincide com serviceStepTime"), config.INVALIDREQUEST)
		return
	}

	bmCreate.BsID = bsID

	//todo bs_bm já cadastrado reativar?

	bmID, err := services.BmService.GetBmIDByCPF(bsID, bmCreate.CPF)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if bmID != 0 {
		config.ResponsePerErr(w, err, config.DOCALREADYREGISTERED)
		return
	}

	bmCreate.StatusID = constants.BMStatusActive

	bmCreate.ID, err = services.BmService.CreateBM(bmCreate)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	for _, service := range bmCreate.Services {
		err = services.BmService.AddBMService(bmCreate.ID, service)
		if err != nil {
			logger.Error("add bs_bm bs_service", zap.String("err", err.Error()))
		}
	}

	for _, weekDay := range bmCreate.WeekDays {
		err = services.BmService.AddBMWeekDay(bmCreate.ID, weekDay)
		if err != nil {
			logger.Error("add bs_bm week day", zap.String("err", err.Error()))
		}
	}

	if bmCreate.Image != "" {
		filePath := utils.GetPathFile(bsID, bmCreate.ID, config.Config.BaseStaticPath, modulePath, "bs_bm", "png")

		err = utils.Base64ImageSave(bmCreate.Image, filePath)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
	}

	resp := response.ID{
		ID: bmCreate.ID,
	}

	config.JSONResponse(resp, http.StatusOK, w)
}

//UpdateBM Update BM
func UpdateBM(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BM, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var bmUpdate entity.BM
	err = json.NewDecoder(r.Body).Decode(&bmUpdate)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if !bmUpdate.ValidateServiceStepTime() {
		config.ResponsePerErr(w, errors.New("serviços duration não coincide com serviceStepTime"), config.INVALIDREQUEST)
		return
	}

	bmUpdate.BsID = bsID
	bmUpdate.StatusID = constants.BMStatusActive

	n, err := services.BmService.GetBmID(bsID, bmUpdate.ID, bmUpdate.StatusID)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if n == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	bmServices, err := services.BmService.GetBmServices(bsID, bmUpdate.ID)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	err = services.BmService.UpdateBM(bmUpdate)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	cntrl := false
	for _, service := range bmUpdate.Services {
		cntrl = false
		for k, s := range bmServices {
			if s.ID == service.ID {
				bmServices = append(bmServices[:k], bmServices[k+1:]...)
				err = services.BmService.UpdateBMService(bmUpdate.ID, service)
				if err != nil {
					logger.Error("add bs_bm bs_service", zap.String("err", err.Error()))
				}
				cntrl = true
			}
		}
		if !cntrl {
			err = services.BmService.AddBMService(bmUpdate.ID, service)
			if err != nil {
				logger.Error("add bs_bm bs_service", zap.String("err", err.Error()))
			}
		}
	}

	for _, s := range bmServices {
		err = services.BmService.DeleteBMService(bmUpdate.ID, s.ID)
		if err != nil {
			logger.Error("remove bs_bm bs_service", zap.String("err", err.Error()))
		}
	}

	err = services.BmService.DeleteBMAllWeekDays(bmUpdate.ID)
	if err != nil {
		//log.Println(err)
	}
	for _, weekDay := range bmUpdate.WeekDays {
		err = services.BmService.AddBMWeekDay(bmUpdate.ID, weekDay)
		if err != nil {
			//log.Println(err)
		}
	}

	if bmUpdate.Image != "" {
		filePath := utils.GetPathFile(bsID, bmUpdate.ID, config.Config.BaseStaticPath, modulePath, "bs_bm", "png")

		err = utils.Base64ImageSave(bmUpdate.Image, filePath)
		if err != nil {
			config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	return
}

func GetBM(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BM, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	bmID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bm, err := services.BmService.GetBM(bsID, bmID, constants.BMStatusActive)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if bm.ID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	bm.Services, err = services.BmService.GetBmServices(bsID, bm.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	bm.WeekDays, err = services.BmService.GetBMWeekDays(bm.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	filePath := utils.GetPathFile(bsID, bm.ID, config.Config.BaseStaticPath, modulePath, "bs_bm", "png")

	//todo responder um avatar padrão fake caso não exista
	bm.Image, err = utils.PathImageToBase64Str(filePath)
	if err != nil {
		//todo logar erro do path
		//logger.Error("erro ao buscar image bs_bm", zap.String("err", err.Error()))
	}

	config.JSONResponse(bm, http.StatusOK, w)
}

func GetBMS(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BM, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	var bsBMResponse response.BSBMResponse

	bsBMResponse.ActivePage, _ = strconv.Atoi(r.URL.Query().Get("page"))
	itemsPerPage, _ := strconv.Atoi(r.URL.Query().Get("itemsPerPage"))
	search := r.URL.Query().Get("search")

	if bsBMResponse.ActivePage < 1 {
		bsBMResponse.ActivePage = 1
	}

	bsBMResponse.Items, bsBMResponse.TotalItems, bsBMResponse.TotalPages, err = services.BmService.GetBMS(bsID, bsBMResponse.ActivePage, itemsPerPage, search)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	config.JSONResponse(bsBMResponse, http.StatusOK, w)
}

func DeleteBM(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BANNER, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	bmID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.BmService.ChangeStatus(bsID, bmID, constants.BMStatusInactive)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
	}

	w.WriteHeader(http.StatusOK)
	return
}

//ADDImageBM func set bs_bm pic
func ADDImageBM(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BM, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	bmID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
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

	filePath := utils.GetPathFile(bsID, bmID, config.Config.BaseStaticPath, modulePath, "bs_bm", "png")

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	err = r.MultipartForm.RemoveAll()
	if err != nil {
		//log.Println(err)
	}

	w.WriteHeader(http.StatusOK)
	return
}

func BMGetImage(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.PROFILE, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	filePath := utils.GetPathFile(bsID, bsID, config.Config.BaseStaticPath, modulePath, "logo", "png")
	f, err := os.Open(filePath)
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
func GetBMServices(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BM, bs.VIEW)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	bmID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	bmServices, err := services.BmService.GetBmServices(bsID, bmID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if bmServices == nil {
		bmServices = make([]entity.Service, 0)
	}
	config.JSONResponse(bmServices, http.StatusOK, w)
}

func AddService(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BM, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	bmID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	var service entity.Service
	err = json.NewDecoder(r.Body).Decode(&service)
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	if service.ID == 0 {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	dbServiceID, err := services.BS.GetServiceID(bsID, service.ID)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	if dbServiceID == 0 {
		config.ResponsePerErr(w, err, config.NOTFOUND)
		return
	}

	err = services.BS.AddBMService(bmID, service)
	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func DeleteBMService(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	bsID, err := services.AuthService.CheckPermission(id, bs.BM, bs.ALL)
	if err != nil || bsID == 0 {
		config.ResponsePerErr(w, err, config.UNAUTHORIZED)
		return
	}

	bmID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}
	serviceID, err := strconv.Atoi(chi.URLParam(r, "service_id"))
	if err != nil {
		config.ResponsePerErr(w, err, config.INVALIDREQUEST)
		return
	}

	err = services.BmService.DeleteBMService(bmID, serviceID)

	if err != nil {
		config.ResponsePerErr(w, err, config.INTERNALSERVERERROR)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
