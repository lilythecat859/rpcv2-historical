```markdown
# RPCv2-Historical Geyser Plugin

## Purpose
Streams every Solana transaction signature 
(and its fee-payer address) into ClickHouse 
in real-time so the RPCv2-Historical service 
can answer `getSignaturesForAddress` queries 
in **<10 ms** instead of **>100 ms**.

## Build (Linux / macOS / WSL)
```bash
cargo build --release
```
Output:
- Linux: `target/release/librpcv2_geyser.so`
- macOS: `target/release/librpcv2_geyser.dylib`

## Configure validator
Create `geyser-config.json`:
```json
{
  "libpath": "/absolute/path/to/target/release/librpcv2_geyser.so",
  "clickhouse_url": "tcp://localhost:9000/rpcv2_hist"
}
```
Start validator:
```bash
solana-validator --geyser-plugin-config geyser-config.json
```

## What it inserts
Table: `sigs_for_address`  
Columns: `address`, `signature`, `slot`, `memo`, `error`, `block_time`  
Inserts are **batched + async**, so validator performance is unaffected.

## Need help?
Open an issue in this repo.
```
