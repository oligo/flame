
function Collector(){

}

Collector.prototype.prepare = function(){
    console.log("Setting up collector...");    
}

Collector.prototype.execute = function () {
    var msg = new Message({rand: Math.floor(Math.random() * 1000)});
    console.log(JSON.stringify(msg));
    emitMessage(msg);
}

Collector.prototype.cleanup = function(){
    console.log("cleaning collector...");
    return;
}

module.exports = new Collector();