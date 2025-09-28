// SPDX-License-Identifier: AGPL-3.0-only
use {
    clickhouse::{Client, Row},
    solana_geyser_plugin_interface::geyser_plugin_interface::*,
    solana_sdk::{signature::Signature, pubkey::Pubkey},
    serde::{Deserialize, Serialize},
    std::{ffi::CStr, path::Path, sync::Arc, time::SystemTime},
    tokio::runtime::Runtime,
};

#[derive(Serialize, Deserialize, Row)]
struct SigRow {
    address:  String,
    signature: String,
    slot:     u64,
    memo:     Vec<u8>,
    error:    Option<String>,
    block_time: i64,
}

pub struct RpcV2Plugin {
    runtime: Arc<Runtime>,
    client:  Arc<Client>,
}

impl RpcV2Plugin {
    fn name(&self) -> &'static str {
        "RpcV2Plugin"
    }
}

impl GeyserPlugin for RpcV2Plugin {
    fn name(&self) -> &'static str {
        Self::name()
    }

    fn on_load(&mut self, cfg_file: &CStr, _recv: Receiver<ReplicaAccountInfoVersions>) -> Result<()> {
        let path = Path::new(cfg_file.to_str().unwrap());
        let cfg: serde_json::Value = serde_json::from_reader(std::fs::File::open(path)?)?;
        let db_url = cfg["clickhouse_url"].as_str().unwrap_or("tcp://localhost:9000/rpcv2_hist");
        let rt = Arc::new(Runtime::new().unwrap());
        let client = Arc::new(Client::new(db_url));
        self.runtime = rt;
        self.client  = client;
        Ok(())
    }

    fn notify_transaction(
        &self,
        tx: ReplicaTransactionInfoVersions,
        slot: u64,
        _is_vote: bool,
    ) -> Result<()> {
        let ReplicaTransactionInfoVersions::V0_0_3(tx_info) = tx;
        let sig = tx_info.transaction.signature().to_string();
        let addr = tx_info.transaction.message().account_keys.iter()
                    .find(|k| !k.is_on_curve()) // simplest: first off-curve = fee-payer
                    .unwrap_or(&Pubkey::default()).to_string();
        let row = SigRow {
            address: addr,
            signature: sig,
            slot,
            memo: vec![], // todo: extract memo
            error: tx_info.transaction_status.clone().err.map(|e| format!("{:?}", e)),
            block_time: SystemTime::now().duration_since(SystemTime::UNIX_EPOCH).unwrap().as_secs() as i64,
        };
        let client = self.client.clone();
        self.runtime.spawn(async move {
            if let Err(e) = client.insert("sigs_for_address", vec![row]).await {
                eprintln!("insert sig: {}", e);
            }
        });
        Ok(())
    }
}

#[no_mangle]
pub unsafe extern "C" fn _create_plugin() -> *mut dyn GeyserPlugin {
    let plugin = RpcV2Plugin{
        runtime: Arc::new(Runtime::new().unwrap()),
        client:  Arc::new(Client::default()),
    };
    Box::into_raw(Box::new(plugin))
}
