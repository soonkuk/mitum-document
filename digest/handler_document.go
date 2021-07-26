package digest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/soonkuk/mitum-data/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/xerrors"
)

func (hd *Handlers) handleDocument(w http.ResponseWriter, r *http.Request) {
	cachekey := cacheKeyPath(r)
	if err := loadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	h, err := parseDocIdFromPath(mux.Vars(r)["documentid"])
	if err != nil {
		hd.problemWithError(w, xerrors.Errorf("invalid document id for document by id: %w", err), http.StatusBadRequest)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleDocumentInGroup(h)
	}); err != nil {
		hd.handleError(w, err)
	} else {
		hd.writeHalBytes(w, v.([]byte), http.StatusOK)

		if !shared {
			hd.writeCache(w, cachekey, time.Hour*30)
		}
	}
}

func (hd *Handlers) handleDocumentInGroup(i currency.Big) ([]byte, error) {
	switch va, found, err := hd.database.Document(i); {
	case err != nil:
		return nil, err
	case !found:
		return nil, util.NotFoundError.Errorf("document value not found")
	default:
		hal, err := hd.buildDocumentHal(va)
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("document:{documentid}", NewHalLink(HandlerPathDocument, nil).SetTemplated())
		hal = hal.AddLink("block:{height}", NewHalLink(HandlerPathBlockByHeight, nil).SetTemplated())

		return hd.enc.Marshal(hal)
	}
}

func (hd *Handlers) handleDocuments(w http.ResponseWriter, r *http.Request) {
	offset := parseOffsetQuery(r.URL.Query().Get("offset"))
	reverse := parseBoolQuery(r.URL.Query().Get("reverse"))

	cachekey := cacheKey(r.URL.Path, stringOffsetQuery(offset), stringBoolQuery("reverse", reverse))
	if err := loadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		i, filled, err := hd.handleDocumentsInGroup(offset, reverse)

		return []interface{}{i, filled}, err
	}); err != nil {
		hd.handleError(w, err)
	} else {
		var b []byte
		var filled bool
		{
			l := v.([]interface{})
			b = l[0].([]byte)
			filled = l[1].(bool)
		}

		hd.writeHalBytes(w, b, http.StatusOK)

		if !shared {
			expire := time.Second * 3
			if filled {
				expire = time.Hour * 30
			}

			hd.writeCache(w, cachekey, expire)
		}
	}
}

func (hd *Handlers) handleDocumentsInGroup(offset string, reverse bool) ([]byte, bool, error) {
	filter, err := buildDocumentsFilterByOffset(offset, reverse)
	if err != nil {
		return nil, false, err
	}

	var vas []Hal
	switch l, e := hd.loadDocumentsHALFromDatabase(filter, reverse); {
	case e != nil:
		return nil, false, e
	case len(l) < 1:
		return nil, false, util.NotFoundError.Errorf("documents not found")
	default:
		vas = l
	}

	h, err := hd.combineURL(HandlerPathDocuments)
	if err != nil {
		return nil, false, err
	}
	hal := hd.buildDocumentsHal(h, vas, offset, reverse)
	if next := nextOffsetOfDocuments(h, vas, reverse); len(next) > 0 {
		hal = hal.AddLink("next", NewHalLink(next, nil))
	}

	b, err := hd.enc.Marshal(hal)
	return b, int64(len(vas)) == hd.itemsLimiter("documents"), err
}

func (hd *Handlers) handleDocumentsByHeight(w http.ResponseWriter, r *http.Request) {
	offset := parseOffsetQuery(r.URL.Query().Get("offset"))
	reverse := parseBoolQuery(r.URL.Query().Get("reverse"))

	cachekey := cacheKey(r.URL.Path, stringOffsetQuery(offset), stringBoolQuery("reverse", reverse))
	if err := loadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	var height base.Height
	switch h, err := parseHeightFromPath(mux.Vars(r)["height"]); {
	case err != nil:
		hd.problemWithError(w, xerrors.Errorf("invalid height found for manifest by height"), http.StatusBadRequest)

		return
	case h <= base.NilHeight:
		hd.problemWithError(w, xerrors.Errorf("invalid height, %v", h), http.StatusBadRequest)
		return
	default:
		height = h
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		i, filled, err := hd.handleDocumentsByHeightInGroup(height, offset, reverse)
		return []interface{}{i, filled}, err
	}); err != nil {
		hd.handleError(w, err)
	} else {
		var b []byte
		var filled bool
		{
			l := v.([]interface{})
			b = l[0].([]byte)
			filled = l[1].(bool)
		}

		hd.writeHalBytes(w, b, http.StatusOK)

		if !shared {
			expire := time.Second * 3
			if filled {
				expire = time.Hour * 30
			}

			hd.writeCache(w, cachekey, expire)
		}
	}
}

func (hd *Handlers) handleDocumentsByHeightInGroup(
	height base.Height,
	offset string,
	reverse bool,
) ([]byte, bool, error) {
	filter, err := buildDocumentsByHeightFilterByOffset(height, offset, reverse)
	if err != nil {
		return nil, false, err
	}

	var vas []Hal
	switch l, e := hd.loadDocumentsHALFromDatabase(filter, reverse); {
	case e != nil:
		return nil, false, e
	case len(l) < 1:
		return nil, false, util.NotFoundError.Errorf("documents not found")
	default:
		vas = l
	}

	h, err := hd.combineURL(HandlerPathDocumentsByHeight, "height", height.String())
	if err != nil {
		return nil, false, err
	}
	hal := hd.buildDocumentsHal(h, vas, offset, reverse)
	if next := nextOffsetOfDocumentsByHeight(h, vas, reverse); len(next) > 0 {
		hal = hal.AddLink("next", NewHalLink(next, nil))
	}

	b, err := hd.enc.Marshal(hal)
	return b, int64(len(vas)) == hd.itemsLimiter("documents"), err
}

func (hd *Handlers) buildDocumentHal(va DocumentValue) (Hal, error) {
	var hal Hal

	h, err := hd.combineURL(HandlerPathDocument, "documentid", va.Document().Info().Index().String())
	if err != nil {
		return nil, err
	}
	hal = NewBaseHal(va, NewHalLink(h, nil))

	h, err = hd.combineURL(HandlerPathBlockByHeight, "height", va.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", NewHalLink(h, nil))

	h, err = hd.combineURL(HandlerPathManifestByHeight, "height", va.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("manifest", NewHalLink(h, nil))

	return hal, nil
}

func (*Handlers) buildDocumentsHal(baseSelf string, vas []Hal, offset string, reverse bool) Hal {
	var hal Hal

	self := baseSelf
	if len(offset) > 0 {
		self = addQueryValue(baseSelf, stringOffsetQuery(offset))
	}
	if reverse {
		self = addQueryValue(self, stringBoolQuery("reverse", reverse))
	}
	hal = NewBaseHal(vas, NewHalLink(self, nil))

	hal = hal.AddLink("reverse", NewHalLink(addQueryValue(baseSelf, stringBoolQuery("reverse", !reverse)), nil))

	return hal
}

func buildDocumentsFilterByOffset(offset string, reverse bool) (bson.M, error) {
	filter := bson.M{}
	if len(offset) > 0 {
		height, index, err := parseOffset(offset)
		if err != nil {
			return nil, err
		}

		if reverse {
			filter["$or"] = []bson.M{
				{"height": bson.M{"$lt": height}},
				{"$and": []bson.M{
					{"height": height},
					{"index": bson.M{"$lt": index}},
				}},
			}
		} else {
			filter["$or"] = []bson.M{
				{"height": bson.M{"$gt": height}},
				{"$and": []bson.M{
					{"height": height},
					{"index": bson.M{"$gt": index}},
				}},
			}
		}
	}

	return filter, nil
}

func buildDocumentsByHeightFilterByOffset(height base.Height, offset string, reverse bool) (bson.M, error) {
	var filter bson.M
	if len(offset) < 1 {
		return bson.M{"height": height}, nil
	}

	index, err := strconv.ParseUint(offset, 10, 64)
	if err != nil {
		return nil, xerrors.Errorf("invalid index of offset: %w", err)
	}

	if reverse {
		filter = bson.M{
			"height": height,
			"index":  bson.M{"$lt": index},
		}
	} else {
		filter = bson.M{
			"height": height,
			"index":  bson.M{"$gt": index},
		}
	}

	return filter, nil
}

func nextOffsetOfDocuments(baseSelf string, vas []Hal, reverse bool) string {
	var nextoffset string
	if len(vas) > 0 {
		va := vas[len(vas)-1].Interface().(DocumentValue)
		nextoffset = buildOffset(va.Height(), va.Document().Info().Index().Uint64())
	}

	if len(nextoffset) < 1 {
		return ""
	}

	next := baseSelf
	if len(nextoffset) > 0 {
		next = addQueryValue(next, stringOffsetQuery(nextoffset))
	}

	if reverse {
		next = addQueryValue(next, stringBoolQuery("reverse", reverse))
	}

	return next
}

func nextOffsetOfDocumentsByHeight(baseSelf string, vas []Hal, reverse bool) string {
	var nextoffset string
	if len(vas) > 0 {
		va := vas[len(vas)-1].Interface().(DocumentValue)
		nextoffset = fmt.Sprintf("%d", va.Document().Info().Index().Uint64())
	}

	if len(nextoffset) < 1 {
		return ""
	}

	next := baseSelf
	if len(nextoffset) > 0 {
		next = addQueryValue(next, stringOffsetQuery(nextoffset))
	}

	if reverse {
		next = addQueryValue(next, stringBoolQuery("reverse", reverse))
	}

	return next
}

func (hd *Handlers) loadDocumentsHALFromDatabase(filter bson.M, reverse bool) ([]Hal, error) {
	var vas []Hal
	if err := hd.database.Documents(
		filter, reverse, hd.itemsLimiter("documents"),
		func(_ currency.Big, va DocumentValue) (bool, error) {
			hal, err := hd.buildDocumentHal(va)
			if err != nil {
				return false, err
			}
			vas = append(vas, hal)

			return true, nil
		},
	); err != nil {
		return nil, err
	} else if len(vas) < 1 {
		return nil, nil
	}

	return vas, nil
}
