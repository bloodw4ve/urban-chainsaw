package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

var (
	re                   = regexp.MustCompile(`(^[-+]{0,1}[0-9]{1,}[.]{0,1}[0-9]{0,})([-+*/]{0,1})([-+]{0,1}[0-9]{1,}[.]{0,1}[0-9]{0,})(=)([?])$`)
	ws                   = regexp.MustCompile(`\s+`)
	errNotMathExpression = fmt.Errorf("not math expression")
	errDivByZero         = fmt.Errorf("impossible to divide by zero")
	errUnknownOperator   = fmt.Errorf("unknown math operator")
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("To start program you have to enter the input file name and the result file name")
		return
	}

	inputFileName := os.Args[1]
	resultFileName := os.Args[2]

	inputFile, err := os.OpenFile(inputFileName, os.O_RDONLY, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer inputFile.Close()

	_ = os.Remove(resultFileName)
	f, err := os.OpenFile(resultFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := bufio.NewReader(inputFile)

	for {
		line, _, err := r.ReadLine()
		if err != nil {
			break
		}
		el, err := splitMathExpression(string(line))
		if err != nil {
			continue
		}

		firstNum, err := strconv.ParseFloat(el[0][1], 64)
		if err != nil {
			f.WriteString(string(line) + "failed to convert " + el[0][1] + "\n")
			continue
		}

		secondNum, err := strconv.ParseFloat(el[0][3], 64)
		if err != nil {
			f.WriteString(string(line) + "failed to convert " + el[0][3] + "\n")
			continue
		}

		operation := el[0][2]
		result, err := calculate(firstNum, secondNum, operation)
		if err != nil {
			f.WriteString(string(line) + " " + err.Error() + "\n")
			continue
		}
		f.WriteString(fmt.Sprintf("%g", firstNum) + operation + fmt.Sprintf("%g", secondNum) + el[0][4] + fmt.Sprintf("%g", result) + "\n")
	}
}

func splitMathExpression(s string) ([][]string, error) {
	trimmed := ws.ReplaceAllString(s, "")
	elements := re.FindAllStringSubmatch(trimmed, -1)
	if elements == nil {
		return make([][]string, 0), errNotMathExpression
	}
	return elements, nil
}

func calculate(firstNum, secondNum float64, operation string) (float64, error) {
	switch {
	case operation == "/" && secondNum == 0:
		return 0, errDivByZero
	case operation == "-":
		return firstNum - secondNum, nil
	case operation == "+":
		return firstNum + secondNum, nil
	case operation == "*":
		return firstNum * secondNum, nil
	case operation == "/":
		return firstNum / secondNum, nil
	}
	return 0, errUnknownOperator
}
