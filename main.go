package main

import (
	"encoding/json"
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	_ "net/http"
	"time"
)

var (
	StatusKey = append([]byte{255, 255}, []byte("/status/json")...)
)

func getFDBStatus(db *fdb.Database) *FDBStatus {
	statusJSON, err := db.ReadTransact(func(tr fdb.ReadTransaction) (interface{}, error) {
		return tr.Get(fdb.Key(StatusKey)).Get()
	})
	if err != nil {
		panic("Can't get FDB status json")
	}
	fmt.Printf("Status JSON: \n%s\n\n", statusJSON)
	var FDBstatus FDBStatus
	err = json.Unmarshal([]byte(statusJSON.([]byte)), &FDBstatus)
	if err != nil {
		panic("Can't deserialize json")
	}
	return &FDBstatus
}

/*
func getFDBStatus(db *fdb.Database) *FDBStatus {
	f, err := os.Open("status.txt")
	if err != nil {
		panic(err)
	}
	statusJSON, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Status JSON: \n%s\n", statusJSON)
	var FDBstatus FDBStatus
	err = json.Unmarshal(statusJSON, &FDBstatus)
	if err != nil {
		panic("Can't deserialize json")
	}
	return &FDBstatus
}*/

func updateMetrics(status *FDBStatus) {
	if status.Cluster.DatabaseAvailable {
		databaseAvailable.Set(1)
	} else {
		databaseAvailable.Set(0)
	}

	if status.Cluster.DatabaseLocked {
		databaseLocked.Set(1)
	} else {
		databaseLocked.Set(0)
	}

	clientsNumber.Set(status.Cluster.Clients.Count)

	for processname, process := range status.Cluster.Processes {
		var (
			commitProxy *DynamicCommitProxyRole
			storage     *DynamicStorageRole
		)
		for _, value := range process.Roles {
			roleType := fmt.Sprintf("%T", value.Value)
			//fmt.Println(roleType)
			switch roleType {
			case "*main.DynamicCommitProxyRole":
				commitProxy = value.Value.(*DynamicCommitProxyRole)
			case "*main.DynamicStorageRole":
				storage = value.Value.(*DynamicStorageRole)
			}
		}
		if commitProxy == nil || storage == nil {
			continue
		}
		//fmt.Printf("Process: %s CommitProxy: role - %s\n", processname, commitProxy.Role)
		//fmt.Printf("Process: %s Storage: role - %s\n", processname, storage.Role)
		commitLatencyStatistics.With(prometheus.Labels{"process": processname, "measure": "median"}).Set(commitProxy.CommitLatencyStatistics.Median)
		commitLatencyStatistics.With(prometheus.Labels{"process": processname, "measure": "p25"}).Set(commitProxy.CommitLatencyStatistics.P25)
		commitLatencyStatistics.With(prometheus.Labels{"process": processname, "measure": "p90"}).Set(commitProxy.CommitLatencyStatistics.P90)
		commitLatencyStatistics.With(prometheus.Labels{"process": processname, "measure": "p95"}).Set(commitProxy.CommitLatencyStatistics.P95)
		commitLatencyStatistics.With(prometheus.Labels{"process": processname, "measure": "p99"}).Set(commitProxy.CommitLatencyStatistics.P99)
		commitLatencyStatistics.With(prometheus.Labels{"process": processname, "measure": "p99.9"}).Set(commitProxy.CommitLatencyStatistics.P99_9)

		readLatencyStatistics.With(prometheus.Labels{"process": processname, "measure": "median"}).Set(storage.ReadLatencyStatistics.Median)
		readLatencyStatistics.With(prometheus.Labels{"process": processname, "measure": "p25"}).Set(storage.ReadLatencyStatistics.P25)
		readLatencyStatistics.With(prometheus.Labels{"process": processname, "measure": "p90"}).Set(storage.ReadLatencyStatistics.P90)
		readLatencyStatistics.With(prometheus.Labels{"process": processname, "measure": "p95"}).Set(storage.ReadLatencyStatistics.P95)
		readLatencyStatistics.With(prometheus.Labels{"process": processname, "measure": "p99"}).Set(storage.ReadLatencyStatistics.P99)
		readLatencyStatistics.With(prometheus.Labels{"process": processname, "measure": "p99.9"}).Set(storage.ReadLatencyStatistics.P99_9)
	}

	for _, value := range status.Cluster.Processes {
		address := value.Address
		classType.With(prometheus.Labels{"address": address, "classType": value.ClassType}).Set(1)

		cpuUsage.With(prometheus.Labels{"address": address}).Set(value.CPU.UsageCores)

		memoryUsedBytes.With(prometheus.Labels{"address": address}).Set(value.Memory.UsedBytes)

		memoryLimitBytes.With(prometheus.Labels{"address": address}).Set(float64(value.Memory.LimitBytes))
		fmt.Printf("address: %s | memory used bytes: %f | memory limit bytes: %d | memory used: %f\n", address, value.Memory.UsedBytes, value.Memory.LimitBytes,
			value.Memory.UsedBytes/float64(value.Memory.LimitBytes))

		diskReads.With(prometheus.Labels{"address": address}).Set(value.Disk.Reads.Counter)

		diskWrites.With(prometheus.Labels{"address": address}).Set(value.Disk.Writes.Counter)

		diskUsed.With(prometheus.Labels{"address": address}).Set(1.0 - float64(value.Disk.FreeBytes)/float64(value.Disk.TotalBytes))
		fmt.Printf("address: %s | free bytes: %d | total bytes: %d | disk used: %f", address, value.Disk.FreeBytes, value.Disk.TotalBytes,
			1.0-float64(value.Disk.FreeBytes)/float64(value.Disk.TotalBytes))

		diskUtilization.With(prometheus.Labels{"address": address}).Set(value.Disk.Busy)

		networkMegabitsReceivedRate.With(prometheus.Labels{"address": address}).Set(value.Network.MegabitsReceived.Hz)

		networkMegabitsSentRate.With(prometheus.Labels{"address": address}).Set(value.Network.MegabitsSent.Hz)

		networkCurrentConnections.With(prometheus.Labels{"address": address}).Set(value.Network.CurrentConnections)

		runloopUtilization.With(prometheus.Labels{"address": address}).Set(value.RunLoopBusy)
	}

	operations.With(prometheus.Labels{"type": "read_requests"}).Set(status.Cluster.Workload.Operations.ReadRequests.Counter)
	operations.With(prometheus.Labels{"type": "reads"}).Set(status.Cluster.Workload.Operations.Reads.Counter)
	operations.With(prometheus.Labels{"type": "writes"}).Set(status.Cluster.Workload.Operations.Writes.Counter)

	transactions.With(prometheus.Labels{"status": "started"}).Set(status.Cluster.Workload.Transactions.Started.Counter)
	transactions.With(prometheus.Labels{"status": "committed"}).Set(status.Cluster.Workload.Transactions.Committed.Counter)
	transactions.With(prometheus.Labels{"status": "conflicted"}).Set(status.Cluster.Workload.Transactions.Conflicted.Counter)

	keysRead.Set(status.Cluster.Workload.Keys.Read.Counter)

	bytesRW.With(prometheus.Labels{"operation": "read"}).Set(status.Cluster.Workload.Bytes.Read.Counter)
	bytesRW.With(prometheus.Labels{"operation": "written"}).Set(status.Cluster.Workload.Bytes.Written.Counter)

	transactionsPerSecondLimit.With(prometheus.Labels{"type": "transactions_limit"}).Set(status.Cluster.Qos.TransactionsPerSecondLimit)
	transactionsPerSecondLimit.With(prometheus.Labels{"type": "batch_transactions_limit"}).Set(status.Cluster.Qos.BatchTransactionsPerSecondLimit)

	releasedTransactionsPerSecond.With(prometheus.Labels{"type": "released_transactions"}).Set(status.Cluster.Qos.ReleasedTransactionsPerSecond)
	releasedTransactionsPerSecond.With(prometheus.Labels{"type": "batch_released_transactions"}).Set(status.Cluster.Qos.BatchReleasedTransactionsPerSecond)
	fmt.Printf("released transactions per second: %f\n", status.Cluster.Qos.ReleasedTransactionsPerSecond)
	fmt.Printf("batch released transactions per second: %f\n", status.Cluster.Qos.BatchReleasedTransactionsPerSecond)

	storageServerDurabilityLag.With(prometheus.Labels{"type": "limiting"}).Set(status.Cluster.Qos.LimitingDurabilityLagStorageServer.Seconds)
	storageServerDurabilityLag.With(prometheus.Labels{"type": "worst"}).Set(status.Cluster.Qos.WorstDurabilityLagStorageServer.Seconds)

	storageServerDataLag.With(prometheus.Labels{"type": "limiting"}).Set(status.Cluster.Qos.LimitingDataLagStorageServer.Seconds)
	storageServerDataLag.With(prometheus.Labels{"type": "worst"}).Set(status.Cluster.Qos.WorstDataLagStorageServer.Seconds)

	storageServerQueue.With(prometheus.Labels{"type": "limiting"}).Set(status.Cluster.Qos.LimitingQueueBytesStorageServer)
	storageServerQueue.With(prometheus.Labels{"type": "worst"}).Set(status.Cluster.Qos.WorstQueueBytesStorageServer)

	logServerQueue.Set(status.Cluster.Qos.WorstQueueBytesLogServer)

	performanceLimitedBy.With(prometheus.Labels{"status": status.Cluster.Qos.PerformanceLimitedBy.Description}).Set(1)
	//performanceLimitedBy.With(prometheus.Labels{"status": "saturated"}).Set(1)

	batchPerformanceLimitedBy.With(prometheus.Labels{"status": status.Cluster.Qos.BatchPerformanceLimitedBy.Description}).Set(1)
	//batchPerformanceLimitedBy.With(prometheus.Labels{"status": "saturated"}).Set(1)

	totalDiskUsedBytes.Set(status.Cluster.Data.TotalDiskUsedBytes)

	totalKVSizeBytes.Set(status.Cluster.Data.TotalKvSizeBytes)

	movingDataInQueueBytes.Set(status.Cluster.Data.MovingData.InQueueBytes)

	movingDataInFlightBytes.Set(status.Cluster.Data.MovingData.InFlightBytes)

	movingDataTotalWrittenBytes.Set(status.Cluster.Data.MovingData.TotalWrittenBytes)

	systemKVSize.Set(status.Cluster.Data.SystemKvSizeBytes)

	averagePartitionSize.Set(status.Cluster.Data.AveragePartitionSizeBytes)

	partitionsCount.Set(status.Cluster.Data.PartitionsCount)

	dataState.With(prometheus.Labels{"status": status.Cluster.Data.State.Name}).Set(1)

	leastOperatingSpaceLogServer.Set(float64(status.Cluster.Data.LeastOperatingSpaceBytesLogServer))

	leastOperatingSpaceStorageServer.Set(float64(status.Cluster.Data.LeastOperatingSpaceBytesStorageServer))
}

func main() {
	fdb.MustAPIVersion(700)

	fmt.Println("Starting fdb-exporter...")

	prometheus.MustRegister(operations)
	prometheus.MustRegister(transactions)
	prometheus.MustRegister(keysRead)
	prometheus.MustRegister(bytesRW)
	prometheus.MustRegister(transactionsPerSecondLimit)
	prometheus.MustRegister(releasedTransactionsPerSecond)
	prometheus.MustRegister(storageServerDurabilityLag)
	prometheus.MustRegister(storageServerDataLag)
	prometheus.MustRegister(storageServerQueue)
	prometheus.MustRegister(logServerQueue)
	prometheus.MustRegister(databaseAvailable)
	prometheus.MustRegister(databaseLocked)
	prometheus.MustRegister(clientsNumber)
	prometheus.MustRegister(commitLatencyStatistics)
	prometheus.MustRegister(readLatencyStatistics)
	prometheus.MustRegister(performanceLimitedBy)
	prometheus.MustRegister(batchPerformanceLimitedBy)
	prometheus.MustRegister(totalDiskUsedBytes)
	prometheus.MustRegister(totalKVSizeBytes)
	prometheus.MustRegister(movingDataInQueueBytes)
	prometheus.MustRegister(movingDataInFlightBytes)
	prometheus.MustRegister(movingDataTotalWrittenBytes)
	prometheus.MustRegister(systemKVSize)
	prometheus.MustRegister(averagePartitionSize)
	prometheus.MustRegister(partitionsCount)
	prometheus.MustRegister(dataState)
	prometheus.MustRegister(leastOperatingSpaceLogServer)
	prometheus.MustRegister(leastOperatingSpaceStorageServer)
	prometheus.MustRegister(classType)
	prometheus.MustRegister(cpuUsage)
	prometheus.MustRegister(memoryUsedBytes)
	prometheus.MustRegister(memoryLimitBytes)
	prometheus.MustRegister(diskReads)
	prometheus.MustRegister(diskWrites)
	prometheus.MustRegister(diskUsed)
	prometheus.MustRegister(diskUtilization)
	prometheus.MustRegister(networkMegabitsReceivedRate)
	prometheus.MustRegister(networkMegabitsSentRate)
	prometheus.MustRegister(networkCurrentConnections)
	prometheus.MustRegister(runloopUtilization)

	db := fdb.MustOpenDefault()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":9090", nil)
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		http.HandleFunc("/liveness", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
		})
		http.HandleFunc("/readiness", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusOK)
		})
		http.ListenAndServe(":8000", nil)
	}()

	for {
		FDBstatus := getFDBStatus(&db)
		updateMetrics(FDBstatus)
		time.Sleep(time.Second)
	}
}
