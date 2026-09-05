package utils

import (
	"bitbucket.org/lyndus/backend/infra/logger"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"go.uber.org/zap"
	"io"
	mathRand "math/rand"
	"net/mail"
	"reflect"
	"time"
	"unicode"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type LevelPermission struct {
	Id        int    `json:"id"         db:"id"`
	Desc      string `json:"desc"       db:"desc"`
	LevelId   int    `json:"levelID"   db:"level_id"`
	LevelName string `json:"levelName" db:"level_name"`
}

type Levels struct {
	Id   int    `json:"id"`
	Desc string `json:"desc"`
}
type Permission struct {
	ID   int    `json:"id"`
	Desc string `json:"desc"`
}

// @important ao criar permissão adicionar a esta lista
var PermissionClientMaster = []LevelPermission{
	LevelPermission{Id: 1, LevelId: 3},
	LevelPermission{Id: 2, LevelId: 3},
	LevelPermission{Id: 3, LevelId: 3},
	LevelPermission{Id: 4, LevelId: 3},
	LevelPermission{Id: 5, LevelId: 3},
	LevelPermission{Id: 6, LevelId: 3},
}

// @important ao criar permissão adicionar a esta lista
var PermissionBsMaster = []LevelPermission{
	LevelPermission{Id: 1, LevelId: 3},
	LevelPermission{Id: 2, LevelId: 3},
	LevelPermission{Id: 3, LevelId: 3},
	LevelPermission{Id: 4, LevelId: 3},
	LevelPermission{Id: 5, LevelId: 3},
	LevelPermission{Id: 6, LevelId: 3},
	LevelPermission{Id: 7, LevelId: 3},
	LevelPermission{Id: 8, LevelId: 3},
	LevelPermission{Id: 9, LevelId: 3},
	LevelPermission{Id: 10, LevelId: 3},
	LevelPermission{Id: 11, LevelId: 3},
	LevelPermission{Id: 12, LevelId: 3},
	LevelPermission{Id: 13, LevelId: 3},
	LevelPermission{Id: 14, LevelId: 3},
}

var AccessLevels = []Levels{
	{1, "Nenhuma"},
	{2, "Visualizar"},
	{3, "Todas"},
}

var (
	cnpjFirstDigitTable  = []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	cnpjSecondDigitTable = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
)

const LAYOUTDATETIME1 = "2006-01-02 15:04"
const LAYOUTDATETIME2 = "02/01/2006 15:04"
const LAYOUTDATE = "2006-01-02"
const LAYOUTDATE2 = "02/01/2006"
const LAYOUTHOUR = "15:04"
const LAYOUTHOURNANO = "15:04:05.999999999"

func Adjust(d int) int {
	if d < 2 {
		return 0
	}
	return 11 - d
}

func OnlyDigits(s string) string {
	buf := bytes.NewBufferString("")
	for _, r := range s {
		if unicode.IsNumber(r) {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

func Limpar(b []byte) string {
	//todo mudar nome limpar
	buf := bytes.NewBufferString("")
	for _, r := range b {
		if unicode.IsNumber(rune(r)) {
			buf.WriteByte(r)
		}
	}
	return buf.String()
}

func Decimalize(b []byte) string {
	buf1 := bytes.NewBufferString("")
	buf2 := bytes.NewBufferString("")
	check := false
	for _, r := range b {
		if rune(r) == '.' || rune(r) == ',' {
			check = true
		}
		if unicode.IsNumber(rune(r)) {
			if check {
				buf2.WriteByte(r)
			} else {
				buf1.WriteByte(r)
			}
		}
	}
	for buf2.Len() < 2 {
		buf2.WriteByte('0')
	}
	buf1.Write(buf2.Bytes())
	return buf1.String()
}

func ValidateCPF(cpf string) bool {
	if len(cpf) != 11 {
		return false
	}
	if cpf[:9] == "111111111" ||
		cpf[:9] == "222222222" ||
		cpf[:9] == "333333333" ||
		cpf[:9] == "444444444" ||
		cpf[:9] == "555555555" ||
		cpf[:9] == "666666666" ||
		cpf[:9] == "777777777" ||
		cpf[:9] == "888888888" ||
		cpf[:9] == "999999999" ||
		cpf[:9] == "123456789" {
		return false
	}

	var n, d int
	for k, v := range cpf[:9] {
		n = int(v - '0')
		d += n * (10 - k)
	}
	d %= 11
	d = Adjust(d)
	if byte('0'+d) != cpf[9] {
		return false
	}
	d = 0
	for k, v := range cpf[:10] {
		n = int(v - '0')
		d += n * (11 - k)
	}
	d %= 11
	d = Adjust(d)
	return byte('0'+d) == cpf[10]
}

func ValidateCNPJ(cnpj string) bool {
	if len(cnpj) != 14 {
		return false
	}
	var n, d int
	for k, v := range cnpj[:12] {
		n = int(v - '0')
		d += n * cnpjFirstDigitTable[k]
	}
	d %= 11
	d = Adjust(d)
	if byte('0'+d) != cnpj[12] {
		return false
	}
	d = 0
	for k, v := range cnpj[:13] {
		n = int(v - '0')
		d += n * cnpjSecondDigitTable[k]
	}

	d %= 11
	d = Adjust(d)
	return byte('0'+d) == cnpj[13]
}

// validar telefone padrão brasileiro 11 9 99999999
func ValidateTelephone(telephone string) bool {
	if len(telephone) < 10 {
		return false
	}
	return true
}

// JWTCustomClaims are custom claims extending default ones.
type JWTCustomClaims struct {
	UserID int `json:"uid"`
	jwt.StandardClaims
}

type ID struct {
	ID int `json:"id"`
}

// OK ok response struct
type OK struct {
	OK string `json:"ok"`
}

// Exist response struct
type Exist struct {
	Exist bool `json:"exist"`
}

type Search struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SearchResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
	//Students  []StudentSimpleResponse `json:"students"`
}

func ComparePasswords(hashedPwd, userPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(userPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false
	}
	return true
}

func HashPwd(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Err create hash password", zap.String("error:", err.Error()))
	}
	return string(hash)
}

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func createHash(key string) []byte {
	hash := md5.Sum([]byte(key))
	dst := make([]byte, hex.EncodedLen(len(hash)))
	hex.Encode(dst, hash[:])
	return dst
}

func Encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher(createHash(passphrase))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		fmt.Println(err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func Decrypt(data []byte, passphrase string) ([]byte, error) {
	key := createHash(passphrase)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

const letters = "abcdefghijklmnopqrstuvwxyz@0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

var seededRand *mathRand.Rand = mathRand.New(mathRand.NewSource(time.Now().UnixNano()))

// RandStringBytesMask generate rand string from length
func RandStringBytesMask(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[seededRand.Intn(len(letters))]
	}
	return string(b)
}

func Assign(target interface{}, object interface{}) {
	t := reflect.ValueOf(target).Elem()
	o := reflect.ValueOf(object).Elem()
	for i := 0; i < o.NumField(); i++ {
		for j := 0; j < t.NumField(); j++ {
			if o.Type().Field(i).Name == t.Type().Field(j).Name {
				t.Field(j).Set(o.Field(i))
			}
		}
	}
}
