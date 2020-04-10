package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type pdaList struct {
	// Note: field names must begin with capital letter for JSON
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	States          []string   `json:"states"`
	InputAlphabet   []string   `json:"inputAlphabet"`
	StackAlphabet   []string   `json:"stackAlphabet"`
	AcceptingStates []string   `json:"acceptingStates"`
	StartState      string     `json:"startState"`
	Transitions     [][]string `json:"transitions"`
	Eos             string     `json:"eos"`

	// Holds the current state.
	CurrentState string

	// Token at the top of the stack.
	CurrentStack string

	// This slice is used to hold the transition states tokens.
	TransitionStack []string

	// This slice is used to hold the token stack.
	TokenStack []string

	// This keeps a count of everytime put method is called
	PutCounter int

	// This keeps a count of everytime is_accepted method is called
	IsAcceptedCount int

	// This keeps a count of everytime peek method is called
	Peek int

	// This keeps a count for everytime a transition  is changed
	TransitionCounter int

	// This keeps a count for everytime current_state method is called
	CurrentStateCounter int

	// This checks if the input is accepted by the PDA
	IsAccepted bool
}

var Pda []pdaList

func showPdas(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllPda")
	json.NewEncoder(w).Encode(Pda)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func resetPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	for i := 0; i < len(Pda); i++ {
		if Pda[i].ID == id {
			fmt.Println("entered the if block")
			Pda[i].TokenStack = []string{}
			Pda[i].CurrentState = Pda[i].StartState
		}
	}
}

func createNewPda(w http.ResponseWriter, r *http.Request) {
	// unmarshal the body of POST request into new PDA struct and append this to our PDA array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	// fmt.Fprintf(w, "%+v", string(reqBody))
	var pda pdaList
	json.Unmarshal(reqBody, &pda)
	// update our global Pda array
	Pda = append(Pda, pda)

	json.NewEncoder(w).Encode(pda)
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/pdas", showPdas)
	myRouter.HandleFunc("/pdas/{id}", createNewPda)
	myRouter.HandleFunc("/pdas/{id}/reset", resetPDA)

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}
