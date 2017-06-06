package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"strings"
	"regexp"
	"bufio"
	"os"
)

func counting_go(url string, c_count chan int)  {
	response, err1 := http.Get(url)
        if err1 != nil {
                fmt.Println("Error: can't open ", url)
		c_count <- 0
		return
        } else {
                defer response.Body.Close()
        }
        doc, err2 := ioutil.ReadAll(response.Body)
        if err2 != nil {
                fmt.Println("Error: can't read ", url)
		c_count <- 0
		return
        }

	var count int = 0
	//Look for "go" in a body of website
        src := string(doc)[strings.Index(string(doc), "<body>"):strings.Index(string(doc), "</body>")]

	//Search of "go" on the website
        re1 := regexp.MustCompile("[ |\n|\t|\r|>][G|g][O|o][ |\n|\t|\r|<]")
        count = count + len(re1.FindAllString(src, -1))

	//Special case: "go" can be the first on the website
        re2 := regexp.MustCompile("^[G|g][O|o][ |\n|\t|\r|<]")
        count = count + len(re2.FindAllString(src, -1))

	//Special case: "go" can be the last on the website
        re3 := regexp.MustCompile("[ |\n|\t|\r|>][G|g][O|o]$")
        count = count + len(re3.FindAllString(src, -1))

        fmt.Println("Count for ", url, ": ", count)
	c_count <- count
}

func worker() {
	var count int = 0
	c_count := make(chan int, 5)
	var url string = ""
	myscanner := bufio.NewScanner(os.Stdin)
        myscanner.Scan()
        url = myscanner.Text()
        for url != "" {
		go counting_go(url, c_count)
		count = count + <-c_count
                myscanner.Scan()
                url = myscanner.Text()
        }
	fmt.Println("Total: ", count)
}

func main() {
	worker()
}
