package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

var (
	endpoint    string
	method      string
	cmdParam    string
	cmdHeader   string
	timeout     int
	headerFlags flagList
	paramFlags  flagList
	headers     map[string]string
	params      map[string]string
	reqFileName string

	trimPrefix string
	trimSuffix string

	client http.Client
)

type flagList []string

func (l *flagList) String() string {
	return ""
}

func (l *flagList) Set(value string) error {
	*l = append(*l, value)
	return nil
}

func init() {
	// Parse flags
	flag.StringVar(&endpoint, "url", "", "Url of web shell")
	flag.StringVar(&cmdParam, "param", "", "Parameter used to send commands")
	flag.StringVar(&cmdHeader, "header", "", "Header used to send commands")
	flag.StringVar(&method, "X", "GET", "HTTP method (GET, POST, PUT, PATCH, DELETE)")
	flag.Var(&headerFlags, "H", "HTTP request headers")
	flag.Var(&paramFlags, "P", "HTTP request body parameters")
	flag.StringVar(&trimPrefix, "trimp", "", "Trim prefix")
	flag.StringVar(&trimSuffix, "trims", "", "Trim suffix")
	flag.IntVar(&timeout, "timeout", 10, "Request timeout in seconds")
	flag.StringVar(&reqFileName, "file", "", "File containing http request to replay")

	flag.Parse()

	// Init headers and params map
	headers = make(map[string]string)
	params = make(map[string]string)

	// If not using request file
	if reqFileName == "" {

		// -url is required
		if endpoint == "" {
			fmt.Println("-url required")
			flag.Usage()
			os.Exit(0)
		}

		// Only valid HTTP methods
		if method != "GET" && method != "POST" && method != "PUT" &&
			method != "PATCH" && method != "DELETE" {
			fmt.Println("Invalid HTTP method. Supported methods are:")
			fmt.Println("\tGET, POST, PUT, PATCH, and DELETE")
			os.Exit(0)
		}
	} else {
		// Open request file
		reqFile, err := os.Open(reqFileName)
		if err != nil {
			fmt.Printf("Could not open request file. %v\n", err)
			os.Exit(1)
		}
		defer reqFile.Close()

		// Read/parse request file
		req, err := http.ReadRequest(bufio.NewReader(reqFile))
		if err != io.EOF && err != nil {
			fmt.Printf("Could not parse request file. %v\n", err)
			os.Exit(1)
		}

		// Use req values
		endpoint = req.Host + req.URL.String()
		method = req.Method
		for k, v := range req.Header {
			for _, vv := range v {
				headers[k] = vv
			}
		}
	}

	// Parse header flags
	for _, h := range headerFlags {
		split := strings.Split(h, ":")
		if len(split) != 2 {
			fmt.Printf("Invalid header: \"%s\"\n", h)
			continue
		}

		headers[split[0]] = split[1]
	}

	// Parse parameter flags
	for _, p := range paramFlags {
		split := strings.Split(p, ":")
		if len(split) != 2 {
			fmt.Printf("Invalid parameter: \"%s\"\n", p)
			continue
		}

		params[split[0]] = split[1]
	}

	client = http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
}

func main() {

	// Validate endpoint url
	host, err := getHost(endpoint)
	if err != nil {
		fmt.Println("Invalid url")
		return
	}

	cyan := color.New(color.FgGreen).SprintFunc()
	prompt := fmt.Sprintf(cyan("%s> "), host)

	l, err := readline.NewEx(&readline.Config{
		Prompt:          prompt,
		HistoryFile:     "/tmp/web-cli.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer l.Close()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		switch {
		// case strings.HasPrefix(line, "get "):
		// 	fmt.Println("Download file")
		// 	fmt.Println(line)
		// case strings.HasPrefix(line, "clear"):
		// 	readline.ClearScreen()
		case line == "exit":
			l.Close()
		default:
			out, err := sendRequest(line)
			if err != nil {
				fmt.Printf("%s: %s", color.RedString("ERROR"), err.Error())
			}

			if trimPrefix != "" {
				i := strings.Index(out, trimPrefix)
				if i > 0 {
					i += len(trimPrefix)
					out = out[i:]
				}
			}
			if trimSuffix != "" {
				i := strings.Index(out, trimSuffix)
				if i > 0 {
					out = out[:i]
				}
			}

			out = strings.TrimSpace(out)
			fmt.Println(out)
		}
	}
}

func sendRequest(cmd string) (string, error) {
	finalURL := endpoint
	var body io.Reader

	// If uploading file...
	if strings.HasPrefix(cmd, "put ") {
		c := strings.Fields(cmd)
		fileName := c[1]

		// Open file for reading
		inFile, err := os.Open(fileName)
		if err != nil {
			return "", err
		}
		defer inFile.Close()

		reader := bufio.NewReader(inFile)
		content, _ := ioutil.ReadAll(reader)

		params["f"] = base64.StdEncoding.EncodeToString(content)

		// Create multipart form
		// b := &bytes.Buffer{}
		// writer := multipart.NewWriter(b)

		// // Add file part
		// part, err := writer.CreateFormFile("f", filepath.Base(fileName))
		// if err != nil {
		// 	return "", err
		// }

		// // Copy file to form body
		// _, err = io.Copy(part, inFile)
		// headers["Content-Type"] = writer.FormDataContentType()

		// body = b
	}

	if method == "GET" {
		data := url.Values{}
		if cmdParam != "" {
			data.Set(cmdParam, cmd)
		} else {
			headers[cmdHeader] = cmd
		}
		for k, v := range params {
			data.Set(k, v)
		}
		if strings.Contains(endpoint, "?") {
			finalURL = fmt.Sprintf("%s&%s", endpoint, data.Encode())
		} else {
			finalURL = fmt.Sprintf("%s?%s", endpoint, data.Encode())
		}
	} else {
		headers["Content-Type"] = "application/x-www-form-urlencoded"
		headers["Accept"] = "*/*"
		data := url.Values{}
		if cmdParam != "" {
			data.Set(cmdParam, cmd)
		} else {
			headers[cmdHeader] = cmd
		}
		for k, v := range params {
			data.Set(k, v)
		}
		body = strings.NewReader(data.Encode())
	}

	// Build HTTP request
	req, err := http.NewRequest(method, finalURL, body)
	if err != nil {
		return "", err
	}
	// Parse headers
	for k, v := range headers {
		req.Header.Add(k, v)
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// If downloading file...
	if strings.HasPrefix(cmd, "get ") {
		if resp.StatusCode == 404 {
			return "", errors.New("file not found")
		} else if resp.StatusCode != 200 {
			fmt.Println(resp.StatusCode)
			rBytes, _ := ioutil.ReadAll(resp.Body)
			response := string(rBytes)
			response = strings.Trim(response, " \n")
			return "", errors.New(response)
		}

		c := strings.Fields(cmd)
		fileName := c[1]
		destPath := fileName
		if len(c) > 2 {
			destPath = c[2]
		} else {
			f := strings.Split(destPath, "\\")
			destPath = f[len(f)-1]
			f = strings.Split(destPath, "/")
			destPath = f[len(f)-1]
		}

		outFile, err := os.Create(destPath)
		if err != nil {
			return "", err
		}
		defer outFile.Close()
		// b, err := ioutil.ReadAll(resp.Body)
		// outFile.WriteString(string(b))
		// return string(b), err

		io.Copy(outFile, resp.Body)

		return fmt.Sprintf("%s downloaded to %s.\n", fileName, destPath), nil
	}

	// Read server response
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(response), nil
}

func getHost(u string) (string, error) {
	t, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	return t.Hostname(), nil
}
