// SPDX-License-Identifier: AGPL-3.0-only
package parquet

import (
	"os"

	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/apache/arrow/go/v14/parquet/file"
	"github.com/apache/arrow/go/v14/parquet/pqarrow"
	"github.com/rpcv2-historical/internal/domain"
)

func ReadBlocks(fname string) ([]domain.Block, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	rdr, err := file.OpenParquetFile(f, true)
	if err != nil {
		return nil, err
	}
	defer rdr.Close()
	mgr := pqarrow.NewFileReader(rdr, pqarrow.ArrowReadProperties{}, memory.NewGoAllocator())
	rec, err := mgr.Read()
	if err != nil {
		return nil, err
	}
	defer rec.Release()
	var out []domain.Block
	slots := rec.Column(0).(*array.Uint64)
	hashes := rec.Column(1).(*array.String)
	parents := rec.Column(2).(*array.Uint64)
	times := rec.Column(3).(*array.Int64)
	counts := rec.Column(4).(*array.Uint32)
	for i := 0; i < int(rec.NumRows()); i++ {
		out = append(out, domain.Block{
			Slot:       domain.Slot(slots.Value(i)),
			Hash:       hashes.ValueStr(i),
			ParentSlot: domain.Slot(parents.Value(i)),
			BlockTime:  times.Value(i),
			TxCount:    counts.Value(i),
		})
	}
	return out, nil
}
