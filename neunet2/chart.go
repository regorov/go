package neunet2

import (
    "io"
    "net/http"
    "net/url"
    "os"
    "strconv"
    "strings"
)

var Colours = []string{
    "FF0000",
    "FFFF00",
    "00FF00",
    "00FFFF",
    "0000FF",
    "FF00FF",
}

type Chart struct {
    Title    string
    Datasets [][]float32
}

func NewChart() (chart *Chart) {
    chart = new(Chart)
    chart.Datasets = make([][]float32, 0)
    return chart
}

func (chart *Chart) Data() (data url.Values) {
    data = make(url.Values)
    data.Add("cht", "lxy")
    data.Add("chs", "300x400")
    data.Add("chtt", chart.Title)
    data.Add("chds", "a")

    nextColour := 0
    strcols := make([]string, len(chart.Datasets))
    strsets := make([]string, len(chart.Datasets))

    for j, dataset := range chart.Datasets {
        strset := make([]string, len(dataset))

        for i := 0; i < len(dataset); i++ {
            strset[i] = strconv.FormatFloat(float64(dataset[i]), 'f', -1, 32)
        }

        strcols[j] = Colours[nextColour]
        strsets[j] = strings.Join(strset, ",")

        nextColour = (nextColour + 1) % len(Colours)
    }

    data.Add("chco", strings.Join(strcols, ","))
    data.Add("chd", "t:"+strings.Join(strsets, "|"))

    return data
}

func (chart *Chart) URL() (url string) {
    return "https://chart.googleapis.com/chart?" + chart.Data().Encode()
}

func (chart *Chart) Download(fname string) (err error) {
    f, err := os.Create(fname)
    if err != nil {
        return err
    }
    defer f.Close()

    resp, err := http.Get(chart.URL())
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    io.Copy(f, resp.Body)

    return nil
}
