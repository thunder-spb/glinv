package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"glinv/pkg/models"

	"github.com/justinas/nosurf"
)

// The serverError helper writes an error message and stack trace to the errorLog,
// then sends a generic 500 Internal Server Error response to the user.
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	frameDepth := 2

	if err := app.errorLog.Output(frameDepth, trace); err != nil {
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description to the user.
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// For consistency, also implement a notFound helper. This is simply a
// convenience wrapper around clientError which sends a 404 Not Found response to the user.
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Create an addDefaultData helper. This takes a pointer to a templateData
// struct, adds the current year to the CurrentYear field, and then returns the pointer and etc.
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CSRFToken = nosurf.Token(r)
	td.AuthenticatedUser = app.authenticatedUser(r)
	td.CurrentYear = time.Now().Year()
	td.Notice = app.session.PopString(r, "notice")

	hosts := app.inventoryHost.GetCountUnapprovedHosts()
	services := app.inventoryService.GetCountUnapprovedServices()
	td.CountUnapproved = hosts + services
	td.Version = "web:0.0.2 api:0.0.1 agent:2:0.0.4-2"

	return td
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serverError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the
	// http.ResponseWriter. If there's an error, call our serverError helper and then return.
	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}

	if _, err := buf.WriteTo(w); err != nil {
		return
	}
}

// authenticatedUser method returns the ID of the current user from the
// session, or zero if the request is from an unauthenticated user.
func (app *application) authenticatedUser(r *http.Request) *models.User {
	//return app.session.GetInt(r, "userID")
	user, ok := r.Context().Value(contextKeyUser).(*models.User)
	if !ok {
		return nil
	}
	return user
}

// // Render Tree Groups for Host
// func GetGVars(groupID int) []*models.InventoryGVar {
// 	return nil
// }

// Tree group model
type Tree struct {
	ID       int
	ParentID int
	Value    string
	Meta     bool
	Checked  []int
	GVars    []*models.InventoryGVar
	Nodes    []*Tree
}

// TreeGroups function fills the Tree structure recursively and
// returns slice to templates. Templates also use recursion for output.
func (t *Tree) TreeGroups(groups []*models.InventoryGroup, groupGVars map[int][]*models.InventoryGVar, checkedIDGroups []int) *Tree {
	tree := &Tree{}
	gvars := []*models.InventoryGVar{}

	// 0-level
	for _, g := range groups {
		if g.ParentID == 0 {
			gvars = groupGVars[g.ID]
			tree = &Tree{ID: g.ID, ParentID: g.ParentID, Value: g.Value, Checked: checkedIDGroups, GVars: gvars, Nodes: nil}
		}
	}

	// 1-level
	for _, g := range groups {
		if g.ParentID == tree.ID {
			gvars = groupGVars[g.ID]
			node := &Tree{ID: g.ID, ParentID: g.ParentID, Value: g.Value, Checked: checkedIDGroups, GVars: gvars, Nodes: nil}
			tree.AddNode(node)
		}
	}

	// 2-level
	for i := range tree.Nodes {
		for _, g := range groups {
			if g.ParentID == tree.Nodes[i].ID {
				gvars = groupGVars[g.ID]
				node := &Tree{ID: g.ID, ParentID: g.ParentID, Value: g.Value, Checked: checkedIDGroups, GVars: gvars, Nodes: nil}
				tree.Nodes[i].AddNode(node)
			}
		}

		// 3-level
		for j := range tree.Nodes[i].Nodes {
			for _, g := range groups {
				if g.ParentID == tree.Nodes[i].Nodes[j].ID {
					gvars = groupGVars[g.ID]
					node := &Tree{ID: g.ID, ParentID: g.ParentID, Value: g.Value, Checked: checkedIDGroups, GVars: gvars, Nodes: nil}
					tree.Nodes[i].Nodes[j].AddNode(node)
				}
			}
		}
	}

	return tree
}

// AddNode ...
func (t *Tree) AddNode(node *Tree) {
	t.Nodes = append(t.Nodes, node)
}
