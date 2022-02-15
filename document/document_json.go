package document

import (
	"encoding/json"

	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/hint"
)

type DocumentJSONPacker struct {
	jsonenc.HintedHead
	DC DocumentData `json:"documentdata"`
}

func (doc Document) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocumentJSONPacker{
		HintedHead: jsonenc.NewHintedHead(doc.Hint()),
		DC:         doc.data,
	})
}

type DocumentJSONUnpacker struct {
	DC json.RawMessage `json:"documentdata"`
}

func (doc *Document) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var dod DocumentJSONUnpacker
	if err := enc.Unmarshal(b, &dod); err != nil {
		return err
	}

	return doc.unpack(enc, dod.DC)
}

type BSDocDataJSONPacker struct {
	jsonenc.HintedHead
	DI DocInfo      `json:"info"`
	OW base.Address `json:"owner"`
	FH FileHash     `json:"filehash"`
	CR DocSign      `json:"creator"`
	TL string       `json:"title"`
	SZ currency.Big `json:"size"`
	SG []DocSign    `json:"signers"`
}

func (doc BSDocData) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(BSDocDataJSONPacker{
		HintedHead: jsonenc.NewHintedHead(doc.Hint()),
		DI:         doc.info,
		OW:         doc.owner,
		FH:         doc.fileHash,
		CR:         doc.creator,
		TL:         doc.title,
		SZ:         doc.size,
		SG:         doc.signers,
	})
}

type BSDocDataJSONUnpacker struct {
	DI json.RawMessage     `json:"info"`
	OW base.AddressDecoder `json:"owner"`
	FH string              `json:"filehash"`
	CR json.RawMessage     `json:"creator"`
	TL string              `json:"title"`
	SZ currency.Big        `json:"size"`
	SG json.RawMessage     `json:"signers"`
}

func (doc *BSDocData) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udoc BSDocDataJSONUnpacker
	if err := enc.Unmarshal(b, &udoc); err != nil {
		return err
	}

	return doc.unpack(enc, udoc.DI, udoc.OW, udoc.FH, udoc.CR, udoc.TL, udoc.SZ, udoc.SG)
}

type UserDataJSONPacker struct {
	jsonenc.HintedHead
	DI DocInfo        `json:"info"`
	OW base.Address   `json:"owner"`
	GD uint           `json:"gold"`
	BG uint           `json:"bankgold"`
	ST UserStatistics `json:"statistics"`
}

func (doc BCUserData) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(UserDataJSONPacker{
		HintedHead: jsonenc.NewHintedHead(doc.Hint()),
		DI:         doc.info,
		OW:         doc.owner,
		GD:         doc.gold,
		BG:         doc.bankgold,
		ST:         doc.statistics,
	})
}

type UserDataJSONUnpacker struct {
	DI json.RawMessage     `json:"info"`
	OW base.AddressDecoder `json:"owner"`
	GD uint                `json:"gold"`
	BG uint                `json:"bankgold"`
	ST json.RawMessage     `json:"statistics"`
}

func (doc *BCUserData) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udoc UserDataJSONUnpacker
	if err := enc.Unmarshal(b, &udoc); err != nil {
		return err
	}

	return doc.unpack(enc, udoc.DI, udoc.OW, udoc.GD, udoc.BG, udoc.ST)
}

type LandDataJSONPacker struct {
	jsonenc.HintedHead
	DI DocInfo      `json:"info"`
	OW base.Address `json:"owner"`
	AD string       `json:"address"`
	AR string       `json:"area"`
	RT string       `json:"renter"`
	AC base.Address `json:"account"`
	RD string       `json:"rentdate"`
	PD uint         `json:"periodday"`
}

func (doc BCLandData) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(LandDataJSONPacker{
		HintedHead: jsonenc.NewHintedHead(doc.Hint()),
		DI:         doc.info,
		OW:         doc.owner,
		AD:         doc.address,
		AR:         doc.area,
		RT:         doc.renter,
		AC:         doc.account,
		RD:         doc.rentdate,
		PD:         doc.periodday,
	})
}

type LandDataJSONUnpacker struct {
	DI json.RawMessage     `json:"info"`
	OW base.AddressDecoder `json:"owner"`
	AD string              `json:"address"`
	AR string              `json:"area"`
	RT string              `json:"renter"`
	AC base.AddressDecoder `json:"account"`
	RD string              `json:"rentdate"`
	PD uint                `json:"periodday"`
}

func (doc *BCLandData) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uld LandDataJSONUnpacker
	if err := enc.Unmarshal(b, &uld); err != nil {
		return err
	}

	return doc.unpack(enc, uld.DI, uld.OW, uld.AD, uld.AR, uld.RT, uld.AC, uld.RD, uld.PD)
}

type VotingDataJSONPacker struct {
	jsonenc.HintedHead
	DI DocInfo           `json:"info"`
	OW base.Address      `json:"owner"`
	RD uint              `json:"round"`
	VT string            `json:"endvotetime"`
	CD []VotingCandidate `json:"candidates"`
	BN string            `json:"bossname"`
	AC base.Address      `json:"account"`
	TM string            `json:"termofoffice"`
}

func (doc BCVotingData) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(VotingDataJSONPacker{
		HintedHead: jsonenc.NewHintedHead(doc.Hint()),
		DI:         doc.info,
		OW:         doc.owner,
		RD:         doc.round,
		VT:         doc.endVoteTime,
		CD:         doc.candidates,
		BN:         doc.bossname,
		AC:         doc.account,
		TM:         doc.termofoffice,
	})
}

type VotingDataJSONUnpacker struct {
	DI json.RawMessage     `json:"info"`
	OW base.AddressDecoder `json:"owner"`
	RD uint                `json:"round"`
	VT string              `json:"endvotetime"`
	CD json.RawMessage     `json:"candidates"`
	BN string              `json:"bossname"`
	AC base.AddressDecoder `json:"account"`
	TM string              `json:"termofoffice"`
}

func (doc *BCVotingData) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uvd VotingDataJSONUnpacker
	if err := enc.Unmarshal(b, &uvd); err != nil {
		return err
	}

	return doc.unpack(enc, uvd.DI, uvd.OW, uvd.RD, uvd.VT, uvd.CD, uvd.BN, uvd.AC, uvd.TM)
}

type HistoryDataJSONPacker struct {
	jsonenc.HintedHead
	DI DocInfo      `json:"info"`
	OW base.Address `json:"owner"`
	NM string       `json:"name"`
	AC base.Address `json:"account"`
	DT string       `json:"date"`
	US string       `json:"usage"`
	AP string       `json:"application"`
}

func (doc BCHistoryData) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(HistoryDataJSONPacker{
		HintedHead: jsonenc.NewHintedHead(doc.Hint()),
		DI:         doc.info,
		OW:         doc.owner,
		NM:         doc.name,
		AC:         doc.account,
		DT:         doc.date,
		US:         doc.usage,
		AP:         doc.application,
	})
}

type HistoryDataJSONUnpacker struct {
	DI json.RawMessage     `json:"info"`
	OW base.AddressDecoder `json:"owner"`
	NM string              `json:"name"`
	AC base.AddressDecoder `json:"account"`
	DT string              `json:"date"`
	US string              `json:"usage"`
	AP string              `json:"application"`
}

func (doc *BCHistoryData) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uhd HistoryDataJSONUnpacker
	if err := enc.Unmarshal(b, &uhd); err != nil {
		return err
	}

	return doc.unpack(enc, uhd.DI, uhd.OW, uhd.NM, uhd.AC, uhd.DT, uhd.US, uhd.AP)
}

type UserStatisticsJSONPacker struct {
	jsonenc.HintedHead
	HP uint `json:"hp"`
	ST uint `json:"strength"`
	AG uint `json:"agility"`
	DX uint `json:"dexterity"`
	CR uint `json:"charisma"`
	IG uint `json:"intelligence"`
	VT uint `json:"vital"`
}

func (us UserStatistics) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(UserStatisticsJSONPacker{
		HintedHead: jsonenc.NewHintedHead(us.Hint()),
		HP:         us.hp,
		ST:         us.strength,
		AG:         us.agility,
		DX:         us.dexterity,
		CR:         us.charisma,
		IG:         us.intelligence,
		VT:         us.vital,
	})
}

type UserStatisticsJSONUnpacker struct {
	HP uint `json:"hp"`
	ST uint `json:"strength"`
	AG uint `json:"agility"`
	DX uint `json:"dexterity"`
	CR uint `json:"charisma"`
	IG uint `json:"intelligence"`
	VT uint `json:"vital"`
}

func (us *UserStatistics) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uus UserStatisticsJSONUnpacker
	if err := enc.Unmarshal(b, &uus); err != nil {
		return err
	}

	return us.unpack(enc, uus.HP, uus.ST, uus.AG, uus.DX, uus.CR, uus.IG, uus.VT)
}

type DocInfoJSONPacker struct {
	jsonenc.HintedHead
	ID DocId     `json:"docid"`
	DT hint.Type `json:"doctype"`
}

func (di DocInfo) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocInfoJSONPacker{
		HintedHead: jsonenc.NewHintedHead(di.Hint()),
		ID:         di.id,
		DT:         di.docType,
	})
}

type DocInfoJSONUnpacker struct {
	ID json.RawMessage `json:"docid"`
	DT string          `json:"docType"`
}

func (di *DocInfo) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udi DocInfoJSONUnpacker
	if err := enc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.ID, udi.DT)
}

type DocSignJSONPacker struct {
	jsonenc.HintedHead
	AD base.Address `json:"address"`
	SC string       `json:"signcode"`
	SG bool         `json:"signed"`
}

func (ds DocSign) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocSignJSONPacker{
		HintedHead: jsonenc.NewHintedHead(ds.Hint()),
		AD:         ds.address,
		SC:         ds.signcode,
		SG:         ds.signed,
	})
}

type DocSignJSONUnpacker struct {
	AD base.AddressDecoder `json:"address"`
	SC string              `json:"signcode"`
	SG bool                `json:"signed"`
}

func (ds *DocSign) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uds DocSignJSONUnpacker
	if err := enc.Unmarshal(b, &uds); err != nil {
		return err
	}

	return ds.unpack(enc, uds.AD, uds.SC, uds.SG)
}

type VotingCandidatesJSONPacker struct {
	jsonenc.HintedHead
	AD base.Address `json:"address"`
	NC string       `json:"nickname"`
	MA string       `json:"manifest"`
}

func (vc VotingCandidate) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(VotingCandidatesJSONPacker{
		HintedHead: jsonenc.NewHintedHead(vc.Hint()),
		AD:         vc.address,
		NC:         vc.nickname,
		MA:         vc.manifest,
	})
}

type VotingCandidatesJSONUnpacker struct {
	AD base.AddressDecoder `json:"address"`
	NC string              `json:"nickname"`
	MA string              `json:"manifest"`
}

func (vc *VotingCandidate) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uvc VotingCandidatesJSONUnpacker
	if err := enc.Unmarshal(b, &uvc); err != nil {
		return err
	}

	return vc.unpack(enc, uvc.AD, uvc.NC, uvc.MA)
}

type DocIdJSONPacker struct {
	jsonenc.HintedHead
	SI string `json:"id"`
}

type DocIdJSONUnpacker struct {
	SI string `json:"id"`
}

func (di BSDocId) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocIdJSONPacker{
		HintedHead: jsonenc.NewHintedHead(di.Hint()),
		SI:         di.s,
	})
}

func (di *BSDocId) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udi DocIdJSONUnpacker
	if err := enc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.SI)
}

func (di UserDocId) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocIdJSONPacker{
		HintedHead: jsonenc.NewHintedHead(di.Hint()),
		SI:         di.s,
	})
}

func (di *UserDocId) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udi DocIdJSONUnpacker
	if err := enc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.SI)
}

func (di LandDocId) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocIdJSONPacker{
		HintedHead: jsonenc.NewHintedHead(di.Hint()),
		SI:         di.s,
	})
}

func (di *LandDocId) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udi DocIdJSONUnpacker
	if err := enc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.SI)
}

func (di VotingDocId) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocIdJSONPacker{
		HintedHead: jsonenc.NewHintedHead(di.Hint()),
		SI:         di.s,
	})
}

func (di *VotingDocId) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udi DocIdJSONUnpacker
	if err := enc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.SI)
}

func (di HistoryDocId) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(DocIdJSONPacker{
		HintedHead: jsonenc.NewHintedHead(di.Hint()),
		SI:         di.s,
	})
}

func (di *HistoryDocId) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udi DocIdJSONUnpacker
	if err := enc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.SI)
}
