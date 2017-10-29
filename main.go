package main

import (
	"github.com/oligo/flame/plugin"

	"github.com/oligo/flame/stream"
)

func main() {
	collectHandler := jsplugin.NewJsCollectorHandler("./collector.js")
	printer := jsplugin.NewJsProcessorHandler("./printer.js")
	collector := stream.NewCollector("textCollector", stream.Config{}, &collectHandler)
	p1 := stream.NewProcessor("printer", &printer,
		stream.NewIndexedChanDemux(4, stream.NewGroupDemuxIndex("char").GroupIndex))
	raw := collector.Execute()
	<-p1.Execute(raw)

}
