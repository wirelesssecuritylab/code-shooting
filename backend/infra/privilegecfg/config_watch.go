package privilegecfg

import (
	"context"

	"code-shooting/infra/logger"

	"github.com/fsnotify/fsnotify"
	"go.uber.org/fx"
)

type WatchCallBackFunc func(sePath string) error

func InvokeConfigWatch(lc fx.Lifecycle, configBasePath string, callbacks ...WatchCallBackFunc) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	err = watcher.Add(configBasePath)
	if err != nil {
		logger.Warnf("error: %s", err.Error())
		return err
	}
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				for {
					select {
					case event, ok := <-watcher.Events:
						if !ok {
							return
						}
						if event.Op&fsnotify.Write == fsnotify.Write {
							logger.Infof("modified file: %s", event.Name)
							for _, callback := range callbacks {
								callback(configBasePath)
							}
						}
					case err, ok := <-watcher.Errors:
						if !ok {
							return
						}
						logger.Warnf("error: %s", err.Error())
					}
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			if watcher != nil {
				watcher.Close()
				watcher = nil
			}
			return nil
		},
	})
	return nil
}
