// Package define implements functions separate of project.
package define

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

// Itoa convet int to string.
var Itoa = strconv.Itoa

var ParseUint = func(x string) uint { xUint, _ := strconv.ParseUint(x, 10, 0); return uint(xUint) }
var ParseFloat = func(x string) float64 { xFloat, _ := strconv.ParseFloat(x, 64); return xFloat }

// Atoi convet string to int.
var Atoi = func(x string) int { xInt, _ := strconv.Atoi(x); return xInt }

// Itob convet int to []byte.
var Itob = func(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// Min return minimal value in args.
func Min[T constraints.Ordered](args ...T) T {
	if len(args) == 0 {
		var zero T
		return zero
	}
	min := args[0]
	for _, arg := range args {
		if arg < min {
			min = arg
		}
	}
	return min
}

// Max return maximal value in args.
func Max[T constraints.Ordered](args ...T) T {
	if len(args) == 0 {
		var zero T
		return zero
	}
	max := args[0]
	for _, arg := range args {
		if arg > max {
			max = arg
		}
	}
	return max
}

// Sum return sum value in args.
func Sum[T constraints.Ordered](args ...T) T {
	var sum T
	for _, arg := range args {
		sum += arg
	}
	return sum
}

// GetToday returns string today date.
func GetToday() string {
	ctime := time.Now()
	cday, cmonth := ctime.Day(), int(ctime.Month())
	day, month := Itoa(cday), Itoa(cmonth)
	if len(day) < 2 {
		day = "0" + day
	}
	if len(month) < 2 {
		month = "0" + month
	}
	return fmt.Sprintf("%v-%v-%v", day, month, ctime.Year())
}

// Contains returns
func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Abs[T constraints.Integer](a T) T {
	if a < 0 {
		return -a
	}
	return a
}

func Pow[T constraints.Integer](a T, b T) float64 {
	res := float64(1)
	var aa float64
	if b < 0 {
		aa = float64(1) / float64(a)
	} else {
		aa = float64(a)
	}
	for i := T(0); i < Abs(b); i++ {
		res *= aa
	}
	return res
}

func Insert[T constraints.Ordered](a []T, index int, value T) []T {
	if index < 0 {
		index = 0
	}
	if len(a) <= index {
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...)
	a[index] = value
	return a
}

// Hash returns hash data
func Hash(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

type Set[T constraints.Ordered] struct {
	items []T
}

func (this *Set[T]) GetItems() []T {
	var tItems []T
	for _, x := range this.items {
		tItems = append(tItems, x)
	}
	return tItems
}

func (this *Set[T]) Add(x T) {
	if !Contains(this.items, x) {
		this.items = append(this.items, x)
	}
}

func (this *Set[T]) Adds(lst ...T) {
	for _, x := range lst {
		this.Add(x)
	}
}

func (this *Set[T]) Get(index int) T {
	var t T
	if index < 0 {
		index += len(this.items)
	}
	if index > len(this.items) || index < 0 {
		return t
	}
	return this.items[index]
}

func (this *Set[T]) Count() int {
	return len(this.items)
}

func GetEncodeFunc(format string) func(io.Writer, image.Image) error {
	switch format {
	case "image/png":
		return png.Encode
	case "png":
		return png.Encode
	case "image/jpeg":
		return func(w io.Writer, i image.Image) error {
			return jpeg.Encode(w, i, nil)
		}
	case "jpeg":
		return func(w io.Writer, i image.Image) error {
			return jpeg.Encode(w, i, nil)
		}
	case "image/jpg":
		return func(w io.Writer, i image.Image) error {
			return jpeg.Encode(w, i, nil)
		}
	case "jpg":
		return func(w io.Writer, i image.Image) error {
			return jpeg.Encode(w, i, nil)
		}
	}
	return nil
}

func GetDecodeFunc(format string) func(io.Reader) (image.Image, error) {
	switch format {
	case "image/png":
		return png.Decode
	case "png":
		return png.Decode
	case "image/jpeg":
		return jpeg.Decode
	case "jpeg":
		return jpeg.Decode
	case "image/jpg":
		return jpeg.Decode
	case "jpg":
		return jpeg.Decode
	}
	return nil
}

func GetImagesFromRequestBody(body []byte, key ...string) ([]image.Image, []string) {
	var data map[string]interface{}
	json.Unmarshal(body, &data)
	if len(key) == 0 {
		key = []string{"images"}
	}
	imagesData := data[key[0]].([]interface{})
	var images []image.Image
	var formats []string
	for i := 0; i < len(imagesData); i++ {
		imgData := imagesData[i].(string)
		coI := strings.Index(imgData, ",")
		rawImage := string(imgData)[coI+1:]
		unbased, _ := base64.StdEncoding.DecodeString(string(rawImage))
		res := bytes.NewReader(unbased)
		format := imgData[strings.Index(imgData, ":")+1 : strings.Index(imgData, ";")]
		img, err := GetDecodeFunc(format)(res)
		if err != nil {
			log.Println("GetImagesFromRequestBody: decode:", err.Error())
			continue
		}
		images = append(images, img)
		formats = append(formats, format)
	}
	return images, formats
}

func ImagesToBytes(images []image.Image, formats []string) [][]byte {
	var res [][]byte
	for i := 0; i < len(images); i++ {
		img := images[i]
		buf := new(bytes.Buffer)
		err := GetEncodeFunc("png")(buf, img)
		if err == nil {
			res = append(res, buf.Bytes())
		}
	}
	return res
}

func IndexOf[T comparable](lst []T, x T) int {
	for i := 0; i < len(lst); i++ {
		if lst[i] == x {
			return i
		}
	}
	return -1
}

func GoToStruct(value reflect.Value) (*reflect.Value, error) {
	switch value.Kind() {
	case reflect.Pointer:
		return GoToStruct(value.Elem())
	case reflect.Interface:
		return GoToStruct(value.Elem())
	case reflect.Struct:
		return &value, nil
	default:
		return nil, fmt.Errorf("GoToStruct: getting value is not struct, is %v", value.Kind())
	}
}

func GetTagField(model any, fieldName string, tagName ...string) string {
	vModel, err := GoToStruct(reflect.ValueOf(model))
	if err != nil {
		return ""
	}
	field, ok := reflect.TypeOf(vModel.Interface()).FieldByName(fieldName)
	if !ok {
		return ""
	}
	tag := field.Tag
	if len(tagName) > 0 && tagName[0] != "" {
		tag = reflect.StructTag(tag.Get(tagName[0]))
	}
	return string(tag)
}

func Check(imodel interface{}, field_name string) (*reflect.Value, error) {
	vModel, err := GoToStruct(reflect.ValueOf(imodel))
	if err != nil {
		return nil, err
	}
	withins := strings.Split(field_name, ".")
	vfield := *vModel
	i := 0
	for i, field_name = range withins {
		if vfield.Kind() != reflect.Struct {
			return nil, fmt.Errorf("Check: %v is not struct, is %v", withins[i-1], vfield.Kind())
		}
		cvfield := vfield.FieldByName(field_name)
		if cvfield.Kind() == reflect.Invalid {
			vfield = reflect.Value{}
			break
		}
		vfield = cvfield
	}
	if vfield.Kind() == reflect.Invalid {
		return nil, fmt.Errorf("Check: field `%v` does not exists", field_name)
	}
	return &vfield, nil
}

func ChangeFieldOfName(imodel interface{}, field_name string, value interface{}) error {
	vfield, err := Check(imodel, field_name)
	if err != nil {
		return err
	}
	if vfield.Kind() == reflect.Invalid {
		return fmt.Errorf("ChangeFieldOfName: Field `%v` does not exists", field_name)
	}
	vvalue := reflect.ValueOf(value)
	if !vvalue.CanConvert(vfield.Type()) {
		return fmt.Errorf("ChangeFieldOfName: `%v` (type %T) not be converted to type %v", value, value, vfield.Type())
	}
	vvalue = vvalue.Convert(vfield.Type())
	vfield.Set(vvalue)
	return nil
}

func CopyMap[T1, T2 comparable](Map map[T1]T2) map[T1]T2 {
	newMap := map[T1]T2{}
	for k, v := range Map {
		newMap[k] = v
	}
	return newMap
}

func CopyMapAny[T1, T2 comparable](Map map[T1]T2) map[string]any {
	newMap := map[string]any{}
	for k, v := range Map {
		newMap[fmt.Sprint(k)] = v
	}
	return newMap
}

func Pop[T comparable](lst []T, index int) ([]T, T) {
	var x T
	if index < 0 || index >= len(lst) {
		return lst, x
	}
	if len(lst) == 1 {
		return []T{}, lst[0]
	}
	x = lst[index]
	t := lst[:index]
	if index != len(lst)-1 {
		lst = append(t, lst[index+1:]...)
	}
	return lst, x
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

type Headers map[string]string
type Data map[string]any

func GetResponse(method, curl string, headers Headers, data Data) (int, []byte, error) {
	dataByte, err := json.Marshal(data)
	if err != nil {
		return 500, nil, fmt.Errorf("GetResponse: marshal data: %v", err)
	}

	req, err := http.NewRequest(method, curl, strings.NewReader(string(dataByte)))
	if err != nil {
		return 500, nil, fmt.Errorf("GetResponse: create request: %v", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 500, nil, fmt.Errorf("GetResponse: get response: %v", err)
	}

	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body, nil
}

func GetJSONResponse(method, curl string, headers Headers, data Data) (int, any, error) {
	status, body, err := GetResponse(method, curl, headers, data)
	if err != nil {
		return status, body, err
	}
	var response any
	if err := json.Unmarshal(body, &response); err != nil {
		if len(body) == 0 {
			return 500, nil, fmt.Errorf("GetJSONResponse: body is empty")
		}
		if err1 := json.Unmarshal([]byte(string(body)[3:]), &response); err1 != nil {
			return 500, nil, fmt.Errorf("GetJSONResponse: unmarshal response:\n error: %v\n body: %v", err1, string(body))
		}
	}

	return status, response, nil
}

/*
Compare returns

	-2 if a == nil || b == nil
	-1 if a < b
	 0 if a == b
	+1 if a > b
*/
func Compare(a, b any) int {
	if a == nil || b == nil {
		return -2
	}
	v := reflect.ValueOf(a)
	u := reflect.ValueOf(b)

	if v.Equal(u) {
		return 0
	}

	switch v.Kind() {
	case reflect.Bool:
		a := v.Bool()
		b := u.Bool()
		if !a && b {
			return -1
		} else if a == b {
			return 0
		}
		return 1
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		a := v.Int()
		b := u.Int()
		if a < b {
			return -1
		} else if a == b {
			return 0
		}
		return 1
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		a := v.Uint()
		b := u.Uint()
		if a < b {
			return -1
		} else if a == b {
			return 0
		}
		return 1
	case reflect.Float32, reflect.Float64:
		a := v.Float()
		b := u.Float()
		if a < b {
			return -1
		} else if a == b {
			return 0
		}
		return 1
	case reflect.String:
		a, b := v.String(), u.String()
		if a < b {
			return -1
		} else if a == b {
			return 0
		}
		return 1
	}
	return -2
}
