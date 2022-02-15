package digest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"go.mongodb.org/mongo-driver/bson"
)

/*
func (hd *Handlers) handleBSDocument(w http.ResponseWriter, r *http.Request) {

	cachekey := CacheKeyPath(r)

	if err := LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	h, err := parseDocIdFromPath(mux.Vars(r)["documentid"])
	if err != nil {
		HTTP2ProblemWithError(w, errors.Errorf("invalid document id for document by id: %q", err), http.StatusBadRequest)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleBSDocumentInGroup(h)
	}); err != nil {
		HTTP2HandleError(w, err)
	} else {
		HTTP2WriteHalBytes(hd.enc, w, v.([]byte), http.StatusOK)

		if !shared {
			HTTP2WriteCache(w, cachekey, time.Second*2)
		}
	}
}
*/

/*
func (hd *Handlers) handleBSDocumentInGroup(i string) ([]byte, error) {
	switch va, found, err := hd.database.BSDocument(i); {
	case err != nil:
		return nil, err
	case !found:
		return nil, util.NotFoundError.Errorf("document value not found")
	default:
		hal, err := hd.buildBSDocumentHal(va)
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("bsdocument:{documentid}", NewHalLink(HandlerPathBSDocument, nil).SetTemplated())
		hal = hal.AddLink("block:{height}", NewHalLink(HandlerPathBlockByHeight, nil).SetTemplated())

		return hd.enc.Marshal(hal)
	}
}
*/

func (hd *Handlers) handleDocuments(w http.ResponseWriter, r *http.Request) {
	limit := parseLimitQuery(r.URL.Query().Get("limit"))
	offset := parseOffsetQuery(r.URL.Query().Get("offset"))
	reverse := parseBoolQuery(r.URL.Query().Get("reverse"))

	cachekey := CacheKey(r.URL.Path, stringOffsetQuery(offset), stringBoolQuery("reverse", reverse))

	if err := LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		i, filled, err := hd.handleDocumentsInGroup(offset, reverse, limit)

		return []interface{}{i, filled}, err
	}); err != nil {
		HTTP2HandleError(w, err)
	} else {
		var b []byte
		var filled bool
		{
			l := v.([]interface{})
			b = l[0].([]byte)
			filled = l[1].(bool)
		}

		HTTP2WriteHalBytes(hd.enc, w, b, http.StatusOK)

		if !shared {
			expire := hd.expireNotFilled
			if len(offset) > 0 && filled {
				expire = time.Hour * 30
			}

			HTTP2WriteCache(w, cachekey, expire)
		}
	}
}

func (hd *Handlers) handleDocumentsInGroup(
	offset string,
	reverse bool,
	l int64,
) ([]byte, bool, error) {
	var limit int64
	if l < 0 {
		limit = hd.itemsLimiter("documents")
	} else {
		limit = l
	}
	filter, err := buildDocumentsFilterByOffset(offset, reverse)
	if err != nil {
		return nil, false, err
	}

	var vas []Hal
	switch l, e := hd.loadDocumentsHALFromDatabase(filter, reverse, limit); {
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
	return b, int64(len(vas)) == limit, err
}

func (hd *Handlers) handleDocument(w http.ResponseWriter, r *http.Request) {

	cachekey := CacheKeyPath(r)

	if err := LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	h, err := parseDocIdFromPath(mux.Vars(r)["documentid"])
	if err != nil {
		HTTP2ProblemWithError(w, errors.Errorf("invalid document id for document by id: %q", err), http.StatusBadRequest)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleDocumentInGroup(h)
	}); err != nil {
		HTTP2HandleError(w, err)
	} else {
		HTTP2WriteHalBytes(hd.enc, w, v.([]byte), http.StatusOK)

		if !shared {
			HTTP2WriteCache(w, cachekey, time.Second*2)
		}
	}
}

func (hd *Handlers) handleDocumentInGroup(i string) ([]byte, error) {
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
		hal = hal.AddLink("bcdocument:{documentid}", NewHalLink(HandlerPathDocument, nil).SetTemplated())
		hal = hal.AddLink("block:{height}", NewHalLink(HandlerPathBlockByHeight, nil).SetTemplated())

		return hd.enc.Marshal(hal)
	}
}

func (hd *Handlers) handleDocumentsByHeight(w http.ResponseWriter, r *http.Request) {
	limit := parseLimitQuery(r.URL.Query().Get("limit"))
	offset := parseOffsetQuery(r.URL.Query().Get("offset"))
	reverse := parseBoolQuery(r.URL.Query().Get("reverse"))

	cachekey := CacheKey(r.URL.Path, stringOffsetQuery(offset), stringBoolQuery("reverse", reverse))

	if err := LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	var height base.Height
	switch h, err := parseHeightFromPath(mux.Vars(r)["height"]); {
	case err != nil:
		HTTP2ProblemWithError(w, errors.Errorf("invalid height found for manifest by height"), http.StatusBadRequest)

		return
	case h <= base.NilHeight:
		HTTP2ProblemWithError(w, errors.Errorf("invalid height, %v", h), http.StatusBadRequest)
		return
	default:
		height = h
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		i, filled, err := hd.handleDocumentsByHeightInGroup(height, offset, reverse, limit)
		return []interface{}{i, filled}, err
	}); err != nil {
		HTTP2HandleError(w, err)
	} else {
		var b []byte
		var filled bool
		{
			l := v.([]interface{})
			b = l[0].([]byte)
			filled = l[1].(bool)
		}

		HTTP2WriteHalBytes(hd.enc, w, b, http.StatusOK)

		if !shared {
			expire := hd.expireNotFilled
			if len(offset) > 0 && filled {
				expire = time.Hour * 30
			}

			HTTP2WriteCache(w, cachekey, expire)
		}
	}
}

func (hd *Handlers) handleDocumentsByHeightInGroup(
	height base.Height,
	offset string,
	reverse bool,
	l int64,
) ([]byte, bool, error) {
	var limit int64
	if l < 0 {
		limit = hd.itemsLimiter("documents")
	} else {
		limit = l
	}
	filter, err := buildDocumentsByHeightFilterByOffset(height, offset, reverse)
	if err != nil {
		return nil, false, err
	}

	var vas []Hal
	switch l, e := hd.loadDocumentsHALFromDatabase(filter, reverse, limit); {
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
	return b, int64(len(vas)) == limit, err
}

/*
func (hd *Handlers) buildBSDocumentHal(va BSDocumentValue) (Hal, error) {
	var hal Hal

	h, err := hd.combineURL(HandlerPathBSDocument, "documentid", va.Document().Info().Index().String())
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
*/
func (hd *Handlers) buildDocumentHal(va DocumentValue) (Hal, error) {
	var hal Hal

	h, err := hd.combineURL(HandlerPathDocument, "documentid", va.Document().DocumentId())
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
		height, documentid, err := parseOffset(offset)
		if err != nil {
			return nil, err
		}

		if reverse {
			filter["$or"] = []bson.M{
				{"height": bson.M{"$lt": height}},
				{"$and": []bson.M{
					{"height": height},
					{"documentid": bson.M{"$lt": documentid}},
				}},
			}
		} else {
			filter["$or"] = []bson.M{
				{"height": bson.M{"$gt": height}},
				{"$and": []bson.M{
					{"height": height},
					{"documentid": bson.M{"$gt": documentid}},
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

	documentid, err := strconv.ParseUint(offset, 10, 64)
	if err != nil {
		return nil, errors.Errorf("invalid index of offset: %q", err)
	}

	if reverse {
		filter = bson.M{
			"height":     height,
			"documentid": bson.M{"$lt": documentid},
		}
	} else {
		filter = bson.M{
			"height":     height,
			"documentid": bson.M{"$gt": documentid},
		}
	}

	return filter, nil
}

func nextOffsetOfDocuments(baseSelf string, vas []Hal, reverse bool) string {
	var nextoffset string
	if len(vas) > 0 {
		va := vas[len(vas)-1].Interface().(DocumentValue)
		nextoffset = buildOffsetByString(va.Height(), va.Document().DocumentId())
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
		nextoffset = fmt.Sprintf("%s", va.Document().DocumentId())
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

func (hd *Handlers) loadDocumentsHALFromDatabase(filter bson.M, reverse bool, limit int64) ([]Hal, error) {
	var vas []Hal

	/*
		if err := hd.database.BSDocuments(
			filter, reverse, limit,
			func(_ currency.Big, va BSDocumentValue) (bool, error) {
				hal, err := hd.buildBSDocumentHal(va)
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
	*/

	if err := hd.database.Documents(
		filter, reverse, limit,
		func(_ string, va DocumentValue) (bool, error) {
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
