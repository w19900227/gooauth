package main

import (
    "net/http"
    "time"
    "strings"
    "testing"
    "io/ioutil"
    "encoding/json"
    "fmt"
)


// func (this *User) Get(request *restful.Request, response *restful.Response) {
//     var r resp

//     token := request.PathParameter("token")
//     key := "user_list:" + token

//     result, err := GetService(key)

//     if result == nil || err != nil {
//         r.Status = "fail"
//         r.Errmsg = "no user"
//         response.WriteEntity(r)
//         return
//     }

//     var user User_format
//     json.Unmarshal(result, &user)

//     r.Status = "ok"
//     r.Data = user
//     response.WriteEntity(r)
//     return 
// }

// func TestGetList222(t *testing.T) {
//     // response, _ := http.Get("http://localhost:9090/user")
//     // datas, _ := ioutil.ReadAll(response.Body)
//     // return datas
//     nilBody, _ := http.NewRequest("POST", "http://www.google.com/search?q=foo", nil)
//     tests := []*http.Request{
//         nilBody,
//         {Method: "GET", URL: nil},
//     }
//     for i, req := range tests {
//         err := req.ParseForm()
//         if req.Form == nil {
//             t.Errorf("%d. Form not initialized, error %v", i, err)
//         }
//         if req.PostForm == nil {
//             t.Errorf("%d. PostForm not initialized, error %v", i, err)
//         }
//     }
// }

type user_format struct {
    Email string `json:"email"`
    Password  string `json:"password"`
    Start_time string `json:"start_time"`
}
type user struct {
    Email string `json:"email"`
    Password  string `json:"password"`
    Start_time time.Time `json:"start_time"`
}
type test_user struct {
    Status string `json:"status"`
    Erron  string `json:"erron"`
    Errmsg string`json:"errmsg"`
    Data   user `json:"data"`
}

type test_users struct {
    Status string `json:"status"`
    Erron  string `json:"erron"`
    Errmsg string`json:"errmsg"`
    Data   []user `json:"data"`
}

var url string = "http://localhost:9090"
var content_type string = "application/json"

func TestCreate(t *testing.T) {
    data := user_format{
        Email: "test12345@gmail.com",
        Password: "test",
        Start_time: "2015-11-05T00:00:00Z"}

    buf, _ := json.Marshal(data)
    body := strings.NewReader(string(buf))

    response, _ := http.NewRequest("POST", url+"/user", body)
    response.Header.Set("Content-Type", content_type)
    client := &http.Client{}
    re, _ := client.Do(response)
    defer re.Body.Close()
    
    datas, _ := ioutil.ReadAll(re.Body)
    
    var user test_user
    json.Unmarshal(datas, &user)
    
    if user.Status == "fail" {
        t.Errorf("response fail error : %v", user)
    }
}

func TestUpdate(t *testing.T) {
    data := user_format{
        Email: "test12345@gmail.com",
        Password: "test12345",
        Start_time: "2015-12-09T00:00:00Z"}
    token := Encryption(data.Email)

    buf, _ := json.Marshal(data)
    body := strings.NewReader(string(buf))

    response, _ := http.NewRequest("PUT", url+"/user/"+token, body)
    response.Header.Set("Content-Type", content_type)
    client := &http.Client{}
    re, _ := client.Do(response)
    defer re.Body.Close()
    
    datas, _ := ioutil.ReadAll(re.Body)
    
    var user test_user
    json.Unmarshal(datas, &user)
    
    if user.Status == "fail" {
        t.Errorf("response fail error : %v", user)
    }
}

func TestGetList(t *testing.T) {
    response, _ := http.Get(url+"/user")
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        t.Errorf("reqpuest body error : %v", err)
    }

    var user test_users
    json.Unmarshal(body, &user)

    if user.Status == "fail" {
        t.Errorf("response fail error : %v", user)
    }
    show(user)
}

func TestGet(t *testing.T) {
    token := Encryption("test12345@gmail.com")

    response, _ := http.Get(url+"/user/"+token)
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        t.Errorf("reqpuest body error : %v", err)
    }

    var user test_user
    json.Unmarshal(body, &user)
    if user.Status == "fail" {
        t.Errorf("response fail error : %v", user)
    }
    
    show(user)
}

func TestCheckToken(t *testing.T) {
    token := Encryption("test12345@gmail.com")

    response, _ := http.NewRequest("GET", url+"/token/"+token, nil)
    response.Header.Set("Content-Type", content_type)
    client := &http.Client{}
    re, _ := client.Do(response)
    defer re.Body.Close()

    datas, _ := ioutil.ReadAll(re.Body)
    var user test_user
    json.Unmarshal(datas, &user)

    if user.Status == "fail" {
        t.Errorf("response fail error : %v", user)
    }

    show(user)
}

func TestLogin(t *testing.T) {
    data := user_format{
        Email: "test12345@gmail.com",
        Password: "test12345"}
    
    buf, _ := json.Marshal(data)
    body := strings.NewReader(string(buf))
    response, _ := http.NewRequest("POST", url+"/user/login", body)
    response.Header.Set("Content-Type", content_type)

    client := &http.Client{}
    re, _ := client.Do(response)
    defer re.Body.Close()

    datas, _ := ioutil.ReadAll(re.Body)
    var user test_user
    json.Unmarshal(datas, &user)

    if user.Status == "fail" {
        t.Errorf("response fail error : %v", user)
    }

    show(user)
}

func TestDelete(t *testing.T) {
    token := Encryption("test12345@gmail.com")

    response, _ := http.NewRequest("DELETE", "http://localhost:9090/user/"+token, nil)
    response.Header.Set("Content-Type", content_type)
    client := &http.Client{}
    re, _ := client.Do(response)
    defer re.Body.Close()
    
    datas, _ := ioutil.ReadAll(re.Body)
    
    var user test_user
    json.Unmarshal(datas, &user)
    
    if user.Status == "fail" {
        t.Errorf("response fail error : %v", user)
    }
}

func TestTimes(t *testing.T) {
    response, _ := http.NewRequest("GET", url+"/token/time", nil)
    response.Header.Set("Content-Type", content_type)
    client := &http.Client{}
    re, _ := client.Do(response)
    defer re.Body.Close()

    datas, _ := ioutil.ReadAll(re.Body)
    var user test_user
    json.Unmarshal(datas, &user)

    if user.Status == "fail" {
        t.Errorf("response fail error : %v", user)
    }
}

func show(result interface{}) {
    js, _ := json.Marshal(result)
    fmt.Println(string(js))
}
