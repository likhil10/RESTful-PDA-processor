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

// Structure of type PdaProcessor.
type PdaProcessor struct {
	// Note: field names must begin with capital letter for JSON
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	States          []string   `json:"states"`
	InputAlphabet   []string   `json:"input_alphabet"`
	StackAlphabet   []string   `json:"stack_alphabet"`
	AcceptingStates []string   `json:"accepting_states"`
	StartState      string     `json:"start_state"`
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
	Is_Accepted int

	// This keeps a count of everytime peek method is called
	Peek int

	// This keeps a count for everytime a transition  is changed
	TransitionCounter int

	// This keeps a count for everytime current_state method is called
	CurrentStateCounter int

	// This checks if the input is accepted by the PDA
	IsAccepted bool
}

type JSONMessage struct {
	cs    string
	qt    []string
	peekK []string
}

// Unmarshals the jsonText string. Returns true if it succeeds.
func (pda *PdaProcessor) open(jsonText string) bool {

	if err := json.Unmarshal([]byte(jsonText), &pda); err != nil {
		check(err)
	}

	// Validate input.
	if len(pda.Name) == 0 || len(pda.States) == 0 || len(pda.InputAlphabet) == 0 ||
		len(pda.StackAlphabet) == 0 || len(pda.AcceptingStates) == 0 || len(pda.StartState) == 0 ||
		len(pda.Transitions) == 0 || len(pda.Eos) == 0 {
		return false
	}

	return true
}

// Sets the CurrentState to StartState and assigns Stack a new empty slice
func reset(pda *PdaProcessor) {
	pda.CurrentState = pda.StartState
	pda.TokenStack = []string{}
}

//  Consumes the token, takes appropriate transition(s)
func put(pda *PdaProcessor, char string) {
	pda.PutCounter += 1
	transitions := pda.Transitions
	transition_length := len(transitions)
	if pda.PutCounter == 1 {
		putForTFirst(pda)
	}

	for j := 1; j < transition_length; j++ {
		t := transitions[j]

		if t[0] == pda.CurrentState && t[1] == char && t[2] == pda.CurrentStack {
			pda.IsAccepted = true
			pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)
			pda.TransitionCounter += 1
			pda.CurrentState = t[3]
			pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)

			if t[4] != "null" {
				push(pda, t[4])
				pda.CurrentStack = t[4]
			} else {
				if len(pda.TokenStack) == 0 {
					pda.IsAccepted = false
					break
				} else {
					pop(pda)
					break
				}
			}
		}

		if len(pda.TokenStack) > 1 {
			pda.CurrentStack = pda.TokenStack[len(pda.TokenStack)-1]
		} else {
			break
		}
	}
}

// Put method for the first transition with no input
func putForTFirst(pda *PdaProcessor) {
	if pda.Transitions[0][0] == pda.CurrentState {
		pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)
		pda.CurrentState = pda.Transitions[0][3]
		push(pda, pda.Transitions[0][4])
		pda.TransitionCounter += 1
	}
}

// Returns True if the PDA was succesfully satisfied
func is_accepted(pda *PdaProcessor) bool {
	pda.Is_Accepted += 1
	if len(pda.TokenStack) == 0 && pda.IsAccepted == true {
		return true
	} else {
		return false
	}
}

// Return up to k stack tokens from the top of the stack (default k=1) without modifying the stack.
func peek(pda *PdaProcessor, k int) []string {
	pda.Peek += 1
	if len(pda.TokenStack) > 0 {
		if len(pda.TokenStack) < k {
			return pda.TokenStack
		} else if len(pda.TokenStack) > k {
			x := len(pda.TokenStack) - (k - 1)
			return pda.TokenStack[x-1:]
		} else if len(pda.TokenStack) == k {
			return pda.TokenStack[:k]
		}
	}
	return pda.TokenStack
}

// Adds an input token to the stack
func push(pda *PdaProcessor, x string) {
	pda.TokenStack = append(pda.TokenStack, x)
}

// Removes an input token from the last of the stack
func pop(pda *PdaProcessor) {
	pda.TokenStack = pda.TokenStack[:len(pda.TokenStack)-1]
}

// A function that calls panic if it detects an error.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Declares the end of string
func eos(pda *PdaProcessor) {
	if len(pda.TransitionStack) > 0 && pda.TransitionStack[0] == "q1" && pda.TransitionStack[len(pda.TransitionStack)-1] == "q4" {
		fmt.Println("pda=", pda.Name, ":method=eos:: Reached the End of String")
	} else {
		fmt.Println("pda=", pda.Name, ":method=eos::Did not reach the end of string but EOS was called.")
	}
}

// Returns the current state
func current_state(pda *PdaProcessor) string {
	pda.CurrentStateCounter += 1
	return pda.CurrentState
}

// Garbage disposal method
func close() {

}

var Pda []PdaProcessor

func showPdas(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllPda")
	//fmt.Fprintf(w, "Here are all the PDAs on the server!")
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
			// Pda[i].Stack = make([]string, 0)
			// Pda[i].CurrentState = pdaList.StartState
		}
	}
}

func createNewPda(w http.ResponseWriter, r *http.Request) {
	// unmarshal the body of POST request into new PDA struct and append this to our PDA array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	// fmt.Fprintf(w, "%+v", string(reqBody))
	var pda PdaProcessor
	json.Unmarshal(reqBody, &pda)
	// update our global Pda array
	Pda = append(Pda, pda)

	json.NewEncoder(w).Encode(pda)
}

func eosPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var pos = vars["position"]
	l, err := strconv.Atoi(pos)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(Pda); i++ {
		if Pda[i].ID == id {
			if len(Pda[i].TokenStack) > l { // Removes tokens after given position (excluding) and calls eos()
				for j := 0; j < (len(Pda[i].TokenStack) - l); j++ {
					pop(&Pda[i])
				}
				eos(&Pda[i])
			} else {
				fmt.Println(w, "Tokens till given position not consumed yet")
			}
		}
	}
}

func isAcceptedPDA(w http.ResponseWriter, r *http.Request) {
	var vars = mux.Vars(r)
	var id = vars["id"]
	var accepted bool

	for i := 0; i < len(Pda); i++ {
		if Pda[i].ID == id {
			accepted = is_accepted(&Pda[i])
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

	for i := 0; i < len(Pda); i++ {
		if Pda[i].ID == id {
			returnStack = peek(&Pda[i], k)
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

	for i := 0; i < len(Pda); i++ {
		if Pda[i].ID == id {
			length = len(Pda[i].TokenStack)
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

	for i := 0; i < len(Pda); i++ {
		if Pda[i].ID == id {
			cs = current_state(&Pda[i])
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

	for i := 0; i < len(Pda); i++ {
		if Pda[i].ID == id {
			message.cs = current_state(&Pda[i])
			// message.qt = <Call function to return queued_tokens>
			message.peekK = peek(&Pda[i], k)
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
