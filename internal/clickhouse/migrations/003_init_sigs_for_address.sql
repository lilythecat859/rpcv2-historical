CREATE TABLE IF NOT EXISTS sigs_for_address
(
    address     String,
    signature   FixedString(88),
    slot        UInt64,
    memo        String,
    error       Nullable(String),
    block_time  Int64,
    updated_on  DateTime DEFAULT now(),
    INDEX addr_idx address TYPE bloom_filter GRANULARITY 1
) ENGINE = MergeTree()
ORDER BY (address, slot DESC);
