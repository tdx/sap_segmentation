CREATE TABLE IF NOT EXISTS sap_segmentation(
    id              serial,
    address_sap_id  varchar(255) UNIQUE,
    adr_segment     varchar(16),
    segment_id      bigint
);
