function SumNum(){
    var total = 0;
}

SumNum.prototype.execute = function(input){  
        this.total += input.get("rand");
        return
};

SumNum.prototype.cleanup = function(){
        console.log("SumNum total:" + this.total);
};

SumNum.prototype.prepare = function(cfg){
        console.log("Setting up SumNum processor...");
};


module.exports = new SumNum();