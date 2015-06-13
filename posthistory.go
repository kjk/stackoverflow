package stackoverflow

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// http://meta.stackexchange.com/questions/2677/database-schema-documentation-for-the-public-data-dump-and-sede?rq=1
const (
	HistoryInitialTitle                 = 1
	HistoryInitialBody                  = 2
	HistoryInitialTags                  = 3
	HistoryEditTitle                    = 4
	HistoryEditBody                     = 5
	HistoyrEditTags                     = 6
	HistoryRollbackTitle                = 7
	HistoryRollbackBody                 = 8
	HistoryRollbackTags                 = 9
	HistoryPostClosed                   = 10
	HistoryPostReopened                 = 11
	HistoryPostDeleted                  = 12
	HistoryPostUndeleted                = 13
	HistoryPostLocked                   = 14
	HistoryPostUnlocked                 = 15
	HistoryCommunityOwned               = 16
	HistoryPostMigrated                 = 17
	HistoryQuestionMerged               = 18
	HistoryQuestionProtected            = 19
	HistoryQuestionUnprotected          = 20
	HistoryPostDisassociated            = 21
	HistoryQuestionUnmerged             = 22
	HistorySuggestedEditApplied         = 24
	HistoryPostTweeted                  = 25
	HistoryCommentDiscussionMovedToChat = 26
	HistoryPostNoticeAdded              = 33
	HistoryPostNoticeRemoved            = 34
	HistoryPostMigratedAway             = 35 // replaces id 17
	HistoryPostMigratedHere             = 36 // replaces id 17
	HistoryPostMergeSource              = 37
	HistoryPostMergeDestination         = 38
)

// PostHistory describes history of a post
type PostHistory struct {
	ID                int
	PostHistoryTypeID int
	PostID            int
	RevisionGUID      string
	CreationDate      time.Time
	UserID            int
	UserDisplayName   string
	// if PostHistoryTypeID is 10, 11, 12, 13, 14, 15, this is JSON
	// with users who voted
	Text string
	// if PostHistoryTypeID is HistoryInitialTags or HistoyrEditTags
	// or HistoryRollbackTags, this is a decoded version of tags
	Tags    []string
	Comment string
}

func decodePostHistoryAttr(attr xml.Attr, h *PostHistory) error {
	var err error
	name := strings.ToLower(attr.Name.Local)
	v := attr.Value
	switch name {
	case "id":
		h.ID, err = strconv.Atoi(v)
	case "posthistorytypeid":
		h.PostHistoryTypeID, err = strconv.Atoi(v)
	case "postid":
		h.PostID, err = strconv.Atoi(v)
	case "revisionguid":
		h.RevisionGUID = v
	case "creationdate":
		h.CreationDate, err = decodeTime(v)
	case "userid":
		h.UserID, err = strconv.Atoi(v)
	case "userdisplayname":
		h.UserDisplayName = v
	case "text":
		h.Text = v
	case "comment":
		h.Comment = v
	default:
		err = fmt.Errorf("unknown post history field: '%s'", name)
	}
	return err
}
func decodePostHistoryRow(t xml.Token, h *PostHistory) error {
	// have been checked before that this is "row" element

	// we reuse the struct, so reset to initial state
	h.RevisionGUID = ""
	h.CreationDate = time.Time{}
	h.UserID = 0
	h.UserDisplayName = ""
	h.Text = ""
	h.Tags = nil
	h.Comment = ""

	e, _ := t.(xml.StartElement)
	for _, attr := range e.Attr {
		err := decodePostHistoryAttr(attr, h)
		if err != nil {
			return err
		}
	}
	switch h.PostHistoryTypeID {
	case HistoryInitialTags, HistoyrEditTags, HistoryRollbackTags:
		if h.Text != "" {
			h.Tags = decodeTags(h.Text)
		}
	}
	return nil
}
