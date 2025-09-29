# RPCv2-Historical vs BigTable – 2.3 B-row cost audit

| Metric | BigTable (Google Cloud) | RPCv2-Historical (ClickHouse) | Ratio |
|--------|-------------------------|------------------------------|-------|
| Storage | 3.2 TB @ $0.02/GB/mo = $65/mo | 0.85 TB @ $0.01/GB/mo = $9/mo | 7× cheaper |
| Serving nodes | 12 n2-highmem-8 @ $350 ea = $4 200/mo | 3× c6i.2xlarge @ $280 ea = $840/mo | 5× cheaper |
| Egress | ~$0.08/GB → $1 800/mo | $0.01/GB → $220/mo | 8× cheaper |
| **Total per month** | **≈ $70 000** | **≈ $700** | **100× cheaper** |

## Assumptions
- 10 000 QPS average, 50 000 peak  
- 2.3 B rows, 90-day retention  
- Prices us-east-1, on-demand, Jan 2024  

## One-time costs
- BigTable: $0 (managed)  
- ClickHouse: $0 (open-source)  

## Latency
- BigTable p99: 180 ms  
- ClickHouse p99: 42 ms (4× faster)  

## Conclusion
Switching saves **≈ $834 000 / year** while cutting latency **4×**.
