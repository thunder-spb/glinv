package main

import (
	"path/filepath"
	"text/template"
	"time"

	"glinv/pkg/forms"
	"glinv/pkg/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
// At the moment it only contains one field, but we'll add more
// to it as the build progresses.
type templateData struct {
	AuthenticatedUser *models.User

	CSRFToken string

	CurrentYear int

	Version string

	Notice string

	Tree *Tree

	Form *forms.Form

	Environment string

	CountUnapproved int

	InventoryHost   *models.InventoryHost
	EditHost        *models.InventoryHost
	EditHosts       []*models.InventoryHost
	InventoryHosts  []*models.InventoryHost
	HostsForApprove []*models.InventoryHost
	HostsForDelete  []*models.InventoryHost
	CountInvHosts   int
	Hostname        string

	InventoryHVar     *models.InventoryHVar
	InventoryHVars    []*models.InventoryHVar
	EditHVars         []*models.InventoryHVar
	CountInvHostsVars int

	InventoryHTag     *models.InventoryHTag
	InventoryHTags    []*models.InventoryHTag
	EditHTags         []*models.InventoryHTag
	CountInvHostsTags int

	InventoryGroup  *models.InventoryGroup
	InventoryGroups []*models.InventoryGroup
	EditGroup       *models.InventoryGroup
	CountInvGroups  int
	SubGroups       bool

	InventoryGVar      *models.InventoryGVar
	InventoryGVars     []*models.InventoryGVar
	EditGVars          []*models.InventoryGVar
	CountInvGroupsVars int

	InventoryService   *models.InventoryService
	EditService        *models.InventoryService
	InventoryServices  []*models.InventoryService
	ServicesForApprove []*models.InventoryService
	ServicesForDelete  []*models.InventoryService
	CountInvServices   int
	Statuses           map[int]int
	Stages             []*models.InventoryService

	ServerAgent                *models.ServerAgent   //TODO
	ServersAgent               []*models.ServerAgent //TODO
	CountServersAgent          int
	TotalCPU                   int
	TotalRAM, UsedRAM, FreeRAM string
	TotalHDD, UsedHDD, FreeHDD string
	Alert                      bool

	SysCtl         []*models.SysCtlServer
	PackagesServer []*models.PackageServer //TODO

	BaseTemplate       *models.BaseTemplate
	BaseTemplates      []*models.BaseTemplate
	CountBaseTemplates int

	BaseTemplateItem  *models.BaseTemplateItem
	BaseTemplatesItem []*models.BaseTemplateItem

	History    *models.History
	HistoryAll []*models.History
}

// Create a humanDate function which returns a nicely formatted string
// representation of a time.Time object.
func formatDate(t time.Time) string {
	localTimezone, _ := time.LoadLocation("Asia/Novosibirsk")

	return t.In(localTimezone).Format("02.01.2006 at 15:04")
}

// Initialize a template.FuncMap object and store it in a global variable. This is
// essentially a string-keyed map which acts as a lookup between the names of our
// custom template functions and the functions themselves.
var functions = template.FuncMap{
	"formatDate": formatDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}

	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the application.
	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	// Loop through the pages one-by-one.
	for _, page := range pages {
		// Extract the file name (like 'index.page.tmpl') from the full file path
		// and assign it to the name variable.
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.
		// ts, err := template.ParseFiles(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'layout' templates to the
		// template set (in our case, it's just the 'base' layout at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'chunk' templates to the
		// template set (in our case, it's just the 'footer' partial at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.chunk.tmpl"))
		if err != nil {
			return nil, err
		}

		// Add the template set to the cache, using the name of the page
		// (like 'index.page.tmpl') as the key.
		cache[name] = ts
	}

	// Return the map.
	return cache, nil
}
