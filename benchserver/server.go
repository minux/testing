package benchserver

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
	"testing"
)

// Overall plan of attack:
// We use an init function to initialize the RPC server if -test.benchserver is set.
// If it's set, we take over the execution and serve the RPC server.
// Pretty trivial, heh? The problem is that demons usually lie in the details:
// 0. we can't use the flag package to parse test.benchserver flag because the user
// might add more options in its init or TestMain. But we do need to register the
// flag so that -help works as intended.
// 1. we can't take over the execution of the process in init, because the package
// under test have not initialized yet; moreover, the package might have a TestMain
// with further initialization code. This is tricky to solve, we solve it by injecting
// a test into main.tests and add a -test.run option.

const (
	benchServerFlag = "test.benchserver"
)

var (
	benchServerAddr = flag.String("test.bench.addr", "localhost:54321", "`host:port` for the JSON-RPC bench server")
	benchServer     = flag.Bool(benchServerFlag, false, "enable the JSON-RPC bench server")
)

func init() {
	_ = inited                      // make sure the benchmarks variable have been initialized
	for i, v := range *benchmarks { // only for debug
		fmt.Printf("%d: %#v\n", i, v)
	}

	// see if server is enabled by querying os.Args
	flagSet := flag.NewFlagSet(os.Args[0], 0)
	flagSet.Usage = func() {}         // suppress usage message
	flagSet.SetOutput(ioutil.Discard) // suppress error messages
	enabled := flagSet.Bool(benchServerFlag, false, "")
	flagSet.Parse(os.Args[1:]) // ignore error, as errors should be reported by the main testing package
	if !*enabled {
		return
	}

	// inject our test
	*tests = append(*tests, testing.InternalTest{Name: "BenchServer", F: BenchServer})

	// add -test.run flag
	os.Args = append(os.Args, "-test.run=^BenchServer$")
}

// Note: can't use t.Logf because the output is buffered by testing package
func logf(format string, args ...interface{}) {
	if testing.Verbose() {
		// can't use the log package, as the user test might be using it
		fmt.Fprintf(os.Stderr, format, args...)
	}
}

func BenchServer(t *testing.T) {
	server := Server{m: make(map[string]InternalBenchmark)}
	for _, b := range *benchmarks {
		server.m[b.Name] = b
	}
	rpc.Register(server)

	l, err := net.Listen("tcp", *benchServerAddr)
	if err != nil {
		t.Fatalf("listening on %v failed: %v", *benchServerAddr, err)
	}
	defer l.Close()

	logf("Listening on %v\n", l.Addr())

	for {
		conn, err := l.Accept()
		if err != nil {
			t.Fatalf("accept failed: %v", err)
		}
		go func() {
			jsonrpc.ServeConn(conn)
			conn.Close()
		}()
	}
}

// Server is the benchmark server state.
type Server struct {
	m map[string]InternalBenchmark
}

type Arg struct {
	Name string `json:"name"`
	N    int    `json:"n"`
}

type Reply struct {
	Result *testing.BenchmarkResult `json:"result,omitempty"`
	Names  []string                 `json:"names,omitempty"`
}

func (s Server) List(args *Arg, reply *Reply) error {
	for _, b := range *benchmarks {
		reply.Names = append(reply.Names, b.Name)
	}
	return nil
}

func (s Server) Run(args *Arg, reply *Reply) error {
	if _, ok := s.m[args.Name]; !ok {
		return fmt.Errorf("no such benchmark: %s", args.Name)
	}
	b := &B{benchmark: s.m[args.Name]}
	logf("Running %s with N=%d\n", args.Name, args.N)
	reply.Result = runN(b, args.N)
	return nil
}

func (s Server) Quit(args *Arg, reply *Reply) error {
	os.Exit(0)
	return nil
}
