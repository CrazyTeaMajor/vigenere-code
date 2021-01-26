package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	numberSize int = 10
	alphaSize  int = 26
	Seed       int = 91
)

var alphaCodeList [alphaSize][alphaSize]int

var codeStr string
var isHelp bool
var modeType, cmdText string
var inputFilename, outputFilename string
var codeLength int

/********************init function***********************/

func initCodeList() {
	for i := 0; i < alphaSize; i++ {
		for j := 0; j < alphaSize; j++ {
			alphaCodeList[i][j] = (i + j) % alphaSize
		}
	}
}

func initCmd() {
	key := flag.String("key", "", "Please input the key!")
	text := flag.String("text", "", "Please input the text!")
	iFile := flag.String("i", "", "The input filename.")
	oFile := flag.String("o", "", "The output filename.")
	mode := flag.String("mode", "", "The operation mode: decode or encode")
	help := flag.Bool("help", false, "Ask for help.")
	h := flag.Bool("h", false, "Ask for help.")
	flag.Parse()
	codeStr = *key
	isHelp = *help || *h
	modeType = *mode
	inputFilename = *iFile
	outputFilename = *oFile
	codeLength = len(codeStr)
	cmdText = *text
}

func mainInit() {
	initCmd()
	initCodeList()
}

/******************encode and other functions*************/

func isLetter(char byte) bool {
	if 'a' <= char && char <= 'z' {
		return true
	} else if 'A' <= char && char <= 'Z' {
		return true
	} else {
		return false
	}
}

func isNumber(char byte) bool {
	if '0' <= char && char <= '9' {
		return true
	} else {
		return false
	}
}

func getPos(char byte) int {
	if isNumber(char) == true {
		return (int)(char - '0')
	} else if isLetter(char) == true {
		if 'a' <= char {
			return int(char - 'a')
		} else {
			return int(char - 'A')
		}
	} else {
		err := errors.New("The letter has something wrong!")
		panic(err)
		return -1
	}
}

func codeLetter(char byte, pos int) byte {
	var tx, ty int
	if modeType == "encode" {
		tx = getPos(codeStr[pos%codeLength])
		ty = getPos(char)
		if 'a' <= char {
			return (byte)('a' + alphaCodeList[tx][ty])
		} else {
			return (byte)('A' + alphaCodeList[tx][ty])
		}
	} else if modeType == "decode" {
		tx = getPos(codeStr[pos%codeLength])
		ty = getPos(char)
		for i := 0; i < alphaSize; i++ {
			if alphaCodeList[tx][i] == ty {
				if 'a' <= char {
					return byte('a' + i)
				} else {
					return byte('A' + i)
				}
			}
		}
	} else {
		err := errors.New("Mode type wrong!")
		panic(err)
	}
	return char
}

func codeNumber(char byte, pos int) byte {
	var tx, ty int
	if modeType == "encode" {
		tx = getPos(codeStr[pos%codeLength])
		ty = getPos(char)
		return byte('0' + alphaCodeList[tx][ty]%numberSize)
	} else if modeType == "decode" {
		tx = getPos(codeStr[pos%codeLength])
		ty = getPos(char)
		for i := 0; i < alphaSize; i++ {
			if alphaCodeList[tx][i]%numberSize == ty {
				return byte('0' + i%numberSize)
			}
		}
	} else {
		err := errors.New("Mode type wrong!")
		panic(err)
	}
	return char
}

func codeText(char byte, pos int) byte {
	var codeChar byte
	if isLetter(char) == true {
		codeChar = codeLetter(char, pos)
	} else if isNumber(char) == true {
		codeChar = codeNumber(char, pos)
	} else {
		codeChar = char
	}
	return codeChar
}

func help() {
	fmt.Println("This is the help manual:(both '-' and '--' can be accepted)")
	fmt.Println("use -help or -h to look at the help manual.")
	fmt.Println("-key xxx: The xxx is your secret key to encode or decode the text.")
	fmt.Println("-text xxx: The xxx is the text you want to encode or decode.")
	fmt.Println("-mode xxx: The xxx is 'encode' or 'decode' to choose what you want to do.")
	fmt.Println("-i xxx: The xxx is the file's path which you want to encode or decode.")
	fmt.Println("-o xxx: The xxx is the file's path and name that you want to save the output result.")
	fmt.Println("Tips:")
	fmt.Println("Please do not include space or enter in xxx.")
	fmt.Println("The file's default path is the progam's path.")
}

func isFileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func mainWork() {
	var textStr string
	if inputFilename != "" {
		file, err := os.Open(inputFilename)
		if err != nil {
			err = errors.New("file error!")
			panic(err)
		}
		defer file.Close()
		var data []byte
		buffer := make([]byte, 1024*1024*1024)
		for {
			tmp, err := file.Read(buffer)
			data = append(data, buffer[:tmp]...)
			if err != io.EOF {
				break
			}
		}
		textStr = string(data)
	} else {
		textStr = cmdText
	}
	var outputText, temp string
	for pos := range textStr {
		temp = fmt.Sprintf("%c", codeText(textStr[pos], pos))
		outputText += temp
	}
	if outputFilename != "" {
		if isFileExist(outputFilename) {
			var opt string
			fmt.Println("The file already exists.Do you want to cover it? [y/n] (y for yes and n for no)")
			fmt.Scanf("%s\n", &opt)
			if opt != "y" {
				fmt.Println("Operation cancelled")
				return
			}
		}
		file, err := os.OpenFile(outputFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0)
		if err != nil {
			err = errors.New("File open error!")
			panic(err)
		}
		defer file.Close()
		_, err = file.WriteString(outputText)
		if err != nil {
			err = errors.New("file write error")
			panic(err)
		}
	} else {
		fmt.Println("The result:")
		fmt.Println(outputText)
	}
}

func main() {
	mainInit()
	if isHelp == true || flag.NFlag() <= 2 {
		help()
	} else {
		mainWork()
	}
}
