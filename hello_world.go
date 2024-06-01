package main

import (
	"fmt"
	fm "fmt" // alias3
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	AAA = iota
)

type Vertex struct {
	X int
	Y int
}
type Vertix struct {
	Lat, Long float64
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fm.Sprintf("%v (%v years)", p.Name, p.Age)
}
func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fm.Println(s)
		fm.Printf("pingping %v\n", i)
	}
}
func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum
}

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

func fibonacci1(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fm.Println("quit")
			return
		}
	}
}

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

func getAlbums(c *gin.Context) {
	var albums = []album{
		{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
		{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
		{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	}
	c.IndentedJSON(http.StatusOK, albums)
}

type ContextMsg struct {
	Role    string
	Content string
}

type ReplaySamples struct {
	choices []string
}

func MockGetReplySamples(userId string, charId string, dialog_context []ContextMsg) (ReplaySamples, error) {
	var replaySamples ReplaySamples
	for _, s := range dialog_context {
		replaySamples.choices = append(replaySamples.choices, s.Content)
	}
	return replaySamples, nil
}

func GenerateUserReplaySample(c *gin.Context) {
	fmt.Println("hehahhah")
	userId := c.PostForm("user_id")
	charId := c.PostForm("character_id")

	if len(userId) < 1 || len(charId) < 1 {
		c.JSON(400, gin.H{"error": "parameter error"})
		return
	}

	contextMsgs := mockContextTracer(userId, charId)
	fmt.Println(len(contextMsgs))
	for i := 0; i < len(contextMsgs); i++ {
		fmt.Println(contextMsgs[i].Content)
	}
	llmResult, err := MockGetReplySamples(userId, charId, contextMsgs)
	for i := 0; i < len(llmResult.choices); i++ {
		fmt.Println(llmResult.choices[i])
	}
	fmt.Println("api run")
	fmt.Println(err)

	c.IndentedJSON(http.StatusOK, llmResult.choices)

}

// 模拟拿到上下文的函数
func mockContextTracer(userId string, charId string) []ContextMsg {

	sampleResult := make([]ContextMsg, 0)

	sampleResult = append(sampleResult, ContextMsg{
		Role:    "user",
		Content: "Hello there",
	})

	sampleResult = append(sampleResult, ContextMsg{
		Role:    "char",
		Content: "Hi",
	})
	sampleResult = append(sampleResult, ContextMsg{
		Role:    "user",
		Content: "Nice to see you!",
	})
	sampleResult = append(sampleResult, ContextMsg{
		Role:    "user",
		Content: "How's going?",
	})

	return sampleResult
}
func main() {
	fm.Println("hello world")
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.POST("/sampleReplys", GenerateUserReplaySample)
	router.Run("localhost:8080")

	//Multi_thread
	// go say("pingping")
	// say("pong")
	// s := []int{7, 2, 4, -9, 8, 0}
	// c := make(chan int)
	// go sum(s[:len(s)/2], c)
	// go sum(s[len(s)/2:], c)
	// x, y := <-c, <-c

	// fm.Println(x, y, x+y)
	// chan1 := make(chan int, 10)
	// go fibonacci(cap(chan1), chan1)
	// for i := range chan1 {
	// 	fmt.Println(i)
	// }

	// chan2 := make(chan int)
	// quit := make(chan int)
	// go func() {
	// 	for i := 0; i < 10; i++ {
	// 		fm.Println(<-chan2)
	// 	}
	// 	quit <- 0
	// }()
	// fibonacci1(chan2, quit)

	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)

	for {
		select {
		case <-tick:
			fm.Println("Thick.")
		case <-boom:
			fm.Println("BOOM!")
			return
		default:
			fmt.Println(".")
			time.Sleep(50 * time.Millisecond)
		}
	}

	//Function related
	// defer fm.Println("hahahahah")
	// fm.Println(AAA)
	// switch os := runtime.GOOS; os {
	// case "darwin":
	// 	fm.Println("OS X.")
	// case "linux":
	// 	fm.Println("Linux.")
	// default:
	// 	fm.Println("%s.\n", os)
	// }

	// today := time.Now().Weekday()
	// switch time.Saturday {
	// case today + 0:
	// 	fm.Println("Today.")
	// case today + 1:
	// 	fm.Println("Tomorrow")
	// default:
	// 	fm.Println("Too far")
	// }
	// primes := [8]int{2, 3, 4, 55, 6, 0, 7, 8}
	// var s []int = primes[1:3]
	// fm.Println(s)
	// s = primes[1:]
	// fm.Println(s)
	// v := Vertex{1, 2}
	// p := &v
	// p.X = 1e9
	// fm.Println(v)
	// aArray := make([]int, 0, 6)
	// fm.Println(aArray)
	// fm.Println(cap(aArray))
	// for i, v := range primes {
	// 	fm.Printf("2**%d = %d\n", i, v)
	// }
	// var mMap map[string]Vertix
	// mMap = make(map[string]Vertix)
	// mMap["haha lab"] = Vertix{
	// 	123, 45,
	// }
	// fm.Println(mMap["haha lab"])
	// val1, ok := mMap["haha"]
	// fm.Println("The value:", val1, "Present?", ok)

	// pos, neg := adder(), adder()
	// for i := 0; i < 10; i++ {
	// 	fm.Println(pos(i), neg(-2*i))
	// }

	// aC := Person{"Arthur Curry", 40}
	// zB := Person{"Zaphod Beeblebrox", 9001}
	// fm.Println(aC, zB)

	//Basic
	// var a int = 5
	// var d int16 = 34
	// var e int32
	// e = int32(d)
	// c := int(a)

	// for i := 0; i < 10; i++ {
	// 	fa := rand.Int()
	// 	fm.Printf("%d / ", fa)
	// }
	// var str string = "This is an example of a string, THE, The"
	// fm.Printf("%t\n", strings.HasPrefix(str, "Th"))
	// s1 := strings.Fields(str)
	// for _, val := range s1 {
	// 	fmt.Printf("%s - ", val)
	// }
	// fm.Println()
	// var sum = 1
	// for sum < 1000 {
	// 	sum += sum
	// }
	// if x < 0{

	// }
	// fm.Println(sum)
	// var i1 = 5
	// fm.Printf("An integer: %d, it's location in memory: %p\n", i1, &i1)
	// fm.Println(c)
	// fm.Println(a)
	// fm.Printf("32 bit int is: %d\n", d)
	// fm.Printf("16 bit int is: %d\n", e)

}

func Uint8FromInt(n int) (uint8, error) {
	if 0 <= n && n <= math.MaxUint8 {
		return uint8(n), nil
	}
	return 0, fm.Errorf("%d is out of the uint8 range", n)
}
