package model

type Segmentation struct {
	AddressSapID string `db:"address_sap_id"`
	AdrSegment   string `db:"adr_segment"`
	SegmentID    int    `db:"segment_id"`
}
