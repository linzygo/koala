package middleware

import (
	"context"
	"fmt"
	"testing"
)

func TestMiddleWare(t *testing.T) {
	middleware1 := func(next HandleFunc) HandleFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			fmt.Println("middleware1 start")
			// n := rand.Intn(20)
			// if n < 5 {
			// 	err = fmt.Errorf("middleware1 this is request is not allow")
			// 	return
			// }
			resp, err = next(ctx, req)
			if err != nil {
				return
			}
			fmt.Println("middleware1 end")
			return
		}
	}

	middleware2 := func(next HandleFunc) HandleFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			fmt.Println("middleware2 start")
			resp, err = next(ctx, req)
			if err != nil {
				return
			}
			fmt.Println("middleware2 end")
			return
		}
	}

	outer := func(next HandleFunc) HandleFunc {
		return func(ctx context.Context, req interface{}) (resp interface{}, err error) {
			fmt.Println("outer start")
			resp, err = next(ctx, req)
			if err != nil {
				return
			}
			fmt.Println("outer end")
			return
		}
	}

	process := func(ctx context.Context, req interface{}) (resp interface{}, err error) {
		fmt.Println("process start")
		fmt.Println("process end")
		return
	}

	chain := Chain(outer, middleware1, middleware2)
	proc := chain(process)
	resp, err := proc(context.Background(), "test")
	fmt.Printf("resp:%#v, err:%v\n", resp, err)
}
