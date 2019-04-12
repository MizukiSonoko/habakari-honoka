// Copyright (C) 2019 MizukiSonoko. All rights reserved.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	bench "golang.org/x/tools/benchmark/parse"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func getToken() (string, error) {
	// ToDo: set env in circle ci environments
	key := os.Getenv("GITHUB_PRIVATE_KEY")
	appID := os.Getenv("GITHUB_APP_ID")
	installationId := os.Getenv("GITHUB_INSTALLATION_ID")

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(
		[]byte(strings.Replace(key, "\\n", "\n", -1)))
	if err != nil {
		panic(err)
	}

	t := jwt.New(jwt.SigningMethodRS256)
	t.Claims = jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 10).Unix(),
		"iss": appID,
	}

	token, err := t.SignedString(privateKey)
	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", installationId), nil)
	if err != nil {
		return "", errors.Wrapf(err, "NewRequest failed")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")
	c := http.DefaultClient
	res, err := c.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "Do failed")
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrapf(err, "ReadAll failed")
	}

	var d map[string]string
	if err := json.Unmarshal(data, &d); err != nil {
		return "", errors.Wrapf(err, "ReadAll failed")
	}
	return d["token"], nil
}

func toMP(num float64) string {
	switch {
	case num < 1e3:
		return fmt.Sprintf("%4.2f ns", num)
	case num < 1e6:
		return fmt.Sprintf("%4.2f μs", num/1e3)
	case num < 1e9:
		return fmt.Sprintf("%4.2f ms", num/1e6)
	default:
		return fmt.Sprintf("%4.2f  s", num/1e9)
	}
}

func sendComment(owner, repo string, issuesId int, text string) error {
	token, err := getToken()
	if err != nil{
		return err
	}
	req, err := http.NewRequest("POST",
		fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d/comments",
			owner, repo, issuesId), nil)
	if err != nil {
		return errors.Wrapf(err, "NewRequest failed")
	}
	fmt.Printf("token:%s\n", token)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/vnd.github.squirrel-girl-preview")
	c := http.DefaultClient
	res, err := c.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Do failed")
	}
	fmt.Printf("res:%v\n", res)
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Printf("res:%v\n", string(data))
	return nil
}

func getIssueId(url string) (int, error) {
	l := strings.Split(url, "/")
	if len(l) == 0 {
		return 0, fmt.Errorf("invalid url(%s)", url)
	}
	ns := l[len(l)-1]
	n, err := strconv.Atoi(ns)
	if err != nil {
		return 0, errors.Wrapf(err, "Atoi failed")
	}
	return n, nil
}

func comment(text string) error {
	owner := os.Getenv("CIRCLE_PROJECT_USERNAME")
	repo := os.Getenv("CIRCLE_PROJECT_REPONAME")
	pullReq := os.Getenv("CIRCLE_PULL_REQUEST")
	if pullReq == ""{
		fmt.Printf("This is not pull reqest")
		return nil
	}
	id, err := getIssueId(pullReq)
	if err != nil {
		return errors.Wrapf(err, "getIssueId failed")
	}
	err = sendComment(owner, repo, id, text)
	if err != nil {
		return errors.Wrapf(err, "sendComment failed")
	}
	return nil
}

func main() {
	flag.Parse()

	var lines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		b, err := bench.ParseLine(text)
		if err != nil {
			continue
		}
		lines = append(lines,
			fmt.Sprintf("%s\t%10.d times\t%s/op\t", b.Name, b.N, toMP(b.NsPerOp)))
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var text string
	for _, l := range lines {
		text += l + "\n"
	}
	err := comment(text)
	if err != nil {
		fmt.Printf("comment failed err:%s\n", err)
		os.Exit(1)
	}
}
