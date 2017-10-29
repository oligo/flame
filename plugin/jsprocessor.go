package jsplugin

import "github.com/oligo/flame/stream"

type JsProcessorHandler struct {
	plugin *JsPlugin
}

func (h *JsProcessorHandler) Prepare(out chan stream.Message) {
	h.plugin.Bind(out)
	h.plugin.Call("prepare")
}

func (h *JsProcessorHandler) Execute(in stream.Message) {
	h.plugin.Call("_execute", in.ToMap())
	// r, ok := rv.(map[string]interface{})
	// if !ok {
	// 	log.Fatalln("Message returned from processor is not of a valid type")
	// }
	// in.Replace(r)
	// h.plugin.Emit(in)
}

func (h *JsProcessorHandler) Cleanup() {
	h.plugin.Call("cleanup")
}

func NewJsProcessorHandler(filename string) JsProcessorHandler {
	plugin := LoadPlugin(filename)
	return JsProcessorHandler{plugin: &plugin}
}
