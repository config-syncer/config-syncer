package v1alpha1

import (
	"fmt"
)

const (
	// GarbdListenPort is the port at which Galera Arbitrator Daemon (garbd) listen
	GarbdListenPort = 4444

	// GarbdXtrabackupSSTMethod is the name of the method or script that is
	// used during a State Snapshot Transfer to Galera Arbitrator Daemon (garbd).
	GarbdXtrabackupSSTMethod = "xtrabackup-v2"

	// GarbdXtrabackupSSTRequestSuffix denotes the suffix of sst request string for xtrabackup
	// used by Galera Arbitrator Daemon (garbd)
	GarbdXtrabackupSSTRequestSuffix = "/xtrabackup_sst//1"
	// GarbdLogFile is the name log file at which Galera Arbitrator Daemon (garbd) puts logs
	GarbdLogFile = "/tmp/garb.log"

	// SOCAT is needed after completing sst by Galera Arbitrator Daemon (garbd)
	// SOCATOptionReUseAddr is the SOCAT reuseaddr option
	SOCATOptionReUseAddr = "reuseaddr"
	// SOCATOptionRetry is the default retry value for `socat` binary
	SOCATOptionRetry = 30
)

// ClusterAddressWithListenOption method returns the galera cluster address with
// the listening option (address at which Galera Cluster listens to connections from
// other nodes) for `--address` option in `garbd`
// Here, ‘?gmcast.listen_addr=tcp://0.0.0.0:4444‘ is an arbitrary listen socket address
// that Galera Arbitrator opens to communicate with the cluster.
// https://galeracluster.com/library/documentation/backup-cluster.html
func (g *GaleraArbitratorConfiguration) ClusterAddressWithListenOption() string {
	if g == nil {
		return ""
	}

	return fmt.Sprintf("%s?gmcast.listen_addr=tcp://0.0.0.0:%d", g.Address, GarbdListenPort)
}

// SSTRequestString method form the sst request string
// for `--sst` option in `garbd`
func (g *GaleraArbitratorConfiguration) SSTRequestString(host string) string {
	if g == nil {
		return ""
	}

	return fmt.Sprintf("%s:%s:%d%s", g.SSTMethod, host, GarbdListenPort, GarbdXtrabackupSSTRequestSuffix)
}

// SOCATOption returns the option string used for `SOCAT` in the
// percona xtradb backup process
func SOCATOption(retry int32) string {
	return fmt.Sprintf("TCP-LISTEN:%d,%s,retry=%d", GarbdListenPort, SOCATOptionReUseAddr, retry)
}
