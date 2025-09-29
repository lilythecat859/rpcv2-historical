
```markdown
# Helm chart for RPCv2-Historical

## Install
```bash
helm install rpcv2-hist helm/rpcv2-hist \
  --set jwtKey=$(base64 -w0 < ed25519.pem) \
  --set image.repository=ghcr.io/YOUR_USERNAME/rpcv2-hist
```

## Upgrade
```bash
helm upgrade rpcv2-hist helm/rpcv2-hist --reuse-values
```

## Defaults
- 2 replicas  
- 1 vCPU / 1 GiB per pod  
- ClickHouse assumed at `tcp://clickhouse:9000/rpcv2_hist`
```
