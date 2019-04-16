
## How to use

> you may config in fluent-bit conf.

```ini
[FILTER]
        Name  cgolib
        Match  *
        golib_so /fluent-bit/lib/cgolib.so
        lib_args -v 5 -log_dir=/tmp/cgolib.log
        ## modules that you must init, split by ','
        init_modules k8s
        ## for k8s init
        #kube_config /fluent-bit/etc/kube-config   // in cluster , no need config
        ## for filter config file
        parse_config_file  /fluent-bit/etc/filter_cgolib.yaml   ## same with ../config/config.yaml
``` 

## modules

### 1. k8s

```ini
parameters:
kube_config          kube client config file
lib_args             init args, init glog
```

## filters
```ini
parameters:
parse_config_file   config file of parseFilter
```
[example](../config/config.yaml)

### matchers

1. matchField

检查key中定义的字段的值是否匹配

    a. 先检查 matchStr 是否为空, 如果不是空, 则根据matchStr精确匹配
    b. 否则根据matchRegex进行正则匹配
    c. 如果revert为true, 表示将结果取反

```yaml
name: empdirlog
key: logpath
matchRegex: .*empty-dir.*
matchStr: ""
revert: false
```

### parses

1. parseRegex

将定义的key的值根据配置的regex进行解析, preserve_key为true表示会保留原key

```yaml
name: getPodUid
key: logpath
preserve_key: true
regex: ".*pods/?(?P<pod_uid>[a-z0-9](?:[-a-z0-9]*[a-z0-9]))?/volumes/.*"
```

2. parseK8sUid

根据k8s的uid来添加k8s相关的标签, key定义从哪个字段取uid值

```yaml
name: getPodInfoByUid
key: pod_uid
```

根据k8s的名称和命名空间来添加k8s相关的标签, key定义从哪个字段取名称, namespaceKey定义从哪个字段取命名空间

3. parseK8sName

```yaml
name: getPodInfoByName
key: pod_name
namespaceKey: namespace_name
```

4. parseJson

将字段按json格式展开, key定义了要展开的字段, preserve_key为true表示会保留原key

```yaml
name: externalLog
key: log
preserve_key: true
```

5. parseTimer

格式化日期, key表示要格式化的字段, time_format 定义了源格式, 会转成 RFC3339 格式

```yaml
name: parseMycatTime
key: time
time_format: 2006-01-02 15:04:05.999
```

### filter

将matcher和parser进行组合。

    a. matchType有两种取值: matchAll表示要满足所有的matchers, matchOne表示matchers中只要满足一个
    b. matchers 中定义需要进行matcher的名称
    c. parsers 中定义需要进行parser的名称, 会保持顺序依次执行 

```yaml
name: emptydirGetUid
matchType: matchAll
matchers:
  - empdirlog
parsers:
  - getPodUid
  - getPodInfoByUid
```

