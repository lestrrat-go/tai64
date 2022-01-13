tai64
=====

```go
import "github.com/lestrrat-go/tai64"
n, err := tai64.ParseNLabel([]byte(`@4000000037c219bf2ef02e94`))

t := n.Time()

n.Format(dst) // Write in format @4000000037c219bf2ef02e94
n.Write(dst) // Write in binary
```

