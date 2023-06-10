package test

import (
	"context"
	"go.uber.org/fx"
)

func StartFxApp(app *fx.App) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.StartTimeout())
	defer cancel()
	return app.Start(ctx)
}

func StopFxApp(app *fx.App) error {
	ctx, cancel := context.WithTimeout(context.Background(), app.StopTimeout())
	defer cancel()
	return app.Stop(ctx)
}
