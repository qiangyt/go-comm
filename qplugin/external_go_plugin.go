package qplugin

import (
	"fmt"
	"reflect"

	"github.com/qiangyt/go-comm/v2/qerr"
	"github.com/qiangyt/go-comm/v2/qio"
	"github.com/qiangyt/go-comm/v2/qlang"
	"github.com/spf13/afero"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type ExternalGoPluginContextT struct {
	Interpreter *interp.Interpreter

	StartFunc *reflect.Value
	StopFunc  *reflect.Value
}

type ExternalGoPluginContext = *ExternalGoPluginContextT

func NewExternalGoPluginContext() ExternalGoPluginContext {
	return &ExternalGoPluginContextT{
		Interpreter: nil,
		StartFunc:   nil,
		StopFunc:    nil,
	}
}

func resolveExternalGoPluginFunc(logger qlang.Logger, Interpreter *interp.Interpreter, funcName string) *reflect.Value {
	r, err := Interpreter.Eval(funcName)
	if err != nil {
		logger.Error(err).Msg("failed to eval " + funcName)
		return nil
	}
	if qlang.IsPrimitiveReflectValue(r) {
		logger.Error(err).Msg(funcName + " is a primitive value instead of a function")
		return nil
	}
	if r.IsNil() {
		logger.Error(err).Msg("symbol not found: " + funcName)
		return nil
	}
	if r.Kind() != reflect.Func {
		logger.Error(err).Msg(funcName + " is not a function")
		return nil
	}

	return &r
}

func (me ExternalGoPluginContext) Init(logger qlang.Logger, fs afero.Fs, codeFile string) {
	logCtx := qlang.NewLogContext(false)
	logCtx.Str("codeFile", codeFile)
	logger = logger.NewSubLogger(logCtx)

	me.Interpreter = interp.New(interp.Options{})
	if err := me.Interpreter.Use(stdlib.Symbols); err != nil {
		panic(qerr.NewSystemError(fmt.Sprintf("use stdlib failed: %s", codeFile), err))
	}

	code := qio.ReadFileTextP(fs, codeFile)
	_, err := me.Interpreter.Eval(code)
	if err != nil {
		panic(qerr.NewSystemError(fmt.Sprintf("eval code: %s", codeFile), err))
	}

	me.StartFunc = resolveExternalGoPluginFunc(logger, me.Interpreter, "plugin.PluginStart")
	me.StopFunc = resolveExternalGoPluginFunc(logger, me.Interpreter, "plugin.PluginStop")
}

func (me ExternalGoPluginContext) GetStartFunc() *reflect.Value {
	return me.StartFunc
}

func (me ExternalGoPluginContext) Start() any {
	if me.StartFunc == nil {
		return ""
	}
	return me.StartFunc.Call([]reflect.Value{})
}

func (me ExternalGoPluginContext) GetStopFunc() *reflect.Value {
	return me.StopFunc
}

func (me ExternalGoPluginContext) Stop() any {
	if me.StopFunc == nil {
		return ""
	}
	return me.StopFunc.Call([]reflect.Value{})
}
