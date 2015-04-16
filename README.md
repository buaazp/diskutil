## Golang Disk Utils

This package is used for go codes to get MegaRaid stat.

### Usage

```
    ds, err := diskutil.NewDiskStatus("/opt/MegaRAID/MegaCli/MegaCli64", 1)
    if err != nil {
        fmt.Fprintf(os.Stderr, "DiskStatus New error: %v\n", err)
        return
    }

    err = ds.Get()
    if err != nil {
        fmt.Fprintf(os.Stderr, "DiskStatus Get error: %v\n", err)
        return
    }

    fmt.Println(ds)
```

Full sample code is in examples.

### GoDoc

[https://godoc.org/github.com/buaazp/diskutil](https://godoc.org/github.com/buaazp/diskutil) 

### Issue

If you meet some problem in your servers, please create a github issue or contact me:

weibo: [@招牌疯子](http://weibo.com/buaazp)

mail: zp@buaa.us


