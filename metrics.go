package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	databaseAvailable = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_available_gauge",
			Help: "Value 1 means database is active, 0 means it is not",
		},
	)

	databaseLocked = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "database_locked_gauge",
			Help: "Value 1 means database is locked, 0 means it is not",
		},
	)

	clientsNumber = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "clients_number_gauge",
			Help: "Shows number of clients",
		},
	)

	commitLatencyStatistics = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "commit_latency_statistics_gauge",
			Help: "Shows commit latency statistics",
		},
		[]string{"process", "measure"},
	)

	readLatencyStatistics = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "read_latency_statistics_gauge",
			Help: "Shows read latency statistics",
		},
		[]string{"process", "measure"},
	)

	operations = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "operations_gauge",
			Help: "Number of certain real-time operations",
		},
		[]string{"type"},
	)

	transactions = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "transactions_gauge",
			Help: "Number of transactions started/committed",
		},
		[]string{"status"},
	)

	keysRead = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "keys_read_gauge",
			Help: "Number of keys read",
		},
	)

	bytesRW = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "bytes_RW_gauge",
			Help: "Number of bytes read/written",
		},
		[]string{"operation"},
	)

	transactionsPerSecondLimit = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "transactions_per_second_limit_gauge",
			Help: "Limit of transactions per second",
		},
		[]string{"type"},
	)

	releasedTransactionsPerSecond = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "released_transactions_per_second_gauge",
			Help: "Number of released transactions per second",
		},
		[]string{"type"},
	)

	storageServerDurabilityLag = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "storage_server_durability_lag_gauge",
			Help: "Storage server durability lag info",
		},
		[]string{"type"},
	)

	storageServerDataLag = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "storage_server_data_lag_gauge",
			Help: "Storage server data lag info",
		},
		[]string{"type"},
	)

	storageServerQueue = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "storage_server_queue_gauge",
			Help: "Storage server queue info",
		},
		[]string{"type"},
	)

	logServerQueue = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "log_server_queue_gauge",
			Help: "Log server worst queue bytes info",
		},
	)

	performanceLimitedBy = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "performance_limited_by_gauge",
			Help: "Shows status of performance limit",
		},
		[]string{"status"},
	)

	batchPerformanceLimitedBy = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "batch_performance_limited_by_gauge",
			Help: "Shows status of batch performance limit",
		},
		[]string{"status"},
	)

	totalDiskUsedBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "total_disk_used_bytes_gauge",
			Help: "Number of total disk used bytes",
		},
	)

	totalKVSizeBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "total_kv_size_bytes_gauge",
			Help: "Number of total KV size bytes",
		},
	)

	movingDataInQueueBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "moving_data_in_queue_bytes_gauge",
			Help: "Number of moving data in queue bytes",
		},
	)

	movingDataInFlightBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "moving_data_in_flight_bytes_gauge",
			Help: "Number of moving data in flight bytes",
		},
	)

	movingDataTotalWrittenBytes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "moving_data_total_written_bytes_gauge",
			Help: "Number of moving data total written bytes",
		},
	)

	systemKVSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "system_kv_size_bytes_gauge",
			Help: "System KV size in bytes",
		},
	)

	averagePartitionSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "average_partition_size_bytes_gauge",
			Help: "Average partition size in bytes",
		},
	)

	partitionsCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "partitions_count_gauge",
			Help: "Number of partitions",
		},
	)

	dataState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "data_state_gauge",
			Help: "Shows data state",
		},
		[]string{"status"},
	)

	leastOperatingSpaceLogServer = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "least_operating_space_bytes_log_server_gauge",
			Help: "Shows least operating space for log server in bytes",
		},
	)

	leastOperatingSpaceStorageServer = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "least_operating_space_bytes_storage_server_gauge",
			Help: "Show least operating space for storage server in bytes",
		},
	)

	//Class
	classType = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "class_type_gauge",
			Help: "Shows the class type of the address",
		},
		[]string{"classType", "address"},
	)

	//CPU
	cpuUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage_gauge",
			Help: "Shows percentage of cpu usage of the address",
		},
		[]string{"address"},
	)

	//Memory + Memory Used
	memoryUsedBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_used_bytes_gauge",
			Help: "Shows used memory used in bytes",
		},
		[]string{"address"},
	)

	//Memory Used
	memoryLimitBytes = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_limit_bytes_gauge",
			Help: "Shows memory limit in bytes",
		},
		[]string{"address"},
	)

	//Disk IO
	diskReads = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_disk_reads_gauge",
			Help: "Shows total number of disk reads",
		},
		[]string{"address"},
	)

	//Disk IO
	diskWrites = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "total_disk_writes_gauge",
			Help: "Shows total number of disk writes",
		},
		[]string{"address"},
	)

	//Disk used
	diskUsed = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_used_gauge",
			Help: "Shows percentage of disk used",
		},
		[]string{"address"},
	)

	//Disk Utilization
	diskUtilization = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_utilization_gauge",
			Help: "Show percentage of disk utilization",
		},
		[]string{"address"},
	)

	//Network IO
	networkMegabitsReceivedRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "network_megabits_received_rate_gauge",
			Help: "Shows the rate of megabits received",
		},
		[]string{"address"},
	)

	//Network IO
	networkMegabitsSentRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "network_megabits_sent_rate_gauge",
			Help: "Shows the rate of megabits sent",
		},
		[]string{"address"},
	)

	//Network Connections
	networkCurrentConnections = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "network_current_connections_gauge",
			Help: "Shows number of current connections",
		},
		[]string{"address"},
	)

	//Run Loop Utilization
	runloopUtilization = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "runloop_utilization_gauge",
			Help: "Shows percentage of runloop utilization",
		},
		[]string{"address"},
	)
)
