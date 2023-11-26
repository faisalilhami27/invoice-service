package helper

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"

	"html/template"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"math"
	"os"
	"reflect"
	"strconv"
)

type PaginationParam struct {
	Count int64       `json:"count"`
	Limit int         `json:"limit"`
	Page  int         `json:"page"`
	Data  interface{} `json:"data"`
}

type PaginationResult struct {
	TotalPage    int         `json:"totalPage"`
	TotalData    int64       `json:"totalData"`
	NextPage     *int        `json:"nextPage,omitempty"`
	PreviousPage *int        `json:"previousPage,omitempty"`
	Page         int         `json:"page"`
	Limit        int         `json:"pageSize"`
	Data         interface{} `json:"data"`
}

func GeneratePagination(params PaginationParam) PaginationResult {
	totalPage := int(math.Ceil(float64(params.Count) / float64(params.Limit)))

	var nextPage, previousPage int
	if params.Page < totalPage {
		nextPage = params.Page + 1
	}

	if params.Page > 1 {
		previousPage = params.Page - 1
	}

	result := PaginationResult{
		TotalPage:    totalPage,
		TotalData:    params.Count,
		NextPage:     &nextPage,
		PreviousPage: &previousPage,
		Page:         params.Page,
		Limit:        params.Limit,
		Data:         params.Data,
	}

	return result
}

func BindFromJSON(dest any, filename, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	v.AddConfigPath(path)
	v.SetConfigName(filename)

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	err = v.Unmarshal(&dest)
	if err != nil {
		log.Errorf("failed to unmarshal config env: %+v\n", err)
		return err
	}

	return nil
}

func BindFromConsul(dest any, endPoint, path string) error {
	v := viper.New()

	v.SetConfigType("json")
	err := v.AddRemoteProvider("consul", endPoint, path)
	if err != nil {
		return err
	}

	err = v.ReadRemoteConfig()
	if err != nil {
		return err
	}

	log.Errorf("using config from consul: %s/%s.\n", endPoint, path)

	err = v.Unmarshal(dest)
	if err != nil {
		log.Errorf("failed to unmarshal config dest: %+v\n", err)
		return err
	}

	err = SetEnvFromConsulKV(v)
	if err != nil {
		log.Errorf("failed to set env from consul: %+v\n", err)
		return err
	}

	return nil
}

func SetEnvFromConsulKV(v *viper.Viper) error {
	env := make(map[string]any)

	err := v.Unmarshal(&env)
	if err != nil {
		log.Errorf("failed to unmarshal config env: %+v\n", err)
		return err
	}

	for k, v := range env {
		var (
			valOf = reflect.ValueOf(v)
			val   string
		)

		switch valOf.Kind() { //nolint:exhaustive
		case reflect.String:
			val = valOf.String()
		case reflect.Int:
			val = strconv.Itoa(int(valOf.Int()))
		case reflect.Uint:
			val = strconv.Itoa(int(valOf.Uint()))
		case reflect.Float64:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Float32:
			val = strconv.Itoa(int(valOf.Float()))
		case reflect.Bool:
			val = strconv.FormatBool(valOf.Bool())
		}

		err = os.Setenv(k, val)
		if err != nil {
			return err
		}
	}

	return nil
}

func GenerateSHA256(inputString string) string {
	hash := sha256.New()
	hash.Write([]byte(inputString))
	hashBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashBytes)
	return hashString
}

func GeneratePDF(ctx context.Context, htmlTemplate string, data any) ([]byte, error) {
	contextPDF, cancel := chromedp.NewContext(
		ctx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	tmpl, err := template.New("htmlTemplate").Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}

	var filledTemplate bytes.Buffer
	if err = tmpl.Execute(&filledTemplate, data); err != nil {
		return nil, err
	}

	htmlContent := filledTemplate.String()

	// Convert HTML to PDF
	var pdf []byte
	if err = chromedp.Run(contextPDF,
		chromedp.Navigate("data:text/html,"+htmlContent),
		chromedp.ActionFunc(func(ctx context.Context) error {
			pdf, _, err = page.PrintToPDF().Do(ctx)
			return err
		}),
	); err != nil {
		return nil, err
	}

	return pdf, nil
}
