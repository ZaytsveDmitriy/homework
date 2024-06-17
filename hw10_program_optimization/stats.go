package hw10programoptimization

import (
	"bufio"
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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	reader := bufio.NewReader(r)

	for i := 0; ; i++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return result, nil
			}
			return result, err
		}
		var user User
		if err = json.Unmarshal(line, &user); err != nil {
			return result, err
		}
		result[i] = user
	}
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat, 100)

	matcher, err := regexp.Compile("@(\\w+\\.)" + domain)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		matched := matcher.FindSubmatch([]byte(user.Email))

		if len(matched) != 0 {
			key := strings.ToLower(string(matched[1]) + domain)
			num := result[key]
			num++
			result[key] = num
		}
	}
	return result, nil
}
