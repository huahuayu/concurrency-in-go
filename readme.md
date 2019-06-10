# concurrency in go
test
## race condition
开发并发程序很容易出错，因为开发者容易陷入一种顺序的思考，错误的认为先写的代码就会先执行，这在写并发程序的时候是不成立的。  

以下程序运行可能有三种结果：  
1. 什么也不打印，第11行先于14行执行，程序运行结束  
2. "the value is 0"被打印，第14和15行先于第11行执行  
3. "the value is 1"被打印，第14行先与11行执行，但是15行在11行后执行  

``` go
package main

import (
	"fmt"
)

func main(){
	var data int

	go func(){
		data++
	}()

	if data == 0{
		fmt.Printf("the value is %v.\n",data)
	}
}

```

加sleep的方式不可取，这样做没有根本性解决问题，race condition依然存在，而且降低系统效率  
``` go
package main

import (
	"fmt"
)

func main(){
	var data int

	go func(){
		data++
	}()

    time.Sleep(1 * time.Second) // 错误示范
    
	if data == 0{
		fmt.Printf("the value is %v.\n",data)
	}
}
```

## 检查race condition
使用`--race`参数可以检查程序是否有race condition  
``` go
go run --race main.go 
```

