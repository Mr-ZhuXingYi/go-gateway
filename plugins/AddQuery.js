
//JS插件 主函数,必须有
function main() {
    //exchange是Go 暴露给JS的对象，里面包含了当前的http request对象
    return function (exchange,key,value) {
        var q = exchange.Request.URL.Query()
        q.Add(key,value)
        exchange.Request.URL.RawQuery=q.Encode()
    }
}

function name() {
    return 'AddQuery'
}

function order() {
    return 4
}

