The code pdaProcessor.go shows the transistions for a PDA.

Input-
The automata file is loaded from a specific file which has the following input format-

1: Current State
2: Input
3: Stack Symbols
4: Initial State 
5: Initial Stack 
6: Push/Pop Operations
7: Current Input Symbol 
8: Current Top of Stack
9: Next State

Currently 2 examples have been included-

1: helloPda.json: 0^n1^n
2: testPDA.json: 1^n0^n

Usage-
To run use -  go run pdaProcessor.go helloPda.json    	// if want to run example helloPda.json
OR
go run pdaProcessor.go testPDA.json    	// if want to run example testPDA.json
You can also run the bash file named PushDownAutomataBashHelloPDA.sh for grammar helloPda.json and PushDownAutomataBashTestPDA.sh for grammar testPDA.json. 
Make sure to give full path of the automata file, or better keep the automata file in the same directory as the program.