
module.exports = {
    execute: function(input){  
        console.log("msg received at: " + JSON.stringify(input.timestamp));  
        if(typeof input === 'object'){
            input.put('printer', "printer-1234");
        }
        emitMessage(input);
    },

    cleanup: function(){
        console.log("Cleaning printer processor");
    },

    prepare: function(cfg){
        console.log("Setting up printer processor...");
    }
};