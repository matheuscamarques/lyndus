package config

import (
	"bufio"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"

	"github.com/BurntSushi/toml"
)

const (
	configPath = "config.toml"
)

var Config Configuration

// Configuration configuraction struct'
type Configuration struct {
	Port              string `toml:"port"`
	PrivKeyPath       string `toml:"priv_key_path"`
	PubKeyPath        string `toml:"pub_key_path"`
	VideosPath        string `toml:"videos_path"`
	PdfsPath          string `toml:"pdfs_path"`
	ImagesCoursePath  string `toml:"images_course_path"`
	ImagesSubjectPath string `toml:"images_subject_path"`
	TemplatePath      string `toml:"template_path"`
	BaseStaticPath    string `toml:"base_static_path"`
	BaseDeletedPath   string `json:"base_deleted_path"`
	TokenExpires      int    `toml:"token_expires"` //in minutes
	SecretHash        string `toml:"secret_hash"`
	AdmPassword       string `toml:"adm_password"`
	PostgresUser      string `toml:"postgres_user"`
	PostgresPass      string `toml:"postgres_pass"`
	PostgresHost      string `toml:"postgres_host"`
	PostgresPort      int    `toml:"postgres_port"`
	PostgresDatabase  string `toml:"postgres_database"`
	BaseServer        string `toml:"base_server,omitempty"`

	MailFrom     string `toml:"mail_from"`
	MailUsername string `toml:"mail_username"`
	MailPassword string `toml:"mail_password"`
	MailSMTPHost string `toml:"mail_smtp_host"`
	MailSMTPPort string `toml:"mail_smtp_port"`

	ReservationCancelationAllowedUpToHours int    `toml:"reservation_cancelation_allowed_up_to_hours"`
	SubscriptionServiceURL                 string `toml:"subscription_service_url"`
}

func ReadConfig() (err error) {
	bs, err := ioutil.ReadFile(configPath)
	if err != nil {
		return
	}
	_, err = toml.Decode(string(bs), &Config)
	return
}

var (
	VerifyKey *rsa.PublicKey
	SignKey   *rsa.PrivateKey
)

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

// Base64Response response default method
func Base64Response(response string, code int, w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(response))
}

// JSONResponse response default method
func JSONResponse(response interface{}, code int, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(json)
}

// Base64ImageResponse response image file
func Base64ImageResponse(path string, code int, w http.ResponseWriter) {

	_, err := os.Stat(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	file, err := os.Open(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var size int64 = fi.Size()
	buf := make([]byte, size)

	fReader := bufio.NewReader(file)
	_, err = fReader.Read(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	w.Write([]byte("data:image/png;base64,"))
	encoder := base64.NewEncoder(base64.StdEncoding, w)

	_, err = encoder.Write(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Must close the encoder when finished to flush any partial blocks.
	// If you comment out the following line, the last partial block "r"
	// won't be encoded.
	err = encoder.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	return

}
