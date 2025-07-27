package transform

import (
	"bufio"
	"io"
	"log"
	"strings"

	"github.com/jstemmer/go-junit-report/v2/junit"
	"github.com/mpontillo/tap13"
)

func must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func readInput(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var s []string
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}

	return s, scanner.Err()
}

// ok 4 - foobar   -> foobar
func cleanseDescription(d string) string {
	if d == "" {
		return ""
	}

	desc := strings.TrimLeft(d, " \n\t\r")
	if len(desc) > 0 && string(desc[0]) == "-" {
		return strings.TrimLeft(desc[1:], " \n\t\r")
	} else {
		return desc
	}
}

func buildTestSuites(tap *tap13.Results) junit.Testsuites {
	var suites junit.Testsuites

	ts := junit.Testsuite{Name: "TAP tests"}

	for _, t := range tap.Tests {
		var tc junit.Testcase
		if desc := cleanseDescription(t.Description); desc != "" {
			tc = junit.Testcase{Name: desc}
		} else {
			tc = junit.Testcase{Name: t.DirectiveText}
		}
		if t.Passed {
			tc.Status = "SUCCESSFUL"
		} else {
			tcres := junit.Result{}

			if len(t.Diagnostics) > 0 {
				tcres.Data = t.Diagnostics[0]
			}

			if t.Failed {
				tcres.Message = "Test failed."
				tc.Failure = &tcres
				tc.Status = "FAILED"
			}
			if t.Skipped {
				tcres.Message = "Test was skipped."
				tc.Skipped = &tcres
			}

		}

		ts.AddTestcase(tc)
	}

	if tap.BailOut {
		var tc = junit.Testcase{Name: "Bailed out."}
		tc.Status = "ABORTED"

		var tcres = junit.Result{}
		if tap.BailOutReason != "" {
			tcres.Message = tap.BailOutReason
		} else {
			tcres.Message = "Test aborted."
		}
		tc.Error = &tcres
		ts.AddTestcase(tc)
	}

	suites.AddSuite(ts)

	return suites
}

func Run(in io.Reader, out io.Writer) {
	inLines, err := readInput(in)
	must(err)

	parsed := tap13.Parse(inLines)
	if !parsed.FoundTapData {
		log.Fatalln("Could not parse TAP (version 13) report!")
	}

	testSuites := buildTestSuites(parsed)
	err = testSuites.WriteXML(out)
	must(err)
}
