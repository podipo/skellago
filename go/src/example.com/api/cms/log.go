package cms

import (
	"time"

	"github.com/coocood/qbs"
)

/*
Log is a log of Entry records (a "weblog", if you will)
*/
type Log struct {
	Id      int64  `json:"id" qbs:"pk"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	Tagline string `json:"tagline"`
	Publish bool   `json:"publish"`
	Image   string `json:"image"`
}

func (*Log) Indexes(indexes *qbs.Indexes) {
	indexes.AddUnique("slug")
}

/*
Entry is a post to a Log
*/
type Entry struct {
	Id      int64     `json:"id" qbs:"pk"`
	LogId   int64     `json:"log-id" qbs:"fk:Log"`
	Log     *Log      `json:"log"`
	Subject string    `json:"subject"`
	Slug    string    `json:"slug"`
	Content string    `json:"content"`
	Publish bool      `json:"publish"`
	Created time.Time `json:"created" qbs:"created"`
	Updated time.Time `json:"updated" qbs:"updated"`
	Issued  time.Time `json:"issued"`
	Image   string    `json:"image"`
}

/*
Tag is a metadata tag for an Entry
*/
type Tag struct {
	Id      int64  `json:"id" qbs:"pk"`
	Name    string `json:"name"`
	EntryId int64  `json:"entry-id" qbs:"fk:Entry"`
	Entry   *Entry `json:"entry"`
}

func (*Tag) Indexes(indexes *qbs.Indexes) {
	indexes.AddUnique("name", "entry_id")
}

func CreateLog(name string, slug string, db *qbs.Qbs) (*Log, error) {
	log := new(Log)
	log.Name = name
	log.Slug = slug
	_, err := db.Save(log)

	if err != nil {
		return nil, err
	}
	return log, nil
}

func UpdateLog(log *Log, db *qbs.Qbs) error {
	_, err := db.Save(log)
	if err != nil {
		return err
	}
	return nil
}

func FindLogs(offset int, limit int, q *qbs.Qbs) ([]*Log, error) {
	var logs []*Log
	err := q.Limit(limit).Offset(offset).FindAll(&logs)
	return logs, err
}

func FindPublicLogs(offset int, limit int, q *qbs.Qbs) ([]*Log, error) {
	var logs []*Log
	err := q.Limit(limit).Offset(offset).WhereEqual("publish", true).FindAll(&logs)
	return logs, err
}

func FindLogBySlug(slug string, db *qbs.Qbs) (*Log, error) {
	return findLogByField("slug", slug, db)
}

func findLogByField(fieldName string, value string, db *qbs.Qbs) (*Log, error) {
	record := new(Log)
	err := db.WhereEqual(fieldName, value).Find(record)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func CreateEntry(log *Log, subject string, slug string, content string, db *qbs.Qbs) (*Entry, error) {
	entry := new(Entry)
	entry.LogId = log.Id
	entry.Log = log
	entry.Subject = subject
	entry.Slug = slug
	entry.Content = content
	_, err := db.Save(entry)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func UpdateEntry(entry *Entry, db *qbs.Qbs) error {
	_, err := db.Save(entry)
	if err != nil {
		return err
	}
	return nil
}

func DeleteEntry(id int64, q *qbs.Qbs) (affected int64, err error) {
	record := new(Entry)
	record.Id = id
	return q.Delete(record)
}

func FindEntryBySlug(slug string, db *qbs.Qbs) (*Entry, error) {
	entry := new(Entry)
	err := db.WhereEqual("entry.slug", slug).Find(entry)
	return entry, err
}

func FindEntry(id int64, db *qbs.Qbs) (*Entry, error) {
	entry := new(Entry)
	err := db.WhereEqual("entry.id", id).Find(entry)
	return entry, err
}

func FindLogEntries(logId int64, offset int, limit int, db *qbs.Qbs) ([]*Entry, error) {
	var entries []*Entry
	err := db.Limit(limit).Offset(offset).WhereEqual("log_id", logId).FindAll(&entries)
	return entries, err
}

func FindPublicLogEntries(logId int64, offset int, limit int, db *qbs.Qbs) ([]*Entry, error) {
	var entries []*Entry
	err := db.Limit(limit).Offset(offset).WhereEqual("log_id", logId).WhereEqual("entry.publish", true).FindAll(&entries)
	return entries, err
}

func CreateTag(entry *Entry, name string, db *qbs.Qbs) (*Tag, error) {
	tag := new(Tag)
	tag.Name = name
	tag.EntryId = entry.Id
	tag.Entry = entry
	_, err := db.Save(tag)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

func DeleteTag(id int64, q *qbs.Qbs) (affected int64, err error) {
	record := new(Tag)
	record.Id = id
	return q.Delete(record)
}

/*
FindTags returns a list of Tag records for a given Entry
*/
func FindTags(q *qbs.Qbs) ([]*Tag, error) {
	var tags []*Tag
	err := q.FindAll(&tags)
	return tags, err
}

/*
FindEntryTags returns a list of Tag records for a given Entry
*/
func FindEntryTags(entryId int64, q *qbs.Qbs) ([]*Tag, error) {
	var tags []*Tag
	err := q.WhereEqual("entry_id", entryId).FindAll(&tags)
	return tags, err
}

/*
FindTaggedEntries returns a list of Entry records matching a tag's Name
*/
func FindTaggedEntries(name string, q *qbs.Qbs) ([]*Entry, error) {
	var records []*Tag
	err := q.WhereEqual("name", name).FindAll(&records)
	if err != nil {
		return nil, err
	}
	results := make([]*Entry, len(records))
	for i, record := range records {
		results[i] = record.Entry
	}
	return results, err
}

func DeleteAllTags(db *qbs.Qbs) error {
	_, err := db.Exec("delete from tag")
	return err
}

func DeleteAllEntries(db *qbs.Qbs) error {
	_, err := db.Exec("delete from entry")
	return err
}

func DeleteAllLogs(db *qbs.Qbs) error {
	_, err := db.Exec("delete from log")
	return err
}
