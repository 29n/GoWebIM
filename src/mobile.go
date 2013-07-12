package main

import (
	"net.http"
)

func main() {

    cli := new(http.Client)
    head := new(http.Head)
    
    head.Set("User-Agent", "")
    
	res, err := http.Get("http://www.google.com/robots.txt")
	if err != nil {
		log.Fatal(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", robots)
}
