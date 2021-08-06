package goziptz

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"strings"
	"testing"
)

func TestZipToLocationAllZips(t *testing.T) {
	t.Skip("slow, run every time the new data is pulled")

	r, err := zlib.NewReader(bytes.NewBuffer(tz))
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		split := strings.SplitN(text, "|", 3)

		t.Run(split[0], func(t *testing.T) {
			l, err := ZipToLocation(split[0])
			if err != nil {
				t.Fatal(err)
			}
			if e, g := split[1], l.String(); e != g {
				t.Fatalf("expected %q, got %q", e, g)
			}
		})
	}
	if err := scanner.Err(); err != nil {
		t.Fatal(err)
	}
}

func TestZipToLocation(t *testing.T) {
	tcases := []struct {
		in  string
		out string
		err string
	}{
		{
			in:  "94606",
			out: "America/Los_Angeles",
		},
		{
			in:  "94606-1234",
			out: "America/Los_Angeles",
		},
		{
			in:  "946061234",
			err: ErrInvalidZip.Error(),
		},
		{
			in:  "asdfafadf",
			err: ErrInvalidZip.Error(),
		},
		{
			in:  "",
			err: ErrInvalidZip.Error(),
		},
		{
			in:  "00600",
			err: ErrLocationNotFound.Error(),
		},
	}

	for _, tcase := range tcases {
		t.Run(tcase.in, func(t *testing.T) {
			out, err := ZipToLocation(tcase.in)
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}
			if errStr != tcase.err {
				t.Fatalf("expected error %q, got error %q", tcase.err, errStr)
			}
			if tcase.err != "" {
				return
			}
			if e, g := tcase.out, out.String(); e != g {
				t.Fatalf("expected location %q, got location %q", e, g)
			}
		})
	}
}
