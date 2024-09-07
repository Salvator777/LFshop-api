package initialize

import (
	"LFshop-api/goods-web/global"
	"LFshop-api/goods-web/proto"
	"LFshop-api/goods-web/utils/otgrpc"
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// 这个方法也是拿到goodsSrvClient
// 调用了第三方库，简化了很多，而且用到了负载均衡
func InitSrvConn() {
	consulInfo := global.ServerConfig.ConsulInfo
	fmt.Println(fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name))
	goodsConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, global.ServerConfig.GoodsSrvInfo.Name),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【goods_srv服务失败】")
	}
	global.GoodsSrvClient = proto.NewGoodsClient(goodsConn)
}
