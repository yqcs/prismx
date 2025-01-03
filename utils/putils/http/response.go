package httputil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/docker/go-units"
)

var (
	MaxBodyRead, _ = units.FromHumanSize("4mb")
)

// DumpResponseIntoBuffer dumps a http response without allocating a new buffer
// for the response body.
func DumpResponseIntoBuffer(resp *http.Response, body bool, buff *bytes.Buffer) (err error) {
	if resp == nil {
		return fmt.Errorf("response is nil")
	}
	save := resp.Body
	savecl := resp.ContentLength

	if !body {
		// For content length of zero. Make sure the body is an empty
		// reader, instead of returning error through failureToReadBody{}.
		if resp.ContentLength == 0 {
			resp.Body = emptyBody
		} else {
			resp.Body = failureToReadBody{}
		}
	} else if resp.Body == nil {
		resp.Body = emptyBody
	} else {
		save, resp.Body, err = drainBody(resp.Body)
		if err != nil {
			return err
		}
	}
	err = resp.Write(buff)
	if err == errNoBody {
		err = nil
	}
	resp.Body = save
	resp.ContentLength = savecl
	return
}

// DrainResponseBody drains the response body and closes it.
func DrainResponseBody(resp *http.Response) {
	defer resp.Body.Close()
	// don't reuse connection and just close if body length is more than 2 * MaxBodyRead
	// to avoid DOS
	_, _ = io.CopyN(io.Discard, resp.Body, 2*MaxBodyRead)
}
