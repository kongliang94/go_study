package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"simple_lb_2/config"
)

/*
轮询（Round Robin）——均匀地分发流量负载，假设所有后端服务都具有同样的处理能力
*/
const (
	Attempts int = iota
	Retry
)

// Backend 定义一个结构体保存后端服务器状态信息
type Backend struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	ReverseProxy *httputil.ReverseProxy
}

// SetAlive 设置服务可用
func (b *Backend) SetAlive(alive bool) {
	// 不同的 goroutine 会同时访问 Backend,使用 RWMutex 来串行化对 Alive 的访问操作
	b.mux.Lock()
	b.Alive = alive
	b.mux.Unlock()
}

// IsAlive 服务可用返回true
func (b *Backend) IsAlive() (alive bool) {
	b.mux.RLock()
	alive = b.Alive
	b.mux.RUnlock()
	return
}

// ServerPool 要一种方式来跟踪所有后端，以及一个计算器变量
type ServerPool struct {
	backends []*Backend
	current  uint64
}


// NextIndex 因为有很多客户端连接到负载均衡器，所以发生竟态条件是不可避免的。
// 为了防止这种情况，我们需要使用 mutex 给 ServerPool 加锁。但这样做对性能会有影响，更何况我们并不是真想要给 ServerPool 加锁，我们只是想要更新计数器。
// 最理想的解决方案是使用原子操作，Go 语言的 atomic 包为此提供了很好的支持
func (s *ServerPool) NextIndex() int {
	// 通过原子操作递增 current 的值，并通过对 slice 的长度取模来获得当前索引值。所以，返回值总是介于 0 和 slice 的长度之间，毕竟我们想要的是索引值，而不是总的计数值
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// GetNextPeer 获取下一个可用服务器
func (s *ServerPool) GetNextPeer() *Backend {
	// 遍历后端列表，找到可用服务器
	next := s.NextIndex()
	// log.Println("next ", next)
	// 从next开始遍历
	l := len(s.backends) + next
	// log.Println("l ", l)
	for i := next; i < l; i++ {
		// 通过取模运算获取索引
		idx := i % len(s.backends)
		// log.Println("idx ", idx)
		//如果找到一个可用服务器
		if s.backends[idx].IsAlive() {
			if i != next {
				// 标记当前可用服务器
				atomic.StoreUint64(&s.current, uint64(idx))
			}

			return s.backends[idx]
		}

	}
	return nil
}

// HealthCheck 被动模式，遍历所有服务并并标记可用状态
func (s *ServerPool) HealthCheck() {
	for _, b := range s.backends {
		status := "up"

		alive := isBackendalive(b.URL)

		b.SetAlive(alive)
		if !alive {
			status = "down"
		}

		log.Printf("%s[%s]\n", b.URL, status)
	}
}

// AddBackend 添加服务到ServerPool
func (s *ServerPool) AddBackend(backend *Backend) {
	s.backends = append(s.backends, backend)
}

// MarkBackendStatus 标记服务状态
func (s *ServerPool) MarkBackendStatus(backendUrl *url.URL, alive bool) {
	for _, b := range s.backends {
		if b.URL.String() == backendUrl.String() {
			b.SetAlive(alive)
			break
		}
	}
}

// GetRetryFromContext 返回重试次数
func GetRetryFromContext(r *http.Request) int {
	if retry, ok := r.Context().Value(Retry).(int); ok {
		return retry
	}
	return 0
}

// GetAttemptsFromContext 返回尝试次数
func GetAttemptsFromContext(r *http.Request) int {
	if attempts, ok := r.Context().Value(Attempts).(int); ok {
		return attempts
	}
	return 1
}

// lb对接收到的请求 进行负载均衡
func lb(w http.ResponseWriter, r *http.Request) {

	// 限制重试次数最大为3
	attempts := GetAttemptsFromContext(r)
	if attempts > 3 {
		log.Printf("%s(%s) Max attempts reached, terminating\n", r.RemoteAddr, r.URL.Path)
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
		return
	}

	peer := serverPool.GetNextPeer()
	log.Println("下一个peer ", peer)
	if peer != nil {
		peer.ReverseProxy.ServeHTTP(w, r)
		return
	}
	http.Error(w, "服务不可用", http.StatusServiceUnavailable)
}

// isBackendalive 被动模式 检测服务可用性，建立tcp连接判断后台服务是否可用
func isBackendalive(u *url.URL) bool {
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", u.Host, timeout)
	if err != nil {
		log.Println("Site unreachable, error: ", err)
		return false
	}
	// 执行完操作后要关闭连接，避免给服务器造成额外的负担，否则服务器会一直维护连接
	_ = conn.Close()
	return true
}

// 每20秒检测一次执行一次健康检测，额外开启一个goroutine去执行此方法
func healthCheck() {
	t := time.NewTicker(time.Second * 20)
	for {
		select {
		// <-t.C 每 20 秒返回一个值，select 会检测到这个事件。在没有 default case 的情况下，select 会一直等待，直到有满足条件的 case 被执行
		case <-t.C:
			log.Println("Starting health check...")
			serverPool.HealthCheck()
			log.Println("Health check completed")
		}
	}
}

var serverPool ServerPool

// 测试simplelb.exe
func main() {
	// 从配置文件读取端口
	port := config.RuntimeViper.GetInt("server.port")

	// 从配置文件读取代理服务
	servers := config.RuntimeViper.GetStringSlice("server.proxy_pass")

	for _, tok := range servers {
		serverURL, err := url.Parse(tok)
		if err != nil {
			log.Fatal(err)
		}

		proxy := httputil.NewSingleHostReverseProxy(serverURL)
		//在处理当前请求时，如果发现当前的后端没有响应，就把它标记为已宕机,
		//在发生错误时，ReverseProxy 会触发 ErrorHandler 回调函数，我们可以利用它来检查故障
		proxy.ErrorHandler = func(writer http.ResponseWriter, request *http.Request, e error) {
			log.Printf("[%s] %s\n", serverURL.Host, e.Error())
			// 从context中获取重试次数
			retries := GetRetryFromContext(request)
			if retries < 3 {
				select {
				case <-time.After(10 * time.Millisecond):
					ctx := context.WithValue(request.Context(), Retry, retries+1)
					proxy.ServeHTTP(writer, request.WithContext(ctx))
				}
				return
			}

			// 3次重试后该服务设置为宕机
			serverPool.MarkBackendStatus(serverURL, false)

			// 同一个请求在尝试了几次后仍然失败，增加计数
			attempts := GetAttemptsFromContext(request)
			log.Printf("%s(%s) Attempting retry %d\n", request.RemoteAddr, request.URL.Path, attempts)
			ctx := context.WithValue(request.Context(), Attempts, attempts+1)
			// 通过lb选择一个新的后端来处理请求
			lb(writer, request.WithContext(ctx))
		}

		serverPool.AddBackend(&Backend{
			URL:          serverURL,
			Alive:        true,
			ReverseProxy: proxy,
		})
		log.Printf("Configured server: %s\n", serverURL)
	}

	//创建一个http server，初始化服务器，并添加处理器
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(lb), // HandlerFunc 传给 http 服务器
	}

	// 开启健康检测
	go healthCheck()

	log.Printf("Load Balancer started at :%d\n", port)
	// 监听服务
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}