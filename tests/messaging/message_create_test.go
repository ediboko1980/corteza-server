package messaging

import (
	"fmt"
	"net/http"
	"testing"

	jsonpath "github.com/steinfletcher/apitest-jsonpath"

	"github.com/cortezaproject/corteza-server/tests/helpers"
)

func TestMessagesCreate(t *testing.T) {
	h := newHelper(t)
	ch := h.repoMakePublicCh()

	rval := struct {
		Response struct {
			ID uint64 `json:"messageID,string"`
		}
	}{}

	h.apiInit().
		Post(fmt.Sprintf("/channels/%d/messages/", ch.ID)).
		JSON(`{"message":"new message"}`).
		Expect(t).
		Status(http.StatusOK).
		Assert(helpers.AssertNoErrors).
		Assert(jsonpath.Present(`$.response.messageID`)).
		Assert(jsonpath.Equal(`$.response.message`, `new message`)).
		End().
		JSON(&rval)

	m := h.repoMsgExistingLoad(rval.Response.ID)
	h.a.Equal(`new message`, m.Message)

}
