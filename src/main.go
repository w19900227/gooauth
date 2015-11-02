package main

import (
    "github.com/emicklei/go-restful"
    "github.com/garyburd/redigo/redis"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "encoding/base64"
    "time"
    "log"
)
import "fmt"

func main() {
    restful.Add(UserRouter())
    restful.Add(TokenRouter())

    // log.Fatal(http.ListenAndServe(":9090", nil))
    log.Fatal(http.ListenAndServe(":31001", nil))
}

func TokenRouter() *restful.WebService {
    service := new(restful.WebService)
    service.
        Path("/token").
        Consumes(restful.MIME_JSON, restful.MIME_JSON).
        Produces(restful.MIME_JSON, restful.MIME_JSON)

    service.Route(service.GET("/{token}").To(CheckToken))
    service.Route(service.GET("/time").To(Times))
 
    return service
}

func UserRouter() *restful.WebService {
    service := new(restful.WebService)
    service.
        Path("/user").
        Consumes(restful.MIME_JSON, restful.MIME_JSON).
        Produces(restful.MIME_JSON, restful.MIME_JSON)

    var user User
    service.Route(service.GET("/{token}").To(user.Get))
    service.Route(service.GET("").To(user.GetList))
    service.Route(service.POST("").To(user.Create))
    service.Route(service.POST("/login").To(user.Login))
    service.Route(service.PUT("/{token}").To(user.Update))
    service.Route(service.DELETE("/{token}").To(user.Delete))
 
    return service
}


type User struct {
}
type resp struct {
    Status string `json:"status"`
    Erron  string `json:"erron"`
    Errmsg string`json:"errmsg"`
    Data   interface{} `json:"data"`
}

type User_format struct {
    Email string `json:"email"`
    Password  string `json:"password"`
    // Start_time string `json:"start_time"`
    Start_time time.Time `json:"start_time"`
    Token string `json:"token"`
}

func NewPool() *redis.Pool {
    return &redis.Pool{
                MaxIdle: 80,
                MaxActive: 12000, // max number of connections
                Dial: func() (redis.Conn, error) {
                        c, err := redis.Dial("tcp", "localhost:6379")
                        // c, err := redis.Dial("tcp", "golangredis:6379")
                        if err != nil {
                                panic(err.Error())
                        }
                        return c, err
                },
        } 
}

func GetService(key string) ([]byte, error) {
    pool := NewPool()
    c := pool.Get()

    result, err := redis.Bytes(c.Do("GET", key))
    return result, err
}

func (this *User) Get(request *restful.Request, response *restful.Response) {
    var r resp

    token := request.PathParameter("token")
    key := "user_list:" + token

    result, err := GetService(key)

    if result == nil || err != nil {
        r.Status = "fail"
        r.Errmsg = "no user"
        response.WriteEntity(r)
        return
    }

    var user User_format
    json.Unmarshal(result, &user)

    r.Status = "ok"
    r.Data = user
    response.WriteEntity(r)
    return 
}

func (this *User) GetList(request *restful.Request, response *restful.Response) {
    var r resp
    pool := NewPool()
    c := pool.Get()

    keys := "user_list:*"
    var user User_format
    var user_array []User_format

    var list_user []string
    result, err := redis.Values(c.Do("KEYS", keys))

    if result == nil || err != nil {
        r.Status = "fail"
        r.Errmsg = "no user"
        response.WriteEntity(r)
        return
    }

    redis.ScanSlice(result, &list_user)
    for _, value := range list_user {
        data, _ := GetService(value)
        json.Unmarshal(data, &user)
        user_array = append(user_array, user)
    }

    r.Status = "ok"
    r.Data = user_array
    response.WriteEntity(r)
    return 
}

func (this *User) Create(request *restful.Request, response *restful.Response) {
    var r resp
    var user User_format

    body, err := ioutil.ReadAll(request.Request.Body)
    if err != nil {
        r.Status = "fail"
        r.Errmsg = "Request Body error"
        response.WriteEntity(r)
        return
    }

    json.Unmarshal(body, &user)
    user.Password = Encryption(user.Password)
    user.Token = Encryption(user.Email)
    user.Start_time = user.Start_time

    value, _ := json.Marshal(user)

    pool := NewPool()
    c := pool.Get()

    key := "user_list:" + user.Token

    result, err := c.Do("SET", key, value)

    if result == nil || err != nil {
        r.Status = "fail"
        r.Errmsg = "no user"
        response.WriteEntity(r)
        return
    }

    r.Status = "ok"
    response.WriteEntity(r)
    return
}

func (this *User) Update(request *restful.Request, response *restful.Response) {
    var r resp
    var user User_format

    token := request.PathParameter("token")
    body, _ := ioutil.ReadAll(request.Request.Body)
    json.Unmarshal(body, &user)
    id, err := Decryption(token)
    if err != nil {
        r.Status = "fail"
        r.Errmsg = "token Decryption error"
        response.WriteEntity(r)
        return
    }
    
    user.Email = id.(string)
    user.Password = Encryption(user.Password)
    user.Token = token
    user.Start_time = user.Start_time

    value, _ := json.Marshal(user)

    pool := NewPool()
    c := pool.Get()
    key := "user_list:" + token

    result, err := c.Do("SET", key, value)

    if result == nil || err != nil {
        r.Status = "fail"
        r.Errmsg = "no user"
        response.WriteEntity(r)
        return
    }

    r.Status = "ok"
    response.WriteEntity(r)
    return
}

func (this *User) Delete(request *restful.Request, response *restful.Response) {
    var r resp
    token := request.PathParameter("token")

    pool := NewPool()
    c := pool.Get()
    key := "user_list:" + token

    result, err := c.Do("DEL", key)

    if result == nil || err != nil {
        r.Status = "fail"
        r.Errmsg = "no user"
        response.WriteEntity(r)
        return
    }

    r.Status = "ok"
    response.WriteEntity(r)
    return 
}

func (this *User) Login(request *restful.Request, response *restful.Response) {
    var r resp
    var data User_format
    body, _ := ioutil.ReadAll(request.Request.Body)
    json.Unmarshal(body, &data)

    pool := NewPool()
    c := pool.Get()

    key := "user_list:" + Encryption(data.Email)
    result, err := redis.Bytes(c.Do("GET", key))

    if result == nil || err != nil {
        r.Status = "fail"
        r.Errmsg = "no user"
        response.WriteEntity(r)
        return
    }

    var user User_format
    json.Unmarshal(result, &user)

    pwd := Encryption(data.Password)

    if user.Password != pwd {
        r.Status = "fail"
        r.Errmsg = "account or password is error"
        response.WriteEntity(r)
        return
    }

    resp := _checkToken(user.Token)

    if resp.Status != "ok" {
        response.WriteEntity(resp)
        return
    }

    r.Status = "ok"
    r.Data = user.Token
    response.WriteEntity(r)
    return 
}

func CheckToken(request *restful.Request, response *restful.Response) {
    token := request.PathParameter("token")

    r := _checkToken(token)
    response.WriteEntity(r)
    return
}

func _checkToken(token string) (resp) {
    var r resp
    
    pool := NewPool()
    c := pool.Get()

    key := "user_list:" + token
    result, err := redis.Bytes(c.Do("GET", key))

    if result == nil || err != nil {
        r.Status = "fail"
        r.Errmsg = "no token"
        return r
    }

    var user User_format
    json.Unmarshal(result, &user)

    if user.Token != token {
        r.Status = "fail"
        r.Errmsg = "token auth issue"
        return r
    }

    layout := "2006-01-02 15:04:05"
    now_time := time.Now()
    // now_time = now_time.Add(time.Hour * 8)
    start_time := user.Start_time
    last_time := user.Start_time.Add(time.Hour * 24)

    now_times := now_time.Format(layout)
    now_time_parse, err := time.Parse(layout, now_times)

    start_times := start_time.Format(layout)
    start_time_parse, err := time.Parse(layout, start_times)
    
    last_times := last_time.Format(layout)
    last_time_parse, err := time.Parse(layout, last_times)

    fmt.Printf("Token:%s, Now:%s, Start:%s, End:%s\n", token, now_time_parse, start_time_parse, last_time_parse)

    if start_time_parse.After(now_time_parse) {
        r.Status = "fail"
        r.Errmsg = "no start"
        return r
    }
    if last_time_parse.Before(now_time_parse) {
        r.Status = "fail"
        r.Errmsg = "expire"
        return r
    }

    r.Status = "ok"
    return r
}

func Times(request *restful.Request, response *restful.Response) {
    var r resp

    now_time := time.Now()

    r.Status = "ok"
    r.Data = now_time
    response.WriteEntity(r)
    return 
}

func Encryption(str string) string {
    data := []byte(str)
    ret := base64.StdEncoding.EncodeToString(data)
    return ret
}

func Decryption(str string) (interface{}, error) {
    data, err := base64.StdEncoding.DecodeString(str)
    if err != nil {
        return nil, err
    }
    return string(data), nil
}

