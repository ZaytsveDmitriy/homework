package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var user User
	const RE = "@(\\w+\\.)"

	matcher, err := regexp.Compile(RE + domain)
	if err != nil {
		return nil, err
	}

	json := jsoniter.ConfigCompatibleWithStandardLibrary
	reader := bufio.NewReader(r)
	stat := make(DomainStat, 100)

	for i := 0; ; i++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, Error(err)
		}
		if err = json.Unmarshal(line, &user); err != nil {
			return nil, Error(err)
		}

		matched := matcher.FindStringSubmatch(user.Email)

		if len(matched) != 0 {
			key := strings.ToLower(matched[1] + domain)
			num := stat[key]
			num++
			stat[key] = num
		}
	}

	return stat, nil
}

func Error(err error) error {
	const message = "get users error: %w"
	return fmt.Errorf(message, err)
}

func GetDomainStatAld(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		matched, err := regexp.Match("\\."+domain, []byte(user.Email))
		if err != nil {
			return nil, err
		}

		if matched {
			num := result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]
			num++
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])] = num
		}
	}
	return result, nil
}
