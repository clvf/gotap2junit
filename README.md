# gotap2junit

Transform TAP version 13 reports to JUnit XML format.

***

## Usage

The program reads the TAP version 13 test report on its stdin and writes the JUnit XML
representation on its stdout.

```bash
$ ./build.sh
$ ./bin/gotap2junit < testdata/todos_should_succeed.tap13

<testsuites tests="4" skipped="2">
        <testsuite name="TAP tests" tests="4" failures="0" errors="0" id="0" skipped="2" time="">
                <testcase name="1 - Creating test program" classname="" status="SUCCESSFUL"></testcase>
                <testcase name="2 - Test program runs, no error" classname="" status="SUCCESSFUL"></testcase>
                <testcase name="3 - infinite loop" classname="" status="SUCCESSFUL">
                        <skipped message="TODO halting problem unsolved"></skipped>
                </testcase>
                <testcase name="4 - infinite loop 2" classname="" status="SUCCESSFUL">
                        <skipped message="TODO halting problem unsolved"></skipped>
                </testcase>
                <system-out><![CDATA[TAP version 13
1..4
ok 1 - Creating test program
ok 2 - Test program runs, no error
not ok 3 - infinite loop # TODO halting problem unsolved
not ok 4 - infinite loop 2 # TODO halting problem unsolved
]]></system-out>
        </testsuite>
</testsuites>
```

## Why?

CI/CD pipelines typically expect the test artifacts to be in some kind of an XML
format.

Sometimes you don't want or you are not in the position to use a test framework
(such as Python Nose, Java JUnit, Go builtin, shUnit2, Perl TAP Harness, etc.) that
would represent the test results as XML.

You can however easily produce TAP (version 13) output of your test cases anytime. In
this case a tool to convert the TAP result to JUnit so your tests can be analysed by
the CI/CD pipeline of your choice is helpful.
