// SPDX-License-Identifier: AGPL-3.0-only
package parquet

import (
	"os"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
	"github.com/apache/arrow/go/v14/parquet"
	"github.com/apache/arrow/go/v14/parquet/file"
	"github.com/apache/arrow/go/v14/parquet/pqarrow"
	"github.com/rpcv2-historical/internal/domain"
)

func WriteBlocks(fname string, blocks []domain.Block) error {
	pool := memory.NewGoAllocator()
	schema := arrow.NewSchema(
		[]arrow.Field{
			{Name: "slot", Type: arrow.PrimitiveTypes.Uint64},
			{Name: "blockhash", Type: arrow.BinaryTypes.String},
			{Name: "parent_slot", Type: arrow.PrimitiveTypes.Uint64},
			{Name: "block_time", Type: arrow.PrimitiveTypes.Int64},
			{Name: "tx_count", Type: arrow.PrimitiveTypes.Uint32},
		}, nil)
	bld := array.NewRecordBuilder(pool, schema)
	defer bld.Release()
	for _, b := range blocks {
		bld.Field(0).(*array.Uint64Builder).Append(uint64(b.Slot))
		bld.Field(1).(*array.StringBuilder).Append(b.Hash)
		bld.Field(2).(*array.Uint64Builder).Append(uint64(b.ParentSlot))
		bld.Field(3).(*array.Int64Builder).Append(b.BlockTime)
		bld.Field(4).(*array.Uint32Builder).Append(b.TxCount)
	}
	rec := bld.NewRecord()
	defer rec.Release()
	f, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	w, err := pqarrow.NewFileWriter(schema, f, parquet.NewWriterProperties(), pqarrow.NewArrowWriterProperties())
	if err != nil {
		return err
	}
	return w.Write(rec)
}
