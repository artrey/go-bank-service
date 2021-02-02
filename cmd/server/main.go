package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/artrey/go-bank-service/pkg/transaction"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	if err := execute(); err != nil {
		os.Exit(1)
	}
}

func execute() (err error) {
	transactionSvc := transaction.NewService()
	data, err := ioutil.ReadFile("transactions.csv")
	if err != nil {
		log.Println(err)
		return err
	}

	reader := csv.NewReader(bytes.NewReader(data))
	records, err := reader.ReadAll()
	if err != nil {
		log.Println(err)
		return err
	}
	transactionSvc.ImportRecords(records)

	listener, err := net.Listen("tcp", "0.0.0.0:9999")
	if err != nil {
		log.Println(err)
		return err
	}
	defer func() {
		if cerr := listener.Close(); cerr != nil {
			log.Println(cerr)
			err = cerr
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(conn, transactionSvc)
	}
}

func handle(conn net.Conn, transactionSvc *transaction.Service) {
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Println(cerr)
		}
	}()

	reader := bufio.NewReader(conn)
	const delim = '\n'
	line, err := reader.ReadString(delim)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("received: %s\n", line)
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		log.Printf("invalid request line %s", line)
		return
	}

	writer := bufio.NewWriter(conn)

	path := parts[1]
	switch path {
	case "/":
		err = writeIndex(writer)
	case "/operations.csv":
		err = writeOperationsCsv(writer, transactionSvc)
	case "/operations.json":
		err = writeOperationsJson(writer, transactionSvc)
	case "/operations.xml":
		err = writeOperationsXml(writer, transactionSvc)
	default:
		err = write404(writer)
	}
	if err != nil {
		log.Println(err)
		return
	}
}

func writeResponse(writer *bufio.Writer, status int, headers []string, content []byte) error {
	const CRLF = "\r\n"

	_, err := writer.WriteString(fmt.Sprintf("HTTP/1.1 %d%s", status, CRLF))
	if err != nil {
		return err
	}

	for _, header := range headers {
		_, err = writer.WriteString(header + CRLF)
		if err != nil {
			return err
		}
	}
	_, err = writer.WriteString(CRLF)
	if err != nil {
		return err
	}

	if content != nil && len(content) > 0 {
		_, err = writer.Write(content)
		if err != nil {
			return err
		}
	}

	return writer.Flush()
}

func write404(writer *bufio.Writer) error {
	return writeResponse(writer, 404, []string{
		"Content-Type: text/html;charset=utf-8",
		"Content-Length: 0",
		"Connection: close",
	}, nil)
}

func writeIndex(writer *bufio.Writer) error {
	username := "Александр"
	balance := "1 000.50"

	page, err := ioutil.ReadFile("web/templates/index.html")
	if err != nil {
		return err
	}
	page = bytes.ReplaceAll(page, []byte("{username}"), []byte(username))
	page = bytes.ReplaceAll(page, []byte("{balance}"), []byte(balance))

	return writeResponse(writer, 200, []string{
		"Content-Type: text/html;charset=utf-8",
		fmt.Sprintf("Content-Length: %d", len(page)),
		"Connection: close",
	}, page)
}

func writeOperationsCsv(writer *bufio.Writer, transactionSvc *transaction.Service) error {
	records := transactionSvc.ExportRecords()

	var contentBuffer bytes.Buffer
	w := csv.NewWriter(&contentBuffer)
	w.Comma = ';'
	err := w.WriteAll(records)
	if err != nil {
		return err
	}
	content := contentBuffer.Bytes()

	return writeResponse(writer, 200, []string{
		"Content-Type: text/csv",
		fmt.Sprintf("Content-Length: %d", len(content)),
		"Connection: close",
	}, content)
}

func writeOperationsJson(writer *bufio.Writer, transactionSvc *transaction.Service) error {
	content, err := json.Marshal(transactionSvc.Transactions())
	if err != nil {
		return err
	}

	return writeResponse(writer, 200, []string{
		"Content-Type: application/json;charset=utf-8",
		fmt.Sprintf("Content-Length: %d", len(content)),
		"Connection: close",
	}, content)
}

func writeOperationsXml(writer *bufio.Writer, transactionSvc *transaction.Service) error {
	data := transaction.Transactions{
		Transactions: transactionSvc.Transactions(),
	}
	content, err := xml.Marshal(data)
	if err != nil {
		return err
	}

	return writeResponse(writer, 200, []string{
		"Content-Type: application/xml",
		fmt.Sprintf("Content-Length: %d", len(content)),
		"Connection: close",
	}, content)
}
