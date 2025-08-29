# Environment variables and flags

Set flags using Go's `flag` package, and complement them with environment 
variables. Flags, if set, take precedence over the environment.

Go's standard `flag.FlagSet` can be initialised as normal and parse errors are 
returned as expected. Environment variables can have custom prefixes, suffixes 
etc. as can the usage displayed with the standard `-h` flag.

```go
fs := flag.NewFlagSet("test", flag.ContinueOnError)
fs.String("db", "", "Database thing")
envflag.ParseFlagSet(fs, envflag.Prefix("APP_"), envflag.UsageSuffixer())
```

produces the following usage

```sh
Usage of test:
  -db string
        Database thing [APP_DB]
```
