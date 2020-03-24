package main

import (
	"fmt"
	"os"
	"io/ioutil"
)

// The main method of the pdaDriver program. This will check the command-line args, import the json,
// and pass the input string and json spec to an instance of pdaProcessor.
func main(){
	// Check to make sure the user has provided a path for the json spec.
	if len(os.Args) < 2{
		fmt.Println("Error: command-line args must include json spec file path")
		os.Exit(0)
	}
	jsonFilename := string(os.Args[1])
	jsonText, err := ioutil.ReadFile(jsonFilename)
	
	// check(err)
	fmt.Println("hello")

	pda := new(PdaProcessor)
	if pda.Open(string(jsonText)){
		fmt.Println(pda)
	} else {
		fmt.Println("Error: could not open json spec")
	}

}

// A function that calls panic if it detects an error.
// func check(e error){
// 	if e != nil{
// 		panic(e)
// 	}
// }

