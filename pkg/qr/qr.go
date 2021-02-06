package qr

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const (
	generateCodeMethod = "create-qr-code/"
)

type Service struct {
	baseUrl string
	client  *http.Client
}

func NewService(baseUrl string, client *http.Client) *Service {
	return &Service{
		baseUrl: baseUrl,
		client:  client,
	}
}

func (s Service) getResponse(ctx context.Context, data string) (*http.Response, error) {
	values := make(url.Values)
	values.Add("data", data)
	reqUrl := fmt.Sprintf("%s/%s?%s", s.baseUrl, generateCodeMethod, values.Encode())

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

	return resp, nil
}

func extractFiletype(contentType string) (string, error) {
	if !strings.HasPrefix(contentType, "image/") {
		err := errors.New("incorrect Content-Type")
		log.Println(err)
		return "", err
	}
	return contentType[6:], nil
}

func (s Service) Encode(ctx context.Context, data string) ([]byte, string, error) {
	resp, err := s.getResponse(ctx, data)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}
	if err = resp.Body.Close(); err != nil {
		log.Println(err)
		return nil, "", err
	}

	filetype, err := extractFiletype(resp.Header.Get("Content-Type"))
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	return body, filetype, nil
}
