package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	domainfile, output, outputtype string
	thread                         int
)

func init() {
	flag.StringVar(&domainfile, "f", "", "打开文件")
	flag.StringVar(&output, "o", "", "输出文件")
	flag.StringVar(&outputtype, "t", "ip", "保存类型")
	flag.IntVar(&thread, "g", 20, "线程")
}

// nslookup

func nslookup(domains <-chan string, results chan<- string, doresults chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	for host := range domains {
		ips, err := net.LookupIP(host)
		if err != nil {
			fmt.Printf("域名%s解析失败\n", host, err)
			continue
		}
		//fmt.Println(host, "=>", ips[:])
		fmt.Printf("成功解析: %s => %v\n", host, ips[:])

		if len(ips) == 1 {
			results <- ips[0].String()
			doresults <- host
		}
	}
}

// 保存文件
func outputf(oresults <-chan string) {
	now := time.Now()
	filenametime := now.Format("20060102_15")

	var ipsresults []string
	for ipre := range oresults {
		ipsresults = append(ipsresults, ipre)
	}

	file, err := os.OpenFile(output+"results_"+filenametime+".txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		fmt.Println("创建文件失败", err)
		return
	}
	defer file.Close()

	newwriter := bufio.NewWriter(file)

	for _, result := range ipsresults {
		_, err := newwriter.WriteString(result + "\n")
		if err != nil {
			fmt.Println("写入失败", err)
			return
		}
	}

	err = newwriter.Flush()
	if err != nil {
		return
	}

	fmt.Println("成功保存" + output + "results_" + filenametime + ".txt")
}

func main() {
	flag.Parse()

	if domainfile == "" {

		fmt.Println(`
  _____   ___    _  __   ___    ____
 / ___/  / _ \  / |/ /  / _ \  / __/
/ /__   / // / /    /  / // / / _/  
\___/  /____/ /_/|_/  /____/ /___/   
                              Version: 1.2.1									   
参数:
  -f  打开文件
  -o  输出
  -t  保存类型
  -g  线程
	
-t 类型:
  ip     以IP地址的方式保存
  domain 以域名的方式保存

		`)
		return
	}

	openfile, err := os.Open(domainfile)
	if err != nil {
		fmt.Println("打开文件失败", err)
		return
	}
	defer openfile.Close()

	var deli = make(map[string]struct{})

	newscanner := bufio.NewScanner(openfile)

	for newscanner.Scan() {
		if newscanner.Text() != "" {
			url := strings.TrimPrefix(newscanner.Text(), "http://")
			url = strings.TrimPrefix(url, "https://")
			deli[url] = struct{}{}
		}
	}
	if err := newscanner.Err(); err != nil {
		return
	}

	var domains = make(chan string, len(deli))
	var results = make(chan string, max(100, len(deli)))
	var doresults = make(chan string, max(100, len(deli)))

	go func() {
		for domain_value := range deli {
			domains <- domain_value
		}
		close(domains)
	}()

	var wg sync.WaitGroup

	for i := 0; i < thread; i++ {
		wg.Add(1)
		go nslookup(domains, results, doresults, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
		close(doresults)
	}()

	// 保存 ip 去重
	var ipres = make(map[string]struct{})
	var ipsresultschan = make(chan string, max(100, len(deli)))

	if output != "" {
		switch outputtype {
		case "ip":
			for results_value := range results {
				ipres[results_value] = struct{}{}
			}
			go func() {
				for ip := range ipres {
					fmt.Println(ip)
					ipsresultschan <- ip
				}
				close(ipsresultschan)
			}()

			outputf(ipsresultschan)

		case "domain":
			outputf(doresults)
		}
	} else {
		for range results {
		}
	}
}
