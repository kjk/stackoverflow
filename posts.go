package stackoverflow

import (
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// http://meta.stackexchange.com/questions/2677/database-schema-documentation-for-the-public-data-dump-and-sede?rq=1
const (
	PostQuestion            = 1
	PostAnswer              = 2
	PostOrphanedTagWiki     = 3
	PostTagWikiExcerpt      = 4
	PostTagWiki             = 5
	PostModeratorNomination = 6
	PostWikiPlaceholder     = 7
	PostPrivilegeWiki       = 8
)

// Post describes a post
type Post struct {
	ID                    int
	PostTypeID            int
	ParentID              int // for PostAnswer
	AcceptedAnswerID      int
	CreationDate          time.Time
	Score                 int
	ViewCount             int
	Body                  string
	OwnerUserID           int
	OwnerDisplayName      string
	LastEditorUserID      int
	LastEditorDisplayName string
	LastEditDate          time.Time
	LastActivitityDate    time.Time
	Title                 string
	Tags                  []string
	AnswerCount           int
	CommentCount          int
	FavoriteCount         int
	CommunityOwnedDate    time.Time
	ClosedDate            time.Time
}

var nTagsToShow = 0

func decodeTags(s string) ([]string, error) {
	// tags are in the format: <foo><bar>
	s = strings.TrimPrefix(s, "<")
	s = strings.TrimSuffix(s, ">")
	tags := strings.Split(s, "><")
	if nTagsToShow > 0 {
		nTagsToShow--
		fmt.Printf("tags: '%s' => %v\n", s, tags)
	}
	return tags, nil
}

func decodePostAttr(attr xml.Attr, p *Post) error {
	var err error
	name := strings.ToLower(attr.Name.Local)
	v := attr.Value
	switch name {
	case "id":
		p.ID, err = strconv.Atoi(v)
	case "parentid":
		p.ParentID, err = strconv.Atoi(v)
	case "posttypeid":
		p.PostTypeID, err = strconv.Atoi(v)
	case "acceptedanswerid":
		p.AcceptedAnswerID, err = strconv.Atoi(v)
	case "creationdate":
		p.CreationDate, err = decodeTime(v)
	case "score":
		p.Score, err = strconv.Atoi(v)
	case "viewcount":
		p.ViewCount, err = strconv.Atoi(v)
	case "body":
		p.Body = v
	case "owneruserid":
		p.OwnerUserID, err = strconv.Atoi(v)
	case "ownerdisplayname":
		p.OwnerDisplayName = v
	case "lasteditoruserid":
		p.LastEditorUserID, err = strconv.Atoi(v)
	case "lasteditordisplayname":
		p.LastEditorDisplayName = v
	case "lasteditdate":
		p.LastEditDate, err = decodeTime(v)
	case "lastactivitydate":
		p.LastActivitityDate, err = decodeTime(v)
	case "title":
		p.Title = v
	case "tags":
		p.Tags, err = decodeTags(v)
	case "answercount":
		p.AnswerCount, err = strconv.Atoi(v)
	case "commentcount":
		p.CommentCount, err = strconv.Atoi(v)
	case "favoritecount":
		p.FavoriteCount, err = strconv.Atoi(v)
	case "communityowneddate":
		p.CommunityOwnedDate, err = decodeTime(v)
	case "closeddate":
		p.ClosedDate, err = decodeTime(v)
	default:
		err = fmt.Errorf("unknown post field: '%s'", name)
	}
	return err
}

func validatePost(p *Post) {
	if p.PostTypeID < 1 || p.PostTypeID > 8 {
		log.Fatalf("invalid PostTypeID: %d\n", p.PostTypeID)
	}
}

func decodePostRow(t xml.Token, p *Post) error {
	// have been checked before that this is "row" element
	e, _ := t.(xml.StartElement)
	for _, attr := range e.Attr {
		err := decodePostAttr(attr, p)
		if err != nil {
			return err
		}
	}
	validatePost(p)
	return nil
}
