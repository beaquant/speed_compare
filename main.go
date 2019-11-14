package main

import (
	"fmt"
	"github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/fcoin"
	"github.com/valyala/fasthttp"
	"math"
	"net/http"
	"strings"
	"time"
)

type State struct {
	Max   time.Duration
	Min   time.Duration
	Ave   time.Duration
	Total time.Duration
	Count int64
}

func (s State) String() string {
	return fmt.Sprintf("Max:%s, Min:%s, Ave:%s, Total:%s, Count:%d", s.Max.String(), s.Min.String(), s.Ave.String(), s.Total.String(), s.Count)
}

var (
	fastHttpState = State{
		Max:   math.MinInt64,
		Min:   math.MaxInt64,
		Ave:   0,
		Total: 0,
		Count: 0,
	}

	httpState = State{
		Max:   math.MinInt64,
		Min:   math.MaxInt64,
		Ave:   0,
		Total: 0,
		Count: 0,
	}
)

func updateTime(state *State, cost time.Duration) {
	if state.Max < cost {
		state.Max = cost
	}
	if state.Min > cost {
		state.Min = cost
	}
	state.Count++
	state.Total += cost
	state.Ave = state.Total / time.Duration(state.Count)
}

func GetDepthFastHttp() {
	fmt.Println("start to test fasthttp")
	uri := fmt.Sprintf("market/depth/L20/%s", strings.ToLower(goex.BTC_USDT.ToSymbol("")))

	n := 1000
	for i := 0; i < n; i++ {
		time.Sleep(200 * time.Millisecond)
		start := time.Now()
		status, _, err := fasthttp.Get(nil, "https://api.fcoin.com/v2/"+uri)
		if err != nil {
			fmt.Println("请求失败:", err.Error())
			continue
		}

		if status != fasthttp.StatusOK {
			fmt.Println("请求没有成功:", status)
			continue
		}

		cost := time.Since(start)
		updateTime(&fastHttpState, cost)
		//fmt.Println(string(resp))
	}
}

func GetDepthHttp() {
	fmt.Println("start to test net/http")
	fcClient := fcoin.NewFCoin(http.DefaultClient, "", "")
	n := 1000
	for i := 0; i < n; i++ {
		time.Sleep(200 * time.Millisecond)
		start := time.Now()
		_, err := fcClient.GetDepth(20, goex.BTC_USDT)
		if err != nil {
			fmt.Println("请求失败:", err.Error())
			continue
		}
		cost := time.Since(start)
		updateTime(&httpState, cost)
	}
}

func main() {
	GetDepthFastHttp()
	GetDepthHttp()
	fmt.Println("\nfasthttp:\n", fastHttpState.String())
	fmt.Println("\nnet/http:\n", httpState.String())
}
