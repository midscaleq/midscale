/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"midscale/midscale/app/controller/http"
	"midscale/midscale/app/data/model/define"

	"github.com/spf13/cobra"
)

var endpointHost string
var serverHost string
var persistentKeepalive int

var listenAddrForRemote string
var listenAddrForLocal string
var listenPort int
var localIP string
var oneMaskLength int
var dNS string

var autoStartTunnel bool

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start a server",
	Long:  `Start a local server and a remote server`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(endpointHost) > 0 {
			define.EndpointHost = endpointHost
		}
		if len(serverHost) < 1 {
			serverHost = "http://" + endpointHost
		}
		if len(serverHost) > 0 {
			define.ServerHost = serverHost
		}
		if persistentKeepalive > 0 {
			define.PersistentKeepalive = persistentKeepalive
		}

		if len(listenAddrForRemote) > 0 {
			define.DefaultListenAddrForRemote = listenAddrForRemote
		}
		if len(listenAddrForLocal) > 0 {
			define.DefaultListenAddrForLocal = listenAddrForLocal
		}
		if listenPort > 0 {
			define.DefaultListenPort = listenPort
		}
		if len(localIP) > 0 {
			define.DefaultLocalIP = localIP
		}
		if oneMaskLength > 0 {
			define.DefaultOneMaskLength = oneMaskLength
		}
		if len(dNS) > 0 {
			define.DefaultDNS = dNS
		}
		if !autoStartTunnel {
			define.AutoStartTunnel = autoStartTunnel
		}

		http.NewHttpServer().Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	serverCmd.Flags().StringVarP(&endpointHost, "endpointHost", "e", "", "setting of wireguard peer's endpoint (required)")
	serverCmd.MarkFlagRequired("endpointHost")
	serverCmd.Flags().StringVar(&serverHost, "serverHost", "", "client sent http request to this host (optional)")
	serverCmd.Flags().IntVar(&persistentKeepalive, "persistentKeepalive", define.PersistentKeepalive, "setting of wireguard peer's persistentKeepalive (optional)")

	serverCmd.Flags().StringVar(&listenAddrForRemote, "listenAddrForRemote", define.DefaultListenAddrForRemote, "listening address and port for client connecting (optional)")
	serverCmd.Flags().StringVar(&listenAddrForLocal, "listenAddrForLocal", define.DefaultListenAddrForLocal, "listening address and port for local connecting (optional)")
	serverCmd.Flags().IntVar(&listenPort, "listenPort", define.DefaultListenPort, "setting of wireguard Interface listen port (optional)")
	serverCmd.Flags().StringVar(&localIP, "localIP", define.DefaultLocalIP, "setting of wireguard Interface local ip (optional)")
	serverCmd.Flags().IntVar(&oneMaskLength, "oneMaskLength", define.DefaultOneMaskLength, "setting of wireguard Interface local ip one mask length (optional)")
	serverCmd.Flags().StringVar(&dNS, "dNS", define.DefaultDNS, "setting of wireguard Interface dns (optional)")

	serverCmd.Flags().BoolVar(&autoStartTunnel, "autoStartTunnel", true, "auto install and start wireguard tunnel (optional)")
}
