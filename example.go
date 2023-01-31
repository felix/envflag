package envflag

import (
	"flag"
	"os"
)

func ExampleFlag() {
	flag.Bool("verbose", false, "Be verbose")
	Parse()
}

func ExampleFlagSet() {
	fs := flag.NewFlagSet("example", flag.ExitOnError)
	fs.Bool("verbose", false, "Be verbose")
	ParseFlagSet(fs, os.Args[1:])
}

func ExamplePrefix() {
	flag.String("db", "", "Database thing")
	Parse(Prefix("APP_"))
}
