# Environment variables and flags

Set flags using the standard library `flag` package, and complement them with 
environment variables.

- All Go's standard `flag` functions are used as normal and parse errors are
  returned as expected.

- Flags, if set, take precedence over the environment.

- Environment variables can have custom prefixes and suffixes.

- Usage output can be displayed with the standard `-h` flag.

- Usage output can be annotated with the configured environment variables.

```go
fs := flag.NewFlagSet("test", flag.ContinueOnError)
fs.String("db", "", "Database thing")
_ = envflag.ParseFlagSet(fs, envflag.Prefix("APP_"), envflag.UsageSuffixer())
```

produces the following usage

```sh
Usage of test:
  -db string
        Database thing [APP_DB]
```
