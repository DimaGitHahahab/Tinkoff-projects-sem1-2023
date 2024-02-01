## Pipeline

`ExecutePipeline`запускает пайплайн,
состоящего из `Stage`.
Каждый `Stage` - функция удовлетворяющая интерфейсу

```go
func Stage1(in In) (out Out) {
    out := make(chan any)
    go func() {
        // do something..
    }()
    return out
}
```