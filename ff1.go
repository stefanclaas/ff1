package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
	"github.com/capitalone/fpe/ff1"
)

func main() {
	keyfile := flag.String("k", "", "Path to the key file (128, 192 or 256 bit, hex encoded)")
	tweakfile := flag.String("t", "", "Path to the tweak file (64 bit, hex encoded)")
	decryptFlag := flag.Bool("d", false, "Specify whether to decrypt the data")

	flag.Parse()

	if *keyfile == "" || *tweakfile == "" {
		fmt.Println("Usage: ff1 -k <keyfile> -t <tweakfile> [-decrypt] < infile > outfile\n       Infile = one long hex line, created with my base16 enc. or xxd.")
		os.Exit(1)
	}

	keyHex, err := readLines(*keyfile)
	if err != nil {
		fmt.Printf("Error reading key file: %v\n", err)
		os.Exit(1)
	}

	tweakHex, err := readLines(*tweakfile)
	if err != nil {
		fmt.Printf("Error reading tweak file: %v\n", err)
		os.Exit(1)
	}

	key, err := hex.DecodeString(string(keyHex))
	if err != nil {
		fmt.Printf("Error decoding key hex: %v\n", err)
		os.Exit(1)
	}

	tweak, err := hex.DecodeString(string(tweakHex))
	if err != nil {
		fmt.Printf("Error decoding tweak hex: %v\n", err)
		os.Exit(1)
	}

	FF1, err := ff1.NewCipher(16, 8, key, tweak) // Changed radix to 16
	if err != nil {
		fmt.Printf("Error creating FF1 cipher: %v\n", err)
		os.Exit(1)
	}

	inputData, err := readStdin()
	if err != nil {
		fmt.Printf("Error reading input data: %v\n", err)
		os.Exit(1)
	}

	var outputData string

	if *decryptFlag {
		// Decrypt if -decrypt flag is specified
		plaintext, err := FF1.Decrypt(inputData)
		if err != nil {
			fmt.Printf("Error decrypting data: %v\n", err)
			os.Exit(1)
	}
	outputData = plaintext
    } else {
        // Encrypt by default
        ciphertext, err := FF1.Encrypt(inputData)
        if err != nil {
            fmt.Printf("Error encrypting data: %v\n", err)
            os.Exit(1)
        }
        outputData = ciphertext
    }

    // Write the output data to the output file
    _, err = os.Stdout.WriteString(outputData + "\n")
    if err != nil {
        fmt.Printf("Error writing output data: %v\n", err)
        os.Exit(1)
    }
}

func readLines(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
	    return "", fmt.Errorf("failed to open file: %s", path)
    }
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var lines []string
	for scanner.Scan() {
	    lines = append(lines, scanner.Text())
    }

	return strings.Join(lines, "\n"), scanner.Err()
}

func readStdin() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)

	var lines []string
	for scanner.Scan() {
	    lines = append(lines, scanner.Text())
    }

	return strings.Join(lines, "\n"), scanner.Err()
}

