package jsplugin

import (
	"github.com/robertkrimen/otto"
)

func MessageExtension() PluginExtension {
	return func(vm *otto.Otto) string {
		const msgObjectSrc = `
		function Message(){
			var _payload = {};
			var _timestamp = new Date().getTime();
			if(arguments.length >0){
				_payload = arguments[0];
				if (arguments[1] !== null){
					_timestamp = arguments[1];
				}
			}

			Object.defineProperty(this, 'payload', {
				get: function(){return _payload;},
				enumerable: true,
				configurable:false
			});
		
			Object.defineProperty(this, 'timestamp', {
				get: function(){return _timestamp;},
				enumerable: true,
				configurable:false
			});
		
			this.put = function(key, value){
				_payload[key] = value;
			};
		
			this.get = function(key){
				return _payload[key];
			};
		}
		`
		vm.Run(msgObjectSrc)

		return "Message"
	}
}
