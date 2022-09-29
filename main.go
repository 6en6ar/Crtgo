package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var url = "https://crt.sh/"
var banner = `

 ______     ______     ______   ______     ______    
/\  ___\   /\  == \   /\__  _\ /\  ___\   /\  __ \   
\ \ \____  \ \  __<   \/_/\ \/ \ \ \__ \  \ \ \/\ \  
 \ \_____\  \ \_\ \_\    \ \_\  \ \_____\  \ \_____\ 
  \/_____/   \/_/ /_/     \/_/   \/_____/   \/_____/ 
                                                     
                                                                                                              
Coded by 6en6ar 3:)

--- pull subdomains from crt.sh

`

type Crt struct {
	Name_value  string
	Common_name string
}

func removeDuplicate(subdomains []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range subdomains {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

var domain = flag.String("d", "example.com", "Domain to be queried")

func usage() {
	fmt.Println(banner)
	fmt.Printf("usage : -d <DOMAIN>\n")
	os.Exit(0)
}
func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NFlag() == 0 {
		usage()
		os.Exit(1)
	}
	fmt.Println(banner)
	var subdomains []string
	var jsonData []Crt

	finalUrl := url + "?q=" + *domain + "&output=json"
	fmt.Println("Fetching from --> " + finalUrl)
	resp, err := http.Get(finalUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &jsonData); err != nil {
		fmt.Println(err)
	}
	for _, i := range jsonData {
		if strings.HasPrefix(i.Name_value, "*") {
			continue
		}
		var sd []string = regexp.MustCompile("\r?\n").Split(i.Name_value, -1)
		for _, s := range sd {
			subdomains = append(subdomains, s)
		}

	}
	var subs = removeDuplicate(subdomains)
	fmt.Printf("Number of subdomains found --> %d", len(subs))
	fmt.Println()
	file, err := os.OpenFile(*domain+"-subdomains.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Printf("failed creating file: %s", err)
		fmt.Println()
	}

	w := bufio.NewWriter(file)

	for _, sd := range subs {
		_, _ = w.WriteString(sd + "\n")
	}

	w.Flush()
	fmt.Println("[ + ] Domains written to --> " + file.Name())
	file.Close()
}
