package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const programTemplate = `// Generated by cc_code_gen
#include <iostream>
#include <cstdint>
#include <vector>

using CellType = wchar_t;

const size_t kMemorySize = 65535;

auto Compute(std::vector<CellType>& memory) -> void ;

auto main() -> int{
	auto memory=std::vector<CellType>();
	memory.resize(kMemorySize);
	Compute(memory);
	return EXIT_SUCCESS;
}

`

func compileBf(input string) (res string, err error) {
	res = programTemplate

	functionCompute := `
auto Compute(std::vector<CellType>& memory) -> void {
	size_t ptr_pos = 0;
`

	lBuckets := 0

	addTabs := func() {
		functionCompute += `	`
		for i := 0; i < lBuckets; i++ {
			functionCompute += `	`
		}
	}

	localPtrCounter := 0
	localValueCounter := 0

	clearLocalPtrCounter := func() {
		if localPtrCounter != 0 {
			addTabs()
			functionCompute += fmt.Sprintf("ptr_pos += (%v); \n", localPtrCounter)
		}
		localPtrCounter = 0
	}

	clearLocalValueCounter := func() {
		if localValueCounter != 0 {
			addTabs()
			functionCompute += fmt.Sprintf("memory[ptr_pos] += (%v); \n", localValueCounter)
		}
		localValueCounter = 0
	}
	for _, ch := range input {
		switch ch {
		case '>':
			clearLocalValueCounter()
			localPtrCounter++
		case '<':
			clearLocalValueCounter()
			localPtrCounter--
		case '+':
			clearLocalPtrCounter()
			localValueCounter++
		case '-':
			clearLocalPtrCounter()
			localValueCounter--
		case ',':
			clearLocalPtrCounter()
			clearLocalValueCounter()
			addTabs()
			functionCompute += "std::wcin >> memory[ptr_pos]; \n"
		case '.':
			clearLocalPtrCounter()
			clearLocalValueCounter()
			addTabs()
			functionCompute += "std::wcout << memory[ptr_pos]; \n"
		case '[':
			clearLocalPtrCounter()
			clearLocalValueCounter()
			addTabs()
			lBuckets++
			functionCompute += "while(memory[ptr_pos]) { \n"
		case ']':
			clearLocalPtrCounter()
			clearLocalValueCounter()
			lBuckets--
			addTabs()
			functionCompute += "} \n"
		}
	}

	functionCompute += `}`

	res += functionCompute

	if lBuckets != 0 {
		err = errors.New("compilation error")
	}
	return
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Printf("Usage: %s filename\n", args[0])
		return
	}
	filename := args[1]
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicf("error reading %s\n", filename)
		return
	}
	program, err := compileBf(string(fileContents))
	if err != nil {
		log.Panicln(err)
		return
	}
	err = ioutil.WriteFile(filename+".cc", []byte(program), 0666)
	if err != nil {
		log.Panicln(err)
	}
}
