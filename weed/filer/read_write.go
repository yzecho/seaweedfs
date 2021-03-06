package filer

import (
	"bytes"
	"fmt"
	"github.com/chrislusf/seaweedfs/weed/pb/filer_pb"
	"github.com/chrislusf/seaweedfs/weed/util"
	"github.com/chrislusf/seaweedfs/weed/wdclient"
	"math"
	"net/http"
)

func ReadEntry(masterClient *wdclient.MasterClient, filerClient filer_pb.SeaweedFilerClient, dir, name string, byteBuffer *bytes.Buffer) error {

	request := &filer_pb.LookupDirectoryEntryRequest{
		Directory: dir,
		Name:      name,
	}
	respLookupEntry, err := filer_pb.LookupEntry(filerClient, request)
	if err != nil {
		return err
	}
	if len(respLookupEntry.Entry.Content) > 0 {
		_, err = byteBuffer.Write(respLookupEntry.Entry.Content)
		return err
	}

	return StreamContent(masterClient, byteBuffer, respLookupEntry.Entry.Chunks, 0, math.MaxInt64)

}

func ReadContent(filerAddress string, dir, name string) ([]byte, error) {

	target := fmt.Sprintf("http://%s%s/%s", filerAddress, dir, name)

	data, _, err := util.Get(target)

	return data, err
}

func SaveAs(host string, port int, dir, name string, contentType string, byteBuffer *bytes.Buffer) error {

	target := fmt.Sprintf("http://%s:%d%s/%s", host, port, dir, name)

	// set the HTTP method, url, and request body
	req, err := http.NewRequest(http.MethodPut, target, byteBuffer)
	if err != nil {
		return err
	}

	// set the request header Content-Type for json
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	util.CloseResponse(resp)

	return nil

}
