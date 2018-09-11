
# K-V操作

---


```golang
import(
    "github.com/coreos/etcd/clientv3"
)
```



```golang
cli, err := clientv3.New(clientv3.Config{
	Endpoints: []string{etcdServer},
})
if err != nil {
	println(err)
}
// key 为 foo
// value 为 aaa
rep, err := cli.Put(ctx, "foo", "aaa")
if err != nil {
	println(err)
}
fmt.Println(rep)
reqs, err := cli.Get(ctx, "foo")
if err != nil {
	println(err)
}
fmt.Println(reqs)
```



