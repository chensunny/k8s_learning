### 概要设计


基于类snowflake方案，设计24位 （36/32进制）的 UUID   
注：32进制对移位处理和性能上有一定的优势


**具体ID格式**
>[g-z random] + random(12) + date(8) + increase(3 random)


#### 关于date

* 使用 Unix 时间戳 
		
		data 8位36进制 对应 毫秒可以支持到 `2059/5/26 1:38:27`

* increase(3 random) 

	 	支持单节点每毫秒 46656

预留覆盖趋势递增的场景:

> [ zxy random] date(8) + increase(3 random)  + random(12) 



### 方案
> random(12) = random(5) + processUnique(7)  

描述：  


启动的时候创建一个N位的进程标示（Process Mark），通过每次都会构造 random（12-N）位的随机数。

```
方案 2启动的进程启动的时候，可以通过调用标准库的 `crypto/rand` 底层是通过    
/dev/urandom 获取内核提供的2个真随机数，其中一个作为进程标示位 (Process Mark），另一个作为一般库的 Seed。这样可以解决性能问题的同时，减少重复概率。
通过 processUnique(N) + random(12-N) 生成12位的随机数。
```


1万进程/线程，进程标示不冲突的概率是 0.998724720285324。
进程标示超过10个重复的概率是: 很低（不会算。。）

当进程标示冲突的场景，random 5位出现冲突的概率同上 

* 如果10个相同 `进程标示` 的进程/线程 一秒一次运行1年不冲突 概率 约为 0.999999999137455
* 如果10个相同 `进程标示` 的进程/线程 1毫秒一次运行1年不冲突 概率约为 0.999999137455695
   
   
   
### quick start


```
go get  github.com/golang-common/uuid

id,err :=  uuid.GenerateC24()
```



### Benchmark


```

goos: darwin
goarch: amd64
pkg: github.com/golang-common/uuid
BenchmarkGenerate-4   	 3000000	       380 ns/op


```


时间主要消耗在:

* 获取当前时间 time.Now()
* 10进制转化36进制 strconv.FormatInt() 



![](https://meitu-test.oss-cn-beijing.aliyuncs.com/WeChat2fe1f47fd5edde057993a659574e4d71.png)
   
   
   