# jingtum-lib-go
jingtum-lib to be used for interacting with jingtum blockchain network。This is the go version。
目录架构：
    |--bin   编译生成可执行文件存放目录
    |--conf  配置文件存放目录
    |--log   日志存放目录
    |--pkg   编译生成包
    |--src   go源码目录
       |--common     公共代码存放目录
            |--github.com   github上第三方库存放目录
            |--goMath       数学运算类公共库
            |--goMisc       其他混合公共函数
            |--goNet        网络公共函数
            |--pkg          编译生成包
       |--jingtumLib 具体代码目录
       |--testLib    测试代码目录
    |--build.sh/build.bat linux/window的build脚本
    |--clean.sh           清空脚本
