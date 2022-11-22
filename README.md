# GEE-X

7天用Go从零实现系列，跟随教程 https://geektutu.com/post/gee.html

## Web框架 - Gee

与原文相比，在部分 API 上细微区别，并且考虑到模板在现代应用中不常用，没有实现模板机制。

中间件的 Next() 很像是基于生成器的协程，值得学习，伪代码基本原理如下：

```
funcs: List[Func]
idx = -1  // 上一次被执行的函数 idx

// 让出执行权，执行接下来的函数
func Next() {
    idx += 1  // 接下来要执行的函数
    // 按照顺序执行
    while idx < len(funcs) {
        funcs[idx]() // 在该函数内调用 Next() 可以实现等待其他函数执行完毕
        idx += 1
    }
}

func A() {
    // DO SOMETHING BEFORE，此时 idx 等于我自己的 idx
    Next() // 等待其他函数执行
    // DO SOMETHING AFTER
}
```

