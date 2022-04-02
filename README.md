# red-trie
An Implement of trie with redcon
使用Redcon redis协议实现的一个Trie树工具
启动后可以使用Redis协议进行简单地操作


```
127.0.0.1:6339> add trie1 "中国" "国家"
OK
127.0.0.1:6339> contain trie1 "中国"
(integer) 1
127.0.0.1:6339> prefix trie1 "中"
(integer) 1
127.0.0.1:6339> prefix trie1 "中国人"
(integer) 0
127.0.0.1:6339> add trie1 "人民" "群体"
OK
127.0.0.1:6339> add trie1 "共和国" "体制"
OK
127.0.0.1:6339> entities trie1 "中国"
1) (integer) 0
2) (integer) 1
3) 国家
4) 中国
127.0.0.1:6339> full trie1 "中华人民共和国"
1) (integer) 2
2) (integer) 3
3) 群体
4) 人民
5) (integer) 4
6) (integer) 6
7) 体制
8) 共和国
```