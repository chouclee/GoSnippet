// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Generating random text: a Markov chain algorithm

Based on the program presented in the "Design and Implementation" chapter
of The Practice of Programming (Kernighan and Pike, Addison-Wesley 1999).
See also Computer Recreations, Scientific American 260, 122 - 125 (1989).

A Markov chain algorithm generates text by creating a statistical model of
potential textual suffixes for a given prefix. Consider this text:

	I am not a number! I am a free man!

Our Markov chain algorithm would arrange this text into this set of prefixes
and suffixes, or "chain": (This table assumes a prefix length of two words.)

	Prefix       Suffix

	"" ""        I
	"" I         am
	I am         a
	I am         not
	a free       man!
	am a         free
	am not       a
	a number!    I
	number! I    am
	not a        number!

To generate text using this table we select an initial prefix ("I am", for
example), choose one of the suffixes associated with that prefix at random
with probability determined by the input statistics ("a"),
and then create a new prefix by removing the first word from the prefix
and appending the suffix (making the new prefix is "am a"). Repeat this process
until we can't find any suffixes for the current prefix or we exceed the word
limit. (The word limit is necessary as the chain table may contain cycles.)

Our version of this program reads text from standard input, parsing it into a
Markov chain, and writes generated text to standard output.
The prefix and output lengths can be specified using the -prefix and -words
flags on the command-line.
*/
package main

import (
	"bufio"
	//"flag"
	"fmt"
	//"io"
	"math/rand"
	"os"
	"strings"
	"time"
	"strconv"
)

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}


type Chain struct {
	chain     map[string]map[string]int
	prefixLen int
}


// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type ChainForGenerate struct {
	chain     map[string][]string
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string]map[string]int), prefixLen}
}

// NewChianForGenerate returns a new ChainForGenerate with prefixes of
// prefixLen words.
func NewChianForGenerate(prefixLen int) *ChainForGenerate {
	return &ChainForGenerate{make(map[string][]string), prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
/*func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}*/

// Generate returns a string of at most n words generated from Chain.
func (c *ChainForGenerate) Generate(n int) string {
	p := make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}
	// Intn returns, as an int, a non-negative pseudo-random number in [0,n)
		next := choices[rand.Intn(len(choices))] 
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

func (c *Chain) Read(filePath string) {
	// open the file, return a Reader
	r, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error: Could not open file " + filePath)
		return 
	}
	defer r.Close()

	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for i, _ := range p {
		p[i] = "\"\""
	}
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		tf, ok := c.chain[key]  // term frequency vector of a certian prefix
		if !ok {				// if this prefix is not in the map 
			tf = make(map[string]int) 
			c.chain[key] = tf
		}
		tf[s]++
		p.Shift(s)
	}
}

func (c *Chain) WriteModel(outfilename string) {
	out, err := os.Create(outfilename) // Create outputfile.
	if err != nil {
		fmt.Println("Error! Couldn't create " + outfilename)
		return
	}
	defer out.Close()

	fmt.Fprintln(out, c.prefixLen)
	for k, v := range c.chain { // iterate over all prefix strings
		fmt.Fprint(out, k + " ")
		for _k, _v := range v { // iterate over all term frequencies
			fmt.Fprint(out, _k + " ")
			fmt.Fprint(out, _v) 
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, "\n")
	}
	out.Close()
}

func ReadModel(modelfile string) *ChainForGenerate {
	var  (
		prefixLen int
		p Prefix
		// string
	)

	r, err := os.Open(modelfile) //open model file
	if err != nil {
		fmt.Println("Error: Could not open file " + modelfile)
		return nil
	}
	defer r.Close()

	// read the prefix length
	scanner := bufio.NewScanner(r)

	if scanner.Scan() {
		var err error
		prefixLen, err = strconv.Atoi(scanner.Text())
		if err != nil {
			fmt.Println("Error: prefix length should be an integer >= 1")
			return nil
		}
	}
	c := NewChianForGenerate(prefixLen) // Initialize a new Chain.

	for scanner.Scan() {
		line := scanner.Text()
		splited := strings.Fields(line)


		p = splited[0:prefixLen]
		// replace all \"\" with empty string
		for i := 0; i < len(p); i++ {
			if p[i] == "\"\"" {
				p[i] = ""
			}
		}
		key := p.String()
		
		for i := prefixLen; i < len(splited); i += 2 {
			frequency, err := strconv.Atoi(splited[i+1])
			if err != nil {
				return nil
			}
			term := splited[i]
			for j := 0; j < frequency; j++ {
				c.chain[key] = append(c.chain[key], term)
			}
		}
		//fmt.Println(c.chain[key])
	}
	return c
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: command should be: mark COMMAND options")
		return
	}
	command := os.Args[1]
	if command == "read" {
		if len(os.Args) < 5 {
			fmt.Println("Error: read command should be: mark read N" + 
				" outfilename infile1 infile2 .... ")
			return
		}

		prefixLen, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Error: prefix length should be an integer >= 1")
			return
		}
		c := NewChain(prefixLen) // Initialize a new Chain.

		outfilename := os.Args[3]

		for i := 3; i < len(os.Args); i++ {
			c.Read(os.Args[i])
		}

		c.WriteModel(outfilename)

	} else if (command == "generate") {
		if len(os.Args) < 4 {
			fmt.Println("Error: generate command should be: mark generate" + 
				" modelfile n")
			return
		}

		c := ReadModel(os.Args[2])

		rand.Seed(time.Now().UnixNano()) // Seed the random number generator.
		text := c.Generate(10) // Generate text.
		fmt.Println(text)
	} else {
		fmt.Println("Error: command should be either read or generate")
		return
	}
	//numWords := flag.Int("words", 100, "maximum number of words to print")
	//prefixLen := flag.Int("prefix", 2, "prefix length in words")

	//flag.Parse()                     // Parse command-line flags.
	//rand.Seed(time.Now().UnixNano()) // Seed the random number generator.

	//c := NewChain(*prefixLen)     // Initialize a new Chain.
	//c.Build(os.Stdin)             // Build chains from standard input.
	//text := c.Generate(*numWords) // Generate text.
	//fmt.Println(text)             // Write text to standard output.
}