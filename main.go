package main

import (
	"encoding/csv"
	"github.com/fatih/color"
	"strconv"
	"strings"

	//"github.com/fatih/color"
	"log"
	"os"
	"sort"
	"time"
)

type Buy struct {
	Oid  string
	Sid  string
	Uid  string
	Time time.Time
}

func uniqueArray(arr []string) []string {
	if arr == nil {
		return nil
	}
	mp := make(map[string]bool)
	var result []string
	for _, v := range arr {
		if !mp[v] {
			result = append(result, v)
		}
		mp[v] = true
	}
	sort.Slice(arr, func(i, j int) bool {
		a, _ := strconv.Atoi(arr[i])
		b, _ := strconv.Atoi(arr[j])
		return a < b
	})
	return result
}

func main() {

	result := make(map[string][]string)

	file, _ := os.Open("order_brush_order.csv")
	data, _ := csv.NewReader(file).ReadAll()
	file.Close()
	data = data[1:]

	sort.Slice(data, func(i, j int) bool {
		a := data[i]
		b := data[j]
		if a[1] != b[1] {
			return a[1] < b[1]
		}
		return a[3] < b[3]
	})

	mp := make(map[string][]Buy)

	for _, v := range data {
		t, _ := time.Parse("2006-01-02 15:04:05", v[3])
		buy := Buy{
			Oid:  v[0],
			Sid:  v[1],
			Uid:  v[2],
			Time: t,
		}
		result[v[1]] = nil
		mp[v[1]] = append(mp[v[1]], buy)
	}

	for sid, event := range mp {

		//if sid != "145777302" {continue}
		//color.Red(sid)
		var queue []Buy
		for _, e := range event {
			//color.Blue("==========")
			if len(queue) == 0 {
				queue = append(queue, e)
			}
			for len(queue) > 0 && e.Time.Sub(queue[0].Time) > time.Hour {
				queue = queue[1:]
			}
			queue = append(queue, e)
			//for _, r := range queue {
			//color.Cyan(r.Time.String())
			//}

			orderNumber := len(queue)
			uniqueBuyer := 0
			buyer := make(map[string]int)
			maximum := 0
			for _, q := range queue {
				if buyer[q.Uid] == 0 {
					uniqueBuyer++
				}
				buyer[q.Uid]++
				if buyer[q.Uid] > maximum {
					maximum = buyer[q.Uid]
				}
			}
			rate := float64(orderNumber) / float64(uniqueBuyer)
			if rate >= 3.0 {
				for uid, btime := range buyer {
					if btime == maximum {
						result[sid] = append(result[sid], uid)
					}
				}
			}
			//color.Green("rate %v, order: %v, buyer: %v", rate, orderNumber, uniqueBuyer)
		}

	}
	number := 0
	output, _ := os.Create("result.csv")
	writer := csv.NewWriter(output)
	//defer output.Close()
	writer.Write([]string{"shopid", "userid"})
	var err error
	log.Println(len(result))
	for sid, r := range result {
		r = uniqueArray(r)
		if r != nil && len(r) > 0 {
			err = writer.Write([]string{sid, strings.Join(r, "&")})
			number++
		} else {
			err = writer.Write([]string{sid, "0"})
		}
		if err != nil {
			log.Println(sid, r)
			log.Fatalln(err)
		}
	}
	writer.Flush()
	output.Close()
	color.Red("Total cheated shop : %v", number)
}
