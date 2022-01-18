package blocksign

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
)

func (it *BaseCreateDocumentsItem) unpack(
	enc encoder.Encoder,
	bfh string,
	bdi currency.Big,
	bsc string,
	btl string,
	bsz currency.Big,
	bsg []base.AddressDecoder,
	bsd []string,
	scid string,
) error {

	signers := make([]base.Address, len(bsg))

	for i := range bsg {
		if a, err := bsg[i].Encode(enc); err != nil {
			return err
		} else {
			signers[i] = a
		}
	}
	it.signers = signers
	it.fileHash = FileHash(bfh)
	it.documentid = bdi
	it.signcode = bsc
	it.title = btl
	it.size = bsz
	it.signcodes = bsd
	it.cid = currency.CurrencyID(scid)

	return nil
}
