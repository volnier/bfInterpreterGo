package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

type Program struct {
	size         int
	instructions []byte
	// array index
	ai int
}

func main() {
	var filename string
	flag.StringVar(&filename, "filename", "main.b", "name of the .b file to execute")
	flag.Parse()
	file, err := ioutil.ReadFile(filename)
	handleErr(err)

	code := string(file)
	// check for wrong characters inside code using regex
	expression := regexp.MustCompile(`\r?\n[^.,\-+<>\[\]]`)
	code = expression.ReplaceAllString(code, "")

	// execute the code
	exec(code)
}

func handleErr(i error) {
	if i != nil {
		log.Panic(i)
	}
}

func exec(code string) {
	var program = new(Program)
	program.size = 30000
	program.instructions = make([]byte, program.size, program.size)
	program.ai = 0

	interpret(program, code)
}

func interpret(program *Program, code string) {
	// loop
	var startLoop = -1
	var endLoop = -1
	var ignore = 0
	var skipEndLoop = 0
	for i, char := range code {
		if ignore == 1 {
			if char == '[' {
				skipEndLoop += 1
			} else if char == ']' {
				if skipEndLoop != 0 {
					skipEndLoop -= 1
					continue
				}
				endLoop = i
				ignore = 0
				// check and correct for []
				if startLoop == endLoop {
					startLoop = -1
					endLoop = -1
					continue
				}
				loop := code[startLoop:endLoop]
				for program.instructions[program.ai] > 0 {
					interpret(program, loop)
				}
			}
			continue
		}

		// instructions

		// - decrement byte at the data pointer, + increment byte at the data pointer, < decrement the data pointer, > increment the data pointer
		// . output byte at the data pointer, , accept and store a byte of input
		// [ if byte at the data pointer = 0 jump to the next command after ], ] if byte at the data pointer != 0 jump back to the next command after [

		if char == '-' {
			program.instructions[program.ai] -= 1
		} else if char == '+' {
			program.instructions[program.ai] += 1
		} else if char == '<' {
			if program.ai == 0 {
				program.ai = program.size - 1
			} else {
				program.ai -= 1
			}
		} else if char == '>' {
			if program.ai == program.size-1 {
				program.ai = 0
			} else {
				program.ai += 1
			}
		} else if char == '.' {
			fmt.Printf("%c", rune(program.instructions[program.ai]))
		} else if char == ',' {
			reader := bufio.NewReader(os.Stdin)
			input, _, err := reader.ReadRune()
			handleErr(err)
			program.instructions[program.ai] = byte(input)
		} else if char == '[' {
			startLoop = i + 1
			ignore = 1
		}
	}
}
