package main

import (
	"log"

	"github.com/oligo/flame/plugin"
)

func loadPlugin() jsplugin.JsPlugin {
	plugin := jsplugin.LoadPlugin("./collector.js")

	plugin.Call("prepare", "{parallism: 12, name: 'js-collector'}")

	v := plugin.Call("execute")
	plugin.Call("cleanup")

	log.Println(v)
	return plugin
}
