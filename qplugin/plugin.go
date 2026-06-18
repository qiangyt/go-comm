package qplugin

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/qiangyt/go-comm/v3/q18n"
	"github.com/qiangyt/go-comm/v3/qlang"
)

type PluginLang = string

const (
	PLUGIN_LANG_GO         = "go"
	PLUGIN_LANG_JAVASCRIPT = "javascript"
	PLUGIN_LANG_SHELL      = "shell"
)

type PluginKind = string

type Plugin interface {
	Name() string
	Kind() PluginKind
	Start(logger qlang.Logger)
	Stop(logger qlang.Logger)
	Version() (major int, minor int)
}

type PluginLoader interface {
	Namespace() string
	Plugins() map[string]Plugin
	Start(logger qlang.Logger) error
	Stop(logger qlang.Logger) error
}

func PluginId(namespace string, name string) string {
	return fmt.Sprintf("%s/%s", namespace, name)
}

func StartPlugin(namespace string, plugin Plugin, logger qlang.Logger) (err error) {
	major, minor := plugin.Version()
	ver := fmt.Sprintf("%d/%d", major, minor)
	pluginId := PluginId(namespace, plugin.Name())

	defer func() {
		if p := recover(); p != nil {
			var err2 error
			var isErr bool
			if err2, isErr = p.(error); isErr {
				err = errors.Wrap(err2, q18n.T("error.plugin.start_failed", map[string]any{
					"PluginId": pluginId,
					"Version":  ver,
					"Cause":    err2,
				}))
			} else {
				err = q18n.LocalizeError("error.plugin.start_failed", map[string]any{
					"PluginId": pluginId,
					"Version":  ver,
					"Cause":    p,
				})
			}
		}
	}()

	logCtx := qlang.NewLogContext(false)
	logCtx.Str("pluginId", pluginId).Str("version", ver)
	subLogger := logger.NewSubLogger(logCtx)

	subLogger.Info().Msg("starting")
	plugin.Start(logger)
	subLogger.Info().Msg("started")

	return err
}

func StopPlugin(namespace string, plugin Plugin, logger qlang.Logger) (err error) {
	major, minor := plugin.Version()
	ver := fmt.Sprintf("%d/%d", major, minor)
	pluginId := PluginId(namespace, plugin.Name())

	defer func() {
		if p := recover(); p != nil {
			var err2 error
			var isErr bool
			if err2, isErr = p.(error); isErr {
				err = errors.Wrap(err2, q18n.T("error.plugin.stop_failed", map[string]any{
					"PluginId": pluginId,
					"Version":  ver,
					"Cause":    err2,
				}))
			} else {
				err = q18n.LocalizeError("error.plugin.stop_failed", map[string]any{
					"PluginId": pluginId,
					"Version":  ver,
					"Cause":    p,
				})
			}
		}
	}()

	logCtx := qlang.NewLogContext(false)
	logCtx.Str("pluginId", pluginId).Str("version", ver)
	subLogger := logger.NewSubLogger(logCtx)

	subLogger.Info().Msg("stopping")
	plugin.Stop(logger)
	subLogger.Info().Msg("stopped")

	return err
}
