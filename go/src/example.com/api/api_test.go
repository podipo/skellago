package main

import (
	"os"
	"strconv"
	"testing"

	"example.com/api/cms"
	. "github.com/chai2010/assert"
	"github.com/coocood/qbs"
	"podipo.com/skellago/be"
)

func TestLogAPI(t *testing.T) {
	be.CreateAndInitDB()
	err := cms.MigrateDB()
	AssertNil(t, err)

	db, err := qbs.GetQbs()
	AssertNil(t, err)
	logger.Print("Deferring TestLogAPI cleanup")
	defer func() {
		be.WipeDB()
		cms.WipeDB()
		db.Close()
		logger.Print("Cleaned up TestLogAPI")
	}()

	testApi, err := NewTestAPI()
	AssertNil(t, err)
	defer testApi.Stop()

	userClient, staffClient, err := be.CreateTestUserAndStaffWithClients(testApi, db)

	list, err := userClient.GetList("/log/")
	AssertNil(t, err)
	AssertNil(t, list.Objects)

	list, err = staffClient.GetList("/log/")
	AssertNil(t, err)
	AssertNil(t, list.Objects)

	log1 := cms.Log{
		Name:    "Blargh",
		Slug:    "blargh",
		Tagline: "Blargh it!",
	}
	log2 := new(cms.Log)
	err = userClient.PostAndReceiveJSON("/log/", &log1, log2)
	AssertNotNil(t, err, "Users should not be able to create logs")
	log2 = new(cms.Log)
	err = staffClient.PostAndReceiveJSON("/log/", &log1, log2)
	AssertNil(t, err)
	err = staffClient.PostAndReceiveJSON("/log/", &log1, log2)
	AssertNotNil(t, err, "Should not allow duplicate slugs")

	list, err = userClient.GetList("/log/")
	AssertNil(t, err)
	AssertNil(t, list.Objects, "User should still see no logs since none are published")

	list, err = staffClient.GetList("/log/")
	AssertNil(t, err)
	AssertNotNil(t, list.Objects, "Staff should now see one log (which is not published)")
	objs := list.Objects.([]interface{})
	AssertEqual(t, 1, len(objs))

	log3 := new(cms.Log)
	err = userClient.GetJSON("/log/"+strconv.FormatInt(log2.Id, 10), log3)
	AssertNotNil(t, err, "User should not see non-published log")
	err = staffClient.GetJSON("/log/"+strconv.FormatInt(log2.Id, 10), log3)
	AssertNil(t, err)
	AssertLogsEqual(t, log2, log3)

	log3.Name = "Neu Blargh"
	log4 := new(cms.Log)
	err = userClient.PutAndReceiveJSON("/log/"+strconv.FormatInt(log3.Id, 10), &log3, log4)
	AssertNotNil(t, err, "Users should not be able to update logs")
	log4 = new(cms.Log)
	err = staffClient.PutAndReceiveJSON("/log/"+strconv.FormatInt(log3.Id, 10), &log3, log4)
	AssertNil(t, err)
	AssertEqual(t, log3.Name, log4.Name)

	log5 := new(cms.Log)
	err = staffClient.GetJSON("/log/"+strconv.FormatInt(log3.Id, 10), log5)
	AssertNil(t, err)
	AssertLogsEqual(t, log3, log5)

	log5.Publish = true
	err = staffClient.PutAndReceiveJSON("/log/"+strconv.FormatInt(log5.Id, 10), &log5, log5)
	AssertNil(t, err, "Users should not be able to update logs")

	list, err = userClient.GetList("/log/")
	AssertNil(t, err)
	AssertNotNil(t, list.Objects, "User should now see the log since it's published")
	objs = list.Objects.([]interface{})
	AssertEqual(t, 1, len(objs))

	list, err = userClient.GetList("/log/" + strconv.FormatInt(log5.Id, 10) + "/entries")
	AssertNil(t, err)
	AssertNil(t, list.Objects, "There should be no entries in this log")

	entry1 := cms.Entry{
		LogId:   log5.Id,
		Subject: "Furst Pohst",
		Slug:    "furst-post",
		Content: "Loohk Ohut, Heah Ih Cohm",
		Publish: false,
	}
	entry2 := new(cms.Entry)
	err = userClient.PostAndReceiveJSON("/log/"+strconv.FormatInt(log5.Id, 10)+"/entries", &entry1, entry2)
	AssertNotNil(t, err)
	err = staffClient.PostAndReceiveJSON("/log/"+strconv.FormatInt(log5.Id, 10)+"/entries", &entry1, entry2)
	AssertNil(t, err)

	imageUrl := "/entry/" + strconv.FormatInt(entry2.Id, 10) + "/image"
	_, err = staffClient.GetFile(imageUrl)
	AssertNotNil(t, err)
	imageFile1, err := be.TempImage(os.TempDir(), 1000, 1000)
	AssertNil(t, err)
	response, err := staffClient.SendFile("PUT", imageUrl, "image", imageFile1)
	AssertNil(t, err)
	AssertEqual(t, 200, response.StatusCode)
	reader, err := staffClient.GetFile(imageUrl)
	AssertNil(t, err)
	AssertNotNil(t, reader)

	list, err = userClient.GetList("/log/" + strconv.FormatInt(log5.Id, 10) + "/entries")
	AssertNil(t, err)
	AssertNil(t, list.Objects, "There should be no entries in this log")

	list, err = staffClient.GetList("/log/" + strconv.FormatInt(log5.Id, 10) + "/entries")
	AssertNil(t, err)
	AssertNotNil(t, list.Objects, "Staff should see the entry")
	objs = list.Objects.([]interface{})
	AssertEqual(t, 1, len(objs))

	entry3 := cms.Entry{
		LogId:   log5.Id,
		Subject: "Secund Pohst",
		Slug:    "secund-post",
		Content: "Loohk Ohut, Heah Ih Cohm Agahn",
		Publish: true,
	}
	entry4 := new(cms.Entry)
	err = staffClient.PostAndReceiveJSON("/log/"+strconv.FormatInt(log5.Id, 10)+"/entries", &entry3, entry4)
	AssertNil(t, err)

	list, err = userClient.GetList("/log/" + strconv.FormatInt(log5.Id, 10) + "/entries")
	AssertNil(t, err)
	AssertNotNil(t, list.Objects, "User should see the entry")
	objs = list.Objects.([]interface{})
	AssertEqual(t, 1, len(objs), "User should see one entry")

	list, err = staffClient.GetList("/log/" + strconv.FormatInt(log5.Id, 10) + "/entries")
	AssertNil(t, err)
	AssertNotNil(t, list.Objects, "Staff should see both entries")
	objs = list.Objects.([]interface{})
	AssertEqual(t, 2, len(objs))

	entry5 := new(cms.Entry)
	err = staffClient.GetJSON("/entry/"+strconv.FormatInt(entry4.Id, 10), entry5)
	AssertNil(t, err, "Could not fetch an entry by id")

	entry2.Publish = true
	entry6 := new(cms.Entry)
	err = staffClient.PutAndReceiveJSON("/entry/"+strconv.FormatInt(entry2.Id, 10), entry2, entry6)
	AssertNil(t, err)
	Assert(t, entry6.Publish, "Did not receive a published entry: %v", entry6)

	entry5 = new(cms.Entry)
	err = staffClient.GetJSON("/entry/"+strconv.FormatInt(entry6.Id, 10), entry5)
	AssertNil(t, err)
	AssertEqual(t, true, entry5.Publish, "After setting, Publish should be set: %v", entry5)

	list, err = userClient.GetList("/log/" + strconv.FormatInt(log5.Id, 10) + "/entries")
	AssertNil(t, err)
	AssertNotNil(t, list.Objects, "User should see the entry")
	objs = list.Objects.([]interface{})
	AssertEqual(t, 2, len(objs), "User should see both entries")

	list, err = staffClient.GetList("/log/" + strconv.FormatInt(log5.Id, 10) + "/entries")
	AssertNil(t, err)
	AssertNotNil(t, list.Objects, "Staff should see both entries")
	objs = list.Objects.([]interface{})
	AssertEqual(t, 2, len(objs))

	err = staffClient.Delete("/entry/" + strconv.FormatInt(entry5.Id, 10))
	AssertNil(t, err)

	list, err = staffClient.GetList("/log/" + strconv.FormatInt(log5.Id, 10) + "/entries")
	AssertNil(t, err)
	AssertNotNil(t, list.Objects, "Staff should see both entries")
	objs = list.Objects.([]interface{})
	AssertEqual(t, 1, len(objs), "One of the entries should have been deleted")

	log6 := cms.Log{
		Name:    "Binky",
		Slug:    "binky",
		Tagline: "Binky it!",
		Publish: true,
	}
	log7 := new(cms.Log)
	err = staffClient.PostAndReceiveJSON("/log/", &log6, log7)
	AssertNil(t, err)
	// Make sure that we're not seeing the other log's entries
	list, err = staffClient.GetList("/log/" + strconv.FormatInt(log7.Id, 10) + "/entries")
	AssertNil(t, err)
	AssertNil(t, list.Objects, "Staff should see no entries for a new log")

	// Make sure that we're not seeing the other log's entries
	list, err = userClient.GetList("/log/" + strconv.FormatInt(log7.Id, 10) + "/entries")
	AssertNil(t, err)
	AssertNil(t, list.Objects, "Users should see no entries for a new log ", list.Objects)
}

func AssertLogsEqual(t *testing.T, log1 *cms.Log, log2 *cms.Log) {
	AssertEqual(t, log1.Id, log2.Id)
	AssertEqual(t, log1.Name, log2.Name)
	AssertEqual(t, log1.Slug, log2.Slug)
	AssertEqual(t, log1.Tagline, log2.Tagline)
	AssertEqual(t, log1.Publish, log2.Publish)
}
