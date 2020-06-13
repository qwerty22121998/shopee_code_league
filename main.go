package main

import (
	"encoding/csv"
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
	sort.Strings(result)
	return result
}

func handle(err error) {
	panic(err)
}

func main() {

	result := make(map[string][]string)

	file, _ := os.Open("order_brush_order.csv")
	data, _ := csv.NewReader(file).ReadAll()
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
	defer output.Close()
	var err error
	for sid, r := range result {
		r = uniqueArray(r)

		//color.Blue("%v",r)
		if len(r) > 0 {
			err = writer.Write([]string{sid, strings.Join(r, "&")})
		} else {
			err = writer.Write([]string{sid, "0"})
		}
		if err != nil {
			log.Println(sid, r)
			panic(err)
		}
		//if r != nil && len(r) != 0 {
		//		//	number++
		//		//}
	}
	log.Println(number)

}
