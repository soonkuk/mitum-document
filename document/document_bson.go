package document

import (
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

func (doc BCUserData) MarshalBSON() ([]byte, error) {
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

type BCUserDataBSONUnpacker struct {
	DI bson.Raw            `bson:"info"`
	US base.AddressDecoder `bson:"owner"`
	GD uint                `bson:"gold"`
	BG uint                `bson:"bankgold"`
	ST bson.Raw            `bson:"statistics"`
}

func (doc *BCUserData) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udoc BCUserDataBSONUnpacker
	if err := enc.Unmarshal(b, &udoc); err != nil {
		return err
	}

	return doc.unpack(enc, udoc.DI, udoc.US, udoc.GD, udoc.BG, udoc.ST)
}

func (doc BCLandData) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(doc.Hint()),
		bson.M{
			"info":      doc.info,
			"owner":     doc.owner,
			"address":   doc.address,
			"area":      doc.area,
			"renter":    doc.renter,
			"account":   doc.account,
			"rentdate":  doc.rentdate,
			"periodday": doc.periodday,
		}),
	)
}

type BCLandDataBSONUnpacker struct {
	DI bson.Raw            `bson:"info"`
	OW base.AddressDecoder `bson:"owner"`
	AD string              `bson:"address"`
	AR string              `bson:"area"`
	RT string              `bson:"renter"`
	AC base.AddressDecoder `bson:"account"`
	RD string              `bson:"rentdate"`
	PD uint                `bson:"periodday"`
}

func (doc *BCLandData) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uld BCLandDataBSONUnpacker
	if err := bsonenc.Unmarshal(b, &uld); err != nil {
		return err
	}

	return doc.unpack(enc, uld.DI, uld.OW, uld.AD, uld.AR, uld.RT, uld.AC, uld.RD, uld.PD)
}

func (doc BCVotingData) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(doc.Hint()),
		bson.M{
			"info":         doc.info,
			"owner":        doc.owner,
			"round":        doc.round,
			"endvotetime":  doc.endVoteTime,
			"candidates":   doc.candidates,
			"bossname":     doc.bossname,
			"account":      doc.account,
			"termofoffice": doc.termofoffice,
		}),
	)
}

type BCVotingDataBSONUnpacker struct {
	DI bson.Raw            `bson:"info"`
	OW base.AddressDecoder `bson:"owner"`
	RD uint                `bson:"round"`
	VT string              `bson:"endvotetime"`
	CD bson.Raw            `bson:"candidates"`
	BN string              `bson:"bossname"`
	AC base.AddressDecoder `bson:"account"`
	TM string              `bson:"termofoffice"`
}

func (doc *BCVotingData) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uvd BCVotingDataBSONUnpacker
	if err := bsonenc.Unmarshal(b, &uvd); err != nil {
		return err
	}

	return doc.unpack(enc, uvd.DI, uvd.OW, uvd.RD, uvd.VT, uvd.CD, uvd.BN, uvd.AC, uvd.TM)
}

func (doc BCHistoryData) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(doc.Hint()),
		bson.M{
			"info":        doc.info,
			"owner":       doc.owner,
			"name":        doc.name,
			"account":     doc.account,
			"date":        doc.date,
			"usage":       doc.usage,
			"application": doc.application,
		}),
	)
}

type BCHistoryDataBSONUnpacker struct {
	DI bson.Raw            `bson:"info"`
	OW base.AddressDecoder `bson:"owner"`
	NM string              `bson:"name"`
	AC base.AddressDecoder `bson:"account"`
	DT string              `bson:"date"`
	US string              `bson:"usage"`
	AP string              `bson:"application"`
}

func (doc *BCHistoryData) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uhd BCHistoryDataBSONUnpacker
	if err := bsonenc.Unmarshal(b, &uhd); err != nil {
		return err
	}

	return doc.unpack(enc, uhd.DI, uhd.OW, uhd.NM, uhd.AC, uhd.DT, uhd.US, uhd.AP)
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
			"nickname": di.nickname,
			"manifest": di.manifest,
		}),
	)
}

type VotingCandidateBSONUnpacker struct {
	AD base.AddressDecoder `bson:"address"`
	NC string              `bson:"nickname"`
	MA string              `bson:"manifest"`
}

func (di *VotingCandidate) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var uvc VotingCandidateBSONUnpacker
	if err := bsonenc.Unmarshal(b, &uvc); err != nil {
		return err
	}

	return di.unpack(enc, uvc.AD, uvc.NC, uvc.MA)
}

func (di UserDocId) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(di.Hint()),
		bson.M{
			"id": di.s,
		}),
	)
}

type DocIdBSONUnpacker struct {
	BI string `bson:"id"`
}

func (di *UserDocId) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udi DocIdBSONUnpacker
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

func (di *LandDocId) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udi DocIdBSONUnpacker
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

func (di *VotingDocId) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udi DocIdBSONUnpacker
	if err := bsonenc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.BI)
}

func (di HistoryDocId) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bsonenc.MergeBSONM(
		bsonenc.NewHintedDoc(di.Hint()),
		bson.M{
			"id": di.s,
		}),
	)
}

func (di *HistoryDocId) UnpackBSON(b []byte, enc *bsonenc.Encoder) error {
	var udi DocIdBSONUnpacker
	if err := bsonenc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.BI)
}
