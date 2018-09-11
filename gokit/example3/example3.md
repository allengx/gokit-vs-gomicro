example3

*中间件的嵌套过程，更近似洋葱皮

    svc = proxyingMiddleware(context.Background(), *proxy, logger)(svc)


*中间件Middleware的方法为endpoint添加一层一层的功能 返回类似都是ServiceMiddleware 的func

    func loggingMiddleware(logger log.Logger) ServiceMiddleware {
        return func(next StringService) StringService {
            return logmw{logger, next}
        }
    }
    func proxyingMiddleware(ctx context.Context, instances string, logger log.Logger) ServiceMiddleware {
        //...
        //...
        //...
        return func(next StringService) StringService {
            return proxymw{ctx, next, retry}
        }
    }
    func instrumentingMiddleware(
        requestCount metrics.Counter,
        requestLatency metrics.Histogram,
        countResult metrics.Histogram,
    ) ServiceMiddleware {
        return func(next StringService) StringService {
            return instrmw{requestCount, requestLatency, countResult, next}
        }
    }


*调用嵌套也变得更加易懂

    var svc StringService
    svc = stringService{}
    svc = proxyingMiddleware(context.Background(), *proxy, logger)(svc)
    svc = loggingMiddleware(logger)(svc)
    svc = instrumentingMiddleware(requestCount, requestLatency, countResult)(svc)


