package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"wechatbot/pkg/log"

	uuid "github.com/gofrs/uuid"
)

func GetUUID() string {
	u4 := uuid.Must(uuid.NewV4()).String()
	return u4
}

func FileNameWithoutExtension(fileName string) string {
	if pos := strings.LastIndexByte(fileName, '.'); pos != -1 {
		return fileName[:pos]
	}
	return fileName
}

func Marshal(object interface{}) (string, error) {
	b, err := json.Marshal(object)
	if err != nil {
		return "", err
	}
	str := string(b)
	return str, nil
}

func Unmarshal(jsonStr string, target interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), target)
	if err != nil {
		log.Warn(err.Error())
		return err
	}
	return nil
}

func RandomNumber(min, max int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := min + rand.Intn(max-min)
	return strconv.Itoa(n)
}

func ContainsInt(s []int, i int) bool {
	for _, a := range s {
		if a == i {
			return true
		}
	}
	return false
}

func ContainsString(s []string, i string) bool {
	for _, a := range s {
		if a == i {
			return true
		}
	}
	return false
}

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func Str2Md5(s string) (r string) {
	h := md5.New()
	h.Write([]byte(s))
	r = hex.EncodeToString(h.Sum(nil))
	return
}

func DownloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func Struct2Map(s interface{}) (m map[string]string, err error) {
	data, err := json.Marshal(s)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &m) // Convert to a map
	return
}
