
// Package classification Gravity Node API.
//
// the purpose of this application is to provide an application
// that is using plain go code to define an API
//
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: node.gravityhub.org:8090
//     BasePath: /
//     Version: 1.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: venlab.dev <shamil@venlab.dev> https://venlab.dev
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: KEY
//          in: header
//
//     Extensions:
//     x-meta-value: value
//     x-meta-array:
//       - value1
//       - value2
//     x-meta-array-obj:
//       - name: obj
//         value: field
//
// swagger:meta
package main

import (
	"github.com/Gravity-Hub-Org/gravity-node-api-mockup/v2/controller"
	"github.com/Gravity-Hub-Org/gravity-node-api-mockup/v2/router"
	"flag"
	"fmt"
	"github.com/Gravity-Hub-Org/gravity-node-api-mockup/v2/utils"
	"net/http"
)

var port, shouldFill string

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func init() {
	flag.StringVar(&port, "port", "8090", "Path to config.toml")
	flag.StringVar(&shouldFill, "fill", "0", "Whether to fill db with mockup data or not")
	flag.Parse()
}

func main () {
	if shouldFill != "0" {
		db := &controller.DBController{ DB: utils.ConnectToPG() }
		go db.PersistMockup()
	}

	http.HandleFunc("/hello", headers)
	http.HandleFunc(router.GetAllNebulas, controller.GetAllNebulas)
	http.HandleFunc(router.GetAllNodes, controller.GetAllNodes)
	http.HandleFunc(router.GetCommonStats, controller.GetCommonStats)
	http.HandleFunc(router.GetNodeRewards, controller.GetNodeRewardsList)
	http.HandleFunc(router.GetNodeActionsHistory, controller.GetNodeActionsHistory)

	fmt.Printf("Listening on port: %s\n", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}