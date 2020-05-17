package middleware

import "context"

// HandleFunc 中间件真正执行业务的函数
// 参数
//   ctx: Context, 用于传值或超时控制
//   req: 请求的参数
// 返回值
//   resp: 处理结果
//   err: error
type HandleFunc func(ctx context.Context, req interface{}) (resp interface{}, err error)

// Middleware 中间件, 使两个业务处理函数关联
// 参数
//   HandleFunc 中间件业务处理函数
// 返回值
//   HandleFunc 中间件业务处理函数
type Middleware func(HandleFunc) HandleFunc

// Chain 把所有中间件串连起来, 组成中间件链, 按传入的顺序遍历中间件
// 参数
//   outer: 是最外层的中间件
//   others 其它中间件
// 返回值
//   Middleware: 串联起来的中间件链
func Chain(outer Middleware, others ...Middleware) Middleware {
	return func(next HandleFunc) HandleFunc {
		// 为了按顺序遍历中间件, 把others最前面的HandleFunc传给outer, 需要反转遍历；
		// 如果不反转遍历, 那么最后next将是others最后一个元素, 传给outer将倒序遍历。
		for i := len(others) - 1; i >= 0; i-- {
			next = others[i](next) // Middleware的作用，使两个HandleFunc关联
		}
		return outer(next)
	}
}
