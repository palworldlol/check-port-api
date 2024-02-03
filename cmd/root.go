package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var checkHost string
var checkPort string
var listenHost string
var listenPort string
var timeout time.Duration

var targetAddress string

const UP = "UP"
const DOWN = "DOWN"

var rootCmd = &cobra.Command{
	Use:   "check-port-api",
	Short: "A really simple API that responds UP or DOWN if a certain address+port combination is listening",
	Long: `Host a really simple API that responds UP or DOWN if it finds a specific address+port is listenin.
This is useful to act as a proxy to check if a certain port is open on a system but its not actually exposed.`,
	PreRun: func(cmd *cobra.Command, args []string) {
		targetAddress = net.JoinHostPort(checkHost, checkPort)
	},
	Run: func(cmd *cobra.Command, args []string) {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// Check if target address+port is listening
			err := ping(targetAddress, timeout)
			if err != nil {
				fmt.Fprint(w, DOWN)
			} else {
				fmt.Fprint(w, UP)
			}
		})

		log.Fatal(http.ListenAndServe(net.JoinHostPort(listenHost, listenPort), nil))
	},
}

// https://github.com/janosgyerik/portping/blob/master/portping.go
// Ping connects to the address on the named network,
// using net.DialTimeout, and immediately closes it.
// It returns the connection error. A nil value means success.
// For examples of valid values of network and address,
// see the documentation of net.Dial
func ping(address string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", address, timeout)
	if conn != nil {
		defer conn.Close()
	}
	return err
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&checkHost, "check-host", "", "The hostname to check")
	rootCmd.Flags().StringVar(&checkPort, "check-port", "", "The port to check")
	rootCmd.Flags().StringVar(&listenHost, "listen-host", "0.0.0.0", "The hostname to listen on")
	rootCmd.Flags().StringVar(&listenPort, "listen-port", "8181", "The port to listen on")
	rootCmd.Flags().DurationVar(&timeout, "timeout", 5*time.Second, "Timeout to check target")

	rootCmd.MarkFlagRequired("check-host")
	rootCmd.MarkFlagRequired("check-port")
}
