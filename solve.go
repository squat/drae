package main

import (
	"bytes"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func Solve(res *http.Response) *http.Response {
	defer res.Body.Close()
	buffer, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	contents := string(buffer)
	//fmt.Println(contents)

	tableRe := regexp.MustCompile(`table = "(.*?)"`)
	table := tableRe.FindStringSubmatch(contents)[1]

	cRe := regexp.MustCompile(`c = (\d*)`)
	c, _ := strconv.ParseInt(cRe.FindStringSubmatch(contents)[1], 0, 64)

	sltRe := regexp.MustCompile(`slt = "(.*?)"`)
	slt := sltRe.FindStringSubmatch(contents)[1]

	s1Re := regexp.MustCompile(`s1 = '(.*?)'`)
	s1 := s1Re.FindStringSubmatch(contents)[1]

	s2Re := regexp.MustCompile(`s2 = '(.*?)'`)
	s2 := s2Re.FindStringSubmatch(contents)[1]

	valueRe := regexp.MustCompile(`value="(.*?)" \+ chlg`)
	value := valueRe.FindStringSubmatch(contents)[1]

	valuesRe := regexp.MustCompile(`value="(.*?)"/`)
	values := valuesRe.FindAllStringSubmatch(contents, -1)

	namesRe := regexp.MustCompile(`input type="hidden" name="(.*?)"`)
	names := namesRe.FindAllStringSubmatch(contents, -1)

	uriRe := regexp.MustCompile(`"POST" action="(.*?)"`)
	uri := uriRe.FindStringSubmatch(contents)

	end := []byte(s2)[0]
	arr := []string{s1, s1, s1, s1}
	var crc int64
	var chlg string

	for i := 0; i < 624; i++ {
		for j := 3; j >= 0; j-- {
			t := []byte(arr[j])[0]
			t++
			arr[j] = string(t)

			if t <= end {
				break
			} else {
				arr[j] = s1
			}
		}
		chlg = strings.Join(arr, "")
		str := chlg + slt
		crc = -1

		for k := 0; k < 12; k++ {
			index := ((crc ^ int64([]byte(str)[k])) & 0x000000FF) * 9
			tmp, _ := strconv.ParseInt(("0x" + table[index:(index+8)]), 0, 64)
			crc = ((crc >> 8) ^ tmp) & 4294967295
			if crc > 2147483647 {
				crc -= 4294967296
			}
		}
		crc = crc ^ -1
		crc = int64(math.Abs(float64(crc)))
		if crc == c {
			break
		}
	}

	fields := make([]string, len(names))

	for i := range fields {
		fields[i] = url.QueryEscape(names[i][1]) + "=" + url.QueryEscape(values[i][1])
	}

	fields[1] = url.QueryEscape(names[1][1]) + "=" + url.QueryEscape(value+chlg+":"+slt+":"+strconv.FormatInt(crc, 10))

	payload := strings.Join(fields[:], "&")
	endpoint, err := url.QueryUnescape(uri[1])
	if err != nil {
		panic(err)
	}

	endpoint = "http://lema.rae.es" + endpoint

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return resp
}
