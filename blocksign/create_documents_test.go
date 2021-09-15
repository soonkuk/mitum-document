package blocksign

import (
	"testing"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
	"github.com/stretchr/testify/suite"
)

type testCreateDocuments struct {
	baseTest
}

func (t *testCreateDocuments) TestNew() {
	spk := key.MustNewBTCPrivatekey()
	rpk := key.MustNewBTCPrivatekey()
	cid := currency.CurrencyID("SHOWME")

	skey, err := currency.NewKey(spk.Publickey(), 50)
	t.NoError(err)
	sgkey, err := currency.NewKey(rpk.Publickey(), 50)
	t.NoError(err)

	// threshold 50, key weight 50
	spkeys := []key.Privatekey{spk}
	keys0, _ := currency.NewKeys([]currency.Key{skey}, 50)
	senderAddr, _ := currency.NewAddressFromKeys(keys0)
	keys1, _ := currency.NewKeys([]currency.Key{sgkey}, 50)
	signerAddr, _ := currency.NewAddressFromKeys(keys1)

	token := util.UUID().Bytes()

	filehash := FileHash("ABCD")
	documentid := currency.NewBig(0)
	signcode0 := "user0"
	title := "title01"
	size := currency.NewBig(555)
	signcode1 := "user1"

	item := NewCreateDocumentsItemSingleFile(filehash, documentid, signcode0, title, size, []base.Address{signerAddr}, []string{signcode1}, cid)
	fact := NewCreateDocumentsFact(token, senderAddr, []CreateDocumentsItem{item})

	var fs []operation.FactSign

	for _, pk := range spkeys {
		sig, err := operation.NewFactSignature(pk, fact, nil)
		t.NoError(err)

		fs = append(fs, operation.NewBaseFactSign(pk.Publickey(), sig))
	}

	op, err := NewCreateDocuments(fact, fs, "")
	t.NoError(err)

	t.NoError(op.IsValid(nil))

	t.Implements((*base.Fact)(nil), op.Fact())
	t.Implements((*operation.Operation)(nil), op)

	ufact := op.Fact().(CreateDocumentsFact)
	t.Equal(filehash, ufact.Items()[0].FileHash())
	t.Equal(signerAddr, ufact.Items()[0].Signers()[0])
}

func (t *testCreateDocuments) TestDuplicatedDocumentId() {
	cid := currency.CurrencyID("SHOWME")
	var items []CreateDocumentsItem

	pk := key.MustNewBTCPrivatekey()
	skey, err := currency.NewKey(pk.Publickey(), 100)
	t.NoError(err)

	skeys, _ := currency.NewKeys([]currency.Key{skey}, 100)
	sender, _ := currency.NewAddressFromKeys(skeys)
	{
		filehash := FileHash("ABCD")
		documentid := currency.NewBig(0)
		signcode0 := "user0"
		title := "title01"
		size := currency.NewBig(555)

		items = append(items, NewCreateDocumentsItemSingleFile(filehash, documentid, signcode0, title, size, []base.Address{}, []string{}, cid))
		items = append(items, NewCreateDocumentsItemSingleFile(filehash, documentid, signcode0, title, size, []base.Address{}, []string{}, cid))
	}

	token := util.UUID().Bytes()

	fact := NewCreateDocumentsFact(token, sender, items)

	sig, err := operation.NewFactSignature(pk, fact, nil)
	t.NoError(err)
	fs := []operation.FactSign{operation.NewBaseFactSign(pk.Publickey(), sig)}

	op, err := NewCreateDocuments(fact, fs, "")
	t.NoError(err)

	err = op.IsValid(nil)
	t.Contains(err.Error(), "duplicated filehash")
}

func TestCreateDocuments(t *testing.T) {
	suite.Run(t, new(testCreateDocuments))
}
