package main

import (
	"bytes"
	"net"
	"sync"
	"testing"

	"github.com/pusher/buddha/tcptest"
	"golang.org/x/net/context"
)

func testFetch(t *testing.T, fixture string) {
	ts := tcptest.NewServer(func(conn net.Conn) {
		defer conn.Close()
		conn.Write([]byte(fixture))
	})
	defer ts.Close()

	ctx := context.Background()
	fetchCnt := 100

	fetchSignal := make(chan struct{}, 1)
	metricCh := make(chan []byte, 0)

	NewFetcher(ctx, "tcp", ts.Addr.String(),
		fetchSignal,
		metricCh,
		fetchCnt,
		"test")

	fetchSignal <- struct{}{}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		metrics := <-metricCh
		metricLines := bytes.Split(metrics, []byte("\n"))
		for _, m := range metricLines {
			if len(m) <= 0 {
				continue
			}
			parts := bytes.Split(m, []byte(" "))
			if len(parts) != 3 {
				t.Errorf("metric format is wrong m:%v", string(m))
			}
			t.Logf("%v", string(m))
		}

		wg.Done()
	}()

	wg.Wait()

}

func Test_Fetcher(t *testing.T) {
	testFetch(t, sample1)
	testFetch(t, sample2)

}

var sample1 = `<?xml version="1.0" encoding="ISO-8859-1" standalone="yes"?>
<!DOCTYPE GANGLIA_XML [
   <!ELEMENT GANGLIA_XML (GRID|CLUSTER|HOST)*>
      <!ATTLIST GANGLIA_XML VERSION CDATA #REQUIRED>
      <!ATTLIST GANGLIA_XML SOURCE CDATA #REQUIRED>
   <!ELEMENT GRID (CLUSTER | GRID | HOSTS | METRICS)*>
      <!ATTLIST GRID NAME CDATA #REQUIRED>
      <!ATTLIST GRID AUTHORITY CDATA #REQUIRED>
      <!ATTLIST GRID LOCALTIME CDATA #IMPLIED>
   <!ELEMENT CLUSTER (HOST | HOSTS | METRICS)*>
      <!ATTLIST CLUSTER NAME CDATA #REQUIRED>
      <!ATTLIST CLUSTER OWNER CDATA #IMPLIED>
      <!ATTLIST CLUSTER LATLONG CDATA #IMPLIED>
      <!ATTLIST CLUSTER URL CDATA #IMPLIED>
      <!ATTLIST CLUSTER LOCALTIME CDATA #REQUIRED>
   <!ELEMENT HOST (METRIC)*>
      <!ATTLIST HOST NAME CDATA #REQUIRED>
      <!ATTLIST HOST IP CDATA #REQUIRED>
      <!ATTLIST HOST LOCATION CDATA #IMPLIED>
      <!ATTLIST HOST TAGS CDATA #IMPLIED>
      <!ATTLIST HOST REPORTED CDATA #REQUIRED>
      <!ATTLIST HOST TN CDATA #IMPLIED>
      <!ATTLIST HOST TMAX CDATA #IMPLIED>
      <!ATTLIST HOST DMAX CDATA #IMPLIED>
      <!ATTLIST HOST GMOND_STARTED CDATA #IMPLIED>
   <!ELEMENT METRIC (EXTRA_DATA*)>
      <!ATTLIST METRIC NAME CDATA #REQUIRED>
      <!ATTLIST METRIC VAL CDATA #REQUIRED>
      <!ATTLIST METRIC TYPE (string | int8 | uint8 | int16 | uint16 | int32 | uint32 | float | double | timestamp) #REQUIRED>
      <!ATTLIST METRIC UNITS CDATA #IMPLIED>
      <!ATTLIST METRIC TN CDATA #IMPLIED>
      <!ATTLIST METRIC TMAX CDATA #IMPLIED>
      <!ATTLIST METRIC DMAX CDATA #IMPLIED>
      <!ATTLIST METRIC SLOPE (zero | positive | negative | both | unspecified) #IMPLIED>
      <!ATTLIST METRIC SOURCE (gmond) 'gmond'>
   <!ELEMENT EXTRA_DATA (EXTRA_ELEMENT*)>
   <!ELEMENT EXTRA_ELEMENT EMPTY>
      <!ATTLIST EXTRA_ELEMENT NAME CDATA #REQUIRED>
      <!ATTLIST EXTRA_ELEMENT VAL CDATA #REQUIRED>
   <!ELEMENT HOSTS EMPTY>
      <!ATTLIST HOSTS UP CDATA #REQUIRED>
      <!ATTLIST HOSTS DOWN CDATA #REQUIRED>
      <!ATTLIST HOSTS SOURCE (gmond | gmetad) #REQUIRED>
   <!ELEMENT METRICS (EXTRA_DATA*)>
      <!ATTLIST METRICS NAME CDATA #REQUIRED>
      <!ATTLIST METRICS SUM CDATA #REQUIRED>
      <!ATTLIST METRICS NUM CDATA #REQUIRED>
      <!ATTLIST METRICS TYPE (string | int8 | uint8 | int16 | uint16 | int32 | uint32 | float | double | timestamp) #REQUIRED>
      <!ATTLIST METRICS UNITS CDATA #IMPLIED>
      <!ATTLIST METRICS SLOPE (zero | positive | negative | both | unspecified) #IMPLIED>
      <!ATTLIST METRICS SOURCE (gmond) 'gmond'>
]>
<GANGLIA_XML VERSION="3.6.0" SOURCE="gmond">
<CLUSTER NAME="cluster1" LOCALTIME="1453984564" OWNER="unspecified" LATLONG="unspecified" URL="unspecified">
<HOST NAME="host1" IP="172.17.0.4" TAGS="" REPORTED="1453984563" TN="0" TMAX="20" DMAX="0" LOCATION="unspecified" GMOND_STARTED="1453885599">
<METRIC NAME="swap_free" VAL="16070712" TYPE="float" UNITS="KB" TN="25" TMAX="180" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="memory"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Amount of available swap memory"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Free Swap Space"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="load_one" VAL="0.09" TYPE="float" UNITS=" " TN="38" TMAX="70" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="load"/>
<EXTRA_ELEMENT NAME="DESC" VAL="One minute load average"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="One Minute Load Average"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="mem_total" VAL="16333692" TYPE="float" UNITS="KB" TN="1029" TMAX="1200" DMAX="0" SLOPE="zero">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="memory"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Total amount of memory displayed in KBs"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Memory Total"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="os_release" VAL="4.2.0-23-generic" TYPE="string" UNITS="" TN="1029" TMAX="1200" DMAX="0" SLOPE="zero">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="system"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Operating system release date"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Operating System Release"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="proc_run" VAL="0" TYPE="uint32" UNITS=" " TN="66" TMAX="950" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="process"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Total number of running processes"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Total Running Processes"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="load_five" VAL="0.25" TYPE="float" UNITS=" " TN="38" TMAX="325" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="load"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Five minute load average"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Five Minute Load Average"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="gexec" VAL="OFF" TYPE="string" UNITS="" TN="129" TMAX="300" DMAX="0" SLOPE="zero">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="core"/>
<EXTRA_ELEMENT NAME="DESC" VAL="gexec available"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Gexec Status"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="disk_free" VAL="303.665" TYPE="double" UNITS="GB" TN="128" TMAX="180" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="disk"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Total free disk space"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Disk Space Available"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="mem_cached" VAL="3554892" TYPE="float" UNITS="KB" TN="25" TMAX="180" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="memory"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Amount of cached memory"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Cached Memory"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="pkts_in" VAL="0.00" TYPE="float" UNITS="packets/sec" TN="129" TMAX="300" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="network"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Packets in per second"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Packets Received"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="bytes_in" VAL="0.00" TYPE="float" UNITS="bytes/sec" TN="129" TMAX="300" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="network"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Number of bytes in per second"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Bytes Received"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="bytes_out" VAL="34.65" TYPE="float" UNITS="bytes/sec" TN="129" TMAX="300" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="network"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Number of bytes out per second"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Bytes Sent"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="swap_total" VAL="16674812" TYPE="float" UNITS="KB" TN="1029" TMAX="1200" DMAX="0" SLOPE="zero">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="memory"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Total amount of swap space displayed in KBs"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Swap Space Total"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="mem_free" VAL="4294852" TYPE="float" UNITS="KB" TN="25" TMAX="180" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="memory"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Amount of available memory"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Free Memory"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="load_fifteen" VAL="0.29" TYPE="float" UNITS=" " TN="38" TMAX="950" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="load"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Fifteen minute load average"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Fifteen Minute Load Average"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="os_name" VAL="Linux" TYPE="string" UNITS="" TN="1029" TMAX="1200" DMAX="0" SLOPE="zero">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="system"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Operating system name"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Operating System"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="boottime" VAL="1453529619" TYPE="uint32" UNITS="s" TN="1029" TMAX="1200" DMAX="0" SLOPE="zero">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="system"/>
<EXTRA_ELEMENT NAME="DESC" VAL="The last time that the system was started"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Last Boot Time"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="cpu_idle" VAL="97.0" TYPE="float" UNITS="%" TN="10" TMAX="90" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="cpu"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Percentage of time that the CPU or CPUs were idle and the system did not have an outstanding disk I/O request"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="CPU Idle"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="cpu_user" VAL="2.7" TYPE="float" UNITS="%" TN="10" TMAX="90" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="cpu"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Percentage of CPU utilization that occurred while executing at the user level"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="CPU User"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="cpu_nice" VAL="0.0" TYPE="float" UNITS="%" TN="10" TMAX="90" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="cpu"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Percentage of CPU utilization that occurred while executing at the user level with nice priority"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="CPU Nice"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="cpu_aidle" VAL="0.0" TYPE="float" UNITS="%" TN="10" TMAX="3800" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="cpu"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Percent of time since boot idle CPU"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="CPU aidle"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="mem_buffers" VAL="644188" TYPE="float" UNITS="KB" TN="25" TMAX="180" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="memory"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Amount of buffered memory"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Memory Buffers"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="cpu_system" VAL="0.3" TYPE="float" UNITS="%" TN="10" TMAX="90" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="cpu"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Percentage of CPU utilization that occurred while executing at the system level"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="CPU System"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="part_max_used" VAL="37.6" TYPE="float" UNITS="%" TN="128" TMAX="180" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="disk"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Maximum percent used for all partitions"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Maximum Disk Space Used"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="disk_total" VAL="486.329" TYPE="double" UNITS="GB" TN="2316" TMAX="1200" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="disk"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Total available disk space"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Total Disk Space"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="mem_shared" VAL="0" TYPE="float" UNITS="KB" TN="25" TMAX="180" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="memory"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Amount of shared memory"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Shared Memory"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="cpu_wio" VAL="0.0" TYPE="float" UNITS="%" TN="10" TMAX="90" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="cpu"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Percentage of time that the CPU or CPUs were idle during which the system had an outstanding disk I/O request"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="CPU wio"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="machine_type" VAL="x86_64" TYPE="string" UNITS="" TN="1029" TMAX="1200" DMAX="0" SLOPE="zero">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="system"/>
<EXTRA_ELEMENT NAME="DESC" VAL="System architecture"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Machine Type"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="proc_total" VAL="1269" TYPE="uint32" UNITS=" " TN="66" TMAX="950" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="process"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Total number of processes"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Total Processes"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="cpu_num" VAL="8" TYPE="uint16" UNITS="CPUs" TN="1029" TMAX="1200" DMAX="0" SLOPE="zero">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="cpu"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Total number of CPUs"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="CPU Count"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="cpu_speed" VAL="3300" TYPE="uint32" UNITS="MHz" TN="1029" TMAX="1200" DMAX="0" SLOPE="zero">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="cpu"/>
<EXTRA_ELEMENT NAME="DESC" VAL="CPU Speed in terms of MHz"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="CPU Speed"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="pkts_out" VAL="0.37" TYPE="float" UNITS="packets/sec" TN="129" TMAX="300" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="network"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Packets out per second"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Packets Sent"/>
</EXTRA_DATA>
</METRIC>
</HOST>
</CLUSTER>
</GANGLIA_XML>
`

var sample2 = `
<GANGLIA_XML VERSION="3.6.0" SOURCE="gmond">
<CLUSTER NAME="cluster1" LOCALTIME="1453984564" OWNER="unspecified" LATLONG="unspecified" URL="unspecified">
<HOST NAME="host1" IP="172.17.0.4" TAGS="" REPORTED="1453984563" TN="0" TMAX="20" DMAX="0" LOCATION="unspecified" GMOND_STARTED="1453885599">
<METRIC NAME="swap_free" VAL="16070712" TYPE="float" UNITS="KB" TN="25" TMAX="180" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="memory"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Amount of available swap memory"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Free Swap Space"/>
</EXTRA_DATA>
</METRIC>
</HOST>
<HOST NAME="host2" IP="172.17.0.4" TAGS="" REPORTED="1453984563" TN="0" TMAX="20" DMAX="0" LOCATION="unspecified" GMOND_STARTED="1453885599">
<METRIC NAME="gexec" VAL="OFF" TYPE="string" UNITS="" TN="129" TMAX="300" DMAX="0" SLOPE="zero">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="core"/>
<EXTRA_ELEMENT NAME="DESC" VAL="gexec available"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Gexec Status"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="disk_free" VAL="303.665" TYPE="double" UNITS="GB" TN="128" TMAX="180" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="disk"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Total free disk space"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Disk Space Available"/>
</EXTRA_DATA>
</METRIC>
</HOST>
<HOST NAME="host3" IP="172.17.0.4" TAGS="" REPORTED="1453984563" TN="0" TMAX="20" DMAX="0" LOCATION="unspecified" GMOND_STARTED="1453885599">
<METRIC NAME="mem_cached" VAL="3554892" TYPE="float" UNITS="KB" TN="25" TMAX="180" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="memory"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Amount of cached memory"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Cached Memory"/>
</EXTRA_DATA>
</METRIC>
<METRIC NAME="bytes_in" VAL="0.00" TYPE="float" UNITS="bytes/sec" TN="129" TMAX="300" DMAX="0" SLOPE="both">
<EXTRA_DATA>
<EXTRA_ELEMENT NAME="GROUP" VAL="network"/>
<EXTRA_ELEMENT NAME="DESC" VAL="Number of bytes in per second"/>
<EXTRA_ELEMENT NAME="TITLE" VAL="Bytes Received"/>
</EXTRA_DATA>
</METRIC>
</HOST>
</CLUSTER>
</GANGLIA_XML>
`
