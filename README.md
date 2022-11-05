# toolkit

## solc 使用

#### erc721

> solc mqy.sol --bin --abi --optimize -o ./output 

#### erc20
> solc token.sol --bin --abi --optimize -o ./output


# jaeger 搭建

```shell
docker run -d  --name jaeger \
-e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
-p 16686:16686 \
-p 14268:14268 \
-p 14269:14269 \
-p 9411:9411 \
-p 6831:6831/udp \
jaegertracing/all-in-one:latest
```