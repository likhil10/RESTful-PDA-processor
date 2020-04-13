package main

import (
	"encoding/json"
	"fmt"
)

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
func put(pda *PdaProcessor, position int, token string) {
	pda.PutCounter++
	var takeToken bool
	transitions := pda.Transitions
	transitionLength := len(transitions)
	pda.CurrentStack = "null"
	pda.IsAccepted = false
	if pda.PutCounter == 1 {
		putForTFirst(pda)
	}
	if position == 1 || position == (pda.LastPosition+1) {
		takeToken = true
		pda.LastPosition = position
	} else if takeToken == false {
		pda.HoldBackToken = append(pda.HoldBackToken, token)
		pda.HoldBackPosition = append(pda.HoldBackPosition, position)
	}
gotoPoint:
	if takeToken {
		takeToken = false
		fmt.Printf("hello sexy")
		for j := 1; j < transitionLength; j++ {
			t := transitions[j]

			if t[0] == pda.CurrentState && t[1] == token && t[2] == pda.CurrentStack {
				pda.IsAccepted = true
				pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)
				pda.TransitionCounter++
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
	if takeToken == false {
		for i := 0; i < len(pda.HoldBackPosition); i++ {
			if pda.HoldBackPosition[i] == (pda.LastPosition + 1) {
				takeToken = true
				token = pda.HoldBackToken[i]
				pda.LastPosition = pda.HoldBackPosition[i]
				pda.HoldBackPosition = append(pda.HoldBackPosition[:i], pda.HoldBackPosition[i+1:]...)
				pda.HoldBackToken = append(pda.HoldBackToken[:i], pda.HoldBackToken[i+1:]...)
				pda.CurrentStack = "null"
				pda.IsAccepted = false
				goto gotoPoint
			}
		}
	}
}

// Put method for the first transition with no input
func putForTFirst(pda *PdaProcessor) {
	if pda.Transitions[0][0] == pda.CurrentState {
		pda.TransitionStack = append(pda.TransitionStack, pda.CurrentState)
		pda.CurrentState = pda.Transitions[0][3]
		push(pda, pda.Transitions[0][4])
		pda.TransitionCounter++
	}
}

// Returns True if the PDA was succesfully satisfied
func isAccepted(pda *PdaProcessor) bool {
	pda.IsAcceptedCount++
	if len(pda.TokenStack) == 0 && pda.IsAccepted == true {
		return true
	}
	return false
}

// Return up to k stack tokens from the top of the stack (default k=1) without modifying the stack.
func peek(pda *PdaProcessor, k int) []string {
	pda.Peek++
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

// Arranges the HoldBAckToken stack in increasing order of their position
func queuedTokens(pda *PdaProcessor) {
	for i := 0; i < len(pda.HoldBackPosition)-1; i++ {
		for j := 1; j < len(pda.HoldBackPosition); j++ {
			if pda.HoldBackPosition[i] > pda.HoldBackPosition[j] {
				tempPosition := pda.HoldBackPosition[i]
				pda.HoldBackPosition[i] = pda.HoldBackPosition[j]
				pda.HoldBackPosition[j] = tempPosition
				tempToken := pda.HoldBackToken[i]
				pda.HoldBackToken[i] = pda.HoldBackToken[j]
				pda.HoldBackToken[j] = tempToken
			}
		}
	}
}

// Returns the current state
func currentState(pda *PdaProcessor) string {
	pda.CurrentStateCounter++
	return pda.CurrentState
}

// Garbage disposal method
func close() {

}
