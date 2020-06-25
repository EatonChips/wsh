package cmd

// import (
// 	"bufio"
// 	"encoding/base64"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/chzyer/readline"
// 	"github.com/fatih/color"
// 	"github.com/spf13/cobra"
// )

// var (
// 	endpoint      string
// 	httpMethod    string
// 	commandParam  string
// 	commandHeader string
// 	headerFlags   []string
// 	paramFlags    []string
// 	timeout       int

// 	headers map[string]string
// 	params  map[string]string

// 	client http.Client
// )

// // connectCmd represents the connect command
// var connectCmd = &cobra.Command{
// 	Use:   "connect",
// 	Short: "Connect to a webshell",
// 	Long: `A longer description that spans multiple lines and likely contains examples
// and usage of using your command. For example:

// Cobra is a CLI library for Go that empowers applications.
// This application is a tool to generate the needed files
// to quickly create a Cobra application.`,
// 	Run: connect,
// 	Args: func(cmd *cobra.Command, args []string) error {
// 		if len(args) < 1 {
// 			return errors.New("url is required")
// 		}
// 		return nil
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(connectCmd)

// 	// Here you will define your flags and configuration settings.

// 	// Cobra supports Persistent Flags which will work for this command
// 	// and all subcommands, e.g.:
// 	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

// 	// Cobra supports local flags which will only run when this command
// 	// is called directly, e.g.:
// 	connectCmd.Flags().StringVar(&commandParam, "param", "c", "Parameter for sending command")
// 	connectCmd.Flags().StringVar(&commandHeader, "header", "", "Header for sending command")
// 	connectCmd.Flags().StringVarP(&httpMethod, "method", "X", "GET", "HTTP httpMethod (GET, POST, PUT, PATCH, DELETE)")
// 	connectCmd.Flags().StringSliceVarP(&headerFlags, "headers", "H", []string{}, "HTTP request headers")
// 	connectCmd.Flags().StringSliceVarP(&paramFlags, "params", "P", []string{}, "HTTP request parameters")
// 	connectCmd.Flags().IntVar(&timeout, "timeout", 10, "Request timeout in seconds")

// }

// func connect(cmd *cobra.Command, args []string) {
// 	endpoint = args[0]

// 	// // Check for invalid http methods
// 	// if httpMethod != "GET" && httpMethod != "POST" && httpMethod != "PUT" &&
// 	// 	httpMethod != "PATCH" && httpMethod != "DELETE" {
// 	// 	fmt.Println("Invalid HTTP httpMethod. Supported httpMethods are:")
// 	// 	fmt.Println("\tGET, POST, PUT, PATCH, and DELETE")
// 	// 	os.Exit(0)
// 	// }

// 	// Parse header flags
// 	headers = make(map[string]string)
// 	for _, h := range headerFlags {
// 		split := strings.Split(h, ":")
// 		if len(split) != 2 {
// 			fmt.Printf("Invalid header: \"%s\"\n", h)
// 			continue
// 		}

// 		headers[split[0]] = split[1]
// 	}

// 	// Parse parameter flags
// 	params = make(map[string]string)
// 	for _, p := range paramFlags {
// 		split := strings.Split(p, ":")
// 		if len(split) != 2 {
// 			fmt.Printf("Invalid parameter: \"%s\"\n", p)
// 			continue
// 		}

// 		params[split[0]] = split[1]
// 	}

// 	client = http.Client{
// 		Timeout: time.Duration(timeout) * time.Second,
// 	}

// 	host, err := getHost(endpoint)
// 	if err != nil {
// 		fmt.Println("Invalid url")
// 		return
// 	}

// 	cyan := color.New(color.FgGreen).SprintFunc()
// 	prompt := fmt.Sprintf(cyan("%s> "), host)

// 	l, err := readline.NewEx(&readline.Config{
// 		Prompt:          prompt,
// 		HistoryFile:     "/tmp/web-cli.tmp",
// 		InterruptPrompt: "^C",
// 		EOFPrompt:       "exit",

// 		HistorySearchFold: true,
// 	})
// 	if err != nil {
// 		fmt.Println(err.Error())
// 		return
// 	}
// 	defer l.Close()

// 	for {
// 		line, err := l.Readline()
// 		if err == readline.ErrInterrupt {
// 			if len(line) == 0 {
// 				break
// 			} else {
// 				continue
// 			}
// 		} else if err == io.EOF {
// 			break
// 		}

// 		switch {
// 		// case strings.HasPrefix(line, "get "):
// 		// 	fmt.Println("Download file")
// 		// 	fmt.Println(line)
// 		// case strings.HasPrefix(line, "clear"):
// 		// 	readline.ClearScreen()
// 		case line == "exit":
// 			l.Close()
// 		default:
// 			out, err := sendRequest(line)
// 			if err != nil {
// 				fmt.Printf("%s: %s", color.RedString("ERROR"), err.Error())
// 			}

// 			// if trimPrefix != "" {
// 			// 	i := strings.Index(out, trimPrefix)
// 			// 	if i > 0 {
// 			// 		i += len(trimPrefix)
// 			// 		out = out[i:]
// 			// 	}
// 			// }
// 			// if trimSuffix != "" {
// 			// 	i := strings.Index(out, trimSuffix)
// 			// 	if i > 0 {
// 			// 		out = out[:i]
// 			// 	}
// 			// }

// 			out = strings.TrimSpace(out)
// 			fmt.Println(out)
// 		}
// 	}
// }

// func sendRequest(cmd string) (string, error) {
// 	finalURL := endpoint
// 	var body io.Reader

// 	// If uploading file...
// 	if strings.HasPrefix(cmd, "put ") {
// 		c := strings.Fields(cmd)
// 		fileName := c[1]

// 		// Open file for reading
// 		inFile, err := os.Open(fileName)
// 		if err != nil {
// 			return "", err
// 		}
// 		defer inFile.Close()

// 		reader := bufio.NewReader(inFile)
// 		content, _ := ioutil.ReadAll(reader)

// 		params["f"] = base64.StdEncoding.EncodeToString(content)

// 		// Create multipart form
// 		// b := &bytes.Buffer{}
// 		// writer := multipart.NewWriter(b)

// 		// // Add file part
// 		// part, err := writer.CreateFormFile("f", filepath.Base(fileName))
// 		// if err != nil {
// 		// 	return "", err
// 		// }

// 		// // Copy file to form body
// 		// _, err = io.Copy(part, inFile)
// 		// headers["Content-Type"] = writer.FormDataContentType()

// 		// body = b
// 	}

// 	if httpMethod == "GET" {
// 		data := url.Values{}
// 		if commandParam != "" {
// 			data.Set(commandParam, cmd)
// 		} else {
// 			headers[commandHeader] = cmd
// 		}
// 		for k, v := range params {
// 			data.Set(k, v)
// 		}
// 		if strings.Contains(endpoint, "?") {
// 			finalURL = fmt.Sprintf("%s&%s", endpoint, data.Encode())
// 		} else {
// 			finalURL = fmt.Sprintf("%s?%s", endpoint, data.Encode())
// 		}
// 	} else {
// 		headers["Content-Type"] = "application/x-www-form-urlencoded"
// 		headers["Accept"] = "*/*"
// 		data := url.Values{}
// 		if commandParam != "" {
// 			data.Set(commandParam, cmd)
// 		} else {
// 			headers[commandHeader] = cmd
// 		}
// 		for k, v := range params {
// 			data.Set(k, v)
// 		}
// 		body = strings.NewReader(data.Encode())
// 	}

// 	// Build HTTP request
// 	req, err := http.NewRequest(httpMethod, finalURL, body)
// 	if err != nil {
// 		return "", err
// 	}
// 	// Parse headers
// 	for k, v := range headers {
// 		req.Header.Add(k, v)
// 	}

// 	// Send request
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	// If downloading file...
// 	if strings.HasPrefix(cmd, "get ") {
// 		if resp.StatusCode == 404 {
// 			return "", errors.New("file not found")
// 		} else if resp.StatusCode != 200 {
// 			fmt.Println(resp.StatusCode)
// 			rBytes, _ := ioutil.ReadAll(resp.Body)
// 			response := string(rBytes)
// 			response = strings.Trim(response, " \n")
// 			return "", errors.New(response)
// 		}

// 		c := strings.Fields(cmd)
// 		fileName := c[1]
// 		destPath := fileName
// 		if len(c) > 2 {
// 			destPath = c[2]
// 		} else {
// 			f := strings.Split(destPath, "\\")
// 			destPath = f[len(f)-1]
// 			f = strings.Split(destPath, "/")
// 			destPath = f[len(f)-1]
// 		}

// 		outFile, err := os.Create(destPath)
// 		if err != nil {
// 			return "", err
// 		}
// 		defer outFile.Close()
// 		// b, err := ioutil.ReadAll(resp.Body)
// 		// outFile.WriteString(string(b))
// 		// return string(b), err

// 		io.Copy(outFile, resp.Body)

// 		return fmt.Sprintf("%s downloaded to %s.\n", fileName, destPath), nil
// 	}

// 	// Read server response
// 	response, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	return string(response), nil
// }

// func getHost(u string) (string, error) {
// 	t, err := url.Parse(u)
// 	if err != nil {
// 		return "", err
// 	}

// 	return t.Hostname(), nil
// }
