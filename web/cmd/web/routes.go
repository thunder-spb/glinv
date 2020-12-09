package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"

	cache "github.com/victorspringer/http-cache"
	"github.com/victorspringer/http-cache/adapter/memory"
)

// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	memcached, err := memory.NewAdapter(
		memory.AdapterWithAlgorithm(memory.LRU),
		memory.AdapterWithCapacity(10000000),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cacheClient, err := cache.NewClient(
		cache.ClientWithAdapter(memcached),
		cache.ClientWithTTL(10*time.Minute),
		cache.ClientWithRefreshKey("opn"),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Middleware chain containing "standard" middleware
	// which will be used for every request application receives.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// Dynamic application routes. For now, this chain
	// will only contain the session middleware.
	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	mux := pat.New()

	mux.Get("/", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.overview)) // Display the home page

	// Inventory Hosts
	mux.Get("/hosts", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showHosts))
	mux.Post("/hosts/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createHost))
	mux.Get("/hosts/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createHostForm))
	mux.Post("/hosts/edit", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.editHost))
	mux.Get("/hosts/edit/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.editHostForm))
	mux.Post("/hosts/delete/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteHost))
	mux.Get("/hosts/api", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showAPIHosts))
	mux.Get("/hosts/ssh", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showSSHConfig))

	// Inventory Tags of Hosts
	mux.Get("/hosts/tags", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showHTags))
	mux.Post("/hosts/tags", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createHTag))
	mux.Get("/hosts/tag/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showHTagForm))
	mux.Post("/hosts/tag/edit/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.editHTag))
	mux.Post("/hosts/tag/delete/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteHTag))

	// Inventory Vars of Hosts
	mux.Get("/hosts/vars", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showHVars))
	mux.Post("/hosts/vars", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createHVar))
	mux.Get("/hosts/var/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showHVarForm))
	mux.Post("/hosts/var/edit/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.editHVar))
	mux.Post("/hosts/var/delete/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteHVar))
	mux.Get("/hosts/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showHost))

	// Inventory Groups
	mux.Get("/groups", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showInventoryGroups))
	mux.Post("/groups/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createInventoryGroup))
	mux.Get("/groups/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createInventoryGroupForm))
	mux.Get("/groups/api", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showAPIGroups))
	mux.Get("/group/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.editInventoryGroupForm))
	mux.Post("/group/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.editInventoryGroup))
	mux.Post("/group/delete/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteInventoryGroup))

	// Inventory Vars of Groups
	mux.Get("/groups/vars", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showGVars))
	mux.Post("/groups/vars", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createGVars))
	mux.Get("/groups/var/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showGVarForm))
	mux.Post("/groups/var/edit/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.editGVar))
	mux.Post("/groups/var/delete/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteGVar))

	// Inventory Services(Environments)
	mux.Get("/services", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showServices))
	mux.Get("/services/host/:hostname", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showServicesHost))
	mux.Post("/services/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createService))
	mux.Get("/services/create", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createServiceForm))
	mux.Post("/services/edit", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.editService))
	mux.Get("/services/edit/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.editServiceForm))
	mux.Post("/services/delete/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteService))
	mux.Get("/services/api", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showAPIServices))
	mux.Get("/services/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showService))

	// Manifests
	mux.Get("/manifests", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showManifests))

	// Soft Info
	// mux.Get("/basesoft", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showBaseSoft))
	// mux.Post("/basesoft", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createBaseSoft))

	// Hard Info
	mux.Get("/servers", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showServers))
	mux.Get("/servers/alert/disable/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.alertDisable))
	mux.Get("/servers/alert/enable/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.alertEnable))
	mux.Get("/servers/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showServer))

	// Base Templates
	mux.Get("/base", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showBaseTpl))
	mux.Post("/base", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createBaseTpl))
	// Hard
	mux.Get("/base/hard/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showBaseTplHard))
	mux.Post("/base/hard/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createBaseTplItem))
	mux.Get("/base/hard/:id/delete/item/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteBaseTplItem))
	mux.Post("/base/hard/edit/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.editBaseTplHard))
	mux.Get("/base/hard/delete/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteBaseTplHard))
	// Package
	mux.Get("/base/package/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showBaseTplPackage))
	mux.Post("/base/package/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createBaseTplItem))
	mux.Get("/base/package/:id/delete/item/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteBaseTplItem))
	mux.Post("/base/package/edit/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.editBaseTplPackage))
	mux.Get("/base/package/delete/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteBaseTplPackage))
	// ResolvConf
	mux.Get("/base/resolvconf/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showBaseTplResolvConf))
	mux.Post("/base/resolvconf/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createBaseTplItem))
	mux.Get("/base/resolvconf/:id/delete/item/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteBaseTplItem))
	mux.Post("/base/resolvconf/edit/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.editBaseTplResolvConf))
	mux.Get("/base/resolvconf/delete/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteBaseTplResolvConf))

	// SysCtl
	mux.Get("/base/sysctl/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showBaseTplSysCtl))
	mux.Post("/base/sysctl/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.createBaseTplItem))
	mux.Get("/base/sysctl/:id/delete/item/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteBaseTplItem))
	mux.Post("/base/sysctl/edit/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.editBaseTplSysCtl))
	mux.Get("/base/sysctl/delete/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.deleteBaseTplSysCtl))

	// Approval
	mux.Get("/approval", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.forApproval))
	mux.Post("/approval/service", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.approveService))
	mux.Post("/approval/host", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.approveHost))
	mux.Post("/delete/host", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.approveDeleteHost))
	mux.Post("/delete/service", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.approveDeleteService))
	mux.Get("/delete/cancel/host/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.approveDeleteCancelHost))
	mux.Get("/delete/cancel/service/:id", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.approveDeleteCancelService))

	// User Authentication
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))                                   // Display the user signup form
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))                                      // Create a new user
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))                                     // Display the user login form
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))                                        // Authenticate and login the user
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.logoutUser)) // Logout the user
	//mux.Post("/agent", http.HandlerFunc(app.doAgent))

	// History
	mux.Get("/history", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showHistory))

	// Stages
	mux.Get("/stages", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showStages))
	mux.Get("/fwd", dynamicMiddleware.ThenFunc(app.showFWD))

	// Help
	mux.Get("/help", cacheClient.Middleware(dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showHelp)))

	// In Dev TODO
	mux.Get("/software", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showInDev))
	mux.Get("/locations", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showInDev))
	mux.Get("/applications", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showInDev))
	mux.Get("/devices", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showInDev))
	mux.Get("/teams", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showInDev))
	mux.Get("/users", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showInDev))
	mux.Get("/documents", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showInDev))
	mux.Get("/reports", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showInDev))
	mux.Get("/administration", dynamicMiddleware.Append(app.requireAuthenticatedUser).ThenFunc(app.showInDev))

	// Assets
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
