package cms

import (
	"encoding/json"
	"net/http"

	"podipo.com/skellago/be"
)

var EntryProperties = []be.Property{
	be.Property{
		Name:        "id",
		Description: "A unique id number",
		DataType:    "int",
	},
	be.Property{
		Name:        "slug",
		Description: "A unique, url fiendly string",
		DataType:    "string",
	},
	be.Property{
		Name:        "subject",
		Description: "The title",
		DataType:    "string",
	},
	be.Property{
		Name:        "content",
		Description: "The body",
		DataType:    "string",
	},
	be.Property{
		Name:        "image",
		Description: "The main image",
		DataType:    "file",
	},
}

var LogProperties = []be.Property{
	be.Property{
		Name:        "id",
		Description: "A unique id number",
		DataType:    "int",
	},
	be.Property{
		Name:        "name",
		Description: "The name of the log",
		DataType:    "string",
	},
	be.Property{
		Name:        "slug",
		Description: "A unique slug for use in URLs",
		DataType:    "string",
	},
	be.Property{
		Name:        "tagline",
		Description: "A short tagline describing the log's contents.",
		DataType:    "string",
		Optional:    true,
	},
	be.Property{
		Name:        "publish",
		Description: "True if the log should be available to the public",
		DataType:    "bool",
	},
	be.Property{
		Name:        "image",
		Description: "An image associated with the log",
		DataType:    "file",
		Optional:    true,
	},
}

var LogsProperties = be.NewAPIListProperties("log")
var LogEntriesProperties = be.NewAPIListProperties("entry")

type LogsResource struct {
}

func NewLogsResource() *LogsResource {
	return &LogsResource{}
}

func (LogsResource) Name() string  { return "logs" }
func (LogsResource) Path() string  { return "/log/" }
func (LogsResource) Title() string { return "Logs" }
func (LogsResource) Description() string {
	return "A list of logs."
}

func (resource LogsResource) Properties() []be.Property {
	return LogsProperties
}

func (resource LogsResource) Get(request *be.APIRequest) (int, interface{}, http.Header) {
	responseHeader := map[string][]string{}
	offset, limit := be.GetOffsetAndLimit(request.Raw.Form)
	var logs []*Log
	var err error
	if request.User.Staff {
		logs, err = FindLogs(offset, limit, request.DB)
	} else {
		logs, err = FindPublicLogs(offset, limit, request.DB)
	}
	if err != nil {
		return 500, be.APIError{
			Id:      "db_error",
			Message: "Database error",
			Error:   err.Error(),
		}, responseHeader
	}
	list := &be.APIList{
		Offset:  offset,
		Limit:   limit,
		Objects: logs,
	}
	return 200, list, responseHeader
}

func (resource LogsResource) Post(request *be.APIRequest) (int, interface{}, http.Header) {
	responseHeader := map[string][]string{}
	if request.User == nil {
		return 401, be.NotLoggedInError, responseHeader
	}
	if request.User.Staff == false {
		return 403, be.ForbiddenError, responseHeader
	}
	log := new(Log)
	err := json.NewDecoder(request.Raw.Body).Decode(&log)
	if err != nil {
		return 400, be.JSONParseError, responseHeader
	}
	newLog, err := CreateLog(log.Name, log.Slug, request.DB)
	if err != nil {
		return 404, be.APIError{
			Id:      "log_creation_error",
			Message: "Could not create the log",
			Error:   err.Error(),
		}, responseHeader
	}
	newLog.Tagline = log.Tagline
	newLog.Publish = log.Publish
	UpdateLog(newLog, request.DB)
	return 200, newLog, responseHeader
}

type LogResource struct {
}

func NewLogResource() *LogResource {
	return &LogResource{}
}

func (LogResource) Name() string  { return "log" }
func (LogResource) Path() string  { return "/log/{slug:[0-9,a-z,-]+}" }
func (LogResource) Title() string { return "Log" }
func (LogResource) Description() string {
	return "A log (some would say blog) contains a series of entries."
}

func (resource LogResource) Properties() []be.Property {
	return LogProperties
}

func (resource LogResource) Get(request *be.APIRequest) (int, interface{}, http.Header) {
	responseHeader := map[string][]string{}
	slug, _ := request.PathValues["slug"]
	log, err := FindLogBySlug(slug, request.DB)
	if err != nil {
		return 404, be.APIError{
			Id:      "no_such_log",
			Message: "No such log: " + slug,
			Error:   err.Error(),
		}, responseHeader
	}
	// Don't show the log if the it isn't published and this isn't a staff request
	if log.Publish == false {
		if request.User == nil {
			return 404, be.NotLoggedInError, responseHeader
		}
		if request.User.Staff == false {
			return 403, be.ForbiddenError, responseHeader
		}
	}
	return 200, log, responseHeader
}

// Put updates the Log
func (resource LogResource) Put(request *be.APIRequest) (int, interface{}, http.Header) {
	responseHeader := map[string][]string{}
	if request.User == nil {
		return 401, be.NotLoggedInError, responseHeader
	}
	if request.User.Staff == false {
		return 403, be.ForbiddenError, responseHeader
	}
	slug, _ := request.PathValues["slug"]
	log, err := FindLogBySlug(slug, request.DB)
	if err != nil {
		return 404, be.APIError{
			Id:      "no_such_log",
			Message: "No such log: " + slug,
			Error:   err.Error(),
		}, responseHeader
	}
	logUpdate := new(Log)
	err = json.NewDecoder(request.Raw.Body).Decode(&logUpdate)
	if err != nil {
		return 400, be.JSONParseError, responseHeader
	}
	log.Name = logUpdate.Name
	log.Publish = logUpdate.Publish
	log.Slug = logUpdate.Slug
	log.Tagline = logUpdate.Tagline
	err = UpdateLog(log, request.DB)
	if err != nil {
		return 404, be.APIError{
			Id:      "bad_log_update",
			Message: "Could not update the log",
			Error:   err.Error(),
		}, responseHeader
	}
	return 200, log, responseHeader
}

// Post creates an Entry record
func (resource LogResource) Post(request *be.APIRequest) (int, interface{}, http.Header) {
	responseHeader := map[string][]string{}
	if request.User == nil {
		return 401, be.NotLoggedInError, responseHeader
	}
	if request.User.Staff == false {
		return 403, be.ForbiddenError, responseHeader
	}

	slug, _ := request.PathValues["slug"]
	log, err := FindLogBySlug(slug, request.DB)
	if err != nil {
		return 404, be.APIError{
			Id:      "no_such_log",
			Message: "No such log: " + slug,
			Error:   err.Error(),
		}, responseHeader
	}

	newEntry := new(Entry)
	err = json.NewDecoder(request.Raw.Body).Decode(&newEntry)
	if err != nil {
		return 400, be.JSONParseError, responseHeader
	}

	entry, err := CreateEntry(log, newEntry.Subject, newEntry.Slug, newEntry.Content, request.DB)
	if err != nil {
		return 404, be.APIError{
			Id:      "entry_creation_error",
			Message: "Could not create that entry",
			Error:   err.Error(),
		}, responseHeader
	}
	entry.Issued = newEntry.Issued
	entry.Publish = newEntry.Publish
	err = UpdateEntry(entry, request.DB)
	if err != nil {
		return 404, be.APIError{
			Id:      "entry_update_error",
			Message: "Could not update that entry",
			Error:   err.Error(),
		}, responseHeader
	}

	return 200, entry, responseHeader
}

type LogEntriesResource struct {
}

func NewLogEntriesResource() *LogEntriesResource {
	return &LogEntriesResource{}
}

func (LogEntriesResource) Name() string  { return "log-entries" }
func (LogEntriesResource) Path() string  { return "/log/{slug:[0-9,a-z,-]+}/entries" }
func (LogEntriesResource) Title() string { return "Log entries" }
func (LogEntriesResource) Description() string {
	return "A list of entries in a log."
}

func (resource LogEntriesResource) Properties() []be.Property {
	return LogEntriesProperties
}

func (resource LogEntriesResource) Get(request *be.APIRequest) (int, interface{}, http.Header) {
	responseHeader := map[string][]string{}

	slug, _ := request.PathValues["slug"]
	log, err := FindLogBySlug(slug, request.DB)
	if err != nil {
		return 404, be.APIError{
			Id:      "no_such_log",
			Message: "No such log: " + slug,
			Error:   err.Error(),
		}, responseHeader
	}

	if log.Publish == false {
		if request.User == nil {
			return 401, be.NotLoggedInError, responseHeader
		}
		if request.User.Staff == false {
			return 403, be.ForbiddenError, responseHeader
		}
	}

	offset, limit := be.GetOffsetAndLimit(request.Raw.Form)
	var entries []*Entry
	if request.User.Staff {
		entries, err = FindLogEntries(log.Id, offset, limit, request.DB)
	} else {
		entries, err = FindPublicLogEntries(log.Id, offset, limit, request.DB)
	}
	if err != nil {
		return 500, be.APIError{
			Id:      "db_error",
			Message: "Database error",
			Error:   err.Error(),
		}, responseHeader
	}
	list := &be.APIList{
		Offset:  offset,
		Limit:   limit,
		Objects: entries,
	}
	return 200, list, responseHeader
}

type EntryResource struct {
}

func NewEntryResource() *EntryResource {
	return &EntryResource{}
}

func (EntryResource) Name() string  { return "entry" }
func (EntryResource) Path() string  { return "/entry/{slug:[0-9,a-z,-]+}" }
func (EntryResource) Title() string { return "Entry" }
func (EntryResource) Description() string {
	return "An entry (aka post) in a log (aka blog)"
}

func (resource EntryResource) Properties() []be.Property {
	return EntryProperties
}

func (resource EntryResource) Get(request *be.APIRequest) (int, interface{}, http.Header) {
	responseHeader := map[string][]string{}
	slug, _ := request.PathValues["slug"]
	entry, err := FindEntryBySlug(slug, request.DB)
	if err != nil {
		return 404, be.APIError{
			Id:      "no_such_entry",
			Message: "No such entry: " + slug,
			Error:   err.Error(),
		}, responseHeader
	}
	// Don't show the entry if the it isn't published and this isn't a staff request
	if entry.Log.Publish == false || entry.Publish == false {
		if request.User == nil {
			return 404, be.NotLoggedInError, responseHeader
		}
		if request.User.Staff == false {
			return 403, be.ForbiddenError, responseHeader
		}
	}
	return 200, entry, responseHeader
}
