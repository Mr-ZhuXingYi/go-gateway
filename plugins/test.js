

function main() {
   return function (exchange,key,value) {
       exchange.Request.Header.Set(key,value)
   }
}

function name() {
     return 'Test'
}
