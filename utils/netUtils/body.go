package netUtils

import (
	"bytes"
	"io"
	"net/http"
)

// CopyRespBody 无损取Body
func CopyRespBody(resp *http.Response) []byte {
	//复制一份body
	if resp != nil && resp.Body != nil {
		bodyBytes, _ := io.ReadAll(resp.Body)
		//返还
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return bodyBytes
	}
	return nil
}

// CopyReqBody 无损取request
func CopyReqBody(req *http.Request) []byte {
	if req.Body != nil {
		data, err := io.ReadAll(req.Body)
		if err != nil {
			return nil
		}
		// bind之前把body写回去
		req.Body = io.NopCloser(bytes.NewBuffer(data))
		return data
	}
	return nil
}
