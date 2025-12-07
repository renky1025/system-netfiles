package middleware

import (
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ThrottledReader 限速读取器
// 使用令牌桶算法控制读取速率
type ThrottledReader struct {
	reader  io.Reader
	limiter *rate.Limiter
}

// NewThrottledReader 创建限速读取器
// bytesPerSecond: 每秒允许读取的字节数, 0表示无限制
func NewThrottledReader(reader io.Reader, bytesPerSecond int64) io.Reader {
	if bytesPerSecond <= 0 {
		return reader // 无限制,直接返回原reader
	}

	// 创建令牌桶
	// Limit: 每秒产生的令牌数 (即bytes/s)
	// Burst: 桶容量,设为速率的1/10或最小64KB,允许短暂突发
	burst := int(bytesPerSecond / 10)
	if burst < 65536 {
		burst = 65536 // 最小64KB burst
	}

	limiter := rate.NewLimiter(rate.Limit(bytesPerSecond), burst)

	return &ThrottledReader{
		reader:  reader,
		limiter: limiter,
	}
}

// Read 实现io.Reader接口,带速率限制
func (t *ThrottledReader) Read(p []byte) (n int, err error) {
	n, err = t.reader.Read(p)
	if n <= 0 {
		return n, err
	}

	// 等待获取足够的令牌
	// 使用Reserve而非Wait,避免阻塞过久
	reservation := t.limiter.ReserveN(time.Now(), n)
	if !reservation.OK() {
		// 请求的令牌数超过桶容量,分批处理
		time.Sleep(time.Duration(n) * time.Second / time.Duration(t.limiter.Limit()))
	} else {
		delay := reservation.Delay()
		if delay > 0 {
			time.Sleep(delay)
		}
	}

	return n, err
}

// ThrottledResponseWriter 限速响应写入器
type ThrottledResponseWriter struct {
	gin.ResponseWriter
	limiter *rate.Limiter
}

// NewThrottledResponseWriter 创建限速响应写入器
func NewThrottledResponseWriter(w gin.ResponseWriter, bytesPerSecond int64) *ThrottledResponseWriter {
	if bytesPerSecond <= 0 {
		return &ThrottledResponseWriter{ResponseWriter: w, limiter: nil}
	}

	burst := int(bytesPerSecond / 10)
	if burst < 65536 {
		burst = 65536
	}

	return &ThrottledResponseWriter{
		ResponseWriter: w,
		limiter:        rate.NewLimiter(rate.Limit(bytesPerSecond), burst),
	}
}

// Write 实现io.Writer接口,带速率限制
func (w *ThrottledResponseWriter) Write(data []byte) (int, error) {
	if w.limiter == nil {
		return w.ResponseWriter.Write(data)
	}

	written := 0
	for written < len(data) {
		// 每次最多写入burst大小的数据块
		chunkSize := len(data) - written
		burst := w.limiter.Burst()
		if chunkSize > burst {
			chunkSize = burst
		}

		// 等待令牌
		reservation := w.limiter.ReserveN(time.Now(), chunkSize)
		if reservation.OK() {
			delay := reservation.Delay()
			if delay > 0 {
				time.Sleep(delay)
			}
		}

		n, err := w.ResponseWriter.Write(data[written : written+chunkSize])
		written += n
		if err != nil {
			return written, err
		}
	}

	return written, nil
}

// DownloadRateLimitMiddleware 下载限速中间件
// getRateLimit: 获取当前请求的速率限制函数
func DownloadRateLimitMiddleware(getRateLimit func(*gin.Context) int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		rateLimit := getRateLimit(c)

		if rateLimit > 0 {
			// 替换ResponseWriter为限速版本
			throttledWriter := NewThrottledResponseWriter(c.Writer, rateLimit)
			c.Writer = throttledWriter
		}

		c.Next()
	}
}

// ServeFileWithRateLimit 带限速的文件下载
func ServeFileWithRateLimit(c *gin.Context, filePath string, fileName string, rateLimit int64) {
	// 设置下载响应头
	c.Header("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	c.Header("Content-Type", "application/octet-stream")

	if rateLimit <= 0 {
		// 无限速,直接使用gin的文件服务
		c.File(filePath)
		return
	}

	// 使用限速Writer
	throttledWriter := NewThrottledResponseWriter(c.Writer, rateLimit)
	http.ServeFile(throttledWriter, c.Request, filePath)
}
