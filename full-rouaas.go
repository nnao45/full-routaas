package main

import (
	"context"
	"fmt"
	api "github.com/osrg/gobgp/api"
	"github.com/osrg/gobgp/config"
	"github.com/osrg/gobgp/packet/bgp"
	"github.com/osrg/gobgp/packet/mrt"
	gobgp "github.com/osrg/gobgp/server"
	"github.com/osrg/gobgp/table"
	log "github.com/sirupsen/logrus"
	grpc "google.golang.org/grpc"
	"io/ioutil"
	"os"
	"strings"
	"io"
	"net"
	"time"
)

type mrtOpts struct {
	OutputDir   string
	FileFormat  string
	Filename    string
	RecordCount int64
	RecordSkip  int64
	QueueSize   int
	Best        bool
	SkipV4      bool
	SkipV6      bool
	NextHop     net.IP
}

func newmrtOpts() mrtOpts {
	return mrtOpts{
		OutputDir:  "./",
		FileFormat: "",
		Filename:   "",
		Best:       false,
		QueueSize:  100000,
		SkipV4:     false,
		SkipV6:     false,
		NextHop:    nil,
	}
}

func main() {

	s := gobgp.NewBgpServer()
	go s.Serve()

	// start grpc api server. this is not mandatory
	// but you will be able to use `gobgp` cmd with this.
	g := api.NewGrpcServer(s, ":50051")
	go g.Serve()

	global := &config.Global{
		Config: config.GlobalConfig{
			As:       65555,
			RouterId: "192.168.0.1",
		},
	}

	if err := s.Start(global); err != nil {
		log.Fatal(err)
	}

	neighbor := &config.Neighbor{
		Config: config.NeighborConfig{
			PeerAs:          65000,
			NeighborAddress: "192.168.0.2",
		},
		EbgpMultihop: config.EbgpMultihop{
			Config: config.EbgpMultihopConfig{
				Enabled:     true,
				MultihopTtl: 255,
			},
		},
		AfiSafis: []config.AfiSafi{
			config.AfiSafi{
				Config: config.AfiSafiConfig{
					AfiSafiName: "ipv4-unicast",
				},
			},
		},
	}

	if err := s.AddNeighbor(neighbor); err != nil {
		log.Error(err)
	}

	timeout := grpc.WithTimeout(time.Second)

	conn, rpcErr := grpc.Dial("localhost:50051", timeout, grpc.WithBlock(), grpc.WithInsecure())
	if rpcErr != nil {
		log.Fatal("GoBGP is probably not running on the local server ... Please start gobgpd process !\n")
		log.Fatal(rpcErr)
		return
	}

	bgpclient := api.NewGobgpApiClient(conn)

	m := newmrtOpts()
	var mErr error
	m.Filename, mErr = findMrt()
	if mErr != nil {
		log.Fatal(mErr)
	}
	log.Info("MRT injection file is ", m.Filename)

	go func() {
		err := injectMrt(bgpclient, m)
		if err != nil {
			fmt.Errorf("failed to add path: %s", err)
			return
		}

		defer log.Info("MRT injection complete!!")

	}()

	select {}

}

func findMrt() (mrtFile string, err error) {
	files, e := ioutil.ReadDir("./")
	if e != nil {
		err = e
	}
	for _, file := range files {
		if strings.Contains(file.Name(), "rib") {
			mrtFile = "./" + file.Name()
			return
		}
	}
	fmt.Errorf("failed to read mib file.", err)
	return
}

func injectMrt(bgpclient api.GobgpApiClient, m mrtOpts) error {

	file, err := os.Open(m.Filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %s", err)
	}

	if m.NextHop != nil && !m.SkipV4 && !m.SkipV6 {
		fmt.Println("You should probably specify either --no-ipv4 or --no-ipv6 when overwriting nexthop, unless your dump contains only one type of routes")
	}

	var idx int64
	if m.QueueSize < 1 {
		return fmt.Errorf("Specified queue size is smaller than 1, refusing to run with unbounded memory usage")
	}

	ch := make(chan []*table.Path, m.QueueSize)
	go func() {

		var peers []*mrt.Peer
		for {
			buf := make([]byte, mrt.MRT_COMMON_HEADER_LEN)
			_, err := file.Read(buf)
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(fmt.Errorf("failed to read: %s", err))
			}

			h := &mrt.MRTHeader{}
			err = h.DecodeFromBytes(buf)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to parse"))
			}

			buf = make([]byte, h.Len)
			_, err = file.Read(buf)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to read"))
			}

			msg, err := mrt.ParseMRTBody(h, buf)
			if err != nil {
				log.Fatal(fmt.Errorf("failed to parse: %s", err))
				continue
			}

			//fmt.Println(msg)
			if msg.Header.Type == mrt.TABLE_DUMPv2 {
				subType := mrt.MRTSubTypeTableDumpv2(msg.Header.SubType)
				switch subType {
				case mrt.PEER_INDEX_TABLE:
					peers = msg.Body.(*mrt.PeerIndexTable).Peers
					continue
				case mrt.RIB_IPV4_UNICAST, mrt.RIB_IPV4_UNICAST_ADDPATH:
					if m.SkipV4 {
						continue
					}
				case mrt.RIB_IPV6_UNICAST, mrt.RIB_IPV6_UNICAST_ADDPATH:
					if m.SkipV6 {
						continue
					}
				case mrt.GEO_PEER_TABLE:
					fmt.Printf("WARNING: Skipping GEO_PEER_TABLE: %s", msg.Body.(*mrt.GeoPeerTable))
				default:
					log.Fatal(fmt.Errorf("unsupported subType: %v", subType))
				}

				if peers == nil {
					log.Fatal(fmt.Errorf("not found PEER_INDEX_TABLE"))
				}

				rib := msg.Body.(*mrt.Rib)
				nlri := rib.Prefix

				paths := make([]*table.Path, 0, len(rib.Entries))

				for _, e := range rib.Entries {
					if len(peers) < int(e.PeerIndex) {
						log.Fatal(fmt.Errorf("invalid peer index: %d (PEER_INDEX_TABLE has only %d peers)\n", e.PeerIndex, len(peers)))
					}
					source := &table.PeerInfo{
						AS: peers[e.PeerIndex].AS,
						ID: peers[e.PeerIndex].BgpId,
					}
					t := time.Unix(int64(e.OriginatedTime), 0)

					switch subType {
					case mrt.RIB_IPV4_UNICAST, mrt.RIB_IPV4_UNICAST_ADDPATH:
						paths = append(paths, table.NewPath(source, nlri, false, e.PathAttributes, t, false))
					default:
						attrs := make([]bgp.PathAttributeInterface, 0, len(e.PathAttributes))
						for _, attr := range e.PathAttributes {
							if attr.GetType() != bgp.BGP_ATTR_TYPE_MP_REACH_NLRI {
								attrs = append(attrs, attr)
							} else {
								a := attr.(*bgp.PathAttributeMpReachNLRI)
								attrs = append(attrs, bgp.NewPathAttributeMpReachNLRI(a.Nexthop.String(), []bgp.AddrPrefixInterface{nlri}))
							}
						}
						paths = append(paths, table.NewPath(source, nlri, false, attrs, t, false))
					}
				}
				if m.NextHop != nil {
					for _, p := range paths {
						p.SetNexthop(m.NextHop)
					}
				}

				if m.Best {
					dst := table.NewDestination(nlri, 0)
					for _, p := range paths {
						dst.AddNewPath(p)
					}
					best, _, _ := dst.Calculate().GetChanges(table.GLOBAL_RIB_NAME, false)
					if best == nil {
						log.Fatal(fmt.Errorf("Can't find the best %v", nlri))
					}
					paths = []*table.Path{best}
				}

				if idx >= m.RecordSkip {
					ch <- paths
				}

				idx += 1
				if idx == m.RecordCount+m.RecordSkip {
					break
				}
			}
		}

		close(ch)
	}()

	bgpmrtclient, err := bgpclient.InjectMrt(context.Background())
	if err != nil {
		return fmt.Errorf("failed to add path: %s", err)
	}

	for paths := range ch {
		var tables []*api.Path
		for _, p := range paths {
			tables = append(tables, api.ToPathApi(p))
		}
		req := &api.InjectMrtRequest{
			Resource: api.Resource_GLOBAL,
			VrfId:    "",
			Paths:    tables,
		}
		err = bgpmrtclient.Send(req)
		if err != nil {
			return fmt.Errorf("failed to send: %s", err)
		}
	}

	if _, err := bgpmrtclient.CloseAndRecv(); err != nil {
		return fmt.Errorf("failed to send: %s", err)
	}
	return nil
}
