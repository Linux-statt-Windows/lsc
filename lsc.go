package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/BurntSushi/toml"
	"github.com/digitalocean/godo"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

type TokenSource struct {
	AccessToken string
}

type tomlConfig struct {
	Token     string
	Log       bool
	DropletID int
}

func main() {
	configPath := getHomeDir() + "/.lsccfg"

	var config tomlConfig
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		fmt.Println(err)
		return
	}

	var cmdOn = &cobra.Command{
		Use:   "on [Droplet ID]",
		Short: "Power On given Droplet.",
		Long:  `Power On given Droplet.`,
		Run: func(cmd *cobra.Command, args []string) {
			PowerOn(strings.Join(args, ""), config)
		},
	}

	var cmdShutdown = &cobra.Command{
		Use:   "shutdown [Droplet ID]",
		Short: "Shutdown given Droplet. (soft)",
		Long:  `Shutdown given Droplet. This is a softer version of cycle.`,
		Run: func(cmd *cobra.Command, args []string) {
			Shutdown(strings.Join(args, ""), config)
		},
	}

	var cmdCycle = &cobra.Command{
		Use:   "cycle [Droplet ID]",
		Short: "PowerCycle given Droplet. Hard version! (turn off) ",
		Long:  `PowerCycle given Droplet. (turn off) This is the hard version. Only use if needed!`,
		Run: func(cmd *cobra.Command, args []string) {
			Cycle(strings.Join(args, ""), config)
		},
	}

	var cmdReboot = &cobra.Command{
		Use:   "reboot [Droplet ID]",
		Short: "Reboot given Droplet. (soft)",
		Long:  `Reboot given Droplet. This is a softer version of cycle.`,
		Run: func(cmd *cobra.Command, args []string) {
			Reboot(strings.Join(args, ""), config)
		},
	}

	var cmdStatus = &cobra.Command{
		Use:   "status [Droplet ID]",
		Short: "Status of given Droplet.",
		Long:  `Status of given Droplet.`,
		Run: func(cmd *cobra.Command, args []string) {
			Status(strings.Join(args, ""), config)
		},
	}

	var rootCmd = &cobra.Command{Use: "lsc"}
	rootCmd.AddCommand(cmdOn, cmdReboot, cmdShutdown, cmdCycle, cmdStatus)
	rootCmd.Execute()
}

/**
 * Gets the full path to the home directory of the current user.
 *
 * @returns the full path to the home directory of the current user as string.
 */
func getHomeDir() string {
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir
}

/**
 * Gets the Token used for authentication
 *
 * @returns *oauth2.Token, error
 */
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

/**
 * Gets the Droplet by a given ID
 *
 * @returns *godo.Droplet, error
 */
func getDropletByID(client *godo.Client, id int) (*godo.Droplet, error) {
	if id < 1 {
		log.Fatal("missing droplet id")
	}

	droplet, _, err := client.Droplets.Get(id)
	return droplet, err
}

/**
 * Writes Output as JSON used for debugging.
 */
func WriteOutput(data interface{}) {
	var output []byte
	var err error

	output, err = json.Marshal(data)
	if err != nil {
		log.Fatalf("JSON Encoding Error: %s", err)
	}

	fmt.Println(string(output))
}

/**
 * Function to PowerOn the Droplet
 */
func PowerOn(arg string, config tomlConfig) {
	tokenSource := &TokenSource{
		AccessToken: config.Token,
	}

	var droplet_ID int
	droplet_ID, _ = strconv.Atoi(arg)

	if droplet_ID == 0 {
		if config.DropletID == 0 {
			log.Fatal("No Droplet ID provided in Config! Use Parameters.")
		} else {
			droplet_ID = config.DropletID
		}
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	action, _, err := client.DropletActions.PowerOn(droplet_ID)
	if err != nil {
		log.Fatal(err)
	}

	if config.Log {
		WriteOutput(action)
	} else {
		fmt.Println("PowerOn the system. The action could take a few moments.")
	}
}

/**
 * Function to Reboot the Droplet
 */
func Reboot(arg string, config tomlConfig) {
	tokenSource := &TokenSource{
		AccessToken: config.Token,
	}

	var droplet_ID int
	droplet_ID, _ = strconv.Atoi(arg)

	if droplet_ID == 0 {
		if config.DropletID == 0 {
			log.Fatal("No Droplet ID provided in Config! Use Parameters.")
		} else {
			droplet_ID = config.DropletID
		}
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	action, _, err := client.DropletActions.Reboot(droplet_ID)
	if err != nil {
		log.Fatal(err)
	}

	if config.Log {
		WriteOutput(action)
	} else {
		fmt.Println("Reboot the system. The action could take a few moments.")
	}
}

/**
 * Function to Shutdown the Droplet
 */
func Shutdown(arg string, config tomlConfig) {
	tokenSource := &TokenSource{
		AccessToken: config.Token,
	}

	var droplet_ID int
	droplet_ID, _ = strconv.Atoi(arg)

	if droplet_ID == 0 {
		if config.DropletID == 0 {
			log.Fatal("No Droplet ID provided in Config! Use Parameters.")
		} else {
			droplet_ID = config.DropletID
		}
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	action, _, err := client.DropletActions.Shutdown(droplet_ID)
	if err != nil {
		log.Fatal(err)
	}

	if config.Log {
		WriteOutput(action)
	} else {
		fmt.Println("Shutdown the system. The action could take a few moments.")
	}
}

/**
 * Function to PowerCycle the Droplet. (hard shutdown)
 */
func Cycle(arg string, config tomlConfig) {
	tokenSource := &TokenSource{
		AccessToken: config.Token,
	}

	var droplet_ID int
	droplet_ID, _ = strconv.Atoi(arg)

	if droplet_ID == 0 {
		if config.DropletID == 0 {
			log.Fatal("No Droplet ID provided in Config! Use Parameters.")
		} else {
			droplet_ID = config.DropletID
		}
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	action, _, err := client.DropletActions.PowerCycle(droplet_ID)
	if err != nil {
		log.Fatal(err)
	}

	if config.Log {
		WriteOutput(action)
	} else {
		fmt.Println("Shutdown the system forcefully. The action could take a few moments.")
	}
}

/**
 * Function show the status of a Droplet
 */
func Status(arg string, config tomlConfig) {
	tokenSource := &TokenSource{
		AccessToken: config.Token,
	}

	var droplet_ID int
	droplet_ID, _ = strconv.Atoi(arg)

	if droplet_ID == 0 {
		if config.DropletID == 0 {
			log.Fatal("No Droplet ID provided in Config! Use Parameters.")
		} else {
			droplet_ID = config.DropletID
		}
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	droplet, err := getDropletByID(client, droplet_ID)
	if err != nil {
		log.Fatal(err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 3, '\t', 0)
	fmt.Fprintln(w, "ID\tName\tStatus\t")
	fmt.Fprintf(w, "%d\t%s\t%s\n", droplet.ID, droplet.Name, droplet.Status)

	w.Flush()
}
