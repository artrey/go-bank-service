package currency

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/artrey/go-bank-service/pkg/currency/dto"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	currenciesMethod = "netology-code/bgo-homeworks/master/10_client/assets/daily.xml"
)

type Service struct {
	baseUrl string
	timeout time.Duration
	client  *http.Client
}

type Currency struct {
	Code  string
	Name  string
	Value int64
}

func NewService(baseUrl string, client *http.Client) *Service {
	return &Service{
		baseUrl: baseUrl,
		client:  client,
	}
}

func (s *Service) getResponseBody(ctx context.Context, method string) ([]byte, error) {
	reqUrl := fmt.Sprintf("%s/%s", s.baseUrl, method)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = resp.Body.Close(); err != nil {
		log.Println(err)
		return nil, err
	}

	return data, nil
}

func (s *Service) extractXmlDTO(ctx context.Context) (*dto.RateListXmlDTO, error) {
	data, err := s.getResponseBody(ctx, currenciesMethod)

	var rateList *dto.RateListXmlDTO
	err = xml.Unmarshal(data, &rateList)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return rateList, nil
}

func (s *Service) extractCurrenciesFromXml(ctx context.Context) ([]Currency, error) {
	rateListDTO, err := s.extractXmlDTO(ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	currencies := make([]Currency, len(rateListDTO.Rates))
	for i := 0; i < len(rateListDTO.Rates); i++ {
		currencies[i].Code = rateListDTO.Rates[i].NumCode
		currencies[i].Name = rateListDTO.Rates[i].Name
		currencies[i].Value = rateListDTO.Rates[i].ValueInCents()
	}
	return currencies, nil
}

func (s *Service) marshalAsJson(currencies []Currency) ([]byte, error) {
	currenciesJson := make([]dto.RateJsonDTO, len(currencies))
	for i := 0; i < len(currencies); i++ {
		currenciesJson[i].Code = currencies[i].Code
		currenciesJson[i].Name = currencies[i].Name
		currenciesJson[i].Value = currencies[i].Value
	}

	data, err := json.Marshal(currenciesJson)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return data, nil
}

func (s *Service) Extract(ctx context.Context, writer io.Writer) (err error) {
	currencies, err := s.extractCurrenciesFromXml(ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	data, err := s.marshalAsJson(currencies)
	if err != nil {
		log.Println(err)
		return err
	}

	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
