package jsplugin

import "github.com/oligo/flame/stream"

type JsCollectorHandler struct {
	plugin *JsPlugin
}

func (h *JsCollectorHandler) Prepare(out chan stream.Message) {
	h.plugin.Bind(out)
	h.plugin.Call("prepare")
}

func (h *JsCollectorHandler) Execute() {
	h.plugin.Call("_execute")
	// if rv != nil {
	// 	msg := stream.NewMessage()
	// 	msg.Replace(rv.(map[string]interface{}))
	// 	log.Printf("collector returns %s", msg.String())
	// 	h.plugin.Emitter(msg)
	// }
}

func (h *JsCollectorHandler) Cleanup() {
	h.plugin.Call("cleanup")
}

func NewJsCollectorHandler(filename string) JsCollectorHandler {
	plugin := LoadPlugin(filename)
	return JsCollectorHandler{plugin: &plugin}
}
