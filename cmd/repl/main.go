// Command repl implements a Read-Execute-Print Loop
// interface for lambda expressions.
//
//   - Input a lambda expression.
//   - Then input a newline to perform a single beta-reduction.
//   - From here you may either:
//     input another newline to perform another beta-reduction on the output, or
//     input a new lambda expression
package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/nightmarlin/lambda"
	"github.com/nightmarlin/lambda/render/text"
)

const title = `} var  ::= (UTF-8 letters, numbers, '_' or '-')
  term ::= var
         | \ variable . term
         | ( term term )
         | ( term )
  generally, there may be any number of spaces between terms
  but <variable>s must have at least one space between them.

} .exit will exit
  "" (empty) will reduce last result
  $ will be replaced with last result

`

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var (
		scan    = bufio.NewScanner(os.Stdin)
		lines   = make(chan string)
		scanErr = make(chan error, 1)
	)
	scan.Split(bufio.ScanLines)
	go func() {
		defer close(lines)
		defer close(scanErr)

		for scan.Scan() {
			lines <- scan.Text()
		}
		if err := scan.Err(); err != nil && !errors.Is(err, io.EOF) {
			scanErr <- err
		}
	}()

	fmt.Print(title)
	var (
		stored    lambda.Term
		storedStr string
		store     = func(t lambda.Term) {
			stored = t
			storedStr, _ = text.RenderString(t)
		}
	)

	for {
		fmt.Print("> ")

		select {
		case line, ok := <-lines:
			if !ok {
				continue
			}

			line = strings.TrimSpace(line)
			if line == "" {
				if stored == nil {
					fmt.Println("input a lambda expression")
					continue
				}
				break // use last result as input
			}

			if line == ".exit" {
				fmt.Println("exiting...")
				return
			}

			if stored != nil {
				line = strings.ReplaceAll(line, "$", storedStr)
			}

			t, err := text.ParseString(line)
			if err != nil {
				fmt.Printf("failed to parse %q: %s\n", line, err.Error())
				continue
			}
			store(t)

		case err, ok := <-scanErr:
			if !ok {
				continue
			}
			fmt.Println("error: ", err.Error())

		case <-ctx.Done():
			fmt.Println("\nexiting...")
			return
		}

		reduced, didReduce := lambda.BetaReduce(stored)
		store(reduced)
		fmt.Print("  ", storedStr)
		if !didReduce {
			fmt.Print("\t[unchanged]")
		}
		fmt.Println()
	}
}
