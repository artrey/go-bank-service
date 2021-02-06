package currency

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"github.com/artrey/go-bank-service/pkg/currency/dto"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Service struct {
	baseUrl string
	timeout time.Duration
	client  *http.Client
}

type Currency struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

func NewService(baseUrl string, timeout time.Duration, client *http.Client) *Service {
	return &Service{
		baseUrl: baseUrl,
		timeout: timeout,
		client:  client,
	}
}

func (s *Service) getResponseBody() ([]byte, error) {
	reqUrl := s.baseUrl

	ctx, _ := context.WithTimeout(context.Background(), s.timeout)
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

func (s *Service) extractDTO() (*dto.RateListDTO, error) {
	data, err := s.getResponseBody()

	var rateList *dto.RateListDTO
	err = xml.Unmarshal(data, &rateList)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return rateList, nil
}

func (s *Service) Extract(writer io.Writer) (err error) {
	rateListDTO, err := s.extractDTO()
	if err != nil {
		log.Println(err)
		return err
	}

	currencies := make([]Currency, len(rateListDTO.Rates))
	for i := 0; i < len(rateListDTO.Rates); i++ {
		currencies[i].Code = rateListDTO.Rates[i].NumCode
		currencies[i].Name = rateListDTO.Rates[i].Name
		currencies[i].Value = rateListDTO.Rates[i].ValueInCents()
	}

	data, err := json.Marshal(currencies)
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
