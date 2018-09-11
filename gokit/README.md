总体架构

1.service
service包含根对象service以及其interface和对象的方法实现


2.transport
其他方法例如解码和编码以及解码编辑结构体和endpoint的创建


3.main
实现方法的调用，路由的设置以及端口的监听

4.otherMiddleware
    a.定义其他功能的中间件，例如流量监控，日志打印，负债均衡等等。。。文件内包含的代码若为在文件出现，请参阅gokit官方文档	https://gokit.io/examples/stringsvc.html
    b.或借鉴官方案例	https://github.com/go-kit/kit/tree/master/examples