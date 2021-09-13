package http

import (
	"fmt"
	"github.com/busgo/pink/conf"
	"github.com/busgo/pink/http/handler"
	"github.com/busgo/pink/node"
	"github.com/busgo/pink/pkg/etcd"
	"github.com/busgo/pink/pkg/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	JobListUrl    = "/job/list"
	JobAddUrl     = "/job/add"
	JobUpdateUrl  = "/job/update"
	JobDeleteUrl  = "/job/delete"
	JobExecuteUrl = "/job/execute"

	// group
	GroupListUrl        = "/group/list"
	GroupDetailsListUrl = "/group/details/list"
	GroupAddUrl         = "/group/add"

	// node
	NodeListUrl = "/node/list"

	// schedule plan
	SchedulePlanListUrl       = "/schedule/plan/list"
	SchedulePlanClientsUrl    = "/schedule/plan/clients"
	ScheduleSnapshotsUrl      = "/schedule/snapshots"
	ScheduleSnapshotDeleteUrl = "/schedule/snapshot/delete"

	// snapshot
	ExecuteSnapshotsUrls            = "/execute/snapshots"
	ExecuteHistorySnapshotsUrls     = "/execute/history/snapshots"
	ExecuteHistorySnapshotDeleteUrl = "/execute/history/snapshot/delete"
)

type Service struct {
	router  *echo.Echo
	etcdCli *etcd.Cli
	node    *node.PinkNode
	h       *handler.PinkWebHandler
	addr    string
}

func NewHttpService(app *conf.AppConf, etcdCli *etcd.Cli, node *node.PinkNode, handler *handler.PinkWebHandler) *Service {

	router := echo.New()
	router.Use(middleware.Recover())
	router.Use(middleware.Logger())
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.POST, echo.GET, echo.PUT},
	}))
	s := &Service{router: router, etcdCli: etcdCli, node: node, h: handler, addr: fmt.Sprintf(":%d", app.Port)}
	s.initRoutes()
	return s
}

// init routes
func (s *Service) initRoutes() {

	// job conf
	s.router.POST(JobListUrl, s.h.JobConfList)
	s.router.POST(JobAddUrl, s.h.AddJobConf)
	s.router.POST(JobUpdateUrl, s.h.UpdateJobConf)
	s.router.POST(JobDeleteUrl, s.h.DeleteJobConf)
	s.router.POST(JobExecuteUrl, s.h.JobExecute)

	// group
	s.router.POST(GroupListUrl, s.h.GroupList)
	s.router.POST(GroupDetailsListUrl, s.h.GroupDetailsList)
	s.router.POST(GroupAddUrl, s.h.AddGroup)

	// node
	s.router.POST(NodeListUrl, s.h.NodeList)

	// schedule plan
	s.router.POST(SchedulePlanListUrl, s.h.SchedulePlanList)
	s.router.POST(SchedulePlanClientsUrl, s.h.SchedulePlanClients)
	s.router.POST(ScheduleSnapshotsUrl, s.h.ScheduleSnapshots)
	s.router.POST(ScheduleSnapshotDeleteUrl, s.h.ScheduleSnapshotDelete)

	// execute
	s.router.POST(ExecuteSnapshotsUrls, s.h.ExecuteSnapshots)
	s.router.POST(ExecuteHistorySnapshotsUrls, s.h.ExecuteHistorySnapshots)
	s.router.POST(ExecuteHistorySnapshotDeleteUrl, s.h.ExecuteHistorySnapshotDelete)
}

func (s *Service) Run() error {

	err := s.node.Run()
	if err != nil {
		log.Errorf("the pink node run fail:%+v", err)
		return err
	}

	log.Infof("the http service listen to %s", s.addr)
	err = s.router.Start(s.addr)
	if err != nil {
		log.Errorf("the http service for  %s start fail:%+v", s.addr, err)
	}
	return err
}
