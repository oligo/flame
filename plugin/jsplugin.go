package jsplugin

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/oligo/flame/stream"
	"github.com/robertkrimen/otto"
)

var pluginCache = make(map[string]JsPlugin)

// JsPlugin is a wrapper around Javascript implemented handler.
// Every plugin is bounded to one otto javascript VM.
type JsPlugin struct {
	Id       string
	Name     string
	Filename string
	vm       *otto.Otto
	Handler  otto.Value
	// Emitter bound to a channel and pushes message to this channel
	Emit       EmitFunc
	extensions []string
}

// EmitFunc accepts a message and pushes the message to the plugin bounded channel
type EmitFunc func(message stream.Message)

// PluginExtension extends Otto by adding Object definitions or helper functions
// and return extension name.
type PluginExtension func(vm *otto.Otto) string

// LoadPlugin load plugin from file
func LoadPlugin(filename string, outChan chan stream.Message) JsPlugin {
	p, exists := pluginCache[filename]
	if exists {
		return p
	}

	loader := NewJsModuleFileLoader(filename)
	vm := otto.New()
	value, err := loader(vm)
	if err != nil {
		log.Fatalf("Loading plugin %s failed. Reason: %s", filename, err)
	}

	if !value.IsObject() {
		log.Fatalf("Loading plugin %s failed. Reason: %s", filename, "Js plugin should be implemented as Object.")
	}

	h := sha1.New()
	io.WriteString(h, filename)
	io.WriteString(h, fmt.Sprintf("%d", time.Now().UnixNano()))

	p = JsPlugin{Id: fmt.Sprintf("%x", h.Sum(nil)),
		vm:       vm,
		Filename: filename,
		Handler:  value,
		Name:     filepath.Base(filename)}
	p.bind(outChan)

	// add extensions to plugin
	p.InjectExtension(MessageExtension())
	p.addProxyMethods()

	pluginCache[filename] = p
	return p
}

func (plugin *JsPlugin) InjectExtension(ext PluginExtension) {
	name := ext(plugin.vm)
	if plugin.extensions == nil {
		plugin.extensions = make([]string, 0)
	}
	plugin.extensions = append(plugin.extensions, name)
	log.Printf("VM extension %s loaded for plugin %s", name, plugin.Name)
}

func (plugin *JsPlugin) bind(out chan stream.Message) {
	plugin.Emit = func(message stream.Message) {
		log.Println(out)
		out <- message
	}
}

/*
Call delegates to the calls to the object's method
value mapping from js to go:
	undefined   -> nil (FIXME?: Should be Value{})
	null        -> nil
	boolean     -> bool
	number      -> A number type (int, float32, uint64, ...)
	string      -> string
	Array       -> []interface{}
	Object      -> map[string]interface{}
*/
func (plugin *JsPlugin) Call(method string, args ...interface{}) interface{} {
	rv, err := plugin.Handler.Object().Call(method, args...)
	if err != nil {
		log.Println(err)
	}

	goValue, _ := rv.Export() // error is always nil

	return goValue
}

func (plugin *JsPlugin) addProxyMethods() {
	// execute and _execute
	const wrapperObjSrc = `
	(function(plugin){		
		if(typeof plugin === 'object'){
			plugin._execute = function(){
				if (arguments.length > 0){
					var msg = arguments[0];
					console.log("msg:" + JSON.stringify(msg));
					if (msg != undefined){
						msg = new Message(msg.payload, msg.timestamp);
						arguments[0] = msg;
					}
				}
				plugin.execute.apply(null, Array.prototype.slice.call(arguments));				
			};
		}
	})
	`
	plugin.vm.Call(wrapperObjSrc, nil, plugin.Handler.Object())

	plugin.vm.Set("emitMessage", func(call otto.FunctionCall) otto.Value {
		jsMsg := call.Argument(0).Object()
		msg := stream.NewMessage()
		if p, err := jsMsg.Get("payload"); err == nil {
			payload, _ := p.Export()
			msg.Replace(payload.(map[string]interface{}))
		}
		if t, err := jsMsg.Get("timestamp"); err == nil {
			timestamp, err := t.ToInteger()
			if err != nil {
				msg.SetTime(timestamp)
			}
		}
		plugin.Emit(msg)
		return otto.UndefinedValue()
	})

}

func (plugin *JsPlugin) String() string {
	var m = map[string]string{
		"Id":       plugin.Id,
		"name":     plugin.Name,
		"filename": plugin.Filename}
	j, _ := json.Marshal(m)
	return string(j)
}

// ModuleLoader loads nodejs-like js module
type ModuleLoader func(vm *otto.Otto) (otto.Value, error)

// NewJsModuleLoader creates a js module loader
func NewJsModuleLoader(source, pwd string) ModuleLoader {
	return func(vm *otto.Otto) (otto.Value, error) {
		source = "(function(module){var require = module.require; var exports = module.exports; var __dirname = module.__dirname;\n" +
			source + "\n})"

		jsRequire := func(call otto.FunctionCall) otto.Value {
			return otto.Value{}
		}

		jsModule, _ := vm.Object("({exports: {}})")
		jsModule.Set("require", jsRequire)
		jsModule.Set("__dirname", pwd)
		jsExports, _ := jsModule.Get("exports")

		// Run the module source
		module, err := vm.Call(source, jsExports, jsModule)
		if err != nil {
			return otto.UndefinedValue(), err
		}

		var moduleValue otto.Value
		if !module.IsUndefined() {
			moduleValue = module
			jsModule.Set("exports", moduleValue)
		} else {
			moduleValue, _ = jsModule.Get("exports")
		}

		return moduleValue, nil
	}
}

// NewJsModuleFileLoader loads a js module from file
func NewJsModuleFileLoader(filename string) ModuleLoader {
	return func(vm *otto.Otto) (otto.Value, error) {
		source, err := ioutil.ReadFile(filename)
		if err != nil {
			return otto.UndefinedValue(), err
		}

		pwd := filepath.Dir(filename)

		return NewJsModuleLoader(string(source), pwd)(vm)
	}
}
