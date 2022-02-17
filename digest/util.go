package digest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/soonkuk/mitum-blocksign/document"
	"github.com/spikeekips/mitum-currency/currency"
	"github.com/spikeekips/mitum/base"
	"github.com/spikeekips/mitum/base/state"
	quicnetwork "github.com/spikeekips/mitum/network/quic"
	"github.com/spikeekips/mitum/util"
	"github.com/spikeekips/mitum/util/encoder"
	jsonenc "github.com/spikeekips/mitum/util/encoder/json"
	"github.com/spikeekips/mitum/util/valuehash"
)

func IsDocumentState(st state.State) (document.DocumentData, bool, error) {
	if !document.IsStateDocumentDataKey(st.Key()) {
		return nil, false, nil
	}

	doc, err := document.StateDocumentDataValue(st)
	if err != nil {
		return nil, false, err
	}
	return doc, true, nil
}

func parseDocIdFromPath(s string) (string, error) {
	s = strings.TrimSpace(s)
	if len(s) < 1 {
		return "", errors.Errorf("empty id")
	}

	//	h, err := document.ParseDocId(s)
	//	if err != nil {
	//		return "", err
	//	}

	return s, nil
}

func IsAccountState(st state.State) (currency.Account, bool, error) {
	if !currency.IsStateAccountKey(st.Key()) {
		return currency.Account{}, false, nil
	}

	ac, err := currency.LoadStateAccountValue(st)
	if err != nil {
		return currency.Account{}, false, err
	}
	return ac, true, nil
}

func IsBalanceState(st state.State) (currency.Amount, bool, error) {
	if !currency.IsStateBalanceKey(st.Key()) {
		return currency.Amount{}, false, nil
	}

	am, err := currency.StateBalanceValue(st)
	if err != nil {
		return currency.Amount{}, false, err
	}
	return am, true, nil
}

func parseHeightFromPath(s string) (base.Height, error) {
	s = strings.TrimSpace(s)

	if len(s) < 1 {
		return base.NilHeight, errors.Errorf("empty height")
	} else if len(s) > 1 && strings.HasPrefix(s, "0") {
		return base.NilHeight, errors.Errorf("invalid height, %q", s)
	}

	return base.NewHeightFromString(s)
}

func parseHashFromPath(s string) (valuehash.Hash, error) {
	s = strings.TrimSpace(s)
	if len(s) < 1 {
		return nil, errors.Errorf("empty hash")
	}

	h := valuehash.NewBytesFromString(s)
	if err := h.IsValid(nil); err != nil {
		return nil, err
	}

	return h, nil
}

func parseLimitQuery(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return int64(-1)
	}
	return n
}

func parseStringQuery(s string) string {
	return strings.TrimSpace(s)
}

func parseOffsetQuery(s string) string {
	return strings.TrimSpace(s)
}

func stringOffsetQuery(offset string) string {
	return fmt.Sprintf("offset=%s", offset)
}

func stringDocumentidQuery(documentid string) string {
	return fmt.Sprintf("documentid=%s", documentid)
}

func stringDoctypeQuery(doctype string) string {
	return fmt.Sprintf("doctype=%s", doctype)
}

func parseBoolQuery(s string) bool {
	return s == "1"
}

func stringBoolQuery(key string, v bool) string { // nolint:unparam
	if v {
		return fmt.Sprintf("%s=1", key)
	}

	return ""
}

func addQueryValue(b, s string) string {
	if len(s) < 1 {
		return b
	}

	if !strings.Contains(b, "?") {
		return b + "?" + s
	}

	return b + "&" + s
}

func HTTP2Stream(enc encoder.Encoder, w http.ResponseWriter, bufsize int, status int) (*jsoniter.Stream, func()) {
	w.Header().Set(HTTP2EncoderHintHeader, enc.Hint().String())
	w.Header().Set("Content-Type", HALMimetype)

	if status != http.StatusOK {
		w.WriteHeader(status)
	}

	stream := jsoniter.NewStream(HALJSONConfigDefault, w, bufsize)
	return stream, func() {
		_ = stream.Flush()
	}
}

func HTTP2NotSupported(w http.ResponseWriter, err error) {
	if err == nil {
		err = quicnetwork.NotSupportedErorr
	}

	HTTP2ProblemWithError(w, err, http.StatusInternalServerError)
}

func HTTP2ProblemWithError(w http.ResponseWriter, err error, status int) {
	HTTP2WritePoblem(w, NewProblemFromError(err), status)
}

func HTTP2WritePoblem(w http.ResponseWriter, pr Problem, status int) {
	if status == 0 {
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", ProblemMimetype)
	w.Header().Set("X-Content-Type-Options", "nosniff")

	var output []byte
	if b, err := jsonenc.Marshal(pr); err != nil {
		output = unknownProblemJSON
	} else {
		output = b
	}

	w.WriteHeader(status)
	_, _ = w.Write(output)
}

func HTTP2WriteHal(enc encoder.Encoder, w http.ResponseWriter, hal Hal, status int) { // nolint:unparam
	stream, flush := HTTP2Stream(enc, w, 1, status)
	defer flush()

	stream.WriteVal(hal)
}

func HTTP2WriteHalBytes(enc encoder.Encoder, w http.ResponseWriter, b []byte, status int) { // nolint:unparam
	w.Header().Set(HTTP2EncoderHintHeader, enc.Hint().String())
	w.Header().Set("Content-Type", HALMimetype)

	if status != http.StatusOK {
		w.WriteHeader(status)
	}

	_, _ = w.Write(b)
}

func HTTP2WriteCache(w http.ResponseWriter, key string, expire time.Duration) {
	if cw, ok := w.(*CacheResponseWriter); ok {
		_ = cw.SetKey(key).SetExpire(expire)
	}
}

func HTTP2HandleError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, util.NotFoundError):
		status = http.StatusNotFound
	case errors.Is(err, quicnetwork.BadRequestError):
		status = http.StatusBadRequest
	case errors.Is(err, quicnetwork.NotSupportedErorr):
		status = http.StatusInternalServerError
	}

	HTTP2ProblemWithError(w, err, status)
}
