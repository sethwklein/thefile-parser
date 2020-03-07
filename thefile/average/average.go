// Command average displays information about average line lengths to enable
// tuning of the constant in file.tokenize.
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sethwklein.net/thefile/storage"
)

func commented(buf []byte, count int) error {
	var buckets []float64
	length := 0
	count = 0
	scanner := bufio.NewScanner(bytes.NewBuffer(buf))
	for scanner.Scan() {
		count++
		length += len(scanner.Text())
		if count >= 100 {
			buckets = append(buckets, float64(length)/float64(count))
			count, length = 0, 0
		}
	}
	fmt.Println("Distribution:")
	/*
		for _, average := range buckets {
			fmt.Printf("%.2f %s\n", average, strings.Repeat("#", int(average)))
		}
	*/
	for i := 0; i < 80 && i < len(buckets)-5; i++ {
		average := 0.0
		for _, l := range buckets[i : i+5] {
			average += l
		}
		average /= 5
		fmt.Printf("%4d %.2f %s\n", i*100, average, strings.Repeat("#", int(average)))
	}

	return nil
}

func mainError() (err error) {
	buf, err := storage.Load()
	if err != nil {
		return err
	}

	whole := 0.0
	{
		var count int

		count = bytes.Count(buf, []byte{'\n'})
		if buf[len(buf)-1] != '\n' {
			count++
		}
		whole = float64(len(buf)) / float64(count)
		// fmt.Printf("Average: %.3f\n", whole)
	}

	{
		count := 0.0
		for i, c := range buf {
			if c != '\n' {
				continue
			}
			count++
			fmt.Printf("%.0f %d %.3f\n", count, i, whole-float64(i)/count)
		}
	}

	return nil
}

func mainCode() int {
	err := mainError()
	if err == nil {
		return 0
	}
	fmt.Fprintf(os.Stderr, "%v: Error: %v\n", filepath.Base(os.Args[0]), err)
	return 1
}

func main() {
	os.Exit(mainCode())
}
