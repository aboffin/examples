package main

import (
	"code.google.com/p/biogo/index/kmerindex"
	"code.google.com/p/biogo/io/seqio/fasta"
	"flag"
	"fmt"
	"os"
)

func main() {
	var (
		in *fasta.Reader
		e  error
	)

	inName := flag.String("in", "", "Filename for input. Defaults to stdin.")
	k := flag.Int("k", 6, "kmer size.")
	chunk := flag.Int("chunk", 1000, "Chunk width.")
	help := flag.Bool("help", false, "Print this usage message.")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	if *inName == "" {
		in = fasta.NewReader(os.Stdin)
	} else if in, e = fasta.NewReaderName(*inName); e != nil {
		fmt.Fprintf(os.Stderr, "Error: %v.", e)
		os.Exit(0)
	}
	defer in.Close()

	for {
		if sequence, err := in.Read(); err != nil {
			os.Exit(0)
		} else {
			index, err := kmerindex.New(*k, sequence)
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			if baseLine, ok := index.NormalisedKmerFrequencies(); ok {
				for i := 0; (i+1)**chunk < sequence.Len(); i++ {
					sub, _ := sequence.Trunc(i**chunk+1, (i+1)**chunk)
					index, err = kmerindex.New(*k, sub)
					if err != nil {
						fmt.Println(err)
						os.Exit(0)
					}
					if chunkFreqs, ok := index.NormalisedKmerFrequencies(); ok {
						fmt.Printf("%s\t%d\t%f\n", sequence.ID, i**chunk, kmerindex.Distance(baseLine, chunkFreqs))
					}
				}
			}
		}
	}
}
