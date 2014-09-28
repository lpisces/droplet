package main

import (
  "fmt"
  "github.com/codegangsta/cli"
  "os"
  "net/http"
  "io/ioutil"
  "strings"
  "strconv"
  //"time"
)

const(
  stockDate = iota
  stockOpen
  stockHigh
  stockLow
  stockClose
  stockVolume
  stockAdjClose
  stockEma250
  stockEma12
  stockEma26
  stockDiff
  stockDea
  stockMacd
  stockEma5
  stockEma10
  stockEma20
  stockEma30
  stockUsable
  stockPos
)


func main() {

  app := cli.NewApp()
  app.Name = "droplet"
  app.Usage = "A analysis tool for Chinese stock"
  app.Version = "0.0.1"

  app.Commands = []cli.Command{
    {
      Name:  "test",
      ShortName: "t",
      Usage: "for test",
      Action: test,
    },
    {
      Name:  "fetch",
      ShortName: "f",
      Usage: "fetch stock history data",
      Action: fetch,
    },
    {
      Name:  "analysis",
      ShortName: "a",
      Usage: "analysis stock history data",
      Action: analysis,
    },
  }

  app.Run(os.Args)

}

func test(c *cli.Context) {
  fmt.Println("Test OK !")
}

func fetch(c *cli.Context) {
  fmt.Println("Test OK !")
}


func prepare_data(code []byte) (map[int][]string){
  _code := "000000"
  if (string(code[0]) == "6") {
    _code = string(code) + ".SS"
  } else {
    _code = string(code) + ".SZ"
  }

  url := "http://ichart.yahoo.com/table.csv?s=" + _code

  resp, _ := http.Get(url)
  body, _ := ioutil.ReadAll(resp.Body)
  defer resp.Body.Close()

  hash := make(map[int][]string)
  cnt := 0
  for _, v := range(strings.Split(string(body), "\n")) {
    cnt++
    if v == "" || cnt == 1{
      continue
    }
    row := strings.Split(v, ",")
    //t, _:= time.Parse("2006-01-02", row[0])
    hash[cnt] = row
  }

  for i := len(hash)-1; i > 1; i-- {
    a := hash[i][stockAdjClose]
    _c, _:= strconv.ParseFloat(a, 64)
    if i >= len(hash) -2 {
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", _c))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", _c))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", _c))
      hash[i] = append(hash[i], "0")
      hash[i] = append(hash[i], "0")
      hash[i] = append(hash[i], "0")
      hash[i] = append(hash[i], "0")
      hash[i] = append(hash[i], "0")
      hash[i] = append(hash[i], "0")
      hash[i] = append(hash[i], "0")
      hash[i] = append(hash[i], "1")
      hash[i] = append(hash[i], "0")
    } else {
      lastEma250, _ := strconv.ParseFloat(strings.Trim(hash[i+1][stockEma250], " "), 64)
      lastEma12, _ := strconv.ParseFloat(strings.Trim(hash[i+1][stockEma12], " "), 64)
      lastEma26, _ := strconv.ParseFloat(strings.Trim(hash[i+1][stockEma26], " "), 64)
      lastEma5, _ := strconv.ParseFloat(strings.Trim(hash[i+1][stockEma5], " "), 64)
      lastEma10, _ := strconv.ParseFloat(strings.Trim(hash[i+1][stockEma10], " "), 64)
      lastEma20, _ := strconv.ParseFloat(strings.Trim(hash[i+1][stockEma20], " "), 64)
      lastEma30, _ := strconv.ParseFloat(strings.Trim(hash[i+1][stockEma30], " "), 64)
      lastDea, _ := strconv.ParseFloat(strings.Trim(hash[i+1][stockDea], " "), 64)
      ema250 := (2.0/(9+1)) * (_c - lastEma250) + lastEma250
      ema12 := (2.0/(12+1)) * (_c - lastEma12) + lastEma12
      ema26 := (2.0/(26+1)) * (_c - lastEma26) + lastEma26
      ema5 := (2.0/(5+1)) * (_c - lastEma5) + lastEma5
      ema10 := (2.0/(10+1)) * (_c - lastEma10) + lastEma10
      ema20 := (2.0/(20+1)) * (_c - lastEma20) + lastEma20
      ema30 := (2.0/(30+1)) * (_c - lastEma30) + lastEma30
      diff := ema12 - ema26
      dea := (2.0/(9+1)) * (diff - lastDea) + lastDea
      macd := 2.0 * (diff -dea)
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", ema250))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", ema12))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", ema26))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", diff))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", dea))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", macd))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", ema5))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", ema10))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", ema20))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", ema30))
      hash[i] = append(hash[i], "1")
      hash[i] = append(hash[i], "0")

      if hash[i][stockHigh] >= hash[i+1][stockHigh] && hash[i][stockLow] <= hash[i+1][stockLow] {
        hash[i+1][stockUsable] = "0"
        if hash[i][stockEma5] >= hash[i+1][stockEma5] {
          hash[i][stockLow] = hash[i+1][stockLow]
        } else {
          hash[i][stockHigh] = hash[i+1][stockHigh]
        }
      }

      if hash[i][stockHigh] <= hash[i+1][stockHigh] && hash[i][stockLow] >= hash[i+1][stockLow] {
        hash[i+1][stockUsable] = "0"
        if hash[i][stockEma5] >= hash[i+1][stockEma5] {
          hash[i][stockHigh] = hash[i+1][stockHigh]
        } else {
          hash[i][stockLow] = hash[i+1][stockLow]
        }
      }

    }
  }

  return hash
}

func analysis(c *cli.Context) {
  if len(c.Args()) < 1 {
    fmt.Println("Stock code should be given .")
  }

  code := []byte(c.Args()[0])

  hash := prepare_data(code)
  for i := len(hash)-1; i > 1; i-- {
    if i >= len(hash) -2 {
      continue
    }

    if hash[i][stockUsable] == "0" {
      continue
    }

    pre := hash[len(hash) - 1]
    n := 1
    pos := i+n
    for {
      if i+n > len(hash) -1 {
        pre = hash[len(hash) - 1]
        pos = len(hash) - 1
        break
      }
      if hash[i+n][stockUsable] != "0" {
        pre = hash[i+n]
        pos = i+n
        break
      }
      n++
    }

    prepre := hash[len(hash) - 2]
    n++
    for {
      if i+n > len(hash) -2 {
        prepre = hash[len(hash) - 2]
        break
      }
      if hash[i+n][stockUsable] != "0" {
        prepre = hash[i+n]
        break
      }
      n++
    }

    preprepre := hash[len(hash) - 3]
    n++
    for {
      if i+n > len(hash) -3 {
        preprepre = hash[len(hash) - 3]
        break
      }
      if hash[i+n][stockUsable] != "0" {
        preprepre = hash[i+n]
        break
      }
      n++
    }

    match := false
    if hash[i][stockLow] >= pre[stockLow] && hash[i][stockHigh] >= pre[stockHigh] && 
       prepre[stockLow] >= pre[stockLow] && prepre[stockHigh] >= pre[stockHigh] &&
       preprepre[stockLow] >= prepre[stockLow] && preprepre[stockHigh] >= prepre[stockHigh] {
       hash[pos][stockPos] = "-1"
       match = true
    }

    if hash[i][stockLow] <= pre[stockLow] && hash[i][stockHigh] <= pre[stockHigh] && 
       prepre[stockLow] <= pre[stockLow] && prepre[stockHigh] <= pre[stockHigh] &&
       preprepre[stockLow] <= prepre[stockLow] && preprepre[stockHigh] <= prepre[stockHigh] {
       hash[pos][stockPos] = "1"
       match = true
    }

    if match {
      n := 1
      for {
        if pos+n > len(hash) -1 {
          break
        }
        if hash[pos+n][stockPos] == "0" {
          n++
          continue
        }
        if hash[pos+n][stockPos] == "1" && hash[pos][stockPos] == "-1"{
          break
        } else if hash[pos+n][stockPos] == "-1" && hash[pos][stockPos] == "1"{
          break
        } else {
          hash[pos+n][stockPos] = "0"
          break
        }
        n++
      }
    }
  }

  for i := len(hash)-1; i > 1; i-- {
    if i >= len(hash) -2 {
      continue
    }

    if hash[i][stockUsable] == "0" {
      continue
    }

    if hash[i][stockPos] != "0" {
      fmt.Println(hash[i])
    }
  }

}
