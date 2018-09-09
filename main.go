package main

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net"
	"os"
)

const (
	opendnsMyIP         = "myip.opendns.com"
	opendnsResolver     = "resolver1.opendns.com"
	vpnServerName       = "vpn"
	doRegionFrankfurt   = "fra1"
	doSize1GB           = "s-1vcpu-1gb"
	doImageUbuntuBionic = "ubuntu-18-04-x64"
	userDataFile        = "user-data.sh"
)

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}

	return token, nil
}

func getExternalIP() string {
	resolver := &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", net.JoinHostPort(opendnsResolver, "53"))
		},
	}

	ip, err := resolver.LookupIPAddr(context.Background(), opendnsMyIP)
	if err != nil {
		fmt.Print(err)
		return ""
	}

	return fmt.Sprintf("%s", ip[0].IP)
}

func start(client *godo.Client) {
	fmt.Printf("Starting server '%s' ...\n", vpnServerName)

	// extIP := getExternalIP()
	userData, err := ioutil.ReadFile(userDataFile)
	if err != nil {
		fmt.Print(err)
	}

	crtReq := &godo.DropletCreateRequest{
		Name:              vpnServerName,
		Region:            doRegionFrankfurt,
		Size:              doSize1GB,
		Image:             godo.DropletCreateImage{Slug: doImageUbuntuBionic},
		SSHKeys:           []godo.DropletCreateSSHKey{{Fingerprint: os.Getenv("DO_SSHKEY")}},
		Backups:           false,
		IPv6:              false,
		PrivateNetworking: false,
		Monitoring:        false,
		UserData:          string(userData),
		Tags:              []string{vpnServerName},
	}

	client.Droplets.Create(context.Background(), crtReq)
}

func stop(client *godo.Client) {
	fmt.Printf("Stopping server '%s' ...\n", vpnServerName)
	client.Droplets.DeleteByTag(context.Background(), vpnServerName)
}

func main() {
	errArg := fmt.Errorf("fatal: invalid argument.\n\nUsage:\n  vpn-server {start|stop}")

	if len(os.Args) == 1 {
		fmt.Println(errArg)
		return
	}

	ts := &tokenSource{AccessToken: os.Getenv("DO_PAT")}
	oauthClient := oauth2.NewClient(context.Background(), ts)
	client := godo.NewClient(oauthClient)

	if os.Args[1] == "start" {
		start(client)
	} else if os.Args[1] == "stop" {
		stop(client)
	} else {
		fmt.Println(errArg)
	}
}
