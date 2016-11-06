package cms

import (
	"fmt"
	"testing"

	. "github.com/chai2010/assert"
	"github.com/coocood/qbs"

	"podipo.com/skellago/be"
)

func TestLog(t *testing.T) {
	be.CreateAndInitDB()
	err := MigrateDB()
	AssertNil(t, err)

	db, err := qbs.GetQbs()
	AssertNil(t, err)
	defer func() {
		be.WipeDB()
		WipeDB()
		db.Close()
	}()

	slug1 := "blargh"
	log1, err := CreateLog("Blargh", slug1, db)
	AssertNil(t, err)
	AssertEqual(t, log1.Slug, slug1)

	log2, err := FindLogBySlug(slug1, db)
	AssertNil(t, err)
	AssertEqual(t, log1.Id, log2.Id)
	AssertEqual(t, slug1, log2.Slug)

	logs1, err := FindLogs(0, 100, db)
	AssertNil(t, err)
	AssertEqual(t, 1, len(logs1))
	logs1, err = FindPublicLogs(0, 100, db)
	AssertNil(t, err)
	AssertEqual(t, 0, len(logs1))

	log2.Name = "Blargh Blargh"
	log2.Publish = true
	err = UpdateLog(log2, db)
	AssertNil(t, err)

	log3, err := FindLogBySlug(slug1, db)
	AssertNil(t, err)
	AssertEqual(t, log1.Id, log3.Id)
	AssertEqual(t, log2.Name, log3.Name)
	Assert(t, log3.Publish)

	logs1, err = FindLogs(0, 100, db)
	AssertNil(t, err)
	AssertEqual(t, 1, len(logs1))
	logs1, err = FindPublicLogs(0, 100, db)
	AssertNil(t, err)
	AssertEqual(t, 1, len(logs1))

	entry1, err := CreateEntry(log3, "Title 1", "title-1", "Content 1", db)
	AssertNil(t, err)
	AssertEqual(t, "Title 1", entry1.Subject)

	entries, err := FindLogEntries(log3.Id, 0, 100, db)
	AssertNil(t, err)
	AssertEqual(t, 1, len(entries))

	entry2, err := CreateEntry(log3, "Title 2", "title-2", "Content 2", db)
	AssertNil(t, err)
	AssertEqual(t, entry2.Subject, "Title 2")
	AssertNotNil(t, entry2.Created)
	AssertNotNil(t, entry2.Updated)
	Assert(t, be.NilTime.Equal(entry2.Issued), fmt.Sprintf("Issued: %v", entry2.Issued))

	entries, err = FindLogEntries(log3.Id, 0, 100, db)
	AssertNil(t, err)
	AssertEqual(t, 2, len(entries))

	tags1, err := FindEntryTags(entry2.Id, db)
	AssertNil(t, err)
	AssertEqual(t, 0, len(tags1))

	tag1, err := CreateTag(entry2, "Taugh", db)
	AssertNil(t, err)
	AssertNotNil(t, tag1)

	tags1, err = FindEntryTags(entry1.Id, db)
	AssertNil(t, err)
	AssertEqual(t, 0, len(tags1))

	tags1, err = FindEntryTags(entry2.Id, db)
	AssertNil(t, err)
	AssertEqual(t, 1, len(tags1))
	AssertEqual(t, tag1.Name, tags1[0].Name)

	tag2, err := CreateTag(entry2, "Anuthur", db)
	AssertNil(t, err)
	AssertNotNil(t, tag1)

	tags1, err = FindEntryTags(entry2.Id, db)
	AssertNil(t, err)
	AssertEqual(t, 2, len(tags1))

	affected, err := DeleteTag(tag1.Id, db)
	AssertNil(t, err)
	Assert(t, 1 == affected)

	tags1, err = FindEntryTags(entry2.Id, db)
	AssertNil(t, err)
	AssertEqual(t, 1, len(tags1))
	AssertEqual(t, tag2.Name, tags1[0].Name)

	entries, err = FindTaggedEntries(tag2.Name, db)
	AssertNil(t, err)
	AssertEqual(t, 1, len(entries), fmt.Sprintf("Found entries: %v", entries))

	_, err = CreateTag(entry1, tag2.Name, db)
	AssertNil(t, err)
	_, err = CreateTag(entry1, tag2.Name, db)
	AssertNotNil(t, err, "Should not allow duplicate tag creation")

	tags1, err = FindEntryTags(entry1.Id, db)
	AssertNil(t, err)
	AssertEqual(t, 1, len(tags1))

	entries, err = FindTaggedEntries(tag2.Name, db)
	AssertNil(t, err)
	AssertEqual(t, 2, len(entries))
}
