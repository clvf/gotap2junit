package transform

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
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

func getTestName(t *tap13.Test) string {
	var b strings.Builder
	desc := t.Description

	if t.TestNumber <= 0 {
		desc = cleanseDescription(t.Description)
	} else {
		b.WriteString(fmt.Sprintf("%s ", strconv.Itoa(t.TestNumber)))
	}

	if desc != "" {
		b.WriteString(desc)
	} else {
		b.WriteString(t.DirectiveText)
	}

	return b.String()
}

func getTestMsg(t *tap13.Test) string {
	if t.DirectiveText != "" {
		return t.DirectiveText
	}

	if t.Failed {
		return "Test failed."
	}

	if t.Skipped {
		return "Test was skipped."
	}

	if t.Todo {
		return "Test is TODO."
	}

	return ""
}

func getTestData(t *tap13.Test) string {
	var b strings.Builder

	if len(t.Diagnostics) > 0 {
		for _, diag := range t.Diagnostics {
			b.WriteString(diag)
			b.WriteString("\n")
		}

		return b.String()
	}

	if len(t.YamlBytes) > 0 {
		b.WriteString("\n---\n")
		b.Write(t.YamlBytes)
		b.WriteString("...\n")

		return b.String()
	}

	return b.String()
}

func buildTestSuites(tap *tap13.Results, tapRaw []string) junit.Testsuites {
	var suites junit.Testsuites

	ts := junit.Testsuite{Name: "TAP tests"}
	ts.SystemOut = &junit.Output{Data: strings.Join(tapRaw, "\n") + "\n"}

	for _, t := range tap.Tests {
		var tc = junit.Testcase{}

		tc.Name = getTestName(&t)

		if t.Passed {
			tc.Status = "SUCCESSFUL"
			ts.AddTestcase(tc)
			continue
		}

		tcres := junit.Result{}
		tcres.Message = getTestMsg(&t)
		tcres.Data = getTestData(&t)

		if t.Failed {
			tc.Failure = &tcres
			tc.Status = "FAILED"
		}

		if t.Skipped {
			tc.Skipped = &tcres
		}

		if t.Todo {
			tc.Skipped = &tcres
			tc.Status = "SUCCESSFUL"
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

	testSuites := buildTestSuites(parsed, inLines)
	err = testSuites.WriteXML(out)
	must(err)
}
