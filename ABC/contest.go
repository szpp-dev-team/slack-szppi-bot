package ABC

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Post struct {
	ID               string `json:"id"`
	StartEpochSecond int    `json:"start_epoch_second"`
	DurationSecond   int64  `json:"duration_second"`
	Title            string `json:"title"`
	RateChange       string `json:"rate_change"`
}

func max(a int, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func Id() int {
	var posts []Post

	resp, _ := http.Get("https://kenkoooo.com/atcoder/resources/contests.json")

	body, _ := io.ReadAll(resp.Body)

	json.Unmarshal(body, &posts)
	res := 1
	for idx := range posts {
		post_id := posts[idx].ID
		if strings.Count(post_id, "abc") > 0 {
			id_num, _ := strconv.Atoi(posts[idx].ID[3:])
			res = max(res, id_num)
		}
	}
	fmt.Println(res)

	defer resp.Body.Close()
	return res
}
