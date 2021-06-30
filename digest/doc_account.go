package digest

import (
	"github.com/soonkuk/mitum-data/blocksign"
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/state"
	mongodbstorage "github.com/spikeekips/mitum/storage/mongodb"
	"github.com/spikeekips/mitum/util/encoder"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"golang.org/x/xerrors"
)

type AccountDoc struct {
	mongodbstorage.BaseDoc
	address string
	height  base.Height
}

func NewAccountDoc(rs AccountValue, enc encoder.Encoder) (AccountDoc, error) {
	b, err := mongodbstorage.NewBaseDoc(nil, rs, enc)
	if err != nil {
		return AccountDoc{}, err
	}

	return AccountDoc{
		BaseDoc: b,
		address: currency.StateAddressKeyPrefix(rs.ac.Address()),
		height:  rs.height,
	}, nil
}

func (doc AccountDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	m["address"] = doc.address
	m["height"] = doc.height

	return bsonenc.Marshal(m)
}

func NewDocumentDoc(rs DocumentValue, enc encoder.Encoder) (AccountDoc, error) {
	b, err := mongodbstorage.NewBaseDoc(nil, rs, enc)
	if err != nil {
		return AccountDoc{}, err
	}

	return AccountDoc{
		BaseDoc: b,
		address: currency.StateAddressKeyPrefix(rs.ac.Address()),
		height:  rs.height,
	}, nil
}

type BalanceDoc struct {
	mongodbstorage.BaseDoc
	st state.State
	am currency.Amount
}

// NewBalanceDoc gets the State of Amount
func NewBalanceDoc(st state.State, enc encoder.Encoder) (BalanceDoc, error) {
	am, err := currency.StateBalanceValue(st)
	if err != nil {
		return BalanceDoc{}, xerrors.Errorf("BalanceDoc needs Amount state: %w", err)
	}

	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return BalanceDoc{}, err
	}

	return BalanceDoc{
		BaseDoc: b,
		st:      st,
		am:      am,
	}, nil
}

func (doc BalanceDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}
	address := doc.st.Key()[:len(doc.st.Key())-len(currency.StateKeyBalanceSuffix)-len(doc.am.Currency())-1]
	m["address"] = address
	m["currency"] = doc.am.Currency().String()
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}

type FileDataDoc struct {
	mongodbstorage.BaseDoc
	st state.State
	sc blocksign.SignCode
	ow base.Address
}

// NewFileDataDoc gets the State of FileData
func NewFileDataDoc(st state.State, enc encoder.Encoder) (FileDataDoc, error) {

	var fd blocksign.FileData
	if i, err := blocksign.StateFileDataValue(st); err != nil {
		return FileDataDoc{}, xerrors.Errorf("FileDataDoc needs FileData state: %w", err)
	} else {
		fd = i
	}

	b, err := mongodbstorage.NewBaseDoc(nil, st, enc)
	if err != nil {
		return FileDataDoc{}, err
	}
	return FileDataDoc{
		BaseDoc: b,
		st:      st,
		sc:      fd.SignCode(),
		ow:      fd.Owner(),
	}, nil
}

func (doc FileDataDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}
	address := doc.st.Key()[:len(doc.st.Key())-len(blocksign.StateKeyFileDataSuffix)]
	m["address"] = address
	m["height"] = doc.st.Height()

	return bsonenc.Marshal(m)
}
