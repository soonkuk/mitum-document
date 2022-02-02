package document

import (
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	bsonenc "github.com/spikeekips/mitum/util/encoder/bson"
	"go.mongodb.org/mongo-driver/bson"
)

func (doc Document) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(doc.Hint()),
		bson.M{
			"documentdata": doc.data,
		}),
	)
}

type DocumentBSONUnpacker struct {
	DC bson.Raw `bson:"documentdata"`
}

func (doc *Document) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var dod DocumentBSONUnpacker
	if err := bsonenc.Unmarshal(b, &dod); err != nil {
		return err
	}

	return doc.unpack(enc, dod.DC)
}

func (doc CityUserData) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(doc.Hint()),
		bson.M{
			"info":       doc.info,
			"owner":      doc.owner,
			"gold":       doc.gold,
			"bankgold":   doc.bankgold,
			"statistics": doc.statistics,
		}),
	)
}

type CityUserDataBSONUnpacker struct {
	DI bson.Raw            `bson:"info"`
	US base.AddressDecoder `bson:"owner"`
	GD currency.Big        `bson:"gold"`
	BG currency.Big        `bson:"bankgold"`
	ST bson.Raw            `bson:"statistics"`
}

func (doc *CityUserData) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udoc CityUserDataBSONUnpacker
	if err := enc.Unmarshal(b, &udoc); err != nil {
		return err
	}

	return doc.unpack(enc, udoc.DI, udoc.US, udoc.GD, udoc.BG, udoc.ST)
}

func (doc CityLandData) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(doc.Hint()),
		bson.M{
			"info":      doc.info,
			"owner":     doc.owner,
			"lender":    doc.lender,
			"starttime": doc.starttime,
			"periodday": doc.periodday,
		}),
	)
}

type CityLandDataBSONUnpacker struct {
	IN bson.Raw            `bson:"info"`
	OW base.AddressDecoder `bson:"owner"`
	LD base.AddressDecoder `bson:"lender"`
	ST string              `bson:"starttime"`
	PD uint                `bson:"periodday"`
}

func (doc *CityLandData) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uld CityLandDataBSONUnpacker
	if err := bsonenc.Unmarshal(b, &uld); err != nil {
		return err
	}

	return doc.unpack(enc, uld.IN, uld.OW, uld.LD, uld.ST, uld.PD)
}

func (doc CityVotingData) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(doc.Hint()),
		bson.M{
			"info":       doc.info,
			"owner":      doc.owner,
			"round":      doc.round,
			"candidates": doc.candidates,
		}),
	)
}

type CityVotingDataBSONUnpacker struct {
	IN bson.Raw            `bson:"info"`
	OW base.AddressDecoder `bson:"owner"`
	RD uint                `bson:"round"`
	CD bson.Raw            `bson:"candidates"`
}

func (doc *CityVotingData) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uld CityVotingDataBSONUnpacker
	if err := bsonenc.Unmarshal(b, &uld); err != nil {
		return err
	}

	return doc.unpack(enc, uld.IN, uld.OW, uld.RD, uld.CD)
}

func (us UserStatistics) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(us.Hint()),
		bson.M{
			"hp":           us.hp,
			"strength":     us.strength,
			"agility":      us.agility,
			"dexterity":    us.dexterity,
			"charisma":     us.charisma,
			"intelligence": us.intelligence,
			"vital":        us.vital,
		}),
	)
}

type UserStatisticsBSONUnpacker struct {
	HP uint `bson:"hp"`
	ST uint `bson:"strength"`
	AG uint `bson:"agility"`
	DX uint `bson:"dexterity"`
	CR uint `bson:"charisma"`
	IG uint `bson:"intelligence"`
	VT uint `bson:"vital"`
}

func (us *UserStatistics) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uus UserStatisticsBSONUnpacker
	if err := bsonenc.Unmarshal(b, &uus); err != nil {
		return err
	}

	return us.unpack(enc, uus.HP, uus.ST, uus.AG, uus.DX, uus.CR, uus.IG, uus.VT)
}

func (di DocInfo) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(di.Hint()),
		bson.M{
			"docid":   di.id,
			"doctype": di.docType,
		}),
	)
}

type DocInfoBSONUnpacker struct {
	ID bson.Raw `bson:"docid"`
	DT string   `bson:"doctype"`
}

func (di *DocInfo) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udi DocInfoBSONUnpacker
	if err := bsonenc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.ID, udi.DT)
}

func (di VotingCandidate) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(di.Hint()),
		bson.M{
			"address":  di.address,
			"manifest": di.manifest,
		}),
	)
}

type VotingCandidateBSONUnpacker struct {
	AD base.AddressDecoder `bson:"address"`
	MA string              `bson:"manifest"`
}

func (di *VotingCandidate) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udi VotingCandidateBSONUnpacker
	if err := bsonenc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.AD, udi.MA)
}

func (di UserDocId) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(di.Hint()),
		bson.M{
			"id": di.s,
		}),
	)
}

type UserDocIdBSONUnpacker struct {
	BI string `bson:"id"`
}

func (di *UserDocId) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udi UserDocIdBSONUnpacker
	if err := bsonenc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.BI)
}

func (di LandDocId) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(di.Hint()),
		bson.M{
			"id": di.s,
		}),
	)
}

type LandDocIdBSONUnpacker struct {
	BI string `bson:"id"`
}

func (di *LandDocId) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udi LandDocIdBSONUnpacker
	if err := bsonenc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.BI)
}

func (di VotingDocId) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(di.Hint()),
		bson.M{
			"id": di.s,
		}),
	)
}

type VotingDocIdBSONUnpacker struct {
	BI string `bson:"id"`
}

func (di *VotingDocId) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udi VotingDocIdBSONUnpacker
	if err := bsonenc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.BI)
}
