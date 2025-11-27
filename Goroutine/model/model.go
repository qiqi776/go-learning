package main

import (
	"fmt"
	"time"
)

func chanOwner() <-chan int {
	results := make(chan int, 5)

	go func() {
		defer close(results)
		for i := 0; i <= 5; i++ {
			results <- i
		}
	}()
	return results
}

func main() {
	//results := chanOwner()
	//for res := range results {
	//	fmt.Println(res)
	//}

	// ------for-select-------
	//for { // 循环：保证 Goroutine 不退出，一直活着
	//	select {
	//	// 多路复用：同时监听多个通道，哪个有消息处理哪个
	//	case <-channelA:
	//		// do something
	//	case <-channelB:
	//		// do something
	//	}
	//}

	// demo1
	//done := make(chan struct{})
	//data := []string{"a", "b", "c", "d", "e"}
	//generator := func(done <-chan struct{}, strings []string) <-chan string {
	//	out := make(chan string)
	//	go func() {
	//		defer close(out)
	//		for _, s := range strings {
	//			select {
	//			case <-done:
	//				return
	//			case out <- s:
	//
	//			}
	//		}
	//	}()
	//	return out
	//}
	//
	//stream := generator(done, data)
	//for i := 0; i < 2; i++ {
	//	fmt.Println(<-stream)
	//}
	//close(done)
	//
	//// demo2
	//// 创建一个可以取消的上下文（代替done）
	//ctx, cancel := context.WithCancel(context.Background())
	//go func() {
	//	workStream := make(chan int)
	//
	//	go func() {
	//		for i := 0; ; i++ {
	//			workStream <- i
	//			time.Sleep(time.Millisecond * 100)
	//		}
	//	}()
	//	// for-select
	//	for {
	//		select {
	//		case task := <-workStream:
	//			fmt.Println("正在处理： %d...\n", task)
	//		case <-time.After(time.Second):
	//			fmt.Println("Worker: 心跳检测 - 我还活着")
	//		case <-ctx.Done():
	//			fmt.Println("Worker: 收到停止信号，清理资源，准备下班！")
	//			return
	//		}
	//	}
	//}()
	//time.Sleep(3 * time.Second)
	//// 发出停止信号
	//fmt.Println("Main: 系统关闭，通知 Worker 下线")
	//cancel()
	//// 等一会看 Worker 的遗言
	//time.Sleep(1 * time.Second)

	//--------default分支-------
	//for {
	//	select {
	//	case req := <-requests:
	//		handle(req)
	//	default:
	//		// 当 requests 通道没数据时，select 不会阻塞，而是立刻走这里
	//		// 用于：
	//		// 1. 轮询 (Polling)
	//		// 2. 占满 CPU (Spin Lock) - 小心使用！
	//		// 3. 尝试性发送/接收
	//	}
	//}

	//----------or-goroutine------------
	//辅助函数：创建一个 n 之后关闭的通道
	//sig := func(after time.Duration) <-chan interface{} {
	//	c := make(chan interface{})
	//	go func() {
	//		defer close(c)
	//		time.Sleep(after)
	//	}()
	//	return c
	//}
	//
	//start := time.Now()
	//
	//// 监听 5 个通道：
	//// 分别在 1小时, 5分钟, 1秒, 1小时, 1分钟 后关闭
	//// 显然，那个 "1秒" 的通道会最先触发
	//<-or(
	//	sig(1*time.Hour),
	//	sig(5*time.Minute),
	//	sig(1*time.Second),
	//	sig(1*time.Hour),
	//	sig(1*time.Minute),
	//)
	//
	//fmt.Printf("Or-channel 完成，耗时: %v\n", time.Since(start))

	////------------错误处理-------------
	//done := make(chan struct{})
	//defer close(done)
	//
	//urls := []string{"https://www.google.com", "https://bad.host", "https://www.baidu.com"}
	//for result := range checkStatus(done, urls...) {
	//	if result.Error != nil {
	//		fmt.Println("Error: %v (on %s)\n", result.Error, result.Url)
	//		continue
	//	}
	//	fmt.Println("Response: %v (on %s)\n", result.Response.Status, result.Url)
	//}

	//// ----------构建流水线-------------
	//done := make(chan struct{})
	//defer close(done)
	//intStream := generator(done, 1, 2, 3, 4)
	//
	//pipeline := add(done, multiply(done, intStream, 2), 1)
	//
	//for v := range pipeline {
	//	fmt.Println(v)
	//}

	//// ------------一些便利的生成器-------------
	//done := make(chan struct{})
	//defer close(done)
	//randNum := func() interface{} { return rand.Intn(100) }
	//pipeline := take(done, repeatFn(done, randNum), 10)
	//
	//for num := range pipeline {
	//	fmt.Println(num)
	//}

	////--------------------扇入扇出------------------------
	//done := make(chan struct{})
	//defer close(done)
	//start := time.Now()
	//
	//randIntStream := Generator(done)
	//// 扇出
	//numWorkers := 10
	//workers := make([]<-chan int, numWorkers)
	//for i := 0; i < numWorkers; i++ {
	//	workers[i] = primerFinder(done, randIntStream)
	//}
	//// 扇入
	//pipeline := fanIn(done, workers...)
	//count := 0
	//for p := range pipeline {
	//	count++
	//	fmt.Printf("%d\t%d\n", count, p)
	//	if count >= 20 {
	//		break
	//	}
	//}
	//fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())

	// ------------------tee-channel-------------------------
	done := make(chan struct{})
	defer close(done)

	// 产生 5 个数据
	myStream := take(done, repeat(done, 1, 2), 5)

	// 分流！
	out1, out2 := tee(done, myStream)

	// 消费者 1 (快)
	go func() {
		for val := range out1 {
			fmt.Printf("⚡ 快速消费者: %v\n", val)
		}
	}()

	// 消费者 2 (慢)
	go func() {
		for val := range out2 {
			time.Sleep(1 * time.Second) // 模拟处理耗时
			fmt.Printf("🐢 慢速消费者: %v\n", val)
		}
	}()

	// 为了演示，主程睡一会等待
	time.Sleep(6 * time.Second)
}

//// ----------防止协程泄露--------------
//// 父子协程联动
//func doWork(done <-chan struct{}, strings <-chan string) <-chan interface{} {
//	completed := make(chan interface{})
//	go func() {
//		defer fmt.Println("doWork 安全退出")
//		defer close(completed)
//
//		for {
//			select {
//			case s := <-strings:
//				// 正常处理
//				fmt.Println(s)
//			case <-done: // 【逃生门】
//				return
//			}
//		}
//	}()
//	return completed
//}
//
//// 使用递归实现or-Channel
//func or(channels ...<-chan interface{}) <-chan interface{} {
//	switch len(channels) {
//	case 0:
//		return nil
//	case 1:
//		return channels[0]
//	}
//	orDone := make(chan interface{})
//
//	go func() {
//		defer close(orDone)
//		switch len(channels) {
//		case 2:
//			select {
//			case <-channels[0]:
//			case <-channels[1]:
//			}
//		default:
//			select {
//			case <-channels[0]:
//			case <-channels[1]:
//			case <-channels[2]:
//			case <-or(append(channels[3:], orDone)...):
//			}
//		}
//	}()
//	return orDone
//}
//
//// ---------错误处理---------
//type Result struct {
//	Error    error
//	Response *http.Response
//	Url      string
//}
//
//func checkStatus(done <-chan struct{}, urls ...string) <-chan Result {
//	results := make(chan Result)
//	go func() {
//		defer close(results)
//
//		for _, url := range urls {
//			var result Result
//			result.Url = url
//			resp, err := http.Get(url)
//
//			result.Error = err
//			result.Response = resp
//
//			select {
//			case <-done:
//				return
//			case results <- result:
//
//			}
//		}
//	}()
//	return results
//}
//
//// ----------构建流水线-----------
//
//// 1. 数据源生成器 (Generator)：把切片变成 Channel 流
//// 作用：将静态数据转化为流动数据
//func generator(done <-chan struct{}, nums ...int) <-chan int {
//	out := make(chan int)
//	go func() {
//		defer close(out)
//		for _, n := range nums {
//			select {
//			case out <- n:
//			case <-done:
//				return
//			}
//		}
//	}()
//	return out
//}
//
//// 2. 处理阶段 A：乘法器
//// 作用：从上游拿数据，乘以 2，发给下游
//func multiply(done <-chan struct{}, inputStream <-chan int, multiplier int) <-chan int {
//	out := make(chan int)
//	go func() {
//		defer close(out)
//		for i := range inputStream {
//			result := i * multiplier
//			select {
//			case out <- result:
//			case <-done:
//				return
//			}
//		}
//	}()
//	return out
//}
//
//// 3. 处理阶段 B：加法器
//// 作用：从上游拿数据，加上 1，发给下游
//func add(done <-chan struct{}, inputStream <-chan int, additive int) <-chan int {
//	out := make(chan int)
//	go func() {
//		defer close(out)
//		for i := range inputStream {
//			result := i + additive
//			select {
//			case out <- result:
//			case <-done:
//				return
//			}
//		}
//	}()
//	return out
//}
//
//// ------------一些便利的生成器-------------
//// repeat: 接收一组值，无限循环地把它们发送到通道里
//func repeat(done <-chan interface{}, values ...interface{}) interface{} {
//	valueStream := make(chan interface{})
//	go func() {
//		defer close(valueStream)
//		for _, value := range values {
//			select {
//			case <-done:
//				return
//			case valueStream <- value:
//			}
//		}
//	}()
//	return valueStream
//}
//
//// repeatFn: 接收一个函数，无限调用它，并把结果发出去
//func repeatFn(done <-chan struct{}, fn func() interface{}) <-chan interface{} {
//	valueStream := make(chan interface{})
//	go func() {
//		defer close(valueStream)
//		for {
//			select {
//			case <-done:
//				return
//			case valueStream <- fn():
//			}
//		}
//	}()
//	return valueStream
//}
//
//// take: 从 valueStream 中读取 num 个数据，然后退出
//func take(done <-chan struct{}, valueStream <-chan interface{}, num int) <-chan interface{} {
//	takeStream := make(chan interface{})
//	go func() {
//		defer close(takeStream)
//		// 只要读满 num 个，循环自动结束
//		for i := 0; i < num; i++ {
//			select {
//			case <-done:
//				return
//			case takeStream <- <-valueStream: // 这一行是从上游读，写给下游
//			}
//		}
//	}()
//	return takeStream
//}
//
//// 泛型写法
//func Repeat[T any](done <-chan struct{}, values ...T) <-chan T {
//	valueStream := make(chan T)
//	go func() {
//		defer close(valueStream)
//		for {
//			for _, value := range values {
//				select {
//				case <-done:
//					return
//				case valueStream <- value:
//				}
//			}
//		}
//	}()
//	return valueStream
//}

//// --------------扇入扇出--------------
//// 1 生成器
//func Generator(done <-chan struct{}) <-chan int {
//	out := make(chan int)
//	go func() {
//		defer close(out)
//		for {
//			select {
//			case out <- rand.Intn(100000):
//			case <-done:
//				return
//			}
//		}
//	}()
//	return out
//}
//
//// 2 耗时函数,判断素数
//func primerFinder(done <-chan struct{}, input <-chan int) <-chan int {
//	out := make(chan int)
//	go func() {
//		defer close(out)
//		for num := range input {
//			isPrime := true
//			if num <= 1 {
//				isPrime = false
//			} else {
//				for i := 2; i*i <= num; i++ {
//					if num%i == 0 {
//						isPrime = false
//						break
//					}
//				}
//			}
//			if isPrime {
//				select {
//				case out <- num:
//				case <-done:
//					return
//				}
//			}
//		}
//	}()
//	return out
//}
//
//// 3. 扇入：把多个通道合并成一个
//func fanIn(done <-chan struct{}, channels ...<-chan int) <-chan int {
//	finalStream := make(chan int)
//	var wg sync.WaitGroup
//	// 定义一个名为 multiplex 的函数，负责搬运数据
//	// 它把一个 channel 的数据搬运到 finalStream
//	multiplex := func(c <-chan int) {
//		defer wg.Done()
//		for i := range c {
//			select {
//			case finalStream <- i:
//			case <-done:
//				return
//			}
//		}
//	}
//
//	// 为每一个输入通道启动一个搬运工协程
//	wg.Add(len(channels))
//	for _, c := range channels {
//		go multiplex(c)
//	}
//
//	// 启动一个后台协程，专门等待所有搬运工下班
//	// 只有所有人都下班了，才能关闭 finalStream
//	go func() {
//		wg.Wait()
//		close(finalStream)
//	}()
//
//	return finalStream
//}
//
//// ---------------or-done-channel--------------------
//func orDone(done <-chan struct{}, c <-chan int) <-chan int {
//	valStream := make(chan int)
//	go func() {
//		defer close(valStream)
//		for {
//			select {
//			case <-done:
//				return
//			case v, ok := <-c:
//				if !ok {
//					return
//				}
//				select {
//				case valStream <- v:
//				case <-done:
//					return
//				}
//			}
//		}
//	}()
//	return valStream
//}

// ---------------tee-chennel------------------
func orDone(done <-chan struct{}, c <-chan int) <-chan int {
	valStream := make(chan int)
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				if !ok {
					return
				}
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

// 简单的 repeat 生成器
func repeat(done <-chan struct{}, values ...int) <-chan int {
	valueStream := make(chan int)
	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

// take 生成器
func take(done <-chan struct{}, valueStream <-chan int, num int) <-chan int {
	takeStream := make(chan int)
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

// ---------------- Tee 函数 (针对 int 类型修改) ----------------
func tee(done <-chan struct{}, in <-chan int) (_, _ <-chan int) {
	out1 := make(chan int)
	out2 := make(chan int)

	go func() {
		defer close(out1)
		defer close(out2)
		for val := range orDone(done, in) {
			// 局部变量覆盖，用于控制 select
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case out1 <- val:
					out1 = nil // 这一路通了，把它屏蔽掉
				case out2 <- val:
					out2 = nil // 这一路通了，把它屏蔽掉
				}
			}
		}
	}()
	return out1, out2
}
