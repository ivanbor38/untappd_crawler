package main

import (
	"fmt"
	. "crawler"
	. "read_config"
)

func main() {

	start_id, depth := ReadConfig()

	Total("https://untappd.com/user/", start_id, depth)
	fmt.Printf("Done")

}
