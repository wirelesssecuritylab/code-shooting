package restserver

import (
	"code-shooting/infra/common"
	"code-shooting/infra/restserver/internal"

	"github.com/pkg/errors"
	"go.uber.org/fx"

	marsconfig "code-shooting/infra/config"
	"code-shooting/infra/config/model"
	"code-shooting/infra/logger"
	"code-shooting/infra/restserver/middleware"
	"code-shooting/infra/restserver/utils"
)

type MultRestServerConfParam struct {
	fx.In
	Config marsconfig.Config
}

func NewMultRestServerConf(param MultRestServerConfParam) (internal.MultRestServerConf, error) {
	utils.Config = param.Config
	multConf, err := GetConfigParser().Parse(param.Config)
	if err != nil {
		return nil, errors.Wrap(err, "parse config")
	}

	return multConf, nil
}

type MultRestServerParam struct {
	fx.In
	Multconf     internal.MultRestServerConf
	Bootstrapers []middleware.IMiddlewareBootstraper `group:"echomiddleware"`
	Opts         []Option                            `optional:"true"`
}

type MultRestServer struct {
	RestServerMap map[string]*RestServer
}

func newMultRestServer() *MultRestServer {
	mRestServer := &MultRestServer{}
	mRestServer.RestServerMap = make(map[string]*RestServer)
	return mRestServer
}

func NewMultRestServer(lc fx.Lifecycle, param MultRestServerParam) (*MultRestServer, error) {

	multRestServer := newMultRestServer()
	for _, itemConf := range param.Multconf {
		serverParam := RestServerParam{}
		serverParam.Conf = itemConf
		serverParam.Bootstrapers = param.Bootstrapers
		serverParam.Opts = param.Opts
		restServer, err := newRestServerModule(lc, serverParam)
		if err != nil {
			return nil, err
		}
		multRestServer.RestServerMap[serverParam.Conf.Name] = restServer
	}
	multRestServer.RegisterConfigChangeEvent()
	return multRestServer, nil
}

func (mult *MultRestServer) GetRestServerByName(name string) (*RestServer, error) {
	if name == "" {
		return nil, errors.New("rest server name is nil")
	}

	server, ok := mult.RestServerMap[name]
	if !ok {
		return nil, errors.Errorf("no rest server with name %s", name)
	}

	if server.server == nil {
		return nil, errors.New("rest server is nil ")
	}
	return server, nil
}

func (mult *MultRestServer) RegisterConfigChangeEvent() {
	marsconfig.RegisterEventHandler(common.ROOT+"."+common.REST_SERVER, func(e []*model.Event) {
		multConf, err := GetConfigParser().Parse(utils.Config)
		if err != nil {
			logger.Error("parse utils.config error: ", err.Error())
			return
		}
		for _, val := range multConf {

			mult.updateServerMiddlewares(val.Name, *val)
		}
	})
}

func (mult *MultRestServer) updateServerMiddlewares(serverName string, serverConfig internal.RestServerConf) {
	server, err := mult.GetRestServerByName(serverName)
	if err != nil {
		return
	}

	server.updateMiddlewares(serverConfig.Middlewares)
}
