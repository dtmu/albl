# albl
This is convert tool from alb access log to xlsx (EXCEL).

これはALBアクセスログをxlsxファイル(EXCEL)にコンバートするCLIツールです。

## Usage (使い方)

```bash
$ ls $GOPATH/src/github.com/dtmu/albl/testFile
123456789012_elasticloadbalancing_us-east-2_app.my-loadbalancer.1234567890abcdef_20140215T2340Z_172.160.001.192_20sg8hgm.log

$ albl -d $GOPATH/src/github.com/dtmu/albl/testFile -n demo.xlsx

$ ls
. .. demo.xlsx
```
