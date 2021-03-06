package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
)

func parse(urlString string, part string, index int, field string) (string, error) {
	var result string

	url, err := url.Parse(urlString)
	if err != nil {
		return result, err
	}

	switch part {
	case "scheme":
		result = url.Scheme

	case "user":
		if url.User != nil {
			result = url.User.Username()
		}

	case "password":
		if url.User == nil {
			break
		}

		pass, hasPassword := url.User.Password()
		if hasPassword {
			result = pass
		}

	case "hostname":
		result = url.Hostname()

	case "port":
		result = url.Port()

	case "path":
		result = url.Path
		if index > -1 {
			paths := strings.Split(result, "/")
			result = paths[index+1]
		}

	case "query":
		result = url.RawQuery
		if field != "" && result != "" {
			result = url.Query().Get(field)
		}

	case "fragment":
		result = url.Fragment

	case "all":
		result = urlString

	default:
		err = errors.New("Please provides valid part name")
	}

	return result, err
}

func usage() {
	appName := "url-parser"
	fmt.Printf("%s\n", appName)
	fmt.Printf("  Parse URL and shows the part of it.\n\n")
	fmt.Println("Usage:")
	fmt.Printf("  %s --part=PART <url>\n", appName)
	fmt.Printf("  %s --part=path [--path-index=INDEX] <url>\n", appName)
	fmt.Printf("  %s --part=query [--query-field=FIELD] <url>\n\n", appName)
	fmt.Println("Options:")
	fmt.Println("  --part=PART          Part of URL to show [default: all].")
	fmt.Println("                       Valid values: all, scheme, user, password,")
	fmt.Println("                       hostname, port, path, query, or fragment.")
	fmt.Println("  --path-index=INDEX   Filter parsed path by index.")
	fmt.Println("  --query-field=FIELD  Filter parsed query string by field name.")
}

func main() {
	flag.Usage = usage

	partPtr := flag.String("part", "all", "Part of URL to show")
	indexPtr := flag.Int("path-index", -1, "Filter parsed path by index")
	fieldPtr := flag.String("query-field", "", "Filter parsed query string by field name")

	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Printf("Please provides URL to parse")
		os.Exit(1)
	}
	urlString := flag.Args()[0]
	if urlString == "-" {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		urlString = strings.Trim(input, "\n\r")
	}

	result, err := parse(urlString, *partPtr, *indexPtr, *fieldPtr)
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
