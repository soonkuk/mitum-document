package document

import (
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util/encoder"
	"github.com/spikeekips/mitum/util/hint"
)

func (doc *Document) unpack(
	enc encoder.Encoder,
	dc []byte,
) error {

	// unpack document info
	if hinter, err := enc.Decode(dc); err != nil {
		return err
	} else if i, ok := hinter.(DocumentData); !ok {
		return errors.Errorf("not Document: %T", hinter)
	} else {
		doc.data = i
	}

	return nil
}

func (doc *BCUserData) unpack(
	enc encoder.Encoder,
	di []byte,
	us base.AddressDecoder,
	gd currency.Big, // gold
	bg currency.Big, // bankgold
	st []byte, // statistics
) error {

	// unpack document info
	if hinter, err := enc.Decode(di); err != nil {
		return err
	} else if i, ok := hinter.(DocInfo); !ok {
		return errors.Errorf("not DocInfo: %T", hinter)
	} else {
		doc.info = i
	}

	a, err := us.Encode(enc)
	if err != nil {
		return err
	}
	doc.owner = a

	doc.gold = gd
	doc.bankgold = bg

	// unpack statistics
	if hinter, err := enc.Decode(st); err != nil {
		return err
	} else if i, ok := hinter.(UserStatistics); !ok {
		return errors.Errorf("not UserStatistics: %T", hinter)
	} else {
		doc.statistics = i
	}

	return nil
}

func (doc *BCLandData) unpack(
	enc encoder.Encoder,
	di []byte,
	ow base.AddressDecoder,
	ad string, // land address
	ar string, // land area
	rt string, // renter nickname
	ac base.AddressDecoder, //renter account address
	rd string, // rentdate
	pd uint, // period day
) error {

	// unpack document info
	if hinter, err := enc.Decode(di); err != nil {
		return err
	} else if i, ok := hinter.(DocInfo); !ok {
		return errors.Errorf("not Document Info: %T", hinter)
	} else {
		doc.info = i
	}

	oa, err := ow.Encode(enc)
	if err != nil {
		return err
	}
	doc.owner = oa

	ra, err := ac.Encode(enc)
	if err != nil {
		return err
	}
	doc.account = ra

	doc.address = ad
	doc.area = ar
	doc.renter = rt
	doc.rentdate = rd
	doc.periodday = pd

	return nil
}

func (doc *BCVotingData) unpack(
	enc encoder.Encoder,
	bdi []byte,
	ow base.AddressDecoder,
	rd uint,
	vt string,
	bcd []byte,
	bn string,
	ac base.AddressDecoder,
	tm string,
) error {

	// unpack document info
	if hinter, err := enc.Decode(bdi); err != nil {
		return err
	} else if i, ok := hinter.(DocInfo); !ok {
		return errors.Errorf("not Document Info: %T", hinter)
	} else {
		doc.info = i
	}

	// decode owner address
	oa, err := ow.Encode(enc)
	if err != nil {
		return err
	}
	doc.owner = oa

	// decode boss account address
	ba, err := ac.Encode(enc)
	if err != nil {
		return err
	}
	doc.account = ba

	doc.round = rd
	doc.endVoteTime = vt

	// unpack candidates
	hits, err := enc.DecodeSlice(bcd)
	if err != nil {
		return err
	}
	candidates := make([]VotingCandidate, len(hits))

	for i := range hits {
		s, ok := hits[i].(VotingCandidate)
		if !ok {
			return errors.Errorf("not VotingCandidate : %T", s)
		}

		candidates[i] = s
	}
	doc.candidates = candidates

	doc.bossname = bn
	doc.termofoffice = tm

	return nil
}

func (doc *BCHistoryData) unpack(
	enc encoder.Encoder,
	bdi []byte,
	ow base.AddressDecoder,
	snm string, // name
	ac base.AddressDecoder, // account address
	sdt string, // date
	sus string, // usage
	sap string, // application
) error {

	// unpack document info
	if hinter, err := enc.Decode(bdi); err != nil {
		return err
	} else if i, ok := hinter.(DocInfo); !ok {
		return errors.Errorf("not Document Info: %T", hinter)
	} else {
		doc.info = i
	}

	oa, err := ow.Encode(enc)
	if err != nil {
		return err
	}
	doc.owner = oa

	ba, err := ac.Encode(enc)
	if err != nil {
		return err
	}
	doc.account = ba

	doc.name = snm
	doc.date = sdt
	doc.usage = sus
	doc.application = sap

	return nil
}

func (us *UserStatistics) unpack(
	enc encoder.Encoder,
	hp,
	st,
	ag,
	dx,
	cr,
	ig,
	vt uint,
) error {

	us.hp = hp
	us.strength = st
	us.agility = ag
	us.dexterity = dx
	us.charisma = cr
	us.intelligence = ig
	us.vital = vt

	return nil
}

func (di *DocInfo) unpack(
	enc encoder.Encoder,
	bi []byte,
	st string,
) error {

	// unpack document info
	if hinter, err := enc.Decode(bi); err != nil {
		return err
	} else if i, ok := hinter.(DocId); !ok {
		return errors.Errorf("not DocId: %T", hinter)
	} else {
		di.id = i
	}

	di.docType = hint.Type(st)

	return nil
}

func (vc *VotingCandidate) unpack(
	enc encoder.Encoder,
	ad base.AddressDecoder,
	ma string,
) error {

	// decode address
	va, err := ad.Encode(enc)
	if err != nil {
		return err
	}
	vc.address = va
	vc.manifest = ma

	return nil
}

func (di *UserDocId) unpack(
	enc encoder.Encoder,
	si string,
) error {

	// unpack document id
	di.s = si

	return nil
}

func (di *LandDocId) unpack(
	enc encoder.Encoder,
	si string,
) error {

	// unpack document id
	di.s = si

	return nil
}

func (di *VotingDocId) unpack(
	enc encoder.Encoder,
	si string,
) error {
	// unpack document id
	di.s = si

	return nil
}

func (di *HistoryDocId) unpack(
	enc encoder.Encoder,
	si string,
) error {
	// unpack document id
	di.s = si

	return nil
}
