package envflag_test

import (
	"flag"
	"os"

	"userspace.com.au/felix/envflag"
)

func ExampleFlag() {
	flag.Bool("verbose", false, "Be verbose")
	envflag.Parse()
}

func ExampleFlagSet() {
	fs := flag.NewFlagSet("example", flag.ExitOnError)
	fs.Bool("verbose", false, "Be verbose")
	_ = envflag.ParseFlagSet(fs, os.Args[1:], envflag.UsageSuffixer())
}

func ExamplePrefix() {
	flag.String("db", "", "Database thing")
	_ = envflag.Parse(envflag.Prefix("APP_"))
}

func ExampleUsageSuffixer() {
	flag.String("db", "", "Database thing")
	_ = envflag.Parse(envflag.Prefix("APP_"))
}
