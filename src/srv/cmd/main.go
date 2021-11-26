package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/src/srv/models"
	"github.com/src/srv/webserver"
)

const drawDataURL = "https://www.lotterywest.wa.gov.au/results/frequency-charts/?soccerpoolsx"
const randomAPIURL = "https://api.random.org/json-rpc/2/invoke"

//Parse command line switches
func getArgs() models.CmdArgs {

	args := models.CmdArgs{}

	flag.StringVar(&args.Game, "game", "oz-lotto", "Game to generate numbers for.")
	flag.StringVar(&args.GameType, "type", "standard", "Standard, System7, System8,...")
	flag.IntVar(&args.DrawOffset, "offset", 1, "Numbers of draws to go back to from highest draw.")
	flag.IntVar(&args.Draws, "draws", 7, "Draws to evaluate from starting draw.")
	flag.BoolVar(&args.UseWebUI, "webui", false, "Load the WebUI using the default browser. Also starts the web server.")
	flag.BoolVar(&args.UseWebserver, "web", false, "Run a webserver to serve the webui.")
	flag.IntVar(&args.Port, "port", 1337, "Specify Port to be used by web server.")
	flag.StringVar(&args.RapiKey, "apikey", "", "Api key for random number web service")
	flag.StringVar(&args.DataURL, "dataurl", drawDataURL, "URL for historic draw info")

	flag.Parse()

	return args

}

func main() {

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

}

func run() error {

	//Handle ctrl-c
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Printf("%s received. Exiting...", sig)
		done <- true
		os.Exit(0)
	}()

	args := getArgs()
	args.RapiURL = randomAPIURL

	cwd, _ := os.Getwd()
	rootDir := filepath.Dir(filepath.ToSlash(cwd))
	fmt.Printf("Current Directory => %s\nRoot Directory => %s\n", cwd, rootDir)

	_, err := os.Stat(fmt.Sprintf("%s/data", rootDir))

	if err != nil {

		//Go up one level and try again
		_, err = os.Stat(fmt.Sprintf("%s/data", filepath.Dir(cwd)))
		if err != nil { //Give up, use current directory
			fmt.Printf("/Data directory not found. Using current directory %s\n", cwd)
		} else {
			cwd = filepath.Dir(cwd)
		}
	}

	//Read and process most current data from website
	err = models.InitGames(fmt.Sprintf("%s/data", rootDir), drawDataURL, args)
	if err != nil {
		return err //log.Fatal(err)
	}

	ctx := webserver.ServerContext{
		Router:             mux.NewRouter(),
		RandomNumberAPIURL: randomAPIURL,
		RandomNumberAPIKEY: args.RapiKey,
		WebServer:          args.UseWebserver,
		WebUI:              args.UseWebUI,
		Port:               args.Port,
		Cwd:                cwd,
		RootDir:            rootDir,
	}

	//Initialise web/server
	_ = ctx.InitWebserver()

	return nil

}
