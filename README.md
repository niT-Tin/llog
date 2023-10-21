# llog

## Install

```bash

```
## Usage

### Simple usage

```go
llog.Ghelper.Debug("key", "value")
llog.Ghelper.Info("key", "value")
llog.Ghelper.Warn("key", "value")
llog.Ghelper.Error("key", "value")
llog.Ghelper.Fatal("key", "value")
```

### Stdlogger

stdlogger是llog最底层默认使用的Log接口的实现，如果不想使用任何额外配置，可以使用默认的DefaultLogger。

```go
package main

import "github.com/niT-Tin/llog"

func main() {
    llog.DefaultLogger.Log(llog.Debug, "key", "value")
    llog.DefaultLogger.Log(llog.Info, "key", "value")
    llog.DefaultLogger.Log(llog.Warn, "key", "value")
    llog.DefaultLogger.Log(llog.Error, "key", "value")
    // 创建stdlogger

	colorMap := make(map[llog.Level][]llog.Color, 1)
    // 设置Info等级的fg与bg颜色
	colorMap[llog.Info] = []llog.Color{llog.FgYellow, llog.BgGreen}
	logger := llog.NewStdLogger(
		llog.WithStdWriter(os.Stdout),
        // 默认开启颜色，需要关闭可设置为false
		llog.WithStdColored(true),
		llog.WithStdTimeZone(time.UTC),
		llog.WithStdTimeFormat("2006-01-02 15:04:05.000"),
		llog.WithStdColors(colorMap),
	)
	logger.Log(llog.Info, "key", "value")
}

```

### Global

llog提供了一个全局logger以及一个全局helper，llog更推荐使用helper已为其相对方便使用。

```go
package main

func main() {
    // Glogger
	llog.Glogger.Log(llog.Debug, "key", "value")
	llog.Glogger.Log(llog.Info, "key", "value")
	llog.Glogger.Log(llog.Warn, "key", "value")
	llog.Glogger.Log(llog.Error, "key", "value")
	llog.Glogger.Log(llog.Fatal, "key", "value")
    // Ghelper
	llog.Ghelper.Debug("key", "value")
	llog.Ghelper.Info("key", "value")
	llog.Ghelper.Warn("key", "value")
	llog.Ghelper.Error("key", "value")
	llog.Ghelper.Fatal("key", "value")

    // 带w后缀的函数不会打印MessageKey的内容，该内容可以在创建helper时设置
	llog.Ghelper.Debugw("key1", "value1", "key2", "value2")
	llog.Ghelper.Infow("key1", "value1", "key2", "value2")
	llog.Ghelper.Warnw("key1", "value1", "key2", "value2")
	llog.Ghelper.Errorw("key1", "value1", "key2", "value2")
	llog.Ghelper.Fatalw("key1", "value1", "key2", "value2")

    // with format
	llog.Ghelper.Debugf("This is a format key1 %s key2 %s", "value1", "value2")
	llog.Ghelper.Infof("This is a format key1 %s key2 %s", "value1", "value2")
	llog.Ghelper.Warnf("This is a format key1 %s key2 %s", "value1", "value2")
	llog.Ghelper.Errorf("This is a format key1 %s key2 %s", "value1", "value2")
	llog.Ghelper.Fatalf("This is a format key1 %s key2 %s", "value1", "value2")
}

```

### Filter

filter可自定义日志需要过滤的内容以及过滤方法。filter也是Log接口的实现，也可以作为logger使用但需要搭配更底层的std使用。

```go
logger := llog.NewStdLogger(
	llog.WithStdWriter(log.Writer()),
)

	customFilter := func(level llog.Level, keyvals ...interface{}) bool {
		if level == llog.Warn {
			for i := 0; i < len(keyvals); i += 2 {
				if keyvals[i] == "password" {
					keyvals[i+1] = "******"
					return true
				}
			}
		}
		return false

	filter := llog.NewFilter(
		logger,
		// password和mobile作为key的kv键值对，v会被***代替
		// 需要注意的是helper内参考上面的MessageKey，可能会将password并不视为key
		// 如果不需要MessageKey可以使用带w后缀的方法
		llog.WithFilterKeys("password", "mobile"),
		// 设置filter级别，在Warn以下的日志不会被打印
		llog.WithLevel(llog.Warn),
		// 对于指定的value进行过滤
		llog.WithFilterValues("world"),
		// 自定义过滤函数
		llog.WithFilterFunc(customFilter),
	)
	// nothing
	filter.Log(llog.Info, "asdfasdfasdfasdf")  // Info等级低于Warn被过滤不被打印
	filter.Log(llog.Warn, "hello", "world")    // world与WithFilterValues内参数匹配，被默认字符串"***"掩盖
	filter.Log(llog.Warn, "password", "world") // password Warn等级与customFilter内匹配并且password也被匹配，world被自定义字符串掩盖

```

### Helper

Helper类的使用和Ghelper使用相同，请参考Ghelper用法。具体配置请查看下面的说明。


```go
helper = llog.NewHelper(
	llog.NewFilter(llog.DefaultLogger,
		llog.WithLevel(llog.Warn),
		llog.WithFilterKeys("password"),
	),
	// 自定义MessageKey
	llog.WithMessageKey("MessageKey"),
	// 自定义字符串打印函数,
	llog.WithSprint(fmt.Sprint),
	llog.WithSprintf(fmt.Sprintf),
)

helper.Infow("nothing")
helper.Warnw("password", "world")
```

## Interface

可通过实现Logger接口自定义相关logger功能，例如加入其他日志库的使用，或者自己封装

```go

type Logger interface {
	Log(l Level, keyvals ...any) error

    // 默认情况下调用最底层logger的callerSkip,如果自己是，最底层则应该设置相关字段，对其进行增加，便于打印正确的日志位置。
	AddCallerSkip(skip int) Logger
	GetCallerSkip() int

    // clone当前logger，如果不是最底层，则可能需要调用底层的logger的Clone方法
	Clone() Logger
}
```
