# Detector

A very simple proof of concept to demonstrate a generic detection engine. Supports YAML rules containing regular expressions to search source code for secrets that have been checked into a repo. All demo source code is held in `./fixtures` which simulates "real" repos containing applications.

There are two examples of vulnerable Python database connections to demonstrate this is a generic detector and not bound to a specific implementation. Examples of Python and Go code that does not match any rules has been included to demonstrate there are no false positives.

Additionally there is a vulnerable Go implementation that holds authentication credentials to AWS. This demonstrates the detector can handle rules for different programming languages.

## Running

`go run main.go`
