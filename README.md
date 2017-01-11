# gitlab

## 主要功能

- 获取指定group下面所有的project
- 获取和统计每个project下面的commits数据，汇总2016年全年的commits数量，最多的10个开发
- 获取每个project的merge request
- 获取每个project的所有review comments,保存到文件


## 修改TOKEN

找到`lib/client.go`, 改成你自己的`private token`和gitlab api url
