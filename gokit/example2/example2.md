example2

*为程序提供额外的中间件
*设置log信息中间件

    type loggingMiddleware struct {
        logger log.Logger
        next   StringService
    }
    // to show logging and instrumentation
    func (mw loggingMiddleware) Uppercase(ctx context.Context, s string) (output string, err error) {
        defer func(begin time.Time) {
            mw.logger.Log(
                "method", "uppercase",
                "input", s,
                "output", output,
                "err", err,
                "took", time.Since(begin),
            )
        }(time.Now())
        output, err = mw.next.Uppercase(ctx, s)
        return
    }


*设置作业量检测中间件

    type instrumentingMiddleware struct {
        requestCount   metrics.Counter
        requestLatency metrics.Histogram
        countResult    metrics.Histogram
        next           StringService
    }
    func (mw instrumentingMiddleware) Uppercase(ctx context.Context, s string) (output string, err error) {
        defer func(begin time.Time) {
            lvs := []string{"method", "uppercase", "error", fmt.Sprint(err != nil)}
            mw.requestCount.With(lvs...).Add(1)
            mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
        }(time.Now())
        output, err = mw.next.Uppercase(ctx, s)
        return
    }


*构建一个含中间件的service

    var svc StringService
    svc = stringService{}
    svc = loggingMiddleware{logger, svc}
    svc = instrumentingMiddleware{requestCount, requestLatency, countResult, svc}

*注：中间件对象都继承于service并实现了对应的接口方法
*效果相当于多态