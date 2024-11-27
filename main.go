package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

func main() {
	// コマンドライン引数チェック
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <log_file> <mode>")
		fmt.Println("<mode>: 'hourly' for hourly access counts, 'ip' for IP address counts")
		return
	}
	logFileName := os.Args[1]
	mode := strings.ToLower(os.Args[2])

	// ログファイルを開く
	file, err := os.Open(logFileName)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// 正規表現で日時とIPアドレスを抽出
	logPattern := regexp.MustCompile(`^([\d.]+) .+ \[(\d{2}/[A-Za-z]{3}/\d{4}:\d{2}):\d{2}:\d{2} .+\]`)
	hourlyCounts := make(map[string]int)
	ipCounts := make(map[string]int)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		match := logPattern.FindStringSubmatch(line)
		if len(match) > 2 {
			ip := match[1]            // IPアドレス
			hour := match[2]          // "dd/mmm/yyyy:hh" フォーマット
			hourlyCounts[hour]++
			ipCounts[ip]++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// モードによる分岐処理
	if mode == "hourly" {
		// 時間順にソート
		var hours []string
		for hour := range hourlyCounts {
			hours = append(hours, hour)
		}
		sort.Slice(hours, func(i, j int) bool {
			timeI, _ := time.Parse("02/Jan/2006:15", hours[i])
			timeJ, _ := time.Parse("02/Jan/2006:15", hours[j])
			return timeI.Before(timeJ)
		})

		// 結果を出力 (時間ごとのアクセス数)
		fmt.Println("Hourly Access Counts:")
		for _, hour := range hours {
			parsedTime, err := time.Parse("02/Jan/2006:15", hour)
			if err != nil {
				fmt.Printf("Error parsing time: %v\n", err)
				continue
			}
			fmt.Printf("%s: %d accesses\n", parsedTime.Format("2006-01-02 15:00"), hourlyCounts[hour])
		}
	} else if mode == "ip" {
		// IPアドレスごとの結果をソート
		var ips []string
		for ip := range ipCounts {
			ips = append(ips, ip)
		}
		sort.Slice(ips, func(i, j int) bool {
			return ipCounts[ips[i]] > ipCounts[ips[j]] // アクセス数で降順ソート
		})

		// 結果を出力 (IPアドレスごとのアクセス数)
		fmt.Println("IP Address Access Counts:")
		for _, ip := range ips {
			fmt.Printf("%s: %d accesses\n", ip, ipCounts[ip])
		}
	} else {
		fmt.Println("Invalid mode. Use 'hourly' or 'ip'.")
	}
}

