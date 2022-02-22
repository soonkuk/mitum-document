package digest

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/util"
	"go.mongodb.org/mongo-driver/bson"
)

func (hd *Handlers) handleDocuments(w http.ResponseWriter, r *http.Request) {
	limit := parseLimitQuery(r.URL.Query().Get("limit"))
	offset := parseOffsetQuery(r.URL.Query().Get("offset"))
	reverse := parseBoolQuery(r.URL.Query().Get("reverse"))
	doctype := parseStringQuery(r.URL.Query().Get("doctype"))

	cachekey := CacheKey(r.URL.Path, stringOffsetQuery(offset), stringBoolQuery("reverse", reverse), stringDoctypeQuery(doctype))

	if err := LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		i, filled, err := hd.handleDocumentsInGroup(offset, reverse, limit, doctype)

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
				expire = time.Minute
			}

			HTTP2WriteCache(w, cachekey, expire)
		}
	}
}

func (hd *Handlers) handleDocumentsInGroup(
	offset string,
	reverse bool,
	l int64,
	doctype string,
) ([]byte, bool, error) {
	var limit int64
	if l < 0 {
		limit = hd.itemsLimiter("documents")
	} else {
		limit = l
	}
	filter, err := buildDocumentsFilterByOffset(offset, reverse, doctype)
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
	if next := nextOffsetOfDocuments(h, vas, doctype, reverse); len(next) > 0 {
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
		hal = hal.AddLink("document:{documentid}", NewHalLink(HandlerPathDocument, nil).SetTemplated())
		hal = hal.AddLink("block:{height}", NewHalLink(HandlerPathBlockByHeight, nil).SetTemplated())

		return hd.enc.Marshal(hal)
	}
}

func (hd *Handlers) handleDocumentsByHeight(w http.ResponseWriter, r *http.Request) {
	limit := parseLimitQuery(r.URL.Query().Get("limit"))
	offset := parseOffsetQuery(r.URL.Query().Get("offset"))
	reverse := parseBoolQuery(r.URL.Query().Get("reverse"))
	doctype := parseStringQuery(r.URL.Query().Get("doctype"))

	cachekey := CacheKey(r.URL.Path, stringOffsetQuery(offset), stringBoolQuery("reverse", reverse), stringDoctypeQuery(doctype))

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
		i, filled, err := hd.handleDocumentsByHeightInGroup(height, offset, reverse, limit, doctype)
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
				expire = time.Minute
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
	doctype string,
) ([]byte, bool, error) {
	var limit int64
	if l < 0 {
		limit = hd.itemsLimiter("documents")
	} else {
		limit = l
	}
	filter, err := buildDocumentsByHeightFilter(height, reverse, doctype)
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
	if next := nextOffsetOfDocuments(h, vas, doctype, reverse); len(next) > 0 {
		hal = hal.AddLink("next", NewHalLink(next, nil))
	}

	b, err := hd.enc.Marshal(hal)
	return b, int64(len(vas)) == limit, err
}

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

func (hd *Handlers) buildDocumentsHal(
	baseSelf string,
	vas []Hal,
	offset string,
	reverse bool) Hal {
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

func nextOffsetOfDocuments(baseSelf string, vas []Hal, doctype string, reverse bool) string {
	var nextoffset, next string

	if len(vas) > 0 {
		va := vas[len(vas)-1].Interface().(DocumentValue)
		nextoffset = buildOffsetHeight(va.Height())
	}

	if len(nextoffset) > 0 {

		next = baseSelf
		next = addQueryValue(next, stringOffsetQuery(nextoffset))

		if len(doctype) > 0 {
			next = addQueryValue(next, stringDoctypeQuery(doctype))
		}

		if reverse {
			next = addQueryValue(next, stringBoolQuery("reverse", reverse))
		}
	}

	return next
}

func (hd *Handlers) loadDocumentsHALFromDatabase(filter bson.D, reverse bool, limit int64) ([]Hal, error) {
	var vas []Hal

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
