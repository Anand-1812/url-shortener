package main

import ( 
	"http"
	"fmt"
)

const PORT = ":8080"

func main()  {	
	serveErr := http.ListenAndServe(PORT, nil)
	if serveErr != nil {
		fmt.Println("server error: %v", error)
	}
}
