## mgo

### 更新

更新 mgo 新版本后，为了简化测试，需要屏蔽掉部分测试文件, 运行命令如下

```
cd mgo
mv auth_test.go auth_test.tgo
mv cluster_test.go cluster_test.tgo
mv export_test.go export_test.tgo
mv gridfs_test.go gridfs_test.tgo
mv queue_test.go queue_test.tgo
mv session_test.go session_test.tgo
mv suite_test.go suite_test.tgo
mv txn/mgo_test.go txn/mgo_test.tgo
mv txn/sim_test.go txn/sim_test.tgo
mv txn/tarjan_test.go txn/tarjan_test.tgo
mv txn/txn_test.go txn/txn_test.tgo
```

### fix https://bugs.launchpad.net/mgo/+bug/1232685

添加了如下两个文件用于修复 LP:1232685

```
mgo/query_apply2.go
mgo/query_apply2_test.go
```
