package sap_segmentation

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/tdx/sap_segmentation/model"
)

type HTTPLoader struct {
	log    *log.Logger
	client *http.Client
	args   *ConnArgs
	offset int
}

type segment struct {
	AddressSapId string `json:"ADDRESS_SAP_ID"`
	AdrSegment   string `json:"ADR_SEGMENT"`
	SegmentID    int    `json:"SEGMENT_ID"`
}

type segmentation struct {
	Segmentation []segment `json:"SEGMENTATION"`
}

//
func NewHTTPLoader(
	log *log.Logger,
	args *ConnArgs) (*HTTPLoader, error) {

	client := &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout: args.ConnectTimeout,
			}).Dial,
		},
		// Timeout: args.ConnectTimeout * 3, // roundTrip timeout
	}

	log.Printf("loader: api endpoint: %s\n", args.Uri)

	return &HTTPLoader{
		log:    log,
		client: client,
		args:   args,
		offset: 1,
	}, nil
}

//
func (hi *HTTPLoader) LoadNext() ([]model.Segmentation, error) {

	if hi.args.Interval > 0 && hi.offset > 1 {
		time.Sleep(hi.args.Interval)
	}

	uri := fmt.Sprintf("%s?p_limit=%d&p_offset=%d",
		hi.args.Uri, hi.args.ImportBatchSize, hi.offset)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	// basic auth
	auth := base64.StdEncoding.EncodeToString([]byte(hi.args.AuthLoginPwd))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("User-Agent", hi.args.UserAgent)

	hi.log.Printf("loader: %s\n", uri)

	resp, err := hi.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// decode response
	var seg segmentation
	if err = json.NewDecoder(resp.Body).Decode(&seg); err != nil {
		return nil, err
	}

	// all segments loaded ?
	if len(seg.Segmentation) == 0 {
		hi.offset = 1
		return nil, nil
	}

	// update next request p_offset
	hi.offset += hi.args.ImportBatchSize
	if hi.offset == hi.args.ImportBatchSize+1 {
		hi.offset--
	}

	// convert result to model type
	msegs := make([]model.Segmentation, 0, len(seg.Segmentation))
	for i := range seg.Segmentation {
		seg := seg.Segmentation[i]
		msegs = append(msegs, model.Segmentation{
			AddressSapID: seg.AddressSapId,
			AdrSegment:   seg.AdrSegment,
			SegmentID:    seg.SegmentID,
		})
	}

	return msegs, nil
}
