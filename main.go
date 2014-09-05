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
)


func main() {

  app := cli.NewApp()
  app.Name = "droplet"
  app.Usage = "A stock analysis tool for Chinese stock"
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


func prepare_data(code []byte) (map[int][]string, [][]string){
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

  var points [][]string
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
    } else {
      lastEma250, _ := strconv.ParseFloat(strings.Trim(hash[i+1][stockEma250], " "), 64)
      lastEma12, _ := strconv.ParseFloat(strings.Trim(hash[i+1][stockEma12], " "), 64)
      lastEma26, _ := strconv.ParseFloat(strings.Trim(hash[i+1][stockEma26], " "), 64)
      //lastDiff, _ := strconv.ParseFloat(strings.Trim(hash[i+1][stockDiff], " "), 64)
      lastDea, _ := strconv.ParseFloat(strings.Trim(hash[i+1][stockDea], " "), 64)
      ema250 := (2.0/(9+1)) * (_c - lastEma250) + lastEma250
      ema12 := (2.0/(12+1)) * (_c - lastEma12) + lastEma12
      ema26 := (2.0/(26+1)) * (_c - lastEma26) + lastEma26
      diff := ema12 - ema26
      dea := (2.0/(9+1)) * (diff - lastDea) + lastDea
      macd := 2.0 * (diff -dea)
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", ema250))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", ema12))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", ema26))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", diff))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", dea))
      hash[i] = append(hash[i], fmt.Sprintf("%6.3f", macd))

      pre_macd, _:= strconv.ParseFloat(strings.Trim(hash[i+1][stockMacd], " "), 64)
      prepre_macd, _:= strconv.ParseFloat(strings.Trim(hash[i+2][stockMacd], " "), 64)
      if macd <= pre_macd && pre_macd >= prepre_macd{
        points = append(points, hash[i])
      }
      if macd >= pre_macd && pre_macd <= prepre_macd{
        points = append(points, hash[i])
      }
/*
      if diff > 0 && macd < 0 && macd > pre_macd && hash[i][stockLow] < hash[i+1][stockLow]{
        points = append(points, hash[i])
      }
*/
    }
  }

  return hash, points
}

func analysis(c *cli.Context) {
  if len(c.Args()) < 1 {
    fmt.Println("Stock code should be given .")
  }

  code := []byte(c.Args()[0])

  _, points:= prepare_data(code)

  for k,v := range(points) {
    if k < 5 {
      continue
    }
    macd, _:= strconv.ParseFloat(strings.Trim(v[stockMacd], " "), 64)
    stockLow, _:= strconv.ParseFloat(strings.Trim(v[stockClose], " "), 64)
    preMacd, _:= strconv.ParseFloat(strings.Trim(points[k-1][stockMacd], " "), 64)
    preStockLow, _:= strconv.ParseFloat(strings.Trim(points[k-1][stockClose], " "), 64)
    //prepre_macd, _:= strconv.ParseFloat(strings.Trim(points[k-2][stockMacd], " "), 64)
    if macd < 0 && preMacd < 0 && (macd > preMacd || v[stockDiff] > points[k-1][stockDiff]) && stockLow < preStockLow{
      fmt.Println("=")
      fmt.Println(points[k-1])
      fmt.Println(v)
      fmt.Println("=")
    }

    //fmt.Println(macd)
  }
}
