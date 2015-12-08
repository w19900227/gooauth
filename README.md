# gooauth

簡易小型會員認證管理

### 安裝golang
    $ wget https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
    $ tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz
    $ export PATH=$PATH:/usr/local/go/bin

### 環境部署
    $ git clone https://github.com/w19900227/gooauth
    $ cd gooauth
    $ export GOPATH=$PWD

### Testing
必須先run主程式，才可進行測試

    $ go run main.go
    $ go test -v 

### get plug
    $ go get github.com/emicklei/go-restful
    $ go get github.com/garyburd/redigo