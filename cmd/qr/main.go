package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/artrey/go-bank-service/pkg/qr"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	if err := execute(); err != nil {
		os.Exit(1)
	}
}

func execute() error {
	const baseUrl = "https://api.qrserver.com/v1"
	const folderName = "qrs"
	const text = "some text for encoding"

	timeoutString, ok := os.LookupEnv("TIMEOUT")
	if !ok {
		err := errors.New("env variable TIMEOUT is not set")
		log.Println(err)
		return err
	}

	timeout, err := strconv.Atoi(timeoutString)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, err = os.Stat(folderName); os.IsNotExist(err) {
		if err = os.Mkdir(folderName, os.ModeDir); err != nil {
			log.Println(err)
			return err
		}
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Millisecond)
	svc := qr.NewService(baseUrl)
	data, filetype, err := svc.Encode(ctx, text)
	if err != nil {
		log.Println(err)
		return err
	}

	filename := fmt.Sprintf("%s.%s", text, filetype)
	relativePath := filepath.Join(folderName, filename)
	if err = ioutil.WriteFile(relativePath, data, os.ModeType); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
