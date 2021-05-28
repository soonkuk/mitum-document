package currency

import (
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
	cid := CurrencyID("SHOWME")

	skey, err := NewKey(spk.Publickey(), 50)
	t.NoError(err)
	rkey, err := NewKey(rpk.Publickey(), 50)
	t.NoError(err)
	//
	documentPubkeys, _ := NewKeys([]Key{skey, rkey}, 100)

	documentPrvkeys := []key.Privatekey{spk, rpk}

	// threshold 50, key weight 50
	keys, _ := NewKeys([]Key{skey}, 50)
	senderAddr, _ := NewAddressFromKeys(keys)

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
	cid := CurrencyID("SHOWME")
	var items []CreateDocumentsItem

	pk := key.MustNewBTCPrivatekey()
	skey, err := NewKey(pk.Publickey(), 100)
	t.NoError(err)

	skeys, _ := NewKeys([]Key{skey}, 100)
	sender, _ := NewAddressFromKeys(skeys)
	{
		pk := key.MustNewBTCPrivatekey()
		key, err := NewKey(pk.Publickey(), 100)
		t.NoError(err)
		keys, err := NewKeys([]Key{key}, 100)
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
	key, err := NewKey(prvk.Publickey(), 100)
	t.NoError(err)
	keys, err := NewKeys([]Key{key}, 100)
	t.NoError(err)
	sender, _ := NewAddressFromKeys(keys)
	owner := sender

	sc := SignCode("ABCD")
	cid := CurrencyID("SHOWME")
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
