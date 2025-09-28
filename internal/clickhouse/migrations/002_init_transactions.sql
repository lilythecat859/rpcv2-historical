CREATE TABLE IF NOT EXISTS transactions
(
    slot        UInt64,
    tx_hash     FixedString(88),
    idx         UInt32,
    meta        String,
    message     String,
    block_time  Int64,
    updated_on  DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (slot, idx);
