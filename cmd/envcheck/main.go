// Command envcheck compares an env file with its template and exits 1 when
// required keys are missing or empty.
//
//	envcheck -template .env.example -env .env
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/PrinceAdvant/envcheck"
)

func main() {
	tplPath := flag.String("template", ".env.example", "template listing every expected key")
	envPath := flag.String("env", ".env", "env file to check")
	strictExtra := flag.Bool("strict", false, "also fail on keys missing from the template")
	flag.Parse()
	tpl, err := envcheck.ParseFile(*tplPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	env, err := envcheck.ParseFile(*envPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
	r := envcheck.Compare(tpl, env)
	for _, k := range r.Missing {
		fmt.Println("missing:", k)
	}
	for _, k := range r.Empty {
		fmt.Println("empty:  ", k)
	}
	for _, k := range r.Extra {
		fmt.Println("extra:  ", k)
	}
	if !r.OK() || (*strictExtra && len(r.Extra) > 0) {
		os.Exit(1)
	}
	fmt.Println("env OK")
}
