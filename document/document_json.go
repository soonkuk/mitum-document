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

type UserDataJSONPacker struct {
	jsonenc.HintedHead
	DI DocInfo        `json:"info"`
	OW base.Address   `json:"owner"`
	GD currency.Big   `json:"gold"`
	BG currency.Big   `json:"bankgold"`
	ST UserStatistics `json:"statistics"`
}

func (doc CityUserData) MarshalJSON() ([]byte, error) {
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
	GD currency.Big        `json:"gold"`
	BG currency.Big        `json:"bankgold"`
	ST json.RawMessage     `json:"statistics"`
}

func (doc *CityUserData) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
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
	LD base.Address `json:"lender"`
	ST string       `json:"starttime"`
	PD uint         `json:"periodday"`
}

func (doc CityLandData) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(LandDataJSONPacker{
		HintedHead: jsonenc.NewHintedHead(doc.Hint()),
		DI:         doc.info,
		OW:         doc.owner,
		LD:         doc.lender,
		ST:         doc.starttime,
		PD:         doc.periodday,
	})
}

type LandDataJSONUnpacker struct {
	DI json.RawMessage     `json:"info"`
	OW base.AddressDecoder `json:"owner"`
	LD base.AddressDecoder `json:"lender"`
	ST string              `json:"starttime"`
	PD uint                `json:"periodday"`
}

func (doc *CityLandData) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uld LandDataJSONUnpacker
	if err := enc.Unmarshal(b, &uld); err != nil {
		return err
	}

	return doc.unpack(enc, uld.DI, uld.OW, uld.LD, uld.ST, uld.PD)
}

type VotingDataJSONPacker struct {
	jsonenc.HintedHead
	DI DocInfo           `json:"info"`
	OW base.Address      `json:"owner"`
	RD uint              `json:"round"`
	CD []VotingCandidate `json:"candidates"`
}

func (doc CityVotingData) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(VotingDataJSONPacker{
		HintedHead: jsonenc.NewHintedHead(doc.Hint()),
		DI:         doc.info,
		OW:         doc.owner,
		RD:         doc.round,
		CD:         doc.candidates,
	})
}

type VotingDataJSONUnpacker struct {
	DI json.RawMessage     `json:"info"`
	OW base.AddressDecoder `json:"owner"`
	RD uint                `json:"round"`
	CD json.RawMessage     `json:"candidates"`
}

func (doc *CityVotingData) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uld VotingDataJSONUnpacker
	if err := enc.Unmarshal(b, &uld); err != nil {
		return err
	}

	return doc.unpack(enc, uld.DI, uld.OW, uld.RD, uld.CD)
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

type VotingCandidatesJSONPacker struct {
	jsonenc.HintedHead
	AD base.Address `json:"address"`
	MA string       `json:"manifest"`
}

func (vc VotingCandidate) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(VotingCandidatesJSONPacker{
		HintedHead: jsonenc.NewHintedHead(vc.Hint()),
		AD:         vc.address,
		MA:         vc.manifest,
	})
}

type VotingCandidatesJSONUnpacker struct {
	AD base.AddressDecoder `json:"address"`
	MA string              `json:"manifest"`
}

func (vc *VotingCandidate) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var uvc VotingCandidatesJSONUnpacker
	if err := enc.Unmarshal(b, &uvc); err != nil {
		return err
	}

	return vc.unpack(enc, uvc.AD, uvc.MA)
}

type UserDocIdJSONPacker struct {
	jsonenc.HintedHead
	SI string `json:"id"`
}

func (di UserDocId) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(UserDocIdJSONPacker{
		HintedHead: jsonenc.NewHintedHead(di.Hint()),
		SI:         di.s,
	})
}

type UserDocIdJSONUnpacker struct {
	SI string `json:"id"`
}

func (di *UserDocId) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udi UserDocIdJSONUnpacker
	if err := enc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.SI)
}

type LandDocIdJSONPacker struct {
	jsonenc.HintedHead
	SI string `json:"id"`
}

func (di LandDocId) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(LandDocIdJSONPacker{
		HintedHead: jsonenc.NewHintedHead(di.Hint()),
		SI:         di.s,
	})
}

type LandDocIdJSONUnpacker struct {
	SI string `json:"id"`
}

func (di *LandDocId) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udi LandDocIdJSONUnpacker
	if err := enc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.SI)
}

type VotingDocIdJSONPacker struct {
	jsonenc.HintedHead
	SI string `json:"id"`
}

func (di VotingDocId) MarshalJSON() ([]byte, error) {
	return jsonenc.Marshal(VotingDocIdJSONPacker{
		HintedHead: jsonenc.NewHintedHead(di.Hint()),
		SI:         di.s,
	})
}

type VotingDocIdJSONUnpacker struct {
	SI string `json:"id"`
}

func (di *VotingDocId) UnpackJSON(b []byte, enc *jsonenc.Encoder) error {
	var udi VotingDocIdJSONUnpacker
	if err := enc.Unmarshal(b, &udi); err != nil {
		return err
	}

	return di.unpack(enc, udi.SI)
}
