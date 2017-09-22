package servicemux

import (
	"fmt"
	"net/http"
	"strings"
)

type service struct {
	action func(http.ResponseWriter, *http.Request)
	ports  []string
}

func (s *service) runAction(res http.ResponseWriter, req *http.Request) bool {

	var port = strings.Split(req.Host, ":")[1]

	for _, val := range s.ports {
		if val == port || val == "*" {
			s.action(res, req)
			return true
		}
	}
	return false
}

var serviceMap = make(map[string]service)

// AddService adds a service to the service list
// Arguments:
// Domain string: domain name in format 'domain.com', Domain '*' is a wildcard for all
// Service func: function runs when there is a connection to this service
// Ports ...string: slice of allowed ports, port '*' is a wildcard for all
func AddService(Domain string, Service func(http.ResponseWriter, *http.Request), Ports ...string) {
	serviceMap[Domain] = service{
		action: Service,
		ports:  Ports,
	}
}

// Run is a function that shound be run on every http query
func Run(res http.ResponseWriter, req *http.Request) {
	if val, ok := serviceMap[strings.Split(req.Host, ":")[0]]; ok {
		if val.runAction(res, req) {
			return
		}
	}

	if val, ok := serviceMap["*"]; ok {
		if val.runAction(res, req) {
			return
		}
	}

	res.WriteHeader(http.StatusBadGateway)
	fmt.Fprint(res, "BAD GATEWAY")
}
