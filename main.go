package main

import (
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	kibanaConfigFile = "./config/kibana.json"
	logPath          = "./logs/go.log"
)

func main() {
	// Setup dashboards in Kibana (not required step)
	if err := setupDashboards(); err != nil {
		fmt.Printf("failed to setup Kibana dashboards, error: %s\n", err.Error())
	}

	os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0666)
	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"stdout", logPath}
	l, err := c.Build()
	if err != nil {
		panic(err)
	}
	i := 0
	for {
		i++
		time.Sleep(time.Second * 3)
		if rand.Intn(10) == 1 {
			l.Error("test error", zap.Error(fmt.Errorf("error because test: %d", i)))
		} else {
			l.Info(fmt.Sprintf("test log: %d", i))
		}
	}
}

// setupDashboards put graphs and dashboards inside kibana
func setupDashboards() error {
	f, err := os.Open(kibanaConfigFile)
	if err != nil {
		return err
	}
	defer f.Close()
	url := os.Getenv("KIBANA_URL") + "/api/kibana/dashboards/import"
	req, err := http.NewRequest("POST", url, f)
	if err != nil {
		return err
	}

	req.Header.Add("Kbn-Xsrf", "true")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("error responce from kibana: %s", string(body))
	}
	return nil
}
