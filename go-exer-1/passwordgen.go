package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
)

type PasswordType int

const (
	Random PasswordType = iota
	AlphaNumeric
	Pin
)

func main() {
	var passwordLength = flag.Int("length", 12, " length of the password (default to 12 if not provided)")
	var includeNumbersFlag = flag.Bool("includeNumbers", false, " A boolean flag indicating whether to include numbers in the password")
	var includeSymbolsFlag = flag.Bool("includeSymbols", false, " A boolean flag indicating whether to include symbols (e.g., !@#$%) in the password")
	var includeUppercaseFlag = flag.Bool("includeUppercase", false, " A boolean flag indicating whether to include uppercase letters in the password")
	var passwordType = flag.Int("type", 0, " 0 = random, 1 = alphanumeric, 2 = pin")
	var generatedPassword string = "password"

	flag.Parse()
	generatedPassword = generatePassword(*passwordLength, *includeNumbersFlag, *includeSymbolsFlag, *includeUppercaseFlag, *passwordType)

	fmt.Printf("Password rules are the following \n\n"+
		"Password Lenght : %d \n"+
		"Password Has Numbers : %t \n"+
		"Password Has Symbols : %t \n"+
		"Password use UpperCase letters : %t \n"+
		"Type : %v \n"+
		"Generate Password : %v \n", *passwordLength, *includeNumbersFlag, *includeSymbolsFlag, *includeUppercaseFlag, *passwordType, generatedPassword)
}

func generatePassword(passwordLength int, includeNumbersFlag bool, includeSymbolsFlag bool, includeUppercaseFlag bool, passwordType int) string {
	var chars string = "abcdefghijklmnopqrstuvwxyz"

	if includeNumbersFlag || passwordType == int(Random) {
		chars += "0123456789"
	}
	if includeSymbolsFlag || passwordType == int(Random) {
		chars += "!@#$%^&*()_+{}[]:;<>,.?/~`"
	}
	if includeUppercaseFlag || passwordType == int(Random) {
		chars += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	if passwordType == int(Pin) {
		chars = "0123456789"
	}

	password := make([]byte, passwordLength)
	for i := 0; i < passwordLength; i++ {
		randomIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		password[i] = chars[randomIndex.Int64()]
	}

	return string(password)
}
