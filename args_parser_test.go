package main

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const (
	programName = "bombardier"
)

func TestInvalidArgsParsing(t *testing.T) {
	expectations := []struct {
		in  []string
		out string
	}{
		{
			[]string{programName},
			"required argument 'url' not provided",
		},
		{
			[]string{programName, "http://google.com", "http://yahoo.com"},
			"unexpected http://yahoo.com",
		},
	}
	for _, e := range expectations {
		p := newKingpinParser()
		if _, err := p.parse(e.in); err == nil ||
			err.Error() != e.out {
			t.Error(err, e.out)
		}
	}
}

func TestUnspecifiedArgParsing(t *testing.T) {
	p := newKingpinParser()
	args := []string{programName, "--someunspecifiedflag"}
	_, err := p.parse(args)
	if err == nil {
		t.Fail()
	}
}

func TestArgsParsing(t *testing.T) {
	ten := uint64(10)
	expectations := []struct {
		in  [][]string
		out config
	}{
		{
			[][]string{{programName, "https://somehost.somedomain"}},
			config{
				numConns:      defaultNumberOfConns,
				timeout:       defaultTimeout,
				headers:       new(headersList),
				method:        "GET",
				url:           "https://somehost.somedomain",
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"-c", "10",
					"-n", strconv.FormatUint(defaultNumberOfReqs, decBase),
					"-t", "10s",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-c10",
					"-n" + strconv.FormatUint(defaultNumberOfReqs, decBase),
					"-t10s",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--connections", "10",
					"--requests", strconv.FormatUint(defaultNumberOfReqs, decBase),
					"--timeout", "10s",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--connections=10",
					"--requests=" + strconv.FormatUint(defaultNumberOfReqs, decBase),
					"--timeout=10s",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:      10,
				timeout:       10 * time.Second,
				headers:       new(headersList),
				method:        "GET",
				numReqs:       &defaultNumberOfReqs,
				url:           "https://somehost.somedomain",
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"--latencies",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-l",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:       defaultNumberOfConns,
				timeout:        defaultTimeout,
				headers:        new(headersList),
				printLatencies: true,
				method:         "GET",
				url:            "https://somehost.somedomain",
				printIntro:     true,
				printProgress:  true,
				printResult:    true,
			},
		},
		{
			[][]string{
				{
					programName,
					"--insecure",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-k",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:      defaultNumberOfConns,
				timeout:       defaultTimeout,
				headers:       new(headersList),
				insecure:      true,
				method:        "GET",
				url:           "https://somehost.somedomain",
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"--key", "testclient.key",
					"--cert", "testclient.cert",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--key=testclient.key",
					"--cert=testclient.cert",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:      defaultNumberOfConns,
				timeout:       defaultTimeout,
				headers:       new(headersList),
				method:        "GET",
				keyPath:       "testclient.key",
				certPath:      "testclient.cert",
				url:           "https://somehost.somedomain",
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"--method", "POST",
					"--body", "reqbody",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--method=POST",
					"--body=reqbody",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-m", "POST",
					"-b", "reqbody",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-mPOST",
					"-breqbody",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:      defaultNumberOfConns,
				timeout:       defaultTimeout,
				headers:       new(headersList),
				method:        "POST",
				body:          "reqbody",
				url:           "https://somehost.somedomain",
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"--header", "One: Value one",
					"--header", "Two: Value two",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-H", "One: Value one",
					"-H", "Two: Value two",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--header=One: Value one",
					"--header=Two: Value two",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns: defaultNumberOfConns,
				timeout:  defaultTimeout,
				headers: &headersList{
					{"One", "Value one"},
					{"Two", "Value two"},
				},
				method:        "GET",
				url:           "https://somehost.somedomain",
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"--rate", "10",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-r", "10",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--rate=10",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-r10",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:      defaultNumberOfConns,
				timeout:       defaultTimeout,
				headers:       new(headersList),
				method:        "GET",
				url:           "https://somehost.somedomain",
				rate:          &ten,
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"--fasthttp",
					"https://somehost.somedomain",
				},
				{
					programName,
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:      defaultNumberOfConns,
				timeout:       defaultTimeout,
				headers:       new(headersList),
				method:        "GET",
				url:           "https://somehost.somedomain",
				clientType:    fhttp,
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"--http1",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:      defaultNumberOfConns,
				timeout:       defaultTimeout,
				headers:       new(headersList),
				method:        "GET",
				url:           "https://somehost.somedomain",
				clientType:    nhttp1,
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"--http2",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:      defaultNumberOfConns,
				timeout:       defaultTimeout,
				headers:       new(headersList),
				method:        "GET",
				url:           "https://somehost.somedomain",
				clientType:    nhttp2,
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"--body-file=testbody.txt",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--body-file", "testbody.txt",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-f", "testbody.txt",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:      defaultNumberOfConns,
				timeout:       defaultTimeout,
				headers:       new(headersList),
				method:        "GET",
				bodyFilePath:  "testbody.txt",
				url:           "https://somehost.somedomain",
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"--stream",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-s",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:      defaultNumberOfConns,
				timeout:       defaultTimeout,
				headers:       new(headersList),
				method:        "GET",
				stream:        true,
				url:           "https://somehost.somedomain",
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:      defaultNumberOfConns,
				timeout:       defaultTimeout,
				headers:       new(headersList),
				method:        "GET",
				url:           "https://somehost.somedomain",
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"--print=r,i,p",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--print", "r,i,p",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-p", "r,i,p",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--print=result,i,p",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--print", "r,intro,p",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-p", "r,i,progress",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:      defaultNumberOfConns,
				timeout:       defaultTimeout,
				headers:       new(headersList),
				method:        "GET",
				url:           "https://somehost.somedomain",
				printIntro:    true,
				printProgress: true,
				printResult:   true,
			},
		},
		{
			[][]string{
				{
					programName,
					"--print=i,r",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--print", "i,r",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-p", "i,r",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--print=intro,r",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--print", "i,result",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-p", "intro,r",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:      defaultNumberOfConns,
				timeout:       defaultTimeout,
				headers:       new(headersList),
				method:        "GET",
				url:           "https://somehost.somedomain",
				printIntro:    true,
				printProgress: false,
				printResult:   true,
			},
		},
	}
	for _, e := range expectations {
		for _, args := range e.in {
			p := newKingpinParser()
			cfg, err := p.parse(args)
			if err != nil {
				t.Error(err)
				continue
			}
			if !reflect.DeepEqual(cfg, e.out) {
				t.Logf("Expected: %#v", e.out)
				t.Logf("Got: %#v", cfg)
				t.Fail()
			}
		}
	}
}

func TestParsePrintSpec(t *testing.T) {
	exps := []struct {
		spec    string
		results [3]bool
		err     error
	}{
		{
			"",
			[3]bool{},
			errEmptyPrintSpec,
		},
		{
			"a,b,c",
			[3]bool{},
			fmt.Errorf("%q is not a valid part of print spec", "a"),
		},
		{
			"i,p,r,i",
			[3]bool{},
			fmt.Errorf(
				"Spec %q has too many parts, at most 3 are allowed", "i,p,r,i",
			),
		},
		{
			"i",
			[3]bool{true, false, false},
			nil,
		},
		{
			"p",
			[3]bool{false, true, false},
			nil,
		},
		{
			"r",
			[3]bool{false, false, true},
			nil,
		},
		{
			"i,p,r",
			[3]bool{true, true, true},
			nil,
		},
	}
	for _, e := range exps {
		var (
			act = [3]bool{}
			err error
		)
		act[0], act[1], act[2], err = parsePrintSpec(e.spec)
		if !reflect.DeepEqual(err, e.err) {
			t.Errorf("For %q, expected err = %q, but got %q",
				e.spec, e.err, err,
			)
			continue
		}
		if !reflect.DeepEqual(e.results, act) {
			t.Errorf("For %q, expected result = %+v, but got %+v",
				e.spec, e.results, act,
			)
		}
	}
}

func TestArgsParsingWithEmptyPrintSpec(t *testing.T) {
	p := newKingpinParser()
	c, err := p.parse(
		[]string{programName, "--print=", "somehost.somedomain"})
	if err == nil {
		t.Fail()
	}
	if c != emptyConf {
		t.Fail()
	}
}
