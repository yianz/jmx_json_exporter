package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"github.com/fatalc/jmx_json_exporter/utils"
	"encoding/json"
	"strings"
	"github.com/fatalc/jmx_json_exporter/collector"
)

// 配置需要监控的数据项 todo：现在仅支持两层嵌套,待改进
var (
	masterConfig = map[string][]string{
		"Hadoop:service=NameNode,name=NameNodeInfo": {
			"Total",
			"Used",
			"Free",
			"NonDfsUsedSpace",
		},
	}
	workerConfig = map[string][]string{
		"Hadoop:service=DataNode,name=JvmMetrics": {
			"MemNonHeapCommittedM",
			"LogWarn",
		},
	}
)

type HadoopCollector struct {
	masterHosts       map[string]string
	nodeHosts         map[string]string
	masterCollectors  map[string]*collector.CommonCollector
	workersCollectors map[string]*collector.CommonCollector
}

func (hc *HadoopCollector) Describe(ch chan<- *prometheus.Desc) {
	for k, v := range hc.masterCollectors {
		log.Printf("Describe of %s", k)
		v.Describe(ch)
	}
	for k, v := range hc.workersCollectors {
		log.Printf("Describe of %s", k)
		v.Describe(ch)
	}
}

//Collect implements the prometheus.Collector interface. 该接口调用来更新数据
func (hc *HadoopCollector) Collect(ch chan<- prometheus.Metric) {
	for k, v := range hc.masterCollectors {
		log.Printf("Collect of %s", k)
		v.Collect(ch)
	}
	for k, v := range hc.workersCollectors {
		log.Printf("Collect of %s", k)
		v.Collect(ch)
	}
}

func NewHadoopCollector(masterHosts map[string]string) *HadoopCollector {
	nodeHosts := getNodeHosts(masterHosts)
	mcs := make(map[string]*collector.CommonCollector, len(masterHosts))
	for _, v := range masterHosts {
		// todo: 使用hostname代替ip
		mcs[v] = collector.NewBeansCollector(v, "master", masterConfig)
	}
	wcs := make(map[string]*collector.CommonCollector, len(nodeHosts))
	for k, v := range nodeHosts {
		wcs[v] = collector.NewBeansCollector(v, k, workerConfig)
	}
	return &HadoopCollector{
		masterHosts:       masterHosts,
		nodeHosts:         nodeHosts,
		masterCollectors:  mcs,
		workersCollectors: wcs,
	}
}

func getNodeHosts(masterHosts map[string]string) map[string]string {
	const protocol = "http://"
	const nameNodeInfo = "Hadoop:service=NameNode,name=NameNodeInfo"
	const path = "/jmx?qry=" + nameNodeInfo
	const liveNodesName = "LiveNodes"
	const infoKey = "infoAddr"

	nodeUrls := make(map[string]string)
	for _, v := range masterHosts {
		liveNodes := utils.JmxJsonBeansParse(utils.Get(protocol + v + path))[nameNodeInfo].Content[liveNodesName].(string)
		nodesJson := make(map[string]interface{})
		json.Unmarshal([]byte(strings.Trim(liveNodes, "/")), &nodesJson)
		for k, v := range nodesJson {
			nodeUrl := v.(map[string]interface{})[infoKey].(string)
			nodeUrls[k] = nodeUrl
		}
	}
	return nodeUrls
}