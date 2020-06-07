package cli

import (
	"fmt"
	"io"
	"net"

	"github.com/spf13/cobra"
	"golang.org/x/net/websocket"
)

func newServerRequestConsoleCommand(cli *CLI) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "request-console [FLAGS] SERVER",
		Short:                 "Request a WebSocket VNC console for a server",
		Args:                  cobra.ExactArgs(1),
		TraverseChildren:      true,
		DisableFlagsInUseLine: true,
		PreRunE:               cli.ensureToken,
		RunE:                  cli.wrap(runServerRequestConsole),
	}
	addOutputFlag(cmd, outputOptionJSON())
	cmd.Flags().BoolP("listen", "l", false, "Listen for a VNC client")
	cmd.Flags().StringP("bind", "b", "localhost:5900", "Bind host/port for VNC port")
	return cmd
}

func consoleProxy(bind, url string) error {
	s, err := net.Listen("tcp", bind)
	if err != nil {
		return err
	}
	defer s.Close()

	fmt.Printf("\nListening for VNC client on %v...\n", bind)

	conn, err := s.Accept()
	if err != nil {
		return err
	}
	fmt.Printf("VNC Client connect, establishing WebSocket connection...\n")
	defer conn.Close()

	ws, err := websocket.Dial(url, "binary", "https://console.hetzner.cloud/")
	if err != nil {
		return err
	}
	ws.PayloadType = websocket.BinaryFrame
	fmt.Printf("Connected!\n")
	defer ws.Close()

	in := make(chan error, 1)
	go func() {
		_, err := io.Copy(conn, ws)
		in <- err
	}()

	out := make(chan error, 1)
	go func() {
		_, err := io.Copy(ws, conn)
		out <- err
	}()

	select {
	case err = <-in:
		return err
	case err = <-out:
		return err
	}
}

func runServerRequestConsole(cli *CLI, cmd *cobra.Command, args []string) error {
	outOpts := outputFlagsForCommand(cmd)
	idOrName := args[0]
	server, _, err := cli.Client().Server.Get(cli.Context, idOrName)
	if err != nil {
		return err
	}
	if server == nil {
		return fmt.Errorf("server not found: %s", idOrName)
	}

	result, _, err := cli.Client().Server.RequestConsole(cli.Context, server)
	if err != nil {
		return err
	}

	if err := cli.ActionProgress(cli.Context, result.Action); err != nil {
		return err
	}

	if outOpts.IsSet("json") {
		return describeJSON(struct {
			WSSURL   string
			Password string
		}{
			WSSURL:   result.WSSURL,
			Password: result.Password,
		})
	}

	fmt.Printf("Console for server %d:\n", server.ID)
	fmt.Printf("WebSocket URL: %s\n", result.WSSURL)
	fmt.Printf("VNC Password: %s\n", result.Password)

	listen, err := cmd.Flags().GetBool("listen")
	if err != nil {
		return err
	}

	if listen {
		bind, err := cmd.Flags().GetString("bind")
		if err != nil {
			return err
		}

		return consoleProxy(bind, result.WSSURL)
	}

	return nil
}
