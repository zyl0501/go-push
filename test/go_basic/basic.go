package main

import (
	"fmt"
	"time"
	"bufio"
	"bytes"
	"strings"
	"io"
	"encoding/json"
	"encoding/binary"
	"crypto/rand"
	"math/big"
)

const (
	Sunday    = 1 + iota //0
	Monday               //1
	Tuesday              //2
	Wednesday            //3
	Thursday             //4
	Friday               //5
	Saturday             //6
	Unknown   = -1
	Unknown2   //-1
)

type oo struct {
inner
ss1 string
ss2 int
ss3 bool
}

type inner struct {
	ss4 string
}

func (i *inner) testMethod(){
	fmt.Println("testMethod is called!!!")
}

func main() {
	//mapTest()
	//enumTest()
	//chanTest()
	//bufReadWrite()
	//jsonTest()
	//appendTest()
	//binaryTest()
	//extendTest()

	//oo1 := new(oo)
	//fmt.Println("ss4无值："+oo1.ss4)
	//oo1.ss4 = "abc"
	//fmt.Println("ss4已赋值"+oo1.ss4)
	//oo1.testMethod()//继承调用
	//oo1.inner.testMethod()//继承调用 这里也可以重写

	//pointTest2()
	//randomTest()

	//go func() {
	//	fmt.Println("go")
	//	time.Sleep(5 * time.Second)
	//	fmt.Println("go End")
	//}()
	//defer fmt.Println("Exit")

	go timeTest(5)
	go timeTest(2)
	time.Sleep(20 * time.Second)
}
func mapTest() {
	var pc map[string]string
	pc = make(map[string]string)
	pc["qingdao"] = "青岛"
	pc["jinan"] = "济南"
	pc["yantai"] = "烟台"
	fmt.Println("size:", len(pc))
	for k := range pc {
		fmt.Println("k:", k)
		delete(pc, k)
	}
	qingdao, ok := pc["qingdao"]
	if ok {
		fmt.Println(qingdao)
	} else {
		fmt.Println("元素不存在")
	}
	fmt.Println("size:", len(pc))
}

func enumTest() {
	fmt.Println("Sunday=", Sunday)
	fmt.Println("Monday=", Monday)
	fmt.Println("Tuesday=", Tuesday)
	fmt.Println("Wedenesday=", Wednesday)
	fmt.Println("Thursday=", Thursday)
	fmt.Println("Friday=", Friday)
	fmt.Println("Saturday=", Saturday)
	fmt.Println("Unknown=", Unknown)
	fmt.Println("Unknown2=", Unknown2)
}

func chanTest() {
	c := make(chan int,2)
	quit := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			d := <-c
			fmt.Println(d)
			time.Sleep(time.Second)
		}
		quit <- 0
	}()
	f1(c, quit)
}

func chanTest2() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	go func() {
		time.Sleep(time.Second * 5)
		ch1 <- 10
		ch2 <- 20
	}()
	select {
	case e1 := <-ch1:
		//如果ch1通道成功读取数据，则执行该case处理语句
		fmt.Printf("1th case is selected. e1=%v", e1)
	case e2 := <-ch2:
		//如果ch2通道成功读取数据，则执行该case处理语句
		fmt.Printf("2th case is selected. e2=%v", e2)
		//default:
		//	//如果上面case都没有成功，则进入default处理流程
		//	fmt.Println("default!.")
	}
}

func chanTest3() {
	//初始化通道
	ch11 := make(chan int, 100)
	sign := make(chan int, 1)
	//给ch11通道写入数据
	for i := 0; i < 100; i++ {
		ch11 <- i
	}
	//关闭ch11通道
	//close(ch11),为了看效果先注释掉
	//单独起一个Goroutine执行select
	go func() {
		var e int
		ok := true
		i := 0
		j := 0
		for {
			temp := func() (timeout chan bool) {
				//经过大约1ms后，该接收语句会从timeout通道接收到一个新元素并赋值给ok,从而恰当地执行了针对单个操作的超时子流程，恰当地结束当前for循环
				timeout = make(chan bool, 1)
				go func() {
					time.Sleep(time.Millisecond * 3) //休息1ms
					j++
					fmt.Println("j.", j)
					timeout <- false
				}()
				i++
				fmt.Println("i.", i)
				return timeout
			}()

			select {
			case e, ok = <-ch11:
				if !ok {
					fmt.Println("End.")
					break
				}
				time.Sleep(time.Millisecond)
				fmt.Printf("ch11 -> %d\n", e)
			case ok, _ = <-temp:
				fmt.Println("Timeout.")
				break
			}
			//终止for循环
			if !ok {
				fmt.Println(" !ok.")
				sign <- 0
				break
			}
			fmt.Println("-------.")
		}
	}()
	//惯用手法，读取sign通道数据，为了等待select的Goroutine执行。
	<-sign
}

func chanTest4() {
	unbufChan := make(chan int)
	sign := make(chan byte, 2)

	go func() {
		for i := 0; i < 10; i++ {
			select {
			case unbufChan <- i:
			case unbufChan <- i + 10:
				//default:
				//	fmt.Println("default!")
			}
			time.Sleep(time.Second)
		}
		close(unbufChan)
		fmt.Println("The channel is closed.")
		sign <- 0
	}()

	go func() {
	loop:
		for {
			select {
			case e, ok := <-unbufChan:
				if !ok {
					fmt.Println("Closed channel.")
					break loop
				}
				fmt.Printf("e: %d\n", e)
				time.Sleep(2 * time.Second)
			}
		}
		sign <- 1
	}()
	<-sign
	<-sign
}

func f1(c, quit chan int) {
	x := 1;
	for {
		select {
		case c <- x:
			x += 1;
			//time.Sleep(time.Second)
		case <-quit:
			fmt.Println("quit")
			return
		default:
			fmt.Println("default")
			//return
		}
	}
}

func bufReadWrite() {
	// 1: 使用bufio.NewReader构造一个reader
	inputReadBuf := strings.NewReader("1234567890")
	reader := bufio.NewReader(inputReadBuf)

	// 2: 使用bufio.NewWriter构造一个writer
	buf := bytes.NewBuffer(make([]byte, 0))
	writer := bufio.NewWriter(buf)

	// 3: 函数Peek函数: 返回缓存的一个Slice(引用,不是拷贝)，引用缓存中前n字节数据
	// > 如果引用的数据长度小于 n，则返回一个错误信息
	// > 如果 n 大于缓存的总大小，则返回 ErrBufferFull
	// 通过Peek的返回值，可以修改缓存中的数据, 但不能修改底层io.Reader中的数据
	b, err := reader.Peek(5)
	if err != nil {
		fmt.Printf("Read data error")
		return
	}
	// 修改第一个字符
	b[0] = 'A'
	// 重新读取
	b, _ = reader.Peek(5)
	writer.Write(b)
	writer.Flush()
	fmt.Println("buf(Changed): ", buf, "\ninputReadBuf(Not Changed): ", inputReadBuf)

	// 4: Read函数, 每次读取一定量的数据, 这个由buf大小觉得, 所以我们可以循环读取数据, 直到Read返回0说明读取数据结束
	for {
		b1 := make([]byte, 3)
		n1, _ := reader.Read(b1)
		if n1 <= 0 {
			break
		}
		fmt.Println(n1, string(b1))
	}

	// 5: ReadByte和UnreadByte函数
	// ReadByte 从 b 中读出一个字节并返回, 如果 b 中无可读数据，则返回一个错误
	// UnreadByte 撤消最后一次读出的字节, 只有最后读出的字节可以被撤消, 无论任何操作，只要有内容被读出，就可以用 UnreadByte 撤消一个字节
	inputReadBuf2 := strings.NewReader("1234567890")
	reader2 := bufio.NewReader(inputReadBuf2)
	// 读一个字节
	b2, _ := reader2.ReadByte()
	fmt.Println(string(b2))
	// Unread一个字节
	reader2.UnreadByte()
	b2, _ = reader2.ReadByte()
	fmt.Println(string(b2))

	// 6: ReadRune和UnreadRune函数, 类似上面两个函数
	// ReadRune读出一个 UTF8 编码的字符并返回编码长度, 如果UTF8序列无法解码出一个正确的Unicode字符, 则只读出b中的一个字节，size 返回 1
	inputReadBuf3 := strings.NewReader("中文1234567890")
	reader3 := bufio.NewReader(inputReadBuf3)
	b3, size, _ := reader3.ReadRune()
	fmt.Println(string(b3), size)
	reader3.UnreadRune()
	b3, size, _ = reader3.ReadRune()
	fmt.Println(string(b3), size)
	// 执行UnreadRune时候, 如果之前一步不是ReadRune, 那么会报错, 看下面
	b33, _ := reader3.ReadByte()
	fmt.Println(string(b33))
	err3 := reader3.UnreadRune()
	if err3 != nil {
		fmt.Println("ERR")
	}

	// 7: 读取缓冲区中数据字节数(只有执行读才会使用到缓冲区, 否则是没有的)
	inputReadBuf4 := strings.NewReader("中文1234567890")
	reader4 := bufio.NewReader(inputReadBuf4)
	// 下面返回0, 因为还没有开始读取, 缓冲区没有数据
	fmt.Println(reader4.Buffered())
	// 下面返回strings的整体长度16(一个人中文是3长度)
	reader4.Peek(1)
	fmt.Println(reader4.Buffered())
	// 下面返回15, 因为readByte已经读取一个字节数据, 所以缓冲区还有15字节
	reader4.ReadByte()
	fmt.Println(reader4.Buffered())
	// 下面的特别有意思: 上面已经读取了一个字节, 想当于是将"中"读取了1/3, 那么如果现在使用readRune读取, 那么
	// 由于无法解析, 那么仅仅读取一个byte, 所以下面的结果很显然
	// 第一次: 无法解析, 那么返回一个byte, 所以输出的是14
	reader4.ReadRune()
	fmt.Println(reader4.Buffered())
	// 第二次读取, 还剩下"中"最后一个字节, 所以也会err, 所以输出13
	reader4.ReadRune()
	fmt.Println(reader4.Buffered())
	// 现在"中"读完了, 那么开始完整读取"文", 这个OK的, 可以解析的, 所以可以读取三字节, 那么剩下10字节
	reader4.ReadRune()
	fmt.Println(reader4.Buffered())

	// 8: ReadSlice查找 delim 并返回 delim 及其之前的所有数据的切片, 该操作会读出数据，返回的切片是已读出数据的"引用"
	// 如果 ReadSlice 在找到 delim 之前遇到错误, 则读出缓存中的所有数据并返回，同时返回遇到error（通常是 io.EOF）
	// 如果 在整个缓存中都找不到 delim，则返回 ErrBufferFull
	// 如果 ReadSlice 能找到 delim，则返回 nil
	// 注意: 因为返回的Slice数据有可能被下一次读写操作修改, 因此大多数操作应该使用 ReadBytes 或 ReadString，它们返回数据copy
	// 不推荐!
	inputReadBuf5 := strings.NewReader("中文123 4567 890")
	reader5 := bufio.NewReader(inputReadBuf5)
	for ; ; {
		b5, err := reader5.ReadSlice(' ')
		fmt.Println(string(b5))
		// 读到最后
		if err == io.EOF {
			break
		}
	}

	// 9: ReadLine 是一个低级的原始的行读取操作, 一般应该使用 ReadBytes('\n') 或 ReadString('\n')
	// ReadLine 通过调用 ReadSlice 方法实现，返回的也是"引用", 回一行数据，不包括行尾标记（\n 或 \r\n）
	// 如果 在缓存中找不到行尾标记，设置 isPrefix 为 true，表示查找未完成
	// 如果 在当前缓存中找到行尾标记，将 isPrefix 设置为 false，表示查找完成
	// 如果 ReadLine 无法获取任何数据，则返回一个错误信息（通常是 io.EOF）
	// 不推荐!
	inputReadBuf6 := strings.NewReader("中文123\n4567\n890")
	reader6 := bufio.NewReader(inputReadBuf6)
	for ; ; {
		l, p, err := reader6.ReadLine()
		fmt.Println(string(l), p, err)
		if err == io.EOF {
			break
		}
	}

	// 10: ReadBytes查找 delim 并读出 delim 及其之前的所有数据
	// 如果 ReadBytes 在找到 delim 之前遇到错误, 则返回遇到错误之前的所有数据，同时返回遇到的错误（通常是 io.EOF）
	// 如果 ReadBytes 找不到 delim 时，err != nil
	// 返回的是数据的copy, 不是引用
	inputReadBuf7 := strings.NewReader("中文123;4567;890")
	reader7 := bufio.NewReader(inputReadBuf7)
	for ; ; {
		line, err := reader7.ReadBytes(';')
		fmt.Println(string(line))
		if err != nil {
			break
		}
	}

	// 11: ReadString返回的是字符串, 不是bytes
	inputReadBuf8 := strings.NewReader("中文123;4567;890")
	reader8 := bufio.NewReader(inputReadBuf8)
	for ; ; {
		line, err := reader8.ReadString(';')
		fmt.Println(line)
		if err != nil {
			break
		}
	}

	//12: Flush函数用于提交数据, 立即更新
	// Available函数返回缓存中的可用空间
	// b10是保存数据的数组, 不是writer的缓冲区, 别搞错了
	b10 := bytes.NewBuffer(make([]byte, 30))
	// 下面会分配4096字节空间缓冲区
	writer10 := bufio.NewWriter(b10)
	writer10.WriteString("1234567890")
	// 此时没有flush, 那么输出的是"", 但是缓冲区使用了10个字节, 那么剩下4086, Buffered()返回的是缓冲区还没有提交的数据, 此处显然是10
	fmt.Println(writer10.Available(), writer10.Buffered(), b10)
	// 下面flush后, 将缓冲区的数据全部写入b10中, 缓冲区被清空, 所以缓冲区变成4096, Buffered()返回的是0, 说明数据被写入
	writer10.Flush()
	fmt.Println(writer10.Available(), writer10.Buffered(), b10)

	// 13: WriteString(...), Write(...), WriteByte(...), WriteRune(...)函数
	// 都是写数据函数
	b11 := bytes.NewBuffer(make([]byte, 1024))
	writer11 := bufio.NewWriter(b11)
	writer11.WriteString("ABC")
	writer11.WriteByte(byte('M'))
	// Rune的意思是: 代表一个字符, 那么需要一次一个字符写入
	writer11.WriteRune(rune('好'))
	writer11.WriteRune(rune('么'))
	writer11.Write([]byte("1234567890"))
	writer11.Flush()
	fmt.Println(b11)

	// 14: WriteTo函数
	inputReadBuf9 := strings.NewReader("中文1234567890")
	reader9 := bufio.NewReader(inputReadBuf9)
	b9 := bytes.NewBuffer(make([]byte, 0))
	reader9.WriteTo(b9)
	fmt.Println(b9)

	// 15: ReadFrom函数
	inputReadBuf15 := strings.NewReader("率哪来的顺丰内部了第三方吧")
	b15 := bytes.NewBuffer(make([]byte, 0))
	writer15 := bufio.NewWriter(b15)
	writer15.ReadFrom(inputReadBuf15)
	fmt.Println(b15)
}

func jsonTest() {
	type Test struct {
		A string
		B string
		C string
	}

	var example []byte = []byte(`{"A": "1", "B": "2", "C": "3"}`)

	out := Test{
		A: "default a",
		B: "default b",
		// default for C will be "", the empty value for a string
	}
	err := json.Unmarshal(example, &out) // <--
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", out)
}

func appendTest() {
	//a := []byte("hello")
	b := []byte("bbbb")
	//s := "world"
	//b = append(a, s...) // use "..." as suffice
	//fmt.Printf("%s",a)
	//fmt.Println(" ")

	append2(b)
	fmt.Printf("%s", b)
}

func append2(buf []byte) {
	buf = append(buf, "2222"...)
}

func binaryTest() {
	b11 := bytes.NewBuffer(make([]byte, 0, 2))
	writer := bufio.NewWriter(b11)
	binary.Write(writer, binary.BigEndian, uint16(23))
	writer.Flush()

	bb := b11.Bytes()
	var result uint16
	binary.Read(bytes.NewReader(bb), binary.BigEndian, &result)
	fmt.Println("result: ", result)
}

func extendTest() {
	baseMessage := BaseMessage{nil, "baseMessage1"}
	byteMessage := ByteBufMessage{&baseMessage, nil,"byteBufMessage1"}
	okMessage := OKMessage{&byteMessage, "okMessage1"}

	//baseMessage2 := BaseMessage{nil, "baseMessage2"}
	//byteMessage2 := ByteBufMessage{&baseMessage2, nil,"byteBufMessage2"}
	//okMessage2 := OKMessage{&byteMessage2, "okMessage2"}

	okMessage.BaseMessageCodec = okMessage
	okMessage.ByteBufMessageCodec = &okMessage

	//baseMessage.BaseMessageCodec = &byteMessage
	//byteMessage.ByteBufMessageCodec = &okMessage
	//okMessage.BaseMessage.BaseMessageCodec = &byteMessage2
	//okMessage.ByteBufMessage.ByteBufMessageCodec = &okMessage2

	//baseMessage2 := BaseMessage{nil, "baseMessage2"}
	//byteMessage2 := ByteBufMessage{&baseMessage2, nil,"byteBufMessage2"}
	//okMessage2 := OKMessage{&byteMessage2, "okMessage2"}
	//baseMessage2.BaseMessageCodec = byteMessage2
	//byteMessage2.ByteBufMessageCodec = okMessage2

	//okMessage.BaseMessage = baseMessage2
	//byteMessage2.ByteBufMessageCodec = okMessage
	//okMessage.BaseMessageCodec = byteMessage2
	//okMessage.BaseMessage.BaseMessageCodec = byteMessage2
	//baseMessage.BaseMessageCodec = byteMessage2

	okMessage.Encode()
}

type IMessage interface {
	Decode()
	Encode()
}

type BaseMessage struct {
	BaseMessageCodec
	tag string
}

type BaseMessageCodec interface {
	DecodeBaseMessage()
	EncodeBaseMessage()
}

func (message BaseMessage) Decode() {
	fmt.Println(message.tag +" BaseMessage Decode")
	message.DecodeBaseMessage()
}

func (message BaseMessage) Encode() {
	fmt.Println(message.tag +" BaseMessage Encode")
	message.BaseMessageCodec.EncodeBaseMessage()
}

type ByteBufMessage struct {
	*BaseMessage
	ByteBufMessageCodec

	tag string
}

type ByteBufMessageCodec interface {
	DecodeByteBufMessage()
	EncodeByteBufMessage()
}

func (message ByteBufMessage) DecodeBaseMessage() {
	fmt.Println(message.tag +" ByteBufMessage DecodeBaseMessage")
	message.ByteBufMessageCodec.DecodeByteBufMessage()
}

func (message ByteBufMessage) EncodeBaseMessage() {
	fmt.Println(message.tag +" ByteBufMessage EncodeBaseMessage")
	message.ByteBufMessageCodec.EncodeByteBufMessage()
}

type OKMessage struct {
	*ByteBufMessage
	tag string
}

func (message *OKMessage) DecodeByteBufMessage() {
	fmt.Println(message.tag + " OKMessage DecodeByteBufMessage")
}

func (message *OKMessage) EncodeByteBufMessage() {
	fmt.Println(message.tag + " OKMessage EncodeByteBufMessage", )
}

func pointTest(){
	//m := MyType{value: "something"}
	//var s Stringer
	//s = m

	m := MyType{value: "something"}
	m2 := MyType2{MyType: m}

	var s Stringer
	s = &m2 // Compile-time error again
	fmt.Println(s)

	//m := MyType{value: "something"}
	//
	//fmt.Println(Stringer(m))
}

type Stringer interface {
	String() string
}

type MyType struct {
	value string
}

type MyType2 struct {
	MyType
}

func (m *MyType) String() string { return m.value }

type INTEGER int
func (a *INTEGER) Less(b INTEGER) bool {
	return *a < b
}
func (a *INTEGER) Add(b INTEGER) {
	*a += b
}
type LessAdder interface {
	Less(INTEGER) bool
	Add(INTEGER)
}
func pointTest2() {
	var a INTEGER = 100
	var b LessAdder = &a
	b.Add(30)
	fmt.Println(a)
}

func randomTest(){
	//rand.Read
	k := make([]byte,32)
	if _,err := rand.Read(k); err != nil {
		fmt.Printf("rand.Read() error : %v \n", err)
	}
	fmt.Printf("rand.Read(): %v \n", k)

	//rand.Int
	rnd, err := rand.Int(rand.Reader,big.NewInt(1000000))
	if err != nil {
		fmt.Printf("rand.Int() error : %v \n", err)
	}
	fmt.Printf("rand.Int() : %v \n", rnd)
}

func timeTest(n int32){
	select {
	case <-time.After(time.Duration(n) * time.Second):
		fmt.Println("after %v ...", n)
	}
}