// Copyright(c) 2026 Beijing Yingfei Networks Technology Co.Ltd.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http: //www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package lib

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

func Recover(module string) {
	if err := recover(); err != nil {
		if err1, ok := err.(error); ok {
			fmt.Println(fmt.Sprintf("%s:Excute recover:%s", module,
				err1.Error()))
		}
	}
}

type BasicAuthInfo struct {
	UserName string
	Password string
	Schema   string
}

// ReadWithRetry 尝试多次读取指定URL的内容，直到成功或达到最大重试次数。
func ReadWithRetry(URLPath string, timeout int, headers map[string]string, retry int, interval int, basicInfo *BasicAuthInfo) ([]byte, error) {
	var resp []byte
	var err error
	for i := 0; i < retry; i++ {
		resp, err = Read(URLPath, timeout, headers, basicInfo)
		if err == nil {
			break
		}
		time.Sleep(time.Microsecond * time.Duration(interval))
	}
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func Read(urlPath string, timeout int, headers map[string]string, basicInfo *BasicAuthInfo) ([]byte, error) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return nil, err
	}

	if basicInfo != nil && basicInfo.UserName != "" {
		req.SetBasicAuth(basicInfo.UserName, basicInfo.Password)
	}

	// Add the necessary headers
	req.Header.Add("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Create an HTTP client with a timeout of 5 seconds
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	if strings.HasPrefix(urlPath, "https") {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	// Defer the recovery of any panics
	defer Recover(fmt.Sprintf("HTTP GET request to %s", urlPath))

	// Execute the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP GET request to %s returned status code %d", urlPath, resp.StatusCode)
	}

	// Read the response body
	return ioutil.ReadAll(resp.Body)
}

func SendFiles(filePaths []string, url string, basicInfo *BasicAuthInfo) (int, []byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, filePath := range filePaths {
		file, err := os.Open(filePath)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to open file: %v", err)
		}
		defer file.Close()

		// 创建一个新的 part for the file
		part, err := writer.CreateFormFile("file", filePath)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to create form file: %v", err)
		}

		// Copy the file data into the part
		_, err = io.Copy(part, file)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to copy file data: %v", err)
		}
	}

	// Close the writer to write the terminating boundary
	err := writer.Close()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to close multipart writer: %v", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create new request: %v", err)
	}

	if basicInfo != nil {
		req.SetBasicAuth(basicInfo.UserName, basicInfo.Password)
	}

	// Set the Content-Type header to multipart/form-data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the request using an HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	return resp.StatusCode, responseBody, err
}

func SendPatchRequest(url string, payload interface{}, basicInfo *BasicAuthInfo) (int, []byte, error) {
	// 将 payload 编码为 JSON
	var jsonData []byte
	var err error
	if payload != nil {
		jsonData, err = json.Marshal(payload)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to marshal payload: %v", err)
		}
	}

	// 创建一个新的 HTTP PATCH 请求
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create new request: %v", err)
	}

	if basicInfo != nil {
		req.SetBasicAuth(basicInfo.UserName, basicInfo.Password)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 使用 HTTP 客户端发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取并打印响应体
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("Error reading response body: %v\n", err)
	}

	return resp.StatusCode, responseBody, nil
}

func SendDeleteRequest(url string, basicInfo *BasicAuthInfo) (int, []byte, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if basicInfo != nil {
		req.SetBasicAuth(basicInfo.UserName, basicInfo.Password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to send request: %v", err)
	}

	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("Error reading response body: %v\n", err)
	}

	return resp.StatusCode, responseBody, nil
}
