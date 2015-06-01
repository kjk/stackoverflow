package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/kjk/stackoverflow"
	"github.com/kjk/u"
)

func usageAndExit() {
	fmt.Printf("usage: tocvs file.xml\n")
	os.Exit(1)
}

type textWriter struct {
	f   *os.File
	pos int
}

func newTextWriter(path string) (*textWriter, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &textWriter{
		f: f,
	}, nil
}

func (w *textWriter) Write(s string) (int, int, error) {
	_, err := w.f.WriteString(s + "\n\n")
	if err != nil {
		return 0, 0, err
	}
	pos := w.pos
	w.pos += len(s) + 2
	return pos, len(s), nil
}

func (w *textWriter) Close() {
	if w.f != nil {
		w.f.Close()
		w.f = nil
	}
}

func userToCsvRecord(u *stackoverflow.User, w *textWriter, rec []string) error {
	aboutPos, aboutLen, err := w.Write(u.AboutMe)
	if err != nil {
		return err
	}
	about := fmt.Sprintf("%d-%d", aboutPos, aboutLen)
	rec[0] = strconv.Itoa(u.ID)
	rec[1] = strconv.Itoa(u.Reputation)
	rec[2] = u.CreationDate.Format(stackoverflow.TimeFormat)
	rec[3] = u.DisplayName
	rec[4] = u.LastAccessDate.Format(stackoverflow.TimeFormat)
	rec[5] = u.WebsiteURL
	rec[6] = u.Location
	rec[7] = about
	rec[8] = strconv.Itoa(u.Views)
	rec[9] = strconv.Itoa(u.UpVotes)
	rec[10] = strconv.Itoa(u.DownVotes)
	rec[11] = strconv.Itoa(u.Age)
	rec[12] = strconv.Itoa(u.AccountID)
	rec[13] = u.ProfileImageURL
	return nil
}

func usersToCsv(path string) error {
	timeStart := time.Now()
	n := 0
	r, err := stackoverflow.NewUsersReaderFromFile(path)
	if err != nil {
		return fmt.Errorf("usersToCsv: NewUsersReader() failed with %s\n", err)
	}

	dir := filepath.Dir(path)
	csvPath := filepath.Join(dir, "users.csv")
	f, err := os.Create(csvPath)
	if err != nil {
		return err
	}
	defer f.Close()
	textPath := filepath.Join(dir, "users.txt")
	textWriter, err := newTextWriter(textPath)
	if err != nil {
		return err
	}
	defer textWriter.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	var rec [14]string
	for r.Next() {
		u := &r.User
		err = userToCsvRecord(u, textWriter, rec[:])
		if err != nil {
			return err
		}
		w.Write(rec[:])
		n++
	}
	if r.Err() != nil {
		return r.Err()
	}
	fmt.Printf("converted %d users in %s\n", n, time.Since(timeStart))
	return nil
}

func main() {
	if len(os.Args) != 2 {
		usageAndExit()
	}
	path := u.ExpandTildeInPath(os.Args[1])
	if !u.PathExists(path) {
		fmt.Printf("file '%s' doesn't exist\n", path)
		usageAndExit()
	}
	ext := strings.ToLower(filepath.Ext(path))
	if ext != ".xml" {
		fmt.Printf("'%s' is not .xml file\n", path)
		usageAndExit()
	}
	name := strings.ToLower(filepath.Base(path))
	var err error
	switch name {
	case "users.xml":
		err = usersToCsv(path)
	default:
		err = fmt.Errorf("'%s' is not a recognized file name\n", path)
	}
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}
