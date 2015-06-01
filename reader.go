package stackoverflow

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

const (
	// TimeFormat is how time is formatted in .xml files
	TimeFormat = "2006-01-02T15:04:05.999999999"

	typeBadges      = "badges"
	typeComments    = "comments"
	typePostHistory = "posthistory"
	typePostLinks   = "postlinks"
	typePosts       = "posts"
	typeTags        = "tags"
	typeUsers       = "users"
	typeVotes       = "votes"
)

// Reader is for iteratively reading records from xml file
type Reader struct {
	r           io.Reader
	d           *xml.Decoder
	typ         string
	User        User
	Post        Post
	Comment     Comment
	Badge       Badge
	Tag         Tag
	PostHistory PostHistory
	PostLink    PostLink
	Vote        Vote
	err         error
	finished    bool
}

// NewBadgesReaderFromFile returns a new reader for Badges.xml file
func NewBadgesReaderFromFile(path string) (*Reader, error) {
	return newReaderFromFile(path, typeBadges)
}

// NewCommentsReaderFromFile returns a new reader for Comments.xml file
func NewCommentsReaderFromFile(path string) (*Reader, error) {
	return newReaderFromFile(path, typeComments)
}

// NewPostHistoryReaderFromFile returns a new reader for PostHistory.xml file
func NewPostHistoryReaderFromFile(path string) (*Reader, error) {
	return newReaderFromFile(path, typePostHistory)
}

// NewPostLinksReaderFromFile returns a new reader for PostLinks.xml file
func NewPostLinksReaderFromFile(path string) (*Reader, error) {
	return newReaderFromFile(path, typePostLinks)
}

// NewPostsReaderFromFile returns a new reader for Posts.xml file
func NewPostsReaderFromFile(path string) (*Reader, error) {
	return newReaderFromFile(path, typePosts)
}

// NewTagsReaderFromFile returns a new reader for Comments.xml file
func NewTagsReaderFromFile(path string) (*Reader, error) {
	return newReaderFromFile(path, typeTags)
}

// NewUsersReaderFromFile returns a new reader for Users.xml file
func NewUsersReaderFromFile(path string) (*Reader, error) {
	return newReaderFromFile(path, typeUsers)
}

// NewVotesReaderFromFile returns a new reader for Votes.xml file
func NewVotesReaderFromFile(path string) (*Reader, error) {
	return newReaderFromFile(path, typeVotes)
}

// NewBadgesReader returns a new reader for Badges.xml file
func NewBadgesReader(r io.Reader) (*Reader, error) {
	return newReader(r, typeBadges)
}

// NewCommentsReader returns a new reader for Comments.xml file
func NewCommentsReader(r io.Reader) (*Reader, error) {
	return newReader(r, typeComments)
}

// NewPostHistoryReader returns a new reader for PostHistory.xml file
func NewPostHistoryReader(r io.Reader) (*Reader, error) {
	return newReader(r, typePostHistory)
}

// NewPostLinksReader returns a new reader for PostLinks.xml file
func NewPostLinksReader(r io.Reader) (*Reader, error) {
	return newReader(r, typePostLinks)
}

// NewPostsReader returns a new reader for Posts.xml file
func NewPostsReader(r io.Reader) (*Reader, error) {
	return newReader(r, typePosts)
}

// NewTagsReader returns a new reader for Comments.xml file
func NewTagsReader(r io.Reader) (*Reader, error) {
	return newReader(r, typeTags)
}

// NewUsersReader returns a new reader for Users.xml file
func NewUsersReader(r io.Reader) (*Reader, error) {
	return newReader(r, typeUsers)
}

// NewVotesReader returns a new reader for Votes.xml file
func NewVotesReader(r io.Reader) (*Reader, error) {
	return newReader(r, typeVotes)
}

func isCharData(t xml.Token) bool {
	_, ok := t.(xml.CharData)
	return ok
}

func isProcInst(t xml.Token) bool {
	_, ok := t.(xml.ProcInst)
	return ok
}

func isStartElement(t xml.Token, name string) bool {
	e, ok := t.(xml.StartElement)
	if !ok {
		return false
	}
	return strings.EqualFold(e.Name.Local, name)
}

func isEndElement(t xml.Token, name string) bool {
	e, ok := t.(xml.EndElement)
	if !ok {
		return false
	}
	return strings.EqualFold(e.Name.Local, name)
}

func getTokenIgnoreCharData(d *xml.Decoder) (xml.Token, error) {
	t, err := d.Token()
	if err != nil {
		return nil, err
	}
	if !isCharData(t) {
		return t, nil
	}
	return d.Token()
}

func decodeTime(s string) (time.Time, error) {
	return time.Parse(TimeFormat, s)
}

func newReader(rd io.Reader, typ string) (*Reader, error) {
	r := &Reader{
		r:   rd,
		d:   xml.NewDecoder(rd),
		typ: typ,
	}
	t, err := getTokenIgnoreCharData(r.d)
	if err != nil {
		r.Close()
		return nil, err
	}
	// skip <?xml ...>
	if isProcInst(t) {
		t, err = getTokenIgnoreCharData(r.d)
		if err != nil {
			r.Close()
			return nil, err
		}
	}
	if !isStartElement(t, r.typ) {
		r.Close()
		return nil, fmt.Errorf("invalid first token '%#v', expected xml.StartElement '%s'", t, r.typ)
	}
	return r, nil
}

func newReaderFromFile(path string, typ string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return newReader(f, typ)
}

// Err returns potential error
func (r *Reader) Err() error {
	return r.err
}

// Close closes a reader
func (r *Reader) Close() {
	if !r.finished && r.r != nil {
		if rc, ok := r.r.(io.ReadCloser); ok {
			rc.Close()
		}
		r.r = nil
		r.finished = true
	}
}

// Next advances to next User record. Returns false on end or
func (r *Reader) Next() bool {
	if r.err != nil || r.finished {
		return false
	}

	defer func() {
		if r.err != nil {
			r.Close()
		}
	}()

	// skip newlines between eleemnts
	t, err := getTokenIgnoreCharData(r.d)
	if err != nil {
		r.err = err
		return false
	}

	if isEndElement(t, "row") {
		t, r.err = getTokenIgnoreCharData(r.d)
		if r.err != nil {
			return false
		}
	}

	if isEndElement(t, r.typ) {
		r.Close()
		return false
	}

	if !isStartElement(t, "row") {
		r.err = fmt.Errorf("unexpected token: %#v, wanted xml.StartElement 'row'", t)
		return false
	}
	switch r.typ {
	case typeBadges:
		r.err = decodeBadgeRow(t, &r.Badge)
	case typeComments:
		r.err = decodeCommentRow(t, &r.Comment)
	case typePosts:
		r.err = decodePostRow(t, &r.Post)
	case typePostHistory:
		r.err = decodePostHistoryRow(t, &r.PostHistory)
	case typePostLinks:
		r.err = decodePostLinkRow(t, &r.PostLink)
	case typeTags:
		r.err = decodeTagRow(t, &r.Tag)
	case typeUsers:
		r.err = decodeUserRow(t, &r.User)
	case typeVotes:
		r.err = decodeVoteRow(t, &r.Vote)
	}
	if r.err != nil {
		return false
	}
	return true
}
