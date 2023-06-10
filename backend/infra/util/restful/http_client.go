package restful

import (
	"bytes"
	"code-shooting/infra/logger"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"time"
)

const (
	HTTP  = "http"
	HTTPS = "https"
)

func Put(url string, data interface{}, contentType string) ([]byte, error) {
	jsonStr, _ := json.Marshal(data)
	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", contentType)
	defer req.Body.Close()
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	logger.Debugf("[HTTP Client PUT] url: %s StatusCode: %d Body: %s", url, resp.StatusCode, string(result))
	if resp.StatusCode >= 300 {
		logger.Errorf("[HTTP Client PUT] url: %s StatusCode: %d Body: %s", url, resp.StatusCode, string(result))
	}
	return result, nil
}

func Post(url string, data interface{}, contentType string) ([]byte, error) {
	jsonStr, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if contentType != "" {
		req.Header.Add("content-type", contentType)
	}

	defer req.Body.Close()
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	//logger.Debugf("[HTTP Client POST] url: %s StatusCode: %d Body: %s",url,resp.StatusCode,string(result))
	if resp.StatusCode >= 300 {
		logger.Errorf("[HTTP Client POST] url: %s StatusCode: %d Body: %s", url, resp.StatusCode, string(result))
	}
	return result, nil
}

func HTTPSPost(url string, data interface{}, contentType string) (content []byte, err error) {
	jsonStr, _ := json.Marshal(data)

	//log.Printf("HTTP Client RequestBody: %s", string(jsonStr))

	//caCert, err := ioutil.ReadFile("conf/server.crt")
	caCertUrl := SetUrl(os.Getenv("openpalette_service_ip"), os.Getenv("openpalette_service_port"), "cert/ca.crt")
	caCert, err := Get(caCertUrl)
	logger.Debugf("[HTTPS Client POST] url: %s, caCert: %v", caCertUrl, string(caCert))
	if err != nil {
		logger.Errorf("[HTTPS Client POST]  Can not get caCert:%v", err)
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	cookieJar, _ := cookiejar.New(nil)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		logger.Errorf("[HTTPS Client POST]  new request error:", err)
		return nil, err
	}
	req.Header.Add("content-type", contentType)
	defer req.Body.Close()
	client := &http.Client{Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				//InsecureSkipVerify: true,
				RootCAs: caCertPool,
			},
		},
		Jar: cookieJar,
	}
	logger.Debugf("[HTTPS Client POST]  client info : %+v", *client)
	resp, err2 := client.Do(req)
	if err2 != nil {
		logger.Errorf("[HTTPS Client POST]  client do error:", err2)
		return nil, err2
	}
	defer resp.Body.Close()
	content, _ = ioutil.ReadAll(resp.Body)
	logger.Debugf("[HTTPS Client POST] url: %s StatusCode: %d Body: %s", url, resp.StatusCode, string(content))
	if resp.StatusCode >= 300 {
		logger.Errorf("[HTTPS Client POST] url: %s StatusCode: %d Body: %s", url, resp.StatusCode, string(content))
	}
	return content, nil
}

func Get(url string) ([]byte, error) {
	client := &http.Client{}
	res, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	//logs.Debugf("[HTTP Client GET]","HTTP Client url %v of ResponseBody: %s",url, string(body))
	if res.StatusCode >= 300 {
		logger.Errorf("[HTTP Client GET] url: %s StatusCode: %d Body: %s", url, res.StatusCode, string(body))
	}
	return body, nil
}

func SetUrl(ip string, port string, pathUrl string) string {
	u, _ := url.Parse("")
	u.Scheme = HTTP
	addr := net.JoinHostPort(ip, port)
	u.Host = addr
	u.Path = path.Join(u.Path, pathUrl)

	return u.String()
}

func SetUrlWithQuery(ip string, port string, pathUrl string, query map[string]string) string {
	u, _ := url.Parse(SetUrl(ip, port, pathUrl))
	q := u.Query()
	for key, value := range query {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func SetHTTPSUrl(ip string, port string, pathUrl string) string {
	u, err := url.Parse("")
	if err != nil {
		logger.Fatal("SET HTTPS URL", err)
	}
	u.Scheme = HTTPS
	addr := net.JoinHostPort(ip, port)
	u.Host = addr
	u.Path = path.Join(u.Path, pathUrl)

	return u.String()
}
