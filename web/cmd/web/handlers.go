package main

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"

	"glinv/pkg/forms"
	"glinv/pkg/models"
)

func (app *application) overview(w http.ResponseWriter, r *http.Request) {

	countInvGroups := app.inventoryGroup.GetCountGroups()
	countInvGroupsVars := app.inventoryGroup.GetCountGroupsVars()

	countInvHosts := app.inventoryHost.GetCountHosts()
	countInvHostsVars := app.inventoryHost.GetCountHostsVars()
	countInvHostsTags := app.inventoryHost.GetCountHostsTags()

	countInvServices := app.inventoryService.GetCountServices()

	countServersAgent := app.serverAgent.GetCountServers()
	countBaseTemplates := app.serverAgent.GetCountBaseTemplates()

	// Use the new render helper.
	app.render(w, r, "overview.page.tmpl", &templateData{
		CountInvGroups:     countInvGroups,
		CountInvGroupsVars: countInvGroupsVars,

		CountInvHosts:     countInvHosts,
		CountInvHostsVars: countInvHostsVars,
		CountInvHostsTags: countInvHostsTags,

		CountInvServices: countInvServices,

		CountServersAgent:  countServersAgent,
		CountBaseTemplates: countBaseTemplates,
	})
}

//
// Hosts
//

// Show
func (app *application) showHosts(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get environment for filtration
	environment := "all" // Init
	for _, value := range r.Form["env"] {
		environment = value
	}

	inventoryHosts, err := app.inventoryHost.GetHostsByEnv(environment)
	if err != nil {
		app.serverError(w, err)
		return
	}

	statuses, err := app.inventoryHost.GetStatuses()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "hosts.page.tmpl", &templateData{
		InventoryHosts: inventoryHosts,
		Statuses:       statuses,
		Environment:    environment,
	})
}

func (app *application) showHost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}
	// Use the HostModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	inventoryHost, invnetoryGroups, inventoryHVars, inventoryHTags, err := app.inventoryHost.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	countServices := app.inventoryService.GetCountServicesByHost(inventoryHost.Hostname)

	// Use the new render helper.
	app.render(w, r, "host.page.tmpl", &templateData{
		InventoryHost:    inventoryHost,
		InventoryGroups:  invnetoryGroups,
		InventoryHVars:   inventoryHVars,
		InventoryHTags:   inventoryHTags,
		CountInvServices: countServices,
	})
}

// Create
func (app *application) createHostForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get environment for filtration
	environment := r.Form["e"][0]

	// GroupName
	inventoryGroups, err := app.inventoryGroup.GetGroupsByEnv(environment)
	if err != nil {
		app.serverError(w, err)
		return
	}

	groups := Tree{}
	tree := groups.TreeGroups(inventoryGroups, nil, nil)

	inventoryHVars, err := app.inventoryHVar.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	htags, err := app.inventoryHTag.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "createHost.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form:            forms.New(nil),
		InventoryGroups: inventoryGroups,
		InventoryHVars:  inventoryHVars,
		InventoryHTags:  htags,
		Environment:     environment,
		Tree:            tree,
	})
}

func (app *application) createHost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("hostname", "ip", "environment")
	form.MaxLength("hostname", 100)
	form.MaxLength("ip", 15)

	// Groups checkbox
	checkGroups := make(map[string]string)
	if len(r.PostForm["check_groups"]) != 0 {
		for _, groups := range r.PostForm["check_groups"] {
			g := strings.Split(groups, " ")
			checkGroups[g[0]] = g[1]
		}
	} else {
		checkGroups["1"] = "all" // if no group is selected
	}

	// Vars checkbox
	checkVars := make(map[string]string)
	for _, vars := range r.PostForm["hvars"] {
		v := strings.Split(vars, " ")
		checkVars[v[0]] = v[1]
	}

	// Tags checkbox
	checkTags := make(map[string]string)
	for _, tags := range r.PostForm["htags"] {
		v := strings.Split(tags, " ")
		checkTags[v[0]] = v[1]
	}

	// Group
	inventoryGroups, err := app.inventoryGroup.GetGroupsByEnv(form.Get("environment"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	groups := Tree{}
	tree := groups.TreeGroups(inventoryGroups, nil, nil)

	// HVars
	inventoryHVars, err := app.inventoryHVar.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// HTags
	htags, err := app.inventoryHTag.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	if !form.Valid() {
		app.render(w, r, "createHost.page.tmpl", &templateData{
			Form:            form,
			InventoryGroups: inventoryGroups,
			InventoryHVars:  inventoryHVars,
			InventoryHTags:  htags,
			Environment:     form.Get("environment"),
			Tree:            tree,
		})

		return
	}

	id, err := app.inventoryHost.Insert(form.Get("hostname"), form.Get("ip"), form.Get("environment"), checkGroups, checkVars, checkTags)
	if err == models.ErrDuplicateEmail {
		form.Errors.Set("hostname", "Hostname is already in use")
		app.render(w, r, "createHost.page.tmpl", &templateData{
			Form: form,
		})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "notice", "Submitted for Approval!")

	http.Redirect(w, r, fmt.Sprintf("/hosts/%d", id), http.StatusSeeOther)
}

// Edit
func (app *application) editHostForm(w http.ResponseWriter, r *http.Request) {
	// Get the host ID
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	editHost, editGroups, editHVars, editHTags, err := app.inventoryHost.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	environment := editHost.Environment

	// GroupName
	inventoryGroups, err := app.inventoryGroup.GetGroupsByEnv(environment)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Creating a list of marked groups for Tree
	var checkedIDGroups []int
	for _, g := range editGroups {
		checkedIDGroups = append(checkedIDGroups, g.ID)
	}

	groups := Tree{}
	tree := groups.TreeGroups(inventoryGroups, nil, checkedIDGroups)

	inventoryHVars, err := app.inventoryHVar.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	inventoryHTags, err := app.inventoryHTag.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "editHost.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form:        forms.New(nil),
		Environment: environment,

		EditHost:  editHost,
		EditHVars: editHVars,
		EditHTags: editHTags,

		InventoryGroups: inventoryGroups,
		InventoryHVars:  inventoryHVars,
		InventoryHTags:  inventoryHTags,

		Tree: tree,
	})
}

func (app *application) editHost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("hostname", "ip")
	form.MaxLength("hostname", 100)
	form.MaxLength("ip", 15)

	// Groups checkbox
	checkGroups := make(map[string]string)
	if len(r.PostForm["check_groups"]) != 0 {
		for _, groups := range r.PostForm["check_groups"] {
			g := strings.Split(groups, " ")
			checkGroups[g[0]] = g[1]
		}
	} else {
		checkGroups["1"] = "all" // if no group is selected
	}

	// Vars checkbox
	checkVars := make(map[string]string)
	for _, vars := range r.PostForm["hvars"] {
		v := strings.Split(vars, " ")
		checkVars[v[0]] = v[1]
	}

	// Tags checkbox
	checkTags := make(map[string]string)
	for _, tags := range r.PostForm["htags"] {
		v := strings.Split(tags, " ")
		checkTags[v[0]] = v[1]
	}

	// Group
	inventoryGroups, err := app.inventoryGroup.GetGroupsByEnv(form.Get("environment"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	groups := Tree{}
	tree := groups.TreeGroups(inventoryGroups, nil, nil)

	// HVars
	inventoryHVars, err := app.inventoryHVar.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// HTags
	inventoryHTags, err := app.inventoryHTag.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	if !form.Valid() {
		app.render(w, r, "editHost.page.tmpl", &templateData{
			Form:            form,
			InventoryGroups: inventoryGroups,
			InventoryHVars:  inventoryHVars,
			InventoryHTags:  inventoryHTags,
			Environment:     form.Get("environment"),
			Tree:            tree,
		})

		return
	}

	id, err := app.inventoryHost.Update(form.Get("id"), form.Get("hostname"), form.Get("ip"), form.Get("description"), checkGroups, checkVars, checkTags)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// History
	var g, v, t string
	for _, value := range checkGroups {
		g += value + " "
	}
	for _, value := range checkVars {
		v += value + " "
	}
	for _, value := range checkTags {
		t += value + " "
	}
	event := fmt.Sprintf("Hostname: <a href='/hosts/edit/%v'>%s</a> \n IP: %s \n Desc.: %v \n Groups: %s  Vars: %s Tags: %s", id, form.Get("hostname"), form.Get("ip"), form.Get("description"), g, v, t)
	app.history.Event(form.Get("userID"), form.Get("userEmail"), form.Get("id"), "host", "edit", event)

	app.session.Put(r, "notice", "Changes have been submitted for approval!")

	http.Redirect(w, r, fmt.Sprintf("/hosts/%d", id), http.StatusSeeOther)
}

// Delete
func (app *application) deleteHost(w http.ResponseWriter, r *http.Request) {
	// Get the host ID
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	inventoryHost, _, _, _, err := app.inventoryHost.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Host For Deletion
	err = app.inventoryHost.MarkHostForDeletion(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	msg := fmt.Sprintf("Host:<strong>%s</strong> IP:<strong>%s</strong> was successfully submitted for deletion approval!", inventoryHost.Hostname, inventoryHost.IP)
	app.session.Put(r, "notice", msg)

	http.Redirect(w, r, fmt.Sprintf("/hosts"), http.StatusSeeOther)
}

// Other
func (app *application) showAPIHosts(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "apiHosts.page.tmpl", &templateData{})
}

func (app *application) showSSHConfig(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "sshConfig.page.tmpl", &templateData{})
}

//
// HVars
//
func (app *application) showHVars(w http.ResponseWriter, r *http.Request) {
	inventoryHVars, err := app.inventoryHVar.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "hvars.page.tmpl", &templateData{
		InventoryHVars: inventoryHVars,
		Form:           forms.New(nil),
	})
}

func (app *application) showHVarForm(w http.ResponseWriter, r *http.Request) {
	// Get the hvar ID
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	inventoryHVar, err := app.inventoryHVar.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	inventoryHosts, err := app.inventoryHost.GetHostsByIDVar(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "hvar.page.tmpl", &templateData{
		InventoryHVar:  inventoryHVar,
		InventoryHosts: inventoryHosts,
		Form:           forms.New(nil),
	})
}

func (app *application) editHVar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, _ := strconv.Atoi(form.Get("idHVar"))

	inventoryHVar, err := app.inventoryHVar.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	inventoryHosts, err := app.inventoryHost.GetHostsByIDVar(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form.Required("valueHVar")
	form.MaxLength("valueHVar", 100)

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "hvar.page.tmpl", &templateData{
			Form:           form,
			InventoryHVar:  inventoryHVar,
			InventoryHosts: inventoryHosts,
		})

		return
	}

	// Update ...
	err = app.inventoryHVar.Update(form.Get("idHVar"), form.Get("valueHVar"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	msg := fmt.Sprintf("The var: <strong>%s</strong>:<strong>%s</strong> change was saved successfully!", form.Get("nameHVar"), form.Get("valueHVar"))
	app.session.Put(r, "notice", msg)

	http.Redirect(w, r, fmt.Sprintf("/hosts/var/%v", form.Get("idHVar")), http.StatusSeeOther)
}

func (app *application) deleteHVar(w http.ResponseWriter, r *http.Request) {
	// Get the hvar ID
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	inventoryHVar, err := app.inventoryHVar.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.inventoryHVar.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	msg := fmt.Sprintf("Tag <strong>%s</strong>:<strong>%s</strong> was successfully deleted!", inventoryHVar.Name, inventoryHVar.Value)
	app.session.Put(r, "notice", msg)

	http.Redirect(w, r, "/hosts/vars", http.StatusSeeOther)
}

func (app *application) createHVar(w http.ResponseWriter, r *http.Request) {
	inventoryHVars, err := app.inventoryHVar.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("value")
	form.MaxLength("value", 100)

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "hvars.page.tmpl", &templateData{
			InventoryHVars: inventoryHVars,
			Form:           form,
		})

		return
	}

	// Because the form data (with type url.Values) has been anonymously embedded
	// in the form.Form struct, can use the Get() method to retrieve
	// the validated value for a particular form field.
	_, err = app.inventoryHVar.Insert(form.Get("type"), form.Get("value"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	value := fmt.Sprintf("New var <strong>%s</strong> added successfully!", form.Get("value"))
	app.session.Put(r, "notice", value)

	http.Redirect(w, r, "/hosts/vars", http.StatusSeeOther)
}

//
// HTags
//
func (app *application) showHTags(w http.ResponseWriter, r *http.Request) {
	htags, err := app.inventoryHTag.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "htags.page.tmpl", &templateData{
		InventoryHTags: htags,
		Form:           forms.New(nil),
	})
}

func (app *application) showHTagForm(w http.ResponseWriter, r *http.Request) {
	// Get the htag ID
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	inventoryHTag, err := app.inventoryHTag.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	inventoryHosts, err := app.inventoryHost.GetHostsByIDTag(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "htag.page.tmpl", &templateData{
		InventoryHTag:  inventoryHTag,
		InventoryHosts: inventoryHosts,
		Form:           forms.New(nil),
	})
}

func (app *application) editHTag(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, _ := strconv.Atoi(form.Get("idHTag"))

	inventoryHTag, err := app.inventoryHTag.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	inventoryHosts, err := app.inventoryHost.GetHostsByIDTag(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form.Required("valueHTag")
	form.MaxLength("valueHTag", 100)

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "htag.page.tmpl", &templateData{
			Form:           form,
			InventoryHTag:  inventoryHTag,
			InventoryHosts: inventoryHosts,
		})

		return
	}

	// Update
	err = app.inventoryHTag.Update(form.Get("idHTag"), form.Get("valueHTag"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	msg := fmt.Sprintf("The tag: <strong>%s</strong> change was saved successfully!", form.Get("valueHTag"))
	app.session.Put(r, "notice", msg)

	http.Redirect(w, r, fmt.Sprintf("/hosts/tag/%v", form.Get("idHTag")), http.StatusSeeOther)
}

func (app *application) deleteHTag(w http.ResponseWriter, r *http.Request) {
	// Get the htag ID
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	inventoryHTag, err := app.inventoryHTag.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.inventoryHTag.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	msg := fmt.Sprintf("Tag <strong>%s</strong> was successfully deleted!", inventoryHTag.Value)
	app.session.Put(r, "notice", msg)

	http.Redirect(w, r, "/hosts/tags", http.StatusSeeOther)
}

func (app *application) createHTag(w http.ResponseWriter, r *http.Request) {
	htags, err := app.inventoryHTag.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("value")
	form.MaxLength("value", 100)

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "htags.page.tmpl", &templateData{
			InventoryHTags: htags,
			Form:           form,
		})

		return
	}

	// Because the form data (with type url.Values) has been anonymously embedded
	// in the form.Form struct, can use the Get() method to retrieve
	// the validated value for a particular form field.
	_, err = app.inventoryHTag.Insert(form.Get("value"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Set("value", "Tag is already in use")
		app.render(w, r, "htags.page.tmpl", &templateData{
			InventoryHTags: htags,
			Form:           form,
		})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	msg := fmt.Sprintf("Tag <strong>%s</strong> added successfully!", form.Get("value"))
	app.session.Put(r, "notice", msg)

	http.Redirect(w, r, "/hosts/tags", http.StatusSeeOther)
}

//
// Groups
//
func (app *application) showInventoryGroups(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get environment for filtration
	environment := "prd" // Init
	for _, value := range r.Form["env"] {
		environment = value
	}

	inventoryGroups, err := app.inventoryGroup.GetGroupsByEnv(environment)
	if err != nil {
		app.serverError(w, err)
		return
	}

	treeGVars, _ := app.inventoryGVar.TreeGVars()

	groups := Tree{}
	tree := groups.TreeGroups(inventoryGroups, treeGVars, nil)

	// Use the new render helper.
	app.render(w, r, "groups.page.tmpl", &templateData{
		Tree:        tree,
		Environment: environment,
	})
}

func (app *application) editInventoryGroupForm(w http.ResponseWriter, r *http.Request) {
	// Get the group ID
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	// Get value of group
	editGroup, err := app.inventoryGroup.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Check subgroups
	subgroups := app.inventoryGroup.CheckSubGroups(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	inventoryGVars, err := app.inventoryGVar.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	editGVars, err := app.inventoryGVar.GetGVarsByGroupID(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	inventoryHosts, err := app.inventoryHost.GetHostsByIDGroup(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "group.page.tmpl", &templateData{
		EditGroup:      editGroup,
		InventoryGVars: inventoryGVars,
		EditGVars:      editGVars,
		SubGroups:      subgroups,
		InventoryHosts: inventoryHosts,
		Form:           forms.New(nil),
	})
}

func (app *application) editInventoryGroup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, _ := strconv.Atoi(form.Get("idGroup"))

	// Get value of group
	editGroup, err := app.inventoryGroup.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Check subgroups
	subgroups := app.inventoryGroup.CheckSubGroups(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	inventoryGVars, err := app.inventoryGVar.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	editGVars, err := app.inventoryGVar.GetGVarsByGroupID(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// GVars checkbox
	gvars := make(map[string]string)
	for _, vars := range r.PostForm["gvars"] {
		v := strings.Split(vars, " ")
		gvars[v[0]] = v[1]
	}

	inventoryHosts, err := app.inventoryHost.GetHostsByIDGroup(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form.Required("valueGroup")
	form.MaxLength("valueGroup", 100)

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "group.page.tmpl", &templateData{
			Form:           form,
			EditGroup:      editGroup,
			InventoryGVars: inventoryGVars,
			EditGVars:      editGVars,
			SubGroups:      subgroups,
			InventoryHosts: inventoryHosts,
		})

		return
	}

	err = app.inventoryGroup.Update(form.Get("idGroup"), form.Get("valueGroup"), gvars)
	if err != nil {
		app.serverError(w, err)
		return
	}

	msg := fmt.Sprintf("The group: <strong>%s</strong> change was saved successfully!", form.Get("valueGroup"))
	app.session.Put(r, "notice", msg)

	http.Redirect(w, r, fmt.Sprintf("/group/%v", form.Get("idGroup")), http.StatusSeeOther)
}

func (app *application) deleteInventoryGroup(w http.ResponseWriter, r *http.Request) {
	// Get the group ID
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	// Get value of group
	group, err := app.inventoryGroup.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Delete Group
	err = app.inventoryGroup.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	msg := fmt.Sprintf("The group: <strong>%v</strong> was delete successfully!", group.Value)
	app.session.Put(r, "notice", msg)

	http.Redirect(w, r, fmt.Sprintf("/groups"), http.StatusSeeOther)
}

func (app *application) createInventoryGroupForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get environment for filtration
	var environment string
	for _, value := range r.Form["e"] {
		environment = value
	}

	inventoryGroups, err := app.inventoryGroup.GetGroupsByEnv(environment)
	if err != nil {
		app.serverError(w, err)
		return
	}

	inventoryGVars, err := app.inventoryGVar.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "createGroup.page.tmpl", &templateData{
		InventoryGroups: inventoryGroups,
		InventoryGVars:  inventoryGVars,
		Environment:     environment,
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) createInventoryGroup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("newgroup")
	form.MaxLength("newgroup", 100)

	// GVars checkbox
	gvars := make(map[string]string)
	for _, vars := range r.PostForm["gvars"] {
		v := strings.Split(vars, " ")
		gvars[v[0]] = v[1]
	}

	inventoryGroups, err := app.inventoryGroup.GetGroupsByEnv(form.Get("environment"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	inventoryGVars, err := app.inventoryGVar.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "createGroup.page.tmpl", &templateData{
			InventoryGroups: inventoryGroups,
			InventoryGVars:  inventoryGVars,
			Environment:     form.Get("environment"),
			Form:            form,
		})

		return
	}

	// Because the form data (with type url.Values) has been anonymously embedded
	// in the form.Form struct, can use the Get() method to retrieve
	// the validated value for a particular form field.
	_, err = app.inventoryGroup.Insert(form.Get("environment"), form.Get("parent"), form.Get("newgroup"), form.Get("description"), gvars)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "notice", "Submitted for Approval!")

	http.Redirect(w, r, "/groups", http.StatusSeeOther)
}

func (app *application) showAPIGroups(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "apiGroups.page.tmpl", &templateData{})
}

//
// GVars
//
func (app *application) showGVars(w http.ResponseWriter, r *http.Request) {
	inventoryGVars, err := app.inventoryGVar.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "gvars.page.tmpl", &templateData{
		InventoryGVars: inventoryGVars,
		Form:           forms.New(nil),
	})
}

func (app *application) showGVarForm(w http.ResponseWriter, r *http.Request) {
	// Get the gvar ID
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	inventoryGVar, err := app.inventoryGVar.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	inventoryGroups, err := app.inventoryGroup.GetGroupsByID(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "gvar.page.tmpl", &templateData{
		InventoryGVar:   inventoryGVar,
		InventoryGroups: inventoryGroups,
		Form:            forms.New(nil),
	})
}

func (app *application) editGVar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	id, _ := strconv.Atoi(form.Get("idGVar"))

	inventoryGVar, err := app.inventoryGVar.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	inventoryGroups, err := app.inventoryGroup.GetGroupsByID(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	form.Required("valueGVar")
	form.MaxLength("valueGVar", 100)

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "gvar.page.tmpl", &templateData{
			Form:            form,
			InventoryGVar:   inventoryGVar,
			InventoryGroups: inventoryGroups,
		})

		return
	}

	err = app.inventoryGVar.Update(form.Get("idGVar"), form.Get("valueGVar"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	msg := fmt.Sprintf("The var: <strong>%s</strong>:<strong>%s</strong> change was saved successfully!", form.Get("nameGVar"), form.Get("valueGVar"))
	app.session.Put(r, "notice", msg)

	http.Redirect(w, r, fmt.Sprintf("/groups/var/%v", form.Get("idGVar")), http.StatusSeeOther)
}

func (app *application) deleteGVar(w http.ResponseWriter, r *http.Request) {
	// Get the gvar ID
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	inventoryGVar, err := app.inventoryGVar.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.inventoryGVar.Delete(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	msg := fmt.Sprintf("Tag <strong>%s</strong>:<strong>%s</strong> was successfully deleted!", inventoryGVar.Name, inventoryGVar.Value)
	app.session.Put(r, "notice", msg)

	http.Redirect(w, r, "/groups/vars", http.StatusSeeOther)
}

func (app *application) createGVars(w http.ResponseWriter, r *http.Request) {
	inventoryGVars, err := app.inventoryGVar.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("value")
	form.MaxLength("value", 100)

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "gvars.page.tmpl", &templateData{
			InventoryGVars: inventoryGVars,
			Form:           form,
		})

		return
	}

	// Because the form data (with type url.Values) has been anonymously embedded
	// in the form.Form struct, can use the Get() method to retrieve
	// the validated value for a particular form field.
	_, err = app.inventoryGVar.Insert(form.Get("type"), form.Get("value"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	value := fmt.Sprintf("New var %s added successfully!", form.Get("value"))
	app.session.Put(r, "notice", value)

	http.Redirect(w, r, "/groups/vars", http.StatusSeeOther)
}

//
// Services
//
func (app *application) showServices(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get environment for filtration
	environment := "all" // Init
	for _, value := range r.Form["env"] {
		environment = value
	}

	inventoryServices, err := app.inventoryService.GetServicesByEnv(environment)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "services.page.tmpl", &templateData{
		InventoryServices: inventoryServices,
		Environment:       environment,
	})
}

func (app *application) showService(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	inventoryService, err := app.inventoryService.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Pass the Notice message to the template.
	app.render(w, r, "service.page.tmpl", &templateData{
		InventoryService: inventoryService,
	})
}

func (app *application) showServicesHost(w http.ResponseWriter, r *http.Request) {
	hostname := r.URL.Query().Get(":hostname")
	fmt.Println(hostname)

	inventoryServices, err := app.inventoryService.GetServicesByHostname(hostname)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "servicesHost.page.tmpl", &templateData{
		Hostname:          hostname,
		InventoryServices: inventoryServices,
	})
}

func (app *application) editServiceForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	inventoryHosts, err := app.inventoryHost.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	editHosts, err := app.inventoryService.GetHostsService(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	editService, err := app.inventoryService.Get(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "editService.page.tmpl", &templateData{
		EditService:    editService,
		EditHosts:      editHosts,
		InventoryHosts: inventoryHosts,
		Form:           forms.New(nil),
	})
}

func (app *application) editService(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("nodes", "techname", "placement", "port")
	form.MaxLength("techname", 100)
	form.PermittedPorts("port", 0, 65535)

	idService, err := strconv.Atoi(form.Get("id"))
	if err != nil || idService < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	inventoryHosts, err := app.inventoryHost.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	editHosts, err := app.inventoryService.GetHostsService(idService)
	if err != nil {
		app.serverError(w, err)
		return
	}

	editService, err := app.inventoryService.Get(idService)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Nodes checkbox
	var nodes []string
	var idHost int
	if len(r.PostForm["nodes"]) != 0 {
		for _, id := range r.PostForm["nodes"] {
			nodes = append(nodes, id)
		}

		idHost, err = strconv.Atoi(nodes[0])
		if err != nil || idHost < 1 {
			app.notFound(w) // Use the notFound() helper.
			return
		}

	}

	if len(nodes) > 1 || len(nodes) == 0 {
		app.session.Put(r, "notice", "<div class='alert alert-danger notice' role='alert'>In edit mode, only one host must be selected, changes are not saved!</div>")
	} else {

		var reg, method string

		if len(form.Get("reg_to_consul")) != 0 {
			reg = "true"
			method = form.Get("method_check_consul")
		} else {
			reg = "false"
		}

		location := form.Get("location")
		kind := form.Get("type")

		title := form.Get("title")
		techname := form.Get("techname")
		domain := form.Get("domain")
		value := form.Get("value")
		placement := form.Get("placement")

		team := form.Get("team")
		resp := form.Get("resp")

		port := form.Get("port")
		description := r.PostForm.Get("description")

		// If the form isn't valid, redisplay the template passing in the
		// form.Form object as the data.
		if !form.Valid() {
			app.render(w, r, "editService.page.tmpl", &templateData{
				InventoryHosts: inventoryHosts,
				EditService:    editService,
				EditHosts:      editHosts,
				Form:           form,
			})

			return
		}
		// Because the form data (with type url.Values) has been anonymously embedded
		// in the form.Form struct, can use the Get() method to retrieve
		// the validated value for a particular form field.
		err = app.inventoryService.Update(kind, techname, domain, placement, team, resp, title, value, location, reg, method, description, port, idHost, idService)
		if err != nil {
			app.serverError(w, err)
			return
		}

		app.session.Put(r, "notice", "<div class='alert alert-warning notice' role='alert'>Submitted for Approval!</div>")
	}

	http.Redirect(w, r, fmt.Sprintf("/services/%d", idService), http.StatusSeeOther)

}

func (app *application) deleteService(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	inventoryService, err := app.inventoryService.Get(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Service For Deletion
	err = app.inventoryService.MarkServiceForDeletion(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	msg := fmt.Sprintf("Service <strong>%s</strong>: <strong>%s</strong> was successfully submitted for deletion approval!", inventoryService.Title, inventoryService.Value)
	app.session.Put(r, "notice", msg)

	http.Redirect(w, r, "/services", http.StatusSeeOther)
}

func (app *application) createServiceForm(w http.ResponseWriter, r *http.Request) {
	inventoryHosts, err := app.inventoryHost.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "createService.page.tmpl", &templateData{
		InventoryHosts: inventoryHosts,
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
	})
}

func (app *application) createService(w http.ResponseWriter, r *http.Request) {
	inventoryHosts, err := app.inventoryHost.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("nodes", "techname", "placement", "port")
	form.MaxLength("techname", 100)
	form.PermittedPorts("port", 0, 65535)

	// Nodes checkbox
	var nodes []string
	for _, id := range r.PostForm["nodes"] {
		nodes = append(nodes, id)
	}

	//hostID := itemNode[0]
	var reg, method string

	if len(form.Get("reg_to_consul")) != 0 {
		reg = "true"
		method = form.Get("method_check_consul")
	} else {
		reg = "false"
	}

	location := form.Get("location")
	kind := form.Get("type")

	title := form.Get("title")
	techname := form.Get("techname")
	domain := form.Get("domain")
	value := form.Get("value")
	placement := form.Get("placement")

	team := form.Get("team")
	resp := form.Get("resp")

	port := form.Get("port")
	description := r.PostForm.Get("description")

	if !form.Valid() {
		app.render(w, r, "createService.page.tmpl", &templateData{
			InventoryHosts: inventoryHosts,
			Form:           form,
		})

		return
	}

	err = app.inventoryService.Insert(nodes, location, kind, title, techname, domain, value, placement, team, resp, reg, method, port, description)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "notice", "Submitted for Approval!")

	http.Redirect(w, r, fmt.Sprintf("/services"), http.StatusSeeOther)
}

func (app *application) showAPIServices(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "apiServices.page.tmpl", &templateData{})
}

//
// Manifests
//
func (app *application) showManifests(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "manifests.page.tmpl", &templateData{})
}

//
// Servers
//
func (app *application) showServers(w http.ResponseWriter, r *http.Request) {
	servers, err := app.serverAgent.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	statuses, err := app.serverAgent.GetStatuses()
	if err != nil {
		app.serverError(w, err)
		return
	}

	var cpu int
	var totalRAM, usedRAM, freeRAM string
	var totalHDD, usedHDD, freeHDD string
	var tr, ur, fr, th, uh, fh float64

	for _, server := range servers {
		var ram []string
		var hdd []string
		if len(server.RAM) != 0 || len(server.HDD) != 0 {
			c, _ := strconv.Atoi(server.NumCPU)
			cpu += c

			ram = strings.Split(server.RAM, "\n")
			ttr, _ := strconv.ParseFloat(ram[0], 64)
			tr += ttr
			uur, _ := strconv.ParseFloat(ram[1], 64)
			ur += uur
			ffr, _ := strconv.ParseFloat(ram[2], 64)
			fr += ffr

			hdd = strings.Split(server.HDD, "\n")
			tth, _ := strconv.ParseFloat(hdd[0], 64)
			th += tth
			uuh, _ := strconv.ParseFloat(hdd[1], 64)
			uh += uuh
			ffh, _ := strconv.ParseFloat(hdd[2], 64)
			fh += ffh
		}
	}
	totalRAM = fmt.Sprintf("%.1f", tr/1024/1024)
	usedRAM = fmt.Sprintf("%.1f", ur/1024/1024)
	freeRAM = fmt.Sprintf("%.1f", fr/1024/1024)

	totalHDD = fmt.Sprintf("%.1f", th)
	usedHDD = fmt.Sprintf("%.1f", uh)
	freeHDD = fmt.Sprintf("%.1f", fh)

	countServersAgent := app.serverAgent.GetCountServers()

	// Use the new render helper.
	app.render(w, r, "servers.page.tmpl", &templateData{
		ServersAgent: servers,
		Statuses:     statuses,

		TotalCPU: cpu,

		TotalRAM: totalRAM,
		UsedRAM:  usedRAM,
		FreeRAM:  freeRAM,

		TotalHDD: totalHDD,
		UsedHDD:  usedHDD,
		FreeHDD:  freeHDD,

		CountServersAgent: countServersAgent,
	})
}

func (app *application) showServer(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}
	// Use the ServerModel object's Get method to retrieve the data for a
	// specific record based on its ID. If no matching record is found,
	// return a 404 Not Found response.
	server, err := app.serverAgent.GetServer(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	sysctl, err := app.serverAgent.GetSysCtl(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	packages, err := app.serverAgent.GetPackages(id)
	if err == models.ErrNoRecord {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	var ram []string
	var totalRAM, usedRAM, freeRAM string
	if len(server.RAM) != 0 {
		ram = strings.Split(server.RAM, "\n")
		t, _ := strconv.ParseFloat(ram[0], 64)
		u, _ := strconv.ParseFloat(ram[1], 64)
		f, _ := strconv.ParseFloat(ram[2], 64)

		totalRAM = fmt.Sprintf("%.1f", t/1024/1024)
		usedRAM = fmt.Sprintf("%.1f", u/1024/1024)
		freeRAM = fmt.Sprintf("%.1f", f/1024/1024)
	}

	var hdd []string
	var totalHDD, usedHDD, freeHDD string
	if len(server.HDD) != 0 {
		hdd = strings.Split(server.HDD, "\n")
		t, _ := strconv.ParseFloat(hdd[0], 64)
		u, _ := strconv.ParseFloat(hdd[1], 64)
		f, _ := strconv.ParseFloat(hdd[2], 64)
		totalHDD = fmt.Sprintf("%.1f", t)
		usedHDD = fmt.Sprintf("%.1f", u)
		freeHDD = fmt.Sprintf("%.1f", f)
	}

	countServices := app.inventoryService.GetCountServicesByHost(server.Hostname)

	alert, err := app.serverAgent.AlertCheck(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the new render helper.
	app.render(w, r, "server.page.tmpl", &templateData{
		ServerAgent:      server,
		SysCtl:           sysctl,
		PackagesServer:   packages,
		CountInvServices: countServices,
		TotalRAM:         totalRAM,
		UsedRAM:          usedRAM,
		FreeRAM:          freeRAM,
		TotalHDD:         totalHDD,
		UsedHDD:          usedHDD,
		FreeHDD:          freeHDD,
		Alert:            alert,
	})
}

func (app *application) alertDisable(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	_, err = app.serverAgent.AlertDisable(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/servers/%d", id), http.StatusSeeOther)

}

func (app *application) alertEnable(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	err = app.serverAgent.AlertEnable(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/servers/%d", id), http.StatusSeeOther)

}

//
// Base Templates
//

// showBaseTpl ...
func (app *application) showBaseTpl(w http.ResponseWriter, r *http.Request) {
	baseTemplates, err := app.baseTemplate.GetAllBaseTpl()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "baseTpl.page.tmpl", &templateData{
		BaseTemplates: baseTemplates,
		Form:          forms.New(nil),
	})
}

// createBaseTpl ...
func (app *application) createBaseTpl(w http.ResponseWriter, r *http.Request) {
	baseTemplates, err := app.baseTemplate.GetAllBaseTpl()
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("value")
	form.MaxLength("value", 100)

	// If the form isn't valid, redisplay the template passing in the
	// form.Form object as the data.
	if !form.Valid() {
		app.render(w, r, "baseTpl.page.tmpl", &templateData{
			BaseTemplates: baseTemplates,
			Form:          form,
		})

		return
	}

	// Because the form data (with type url.Values) has been anonymously embedded
	// in the form.Form struct, can use the Get() method to retrieve
	// the validated value for a particular form field.
	id, err := app.baseTemplate.InsertBaseTpl(form.Get("type"), form.Get("value"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	value := fmt.Sprintf("New base template <strong>%s</strong> type <strong>%s</strong> added successfully! Please add parameters", form.Get("value"), form.Get("type"))
	app.session.Put(r, "notice", value)

	http.Redirect(w, r, fmt.Sprintf("/base/%s/%d", form.Get("type"), id), http.StatusSeeOther)
}

//
// Hard base
//

// editBaseTplHard ...
func (app *application) editBaseTplHard(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("valueTpl")

	baseTemplate, err := app.baseTemplate.GetBaseTpl(form.Get("idTpl"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	baseTemplatesItem, err := app.baseTemplate.GetAllBaseTplHardItem(form.Get("idTpl"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	if !form.Valid() {
		app.render(w, r, "baseTplItem.page.tmpl", &templateData{
			BaseTemplate:      baseTemplate,
			BaseTemplatesItem: baseTemplatesItem,
			Form:              form,
		})

		return
	}

	err = app.baseTemplate.EditBaseTpl(form.Get("idTpl"), form.Get("valueTpl"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/base/%s/%s", form.Get("typeTpl"), form.Get("idTpl")), http.StatusSeeOther)
}

// deleteBaseTplHard ...
func (app *application) deleteBaseTplHard(w http.ResponseWriter, r *http.Request) {
	p, err := url.Parse(r.URL.Path)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.baseTemplate.DeleteBaseTplHard(path.Base(p.Path))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/base", http.StatusSeeOther)
}

// showBaseTplHard ...
func (app *application) showBaseTplHard(w http.ResponseWriter, r *http.Request) {
	p, err := url.Parse(r.URL.Path)
	if err != nil {
		app.serverError(w, err)
		return
	}

	baseTemplate, err := app.baseTemplate.GetBaseTpl(path.Base(p.Path))
	if err != nil {
		app.serverError(w, err)
		return
	}

	baseTemplatesItem, err := app.baseTemplate.GetAllBaseTplHardItem(path.Base(p.Path))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "baseTplItem.page.tmpl", &templateData{
		BaseTemplate:      baseTemplate,
		BaseTemplatesItem: baseTemplatesItem,
		Form:              forms.New(nil),
	})
}

//
// Package base
//

// editBaseTplPackage ...
func (app *application) editBaseTplPackage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("valueTpl")

	baseTemplate, err := app.baseTemplate.GetBaseTpl(form.Get("idTpl"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	baseTemplatesItem, err := app.baseTemplate.GetAllBaseTplPackageItem(form.Get("idTpl"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	if !form.Valid() {
		app.render(w, r, "baseTplItem.page.tmpl", &templateData{
			BaseTemplate:      baseTemplate,
			BaseTemplatesItem: baseTemplatesItem,
			Form:              form,
		})

		return
	}

	err = app.baseTemplate.EditBaseTpl(form.Get("idTpl"), form.Get("valueTpl"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/base/%s/%s", form.Get("typeTpl"), form.Get("idTpl")), http.StatusSeeOther)
}

// deleteBaseTplPackage ...
func (app *application) deleteBaseTplPackage(w http.ResponseWriter, r *http.Request) {
	p, err := url.Parse(r.URL.Path)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.baseTemplate.DeleteBaseTplPackage(path.Base(p.Path))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/base", http.StatusSeeOther)
}

// showBaseTplPackage ...
func (app *application) showBaseTplPackage(w http.ResponseWriter, r *http.Request) {
	p, err := url.Parse(r.URL.Path)
	if err != nil {
		app.serverError(w, err)
		return
	}

	baseTemplate, err := app.baseTemplate.GetBaseTpl(path.Base(p.Path))
	if err != nil {
		app.serverError(w, err)
		return
	}

	baseTemplatesItem, err := app.baseTemplate.GetAllBaseTplPackageItem(path.Base(p.Path))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "baseTplItem.page.tmpl", &templateData{
		BaseTemplate:      baseTemplate,
		BaseTemplatesItem: baseTemplatesItem,
		Form:              forms.New(nil),
	})
}

//
// ResolvConf base
//

// editBaseTplResolvConf ...
func (app *application) editBaseTplResolvConf(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("valueTpl")

	baseTemplate, err := app.baseTemplate.GetBaseTpl(form.Get("idTpl"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	baseTemplatesItem, err := app.baseTemplate.GetAllBaseTplResolvConfItem(form.Get("idTpl"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	if !form.Valid() {
		app.render(w, r, "baseTplItem.page.tmpl", &templateData{
			BaseTemplate:      baseTemplate,
			BaseTemplatesItem: baseTemplatesItem,
			Form:              form,
		})

		return
	}

	err = app.baseTemplate.EditBaseTpl(form.Get("idTpl"), form.Get("valueTpl"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/base/%s/%s", form.Get("typeTpl"), form.Get("idTpl")), http.StatusSeeOther)
}

// deleteBaseTplResolvConf ...
func (app *application) deleteBaseTplResolvConf(w http.ResponseWriter, r *http.Request) {
	p, err := url.Parse(r.URL.Path)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.baseTemplate.DeleteBaseTplResolvConf(path.Base(p.Path))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/base", http.StatusSeeOther)
}

// showBaseTplResolvConf ...
func (app *application) showBaseTplResolvConf(w http.ResponseWriter, r *http.Request) {
	p, err := url.Parse(r.URL.Path)
	if err != nil {
		app.serverError(w, err)
		return
	}

	baseTemplate, err := app.baseTemplate.GetBaseTpl(path.Base(p.Path))
	if err != nil {
		app.serverError(w, err)
		return
	}

	baseTemplatesItem, err := app.baseTemplate.GetAllBaseTplResolvConfItem(path.Base(p.Path))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "baseTplItem.page.tmpl", &templateData{
		BaseTemplate:      baseTemplate,
		BaseTemplatesItem: baseTemplatesItem,
		Form:              forms.New(nil),
	})
}

//
// SysCtl base
//

// editBaseTplResolvConf ...
func (app *application) editBaseTplSysCtl(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("valueTpl")

	baseTemplate, err := app.baseTemplate.GetBaseTpl(form.Get("idTpl"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	baseTemplatesItem, err := app.baseTemplate.GetAllBaseTplSysCtlItem(form.Get("idTpl"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	if !form.Valid() {
		app.render(w, r, "baseTplItem.page.tmpl", &templateData{
			BaseTemplate:      baseTemplate,
			BaseTemplatesItem: baseTemplatesItem,
			Form:              form,
		})

		return
	}

	err = app.baseTemplate.EditBaseTpl(form.Get("idTpl"), form.Get("valueTpl"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/base/%s/%s", form.Get("typeTpl"), form.Get("idTpl")), http.StatusSeeOther)
}

// deleteBaseTplResolvConf ...
func (app *application) deleteBaseTplSysCtl(w http.ResponseWriter, r *http.Request) {
	p, err := url.Parse(r.URL.Path)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.baseTemplate.DeleteBaseTplSysCtl(path.Base(p.Path))
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, "/base", http.StatusSeeOther)
}

// showBaseTplResolvConf ...
func (app *application) showBaseTplSysCtl(w http.ResponseWriter, r *http.Request) {
	p, err := url.Parse(r.URL.Path)
	if err != nil {
		app.serverError(w, err)
		return
	}

	baseTemplate, err := app.baseTemplate.GetBaseTpl(path.Base(p.Path))
	if err != nil {
		app.serverError(w, err)
		return
	}

	baseTemplatesItem, err := app.baseTemplate.GetAllBaseTplSysCtlItem(path.Base(p.Path))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "baseTplItem.page.tmpl", &templateData{
		BaseTemplate:      baseTemplate,
		BaseTemplatesItem: baseTemplatesItem,
		Form:              forms.New(nil),
	})
}

//
// Common base
//

// createBaseTplItem ...
func (app *application) createBaseTplItem(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create a new forms.Form struct containing the POSTed data from the
	// form, then use the validation methods to check the content.
	form := forms.New(r.PostForm)
	form.Required("title")
	form.Required("value")

	baseTemplate, err := app.baseTemplate.GetBaseTpl(form.Get("id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	if form.Get("type") == "hard" {
		baseTemplatesItem, err := app.baseTemplate.GetAllBaseTplHardItem(form.Get("id"))
		if err != nil {
			app.serverError(w, err)
			return
		}

		if !form.Valid() {
			app.render(w, r, "baseTplItem.page.tmpl", &templateData{
				BaseTemplate:      baseTemplate,
				BaseTemplatesItem: baseTemplatesItem,
				Form:              form,
			})

			return
		}

		_, err = app.baseTemplate.InsertBaseTplHardItem(form.Get("id"), form.Get("title"), form.Get("value"))
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	if form.Get("type") == "package" {
		baseTemplatesItem, err := app.baseTemplate.GetAllBaseTplPackageItem(form.Get("id"))
		if err != nil {
			app.serverError(w, err)
			return
		}

		if !form.Valid() {
			app.render(w, r, "baseTplItem.page.tmpl", &templateData{
				BaseTemplate:      baseTemplate,
				BaseTemplatesItem: baseTemplatesItem,
				Form:              form,
			})

			return
		}

		_, err = app.baseTemplate.InsertBaseTplPackageItem(form.Get("id"), form.Get("title"), form.Get("value"))
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	if form.Get("type") == "resolvconf" {
		baseTemplatesItem, err := app.baseTemplate.GetAllBaseTplResolvConfItem(form.Get("id"))
		if err != nil {
			app.serverError(w, err)
			return
		}

		if !form.Valid() {
			app.render(w, r, "baseTplItem.page.tmpl", &templateData{
				BaseTemplate:      baseTemplate,
				BaseTemplatesItem: baseTemplatesItem,
				Form:              form,
			})

			return
		}

		_, err = app.baseTemplate.InsertBaseTplResolvConfItem(form.Get("id"), form.Get("title"), form.Get("value"))
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	if form.Get("type") == "sysctl" {
		baseTemplatesItem, err := app.baseTemplate.GetAllBaseTplSysCtlItem(form.Get("id"))
		if err != nil {
			app.serverError(w, err)
			return
		}

		if !form.Valid() {
			app.render(w, r, "baseTplItem.page.tmpl", &templateData{
				BaseTemplate:      baseTemplate,
				BaseTemplatesItem: baseTemplatesItem,
				Form:              form,
			})

			return
		}

		_, err = app.baseTemplate.InsertBaseTplSysCtlItem(form.Get("id"), form.Get("title"), form.Get("value"))
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	value := fmt.Sprintf("Added a new parameter %s = %s type %s", form.Get("title"), form.Get("value"), form.Get("type"))
	app.session.Put(r, "notice", value)

	http.Redirect(w, r, fmt.Sprintf("/base/%s/%s", form.Get("type"), form.Get("id")), http.StatusSeeOther)
}

// deleteBaseTplItem ...
func (app *application) deleteBaseTplItem(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	u := strings.Split(url, "/")
	typeTpl := u[2]
	idTpl := u[3]
	idItem := u[6]

	// Hard
	if typeTpl == "hard" {
		err := app.baseTemplate.DeleteBaseTplHardItem(idTpl, idItem)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	// Package
	if typeTpl == "package" {
		err := app.baseTemplate.DeleteBaseTplPackageItem(idTpl, idItem)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	// ResolvConf
	if typeTpl == "resolvconf" {
		err := app.baseTemplate.DeleteBaseTplResolvConfItem(idTpl, idItem)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	// SysCtl
	if typeTpl == "sysctl" {
		err := app.baseTemplate.DeleteBaseTplSysCtlItem(idTpl, idItem)
		if err != nil {
			app.serverError(w, err)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/base/%s/%s", typeTpl, idTpl), http.StatusSeeOther)
}

//
// Approval
//
func (app *application) forApproval(w http.ResponseWriter, r *http.Request) {
	hostsForApprove, err := app.inventoryHost.GetAllForApprove()
	if err != nil {
		app.serverError(w, err)
		return
	}

	hostsForDelete, err := app.inventoryHost.GetAllForDelete()
	if err != nil {
		app.serverError(w, err)
		return
	}

	servicesForApprove, err := app.inventoryService.GetAllForApprove()
	if err != nil {
		app.serverError(w, err)
		return
	}

	servicesForDelete, err := app.inventoryService.GetAllForDelete()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "approval.page.tmpl", &templateData{
		HostsForApprove:    hostsForApprove,
		HostsForDelete:     hostsForDelete,
		ServicesForApprove: servicesForApprove,
		ServicesForDelete:  servicesForDelete,
	})
}

func (app *application) approveService(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	err = app.inventoryService.ApproveService(form.Get("id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "notice", "Approved successfully")
	http.Redirect(w, r, "/approval", http.StatusSeeOther)
}

func (app *application) approveDeleteService(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	err = app.inventoryService.Delete(form.Get("id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "notice", "Delete service successfully!")
	http.Redirect(w, r, "/services", http.StatusSeeOther)
}

func (app *application) approveDeleteCancelService(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	err = app.inventoryService.CancelDeletion(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "notice", "Cancel deletion service successfully!")
	http.Redirect(w, r, "/services", http.StatusSeeOther)
}

func (app *application) approveHost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	err = app.inventoryHost.ApproveHost(form.Get("id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "notice", "Approved successfully")
	http.Redirect(w, r, "/approval", http.StatusSeeOther)
}

func (app *application) approveDeleteHost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	err = app.inventoryHost.Delete(form.Get("id"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "notice", "Delete host successfully!")
	http.Redirect(w, r, "/hosts", http.StatusSeeOther)
}

func (app *application) approveDeleteCancelHost(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w) // Use the notFound() helper.
		return
	}

	err = app.inventoryHost.CancelDeletion(id)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.session.Put(r, "notice", "Cancel deletion host successfully!")
	http.Redirect(w, r, "/hosts", http.StatusSeeOther)
}

//
// User Authentication
//
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents using the form helper.
	form := forms.New(r.PostForm)
	form.Required("username", "email", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 8)

	// If there are any errors, redisplay the signup form.
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})

		return
	}

	// Create a new user record in the database. If the email already
	// exists add an error message to the form and re-display it.
	err = app.users.Insert(form.Get("username"), form.Get("email"), form.Get("password"))
	if err == models.ErrDuplicateEmail {
		form.Errors.Set("email", "Address is already in use")
		app.render(w, r, "signup.page.tmpl", &templateData{
			Form: form,
		})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Otherwise add a confirmation notice message to the session
	// confirming that signup worked and asking them to log in.
	app.session.Put(r, "notice", "Your signup was successful. Please log in.")

	// Redirect the user to the login page.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check credentials are valid. If not, add a generic error message
	// to the form failures map and re-display the login page.
	form := forms.New(r.PostForm)
	id, err := app.users.Auth(form.Get("email"), form.Get("password"))
	if err == models.ErrInvalidCredentials {
		form.Errors.Set("generic", "E-mail or Password is incorrect")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}

	// Add the ID of the current user to the session, so that they are now 'logged in'.
	app.session.Put(r, "userID", id)

	// Add the userRole of the current user to the session.
	//app.session.Put(r, "userRole", userRole)

	// Redirect the user to overview page.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// Remove the userID from the session data so that the user is 'logged out'.
	app.session.Remove(r, "userID")
	// Add a notice message to the session to confirm to the user that logged out.
	app.session.Put(r, "notice", "You've been logged out successfully!")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

//
// History
//
func (app *application) showHistory(w http.ResponseWriter, r *http.Request) {
	historyAll, err := app.history.GetAll()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "history.page.tmpl", &templateData{
		HistoryAll: historyAll,
	})

}

//
// Stages
//
func (app *application) showStages(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get environment for filtration
	environment := "all" // Init
	for _, value := range r.Form["env"] {
		environment = value
	}

	stages, err := app.inventoryService.GetServicesByEnv(environment)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "stages.page.tmpl", &templateData{
		Stages:      stages,
		Environment: environment,
	})
}

func (app *application) showFWD(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get environment for filtration
	environment := "all" // Init
	for _, value := range r.Form["env"] {
		environment = value
	}

	stages, err := app.inventoryService.GetServicesByEnv(environment)
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "fwd.page.tmpl", &templateData{
		Stages:      stages,
		Environment: environment,
	})
}

//
// Help
//
func (app *application) showHelp(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "help.page.tmpl", &templateData{})
}

//
// In Dev
//
func (app *application) showInDev(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "inDev.page.tmpl", &templateData{})
}
