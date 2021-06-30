package blocksign

import (
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/key"
	"github.com/spikeekips/mitum/base/operation"
	"github.com/spikeekips/mitum/util"
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
	rkey, err := currency.NewKey(rpk.Publickey(), 50)
	t.NoError(err)
	//
	documentPubkeys, _ := currency.NewKeys([]currency.Key{skey, rkey}, 100)

	documentPrvkeys := []key.Privatekey{spk, rpk}

	// threshold 50, key weight 50
	keys, _ := currency.NewKeys([]currency.Key{skey}, 50)
	senderAddr, _ := currency.NewAddressFromKeys(keys)

	token := util.UUID().Bytes()

	sc := SignCode("ABCD")

	item := NewCreateDocumentsItemSingleFile(documentPubkeys, sc, senderAddr, cid)
	fact := NewCreateDocumentsFact(token, senderAddr, []CreateDocumentsItem{item})

	var fs []operation.FactSign

	for _, pk := range documentPrvkeys {
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
	t.Equal(sc, ufact.Items()[0].SignCode())
	t.Equal(senderAddr, ufact.Items()[0].Owner())
}

func (t *testCreateDocuments) TestDuplicatedKeys() {
	cid := currency.CurrencyID("SHOWME")
	var items []CreateDocumentsItem

	pk := key.MustNewBTCPrivatekey()
	skey, err := currency.NewKey(pk.Publickey(), 100)
	t.NoError(err)

	skeys, _ := currency.NewKeys([]currency.Key{skey}, 100)
	sender, _ := currency.NewAddressFromKeys(skeys)
	{
		pk := key.MustNewBTCPrivatekey()
		key, err := currency.NewKey(pk.Publickey(), 100)
		t.NoError(err)
		keys, err := currency.NewKeys([]currency.Key{key}, 100)
		t.NoError(err)

		sc := SignCode("ABCD")

		items = append(items, NewCreateDocumentsItemSingleFile(keys, sc, sender, cid))
		items = append(items, NewCreateDocumentsItemSingleFile(keys, sc, sender, cid))
	}

	token := util.UUID().Bytes()

	fact := NewCreateDocumentsFact(token, sender, items)

	sig, err := operation.NewFactSignature(pk, fact, nil)
	t.NoError(err)
	fs := []operation.FactSign{operation.NewBaseFactSign(pk.Publickey(), sig)}

	op, err := NewCreateDocuments(fact, fs, "")
	t.NoError(err)

	err = op.IsValid(nil)
	t.Contains(err.Error(), "duplicated acocunt Keys found")
}

func (t *testCreateDocuments) TestSameWithSender() {
	prvk := key.MustNewBTCPrivatekey()
	key, err := currency.NewKey(prvk.Publickey(), 100)
	t.NoError(err)
	keys, err := currency.NewKeys([]currency.Key{key}, 100)
	t.NoError(err)
	sender, _ := currency.NewAddressFromKeys(keys)
	owner := sender

	sc := SignCode("ABCD")
	cid := currency.CurrencyID("SHOWME")
	items := []CreateDocumentsItem{NewCreateDocumentsItemSingleFile(keys, sc, owner, cid)}

	token := util.UUID().Bytes()

	fact := NewCreateDocumentsFact(token, sender, items)

	sig, err := operation.NewFactSignature(prvk, fact, nil)
	t.NoError(err)
	fs := []operation.FactSign{operation.NewBaseFactSign(prvk.Publickey(), sig)}

	op, err := NewCreateDocuments(fact, fs, "")
	t.NoError(err)

	err = op.IsValid(nil)
	t.Contains(err.Error(), "target document address is same with sender")
}
