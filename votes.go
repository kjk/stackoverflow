package stackoverflow

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// http://blog.stackoverflow.com/2009/06/stack-overflow-creative-commons-data-dump/#comment-24147
const (
	VoteAcceptedByOriginator = 1
	VoteUpMod                = 2
	VoteDownMod              = 3
	VoteOffensive            = 4
	VoteFavorite             = 5
	VoteClose                = 6
	VoteReopen               = 7
	VoteBountyStart          = 8
	VoteBountyClose          = 9
	VoteDeletion             = 10
	VoteUndeletion           = 11
	VoteSpam                 = 12
	VoteInformModerator      = 13
)

// Vote describes a vote
type Vote struct {
	ID           int
	PostID       int
	VoteTypeID   int
	UserID       int
	BountyAmount int
	CreationDate time.Time
}

func decodeVoteAttr(attr xml.Attr, vote *Vote) error {
	var err error
	name := strings.ToLower(attr.Name.Local)
	v := attr.Value
	switch name {
	case "id":
		vote.ID, err = strconv.Atoi(v)
	case "postid":
		vote.PostID, err = strconv.Atoi(v)
	case "userid":
		vote.UserID, err = strconv.Atoi(v)
	case "votetypeid":
		vote.VoteTypeID, err = strconv.Atoi(v)
	case "bountyamount":
		vote.BountyAmount, err = strconv.Atoi(v)
	case "creationdate":
		vote.CreationDate, err = decodeTime(v)
	default:
		err = fmt.Errorf("unknown vote field: '%s'", name)
	}
	return err
}

func decodeVoteRow(t xml.Token, vote *Vote) error {
	// have been checked before that this is "row" element
	e, _ := t.(xml.StartElement)
	for _, attr := range e.Attr {
		err := decodeVoteAttr(attr, vote)
		if err != nil {
			return err
		}
	}
	return nil
}
