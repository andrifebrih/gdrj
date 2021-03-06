package main

import (
	"eaciit/gdrj/web/controller"
	"eaciit/gdrj/web/installation"
	"eaciit/gdrj/web/model"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"net/http"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	server *knot.Server
)

func main() {
	runtime.GOMAXPROCS(4)
	gocore.ConfigPath = controller.GDRJ_CONFIG_PATH

	server = new(knot.Server)

	port := new(gocore.Ports)
	port.ID = "port"
	if err := port.GetPort(); err != nil {
		toolkit.Printf("Error get port: %s \n", err.Error())
	}
	if port.Port == 0 {
		if err := setup.Port(port); err != nil {
			toolkit.Printf("Error set port: %s \n", err.Error())
		}
	}

	webController := controller.CreateWebController(server)

	server.Address = toolkit.Sprintf("localhost:%v", toolkit.ToString(port.Port))
	server.RouteStatic("res", filepath.Join(controller.AppBasePath, "web", "assets"))
	server.RouteStatic("image", filepath.Join(controller.AppBasePath, "web", "assets", "img"))
	server.Register(webController, "")
	server.Register(controller.CreateLoginController(server), "")
	server.Register(controller.CreateDataBrowserController(server), "")
	server.Register(controller.CreateUploadDataController(server), "")
	server.Register(controller.CreateReportController(server), "")
	server.Register(controller.CreateAllocationFlowController(server), "")
	server.Register(controller.CreateAdminisrationController(server), "")
	server.Register(controller.CreateSessionController(server), "")
	server.Register(controller.CreateUserController(server), "")
	server.Register(controller.CreateGroupController(server), "")
	server.Register(controller.CreateOrganizationController(server), "")

	server.Route("/", func(r *knot.WebContext) interface{} {
		regex := regexp.MustCompile("/web/report/[a-zA-Z0-9_]+(/.*)?$")
		rURL := r.Request.URL.String()

		if regex.MatchString(rURL) {
			args := strings.Split(strings.Replace(rURL, "/web/report/", "", -1), "/")
			return webController.PageReport(r, args)
		}

		sessionid := r.Session("sessionid", "")
		if sessionid == "" {
			http.Redirect(r.Writer, r.Request, "/web/login", 301)
		} else {
			http.Redirect(r.Writer, r.Request, "/web/databrowser", 301)
		}

		return true
	})

	if err := setAclDatabase(); err != nil {
		toolkit.Printf("Error set acl database : %s \n", err.Error())
	}

	server.Listen()
}

func setAclDatabase() error {
	if err := gocore.InitialSetDatabase(); err != nil {
		return err
	}

	if gocore.GetConfig("default_username") == nil {
		gocore.SetConfig("default_username", "eaciit")
		gocore.SetConfig("default_password", "Password.1")
	}

	if err := gocore.PrepareDefaultUser(); err != nil {
		return err
	}
	return nil
}
