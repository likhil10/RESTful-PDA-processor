package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// PdaProcessor structure.
type PdaProcessor struct {
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

	// for storing the positions
	HoldBackPosition []int

	// for the tokens that were not consumed and held back due to postion
	HoldBackToken []string

	// to store the position of the token last consumed
	LastPosition int
}

// TokenList struct
type TokenList struct {
	// takes the tokens from the http request body
	Tokens string `json:"tokens"`
}

// JSONMessage Structure
type JSONMessage struct {
	curState    string
	quToken    []string
	peekK []string
}

var pdaArr []PdaProcessor
var tokenArr []TokenList
var positionArr []int

func showPdas(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllPda")
	json.NewEncoder(w).Encode(pdaArr)
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func resetPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			fmt.Println("entered the if block")
			pdaArr[i].TokenStack = []string{}
			pdaArr[i].CurrentState = pdaArr[i].StartState
		}
	}
}

func createNewPda(w http.ResponseWriter, r *http.Request) {
	// unmarshal the body of PUT request into new PDA struct and append this to our PDA array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	params := mux.Vars(r)
	var pda PdaProcessor
	json.Unmarshal(reqBody, &pda)
	if len(pdaArr) > 0 {
		for i := 0; i < len(pdaArr); i++ {
			if pdaArr[i].ID == params["id"] {
				panic("This PDA already exists")
			} else {
				// update our global pdaArr array
				pdaArr = append(pdaArr, pda)
			}
		}
	} else {
		// update our global pdaArr array
		pdaArr = append(pdaArr, pda)
	}
}

func putPda(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	params := mux.Vars(r)
	var token TokenList
	position := params["position"]
	id := params["id"]
	json.Unmarshal(reqBody, &token)
	positionInt, err := strconv.Atoi(position)
	if err == nil {
		for i := 0; i < len(pdaArr); i++ {
			if pdaArr[i].ID == id {
				put(&pdaArr[i], positionInt, token.Tokens)
				break
			} else {
				fmt.Printf("Error finding given PDA")
			}
		}
	}
}

func getTokens(w http.ResponseWriter, r *http.Request) {
	// parse the path parameters
	vars := mux.Vars(r)
	// extract the `id` of the pda we wish to delete
	id := vars["id"]

	// we then need to loop through all our pdas
	for _, pda := range pdaArr {
		// if our id path parameter matches one of the pdas
		if pda.ID == id {
			// call the queueTokens() function
			queuedTokens(&pda)
			json.NewEncoder(w).Encode(pda.HoldBackToken)
			break
		}
	}
}

func deletePda(w http.ResponseWriter, r *http.Request) {
	// parse the path parameters
	vars := mux.Vars(r)
	// extract the `id` of the pda we wish to delete
	id := vars["id"]

	// we then need to loop through all our pdas
	for index, pda := range pdaArr {
		// if our id path parameter matches one of the pdas
		if pda.ID == id {
			// updates our pdaArray array to remove the pda
			pdaArr = append(pdaArr[:index], pdaArr[index+1:]...)
			break
		}
	}

}

func eosPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var pos = vars["position"]
	l, err := strconv.Atoi(pos)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			if len(pdaArr[i].TokenStack) == l {
				eos(&pdaArr[i])
			} else {
				fmt.Println(w, "Incorrect Position Value Passed")
			}
		}
	}
}

func isAcceptedPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var accepted bool

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			accepted = isAccepted(&pdaArr[i])
		}
	}
	json.NewEncoder(w).Encode(accepted)
	//fmt.Fprintln(w, strconv.FormatBool(accepted))
	return
}

func stackTopPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var kStr = vars["k"]
	k, err := strconv.Atoi(kStr)
	if err != nil {
		log.Fatal(err)
	}

	var returnStack []string

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			returnStack = peek(&pdaArr[i], k)
		}
	}
	json.NewEncoder(w).Encode(returnStack)
	//fmt.Fprintln(w, returnStack)
	return
}

func stackLenPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var length int

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			length = len(pdaArr[i].TokenStack)
		}
	}
	json.NewEncoder(w).Encode(length)
	//fmt.Fprintln(w, length)
	return
}

func statePDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var cs string

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			cs = currentState(&pdaArr[i])
		}
	}
	json.NewEncoder(w).Encode(cs)
	//fmt.Fprintln(w, cs)
	return
}

func snapshotPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var kStr = vars["k"]

	// Convert string k to int
	k, err := strconv.Atoi(kStr)
	if err != nil {
		log.Fatal(err)
	}

	var message JSONMessage

	for i := 0; i < len(pdaArr); i++ {
		if pdaArr[i].ID == id {
			message.curState = currentState(&pdaArr[i])
			message.quToken = pdaArr[i].HoldBackToken
			message.peekK = peek(&pdaArr[i], k)
		}
	}
	json.NewEncoder(w).Encode(message)
	return
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homepage)
	myRouter.HandleFunc("/pdas", showPdas).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}", createNewPda).Methods("PUT")
	myRouter.HandleFunc("/pdas/{id}/reset", resetPDA).Methods("PUT")
	myRouter.HandleFunc("/pdas/{id}/tokens/{position}", putPda).Methods("PUT")
	myRouter.HandleFunc("/pdas/{id}/tokens", getTokens).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/delete", deletePda).Methods("DELETE")

	//Rhea's APIs

	myRouter.HandleFunc("/pdas/{id}/eos/{position}", eosPDA).Methods("PUT")
	myRouter.HandleFunc("/pdas/{id}/is_accepted", isAcceptedPDA).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/stack/top/{k}", stackTopPDA).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/stack/len", stackLenPDA).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/state", statePDA).Methods("GET")
	myRouter.HandleFunc("/pdas/{id}/snapshot/{k}", snapshotPDA).Methods("GET")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	handleRequests()
}
