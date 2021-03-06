package util

import (
	"bytes"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	log "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/fsouza/go-dockerclient"

	. "github.com/eris-ltd/eris-cli/Godeps/_workspace/src/github.com/eris-ltd/common/go/common"
)

// Docker Client initialization
var DockerClient *docker.Client

func DockerConnect(verbose bool, machName string) { // TODO: return an error...?
	var err error
	var dockerHost string
	var dockerCertPath string
	if runtime.GOOS == "linux" {
		// this means we aren't gonna use docker-machine (kind of)
		if (machName == "eris" || machName == "default") && (os.Getenv("DOCKER_HOST") == "" && os.Getenv("DOCKER_CERT_PATH") == "") {
			//if os.Getenv("DOCKER_HOST") == "" && os.Getenv("DOCKER_CERT_PATH") == "" {
			endpoint := "unix:///var/run/docker.sock"

			log.WithField("=>", endpoint).Debug("Checking Linux Docker socket")
			u, _ := url.Parse(endpoint)
			_, err := net.Dial(u.Scheme, u.Path)
			if err != nil {
				IfExit(fmt.Errorf("%v\n", mustInstallError()))
			}
			log.WithField("=>", endpoint).Debug("Connecting to Docker")
			DockerClient, err = docker.NewClient(endpoint)
			if err != nil {
				IfExit(fmt.Errorf("%v\n", mustInstallError()))
			}
		} else {
			log.WithFields(log.Fields{
				"host":      os.Getenv("DOCKER_HOST"),
				"cert path": os.Getenv("DOCKER_CERT_PATH"),
			}).Debug("Getting connection details from environment")
			log.WithField("machine", machName).Debug("Getting connection details from Docker Machine")
			dockerHost, dockerCertPath, err = getMachineDeets(machName)
			if err != nil {
				IfExit(fmt.Errorf("Error getting Docker Machine details for connection via TLS.\nERROR =>\t\t\t%v\n\nEither re-run the command without a machine or correct your machine name.\n", err))
			}

			log.WithFields(log.Fields{
				"host":      dockerHost,
				"cert path": dockerCertPath,
			}).Debug()

			if err := connectDockerTLS(dockerHost, dockerCertPath); err != nil {
				IfExit(fmt.Errorf("Error connecting to Docker via TLS.\nERROR =>\t\t\t%v\n", err))
			}

			setIPFSHostViaDockerHost(dockerHost)
		}

		log.Debug("Successfully connected to Docker daemon")

	} else {
		log.WithFields(log.Fields{
			"host":      os.Getenv("DOCKER_HOST"),
			"cert path": os.Getenv("DOCKER_CERT_PATH"),
		}).Debug("Getting connection details from environment")
		log.WithField("machine", machName).Debug("Getting connection details from Docker Machine")

		dockerHost, dockerCertPath, err = getMachineDeets(machName) // machName is "eris" by default

		if err != nil {
			log.Debug("Could not connect to Eris Docker Machine")
			log.Errorf("Trying %q Docker Machine: %v", "default", err)
			dockerHost, dockerCertPath, err = getMachineDeets("default") // during toolbox setup this is the machine that is created
			if err != nil {
				log.Debugf("Could not connect to %q Docker Machine", "default")
				log.Debugf("Error: %v", err)
				log.Debug("Trying to set up new machine")
				if e2 := CheckDockerClient(); e2 != nil {
					IfExit(fmt.Errorf("%v\n", e2))
				}
				dockerHost, dockerCertPath, _ = getMachineDeets("eris")
			}

		}

		log.WithFields(log.Fields{
			"host":      dockerHost,
			"cert path": dockerCertPath,
		}).Debug()

		if err := connectDockerTLS(dockerHost, dockerCertPath); err != nil {
			IfExit(fmt.Errorf("Error connecting to Docker Backend via TLS.\nERROR =>\t\t\t%v\n", err))
		}
		log.Debug("Successfully connected to Docker daemon")

		setIPFSHostViaDockerHost(dockerHost)
	}
}

func CheckDockerClient() error {
	if runtime.GOOS == "linux" {
		return nil
	}

	var input string
	dockerHost, dockerCertPath := popHostAndPath()

	if dockerCertPath == "" || dockerHost == "" {
		driver := "virtualbox" // when we use agents we'll wanna turn this driver into a flag

		if runtime.GOOS == "windows" {
			if err := prepWin(); err != nil {
				return fmt.Errorf("Could not add ssh.exe to PATH.\nError:%v\n", err)
			}
		}

		if _, _, err := getMachineDeets("default"); err == nil {
			fmt.Println("A Docker Machine VM exists, which Eris can use")
			fmt.Println("However, our marmots recommend that you have a VM dedicated to Eris dev-ing")
			fmt.Print("Would you like the marmots to create a machine for you? (y/n): ")
			fmt.Scanln(&input)

			if input == "Y" || input == "y" || input == "YES" || input == "Yes" || input == "yes" {
				log.Debug("The marmots will create an Eris machine")
				if err := setupErisMachine(driver); err != nil {
					return err
				}

				log.WithField("driver", driver).Info("New Docker Machine created")
				log.Info("Getting proper environment variables")

				if _, _, err := getMachineDeets("eris"); err != nil {
					return err
				}
			} else {
				log.Info("No Eris Docker Machine will be created")
			}

		} else {

			fmt.Println("The marmots could not find a Docker Machine VM they could connect to")
			fmt.Println("Our marmots recommend that you have a VM dedicated to eris dev-ing")
			fmt.Print("Would you like the marmots to create a machine for you? (y/n): ")
			fmt.Scanln(&input)

			if input == "Y" || input == "y" || input == "YES" || input == "Yes" || input == "yes" {
				log.Warn("The marmots will create an Eris machine")
				if err := setupErisMachine(driver); err != nil {
					return err
				}

				log.WithField("driver", driver).Info("New Docker Machine created")
				log.Info("Getting proper environment variables")
				if _, _, err := getMachineDeets("eris"); err != nil {
					return err
				}
			}

		}
	}

	log.Info("Docker client connected")
	return nil
}

func getMachineDeets(machName string) (string, string, error) {
	var out = new(bytes.Buffer)
	var out2 = new(bytes.Buffer)

	noConnectError := fmt.Errorf("Could not evaluate the env vars for the %s docker-machine.\n", machName)
	dHost, dPath := popHostAndPath()

	if (dHost != "" && dPath != "") && (machName == "eris" || machName == "default") {
		return dHost, dPath, nil
	}

	// TODO: when go-dockerclient adds machine API endpoints use those instead.
	log.WithField("machine", machName).Debug("Querying Docker Machine URL")
	cmd := exec.Command("docker-machine", "url", machName)
	cmd.Stdout = out
	if err := cmd.Run(); err != nil {
		return "", "", fmt.Errorf("%vError:\t%v\n", noConnectError, err)
	}
	dHost = strings.TrimSpace(out.String())
	log.WithField("host", dHost).Debug()

	// TODO: when go-dockerclient adds machine API endpoints use those instead.
	log.WithField("machine", machName).Debug("Querying Docker Machine cert path")
	cmd2 := exec.Command("docker-machine", "inspect", machName, "--format", "{{.HostOptions.AuthOptions.ServerCertPath}}")
	cmd2.Stdout = out2
	//cmd2.Stderr = os.Stderr
	if err := cmd2.Run(); err != nil {
		return "", "", fmt.Errorf("%vError:\t%v\n", noConnectError, err)
	}
	dPath = out2.String()
	dPath = strings.Replace(dPath, "'", "", -1)
	dPath = filepath.Dir(dPath)
	log.WithField("cert path", dPath).Debug()

	if dPath == "" || dHost == "" {
		return "", "", noConnectError
	}

	log.Info("Querying host and user have access to the right files for TLS connection to Docker")
	if err := checkKeysAndCerts(dPath); err != nil {
		return "", "", err
	}
	log.Debug("Certificate files look good")

	// technically, do not *have* to do this, but it will make repetitive tasks faster
	log.Debug("Setting environment variables for quick future development")
	os.Setenv("DOCKER_HOST", dHost)
	os.Setenv("DOCKER_CERT_PATH", dPath)
	os.Setenv("DOCKER_TLS_VERIFY", "1")
	os.Setenv("DOCKER_MACHINE_NAME", machName)

	log.WithField("machine", machName).Debug("Finished getting machine details")
	return dHost, dPath, nil
}

func DockerClientVersion() (float64, error) {
	verR, err := DockerClient.Version()
	if err != nil {
		return 0, err
	}

	verD := verR.Get("Version")
	v := strings.Split(verD, ".")
	v = v[0:2] // we only want 1.8 so we can marshal that into a float64

	return strconv.ParseFloat(strings.Join(v, "."), 10)
}

func DockerAPIVersion() (float64, error) {
	verR, err := DockerClient.Version()
	if err != nil {
		return 0, err
	}

	verD := verR.Get("APIVersion")
	return strconv.ParseFloat(verD, 10)
}

func setupErisMachine(driver string) error {
	cmd := "docker-machine"
	args := []string{"status", "eris"}
	if err := exec.Command(cmd, args...).Run(); err == nil {
		// if err == nil this means the machine is created. if err != nil that means machine doesn't exist.
		log.Debug("Eris Docker Machine exists. Starting")
		return startErisMachine()
	}
	log.Debug("Eris Docker Machine doesn't exist")

	return createErisMachine(driver)
}

func createErisMachine(driver string) error {
	log.Warn("Creating the Eris Docker Machine")
	log.Warn("This will take some time, please feel free to go feed your marmot")
	log.WithField("driver", driver).Debug()
	cmd := "docker-machine"
	args := []string{"create", "--driver", driver, "eris"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		log.Debugf("There was an error creating the Eris Docker Machine: %v", err)
		return mustInstallError()
	}
	log.Debug("Eris Docker Machine created")

	return startErisMachine()
}

func startErisMachine() error {
	log.Info("Starting Eris Docker Machine")
	cmd := "docker-machine"
	args := []string{"start", "eris"}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		return fmt.Errorf("There was an error starting the newly created docker-machine.\nError:\t%v\n", err)
	}
	log.Debug("Eris Docker Machine started")

	return nil
}

func connectDockerTLS(dockerHost, dockerCertPath string) error {
	var err error

	log.WithFields(log.Fields{
		"host":      dockerHost,
		"cert path": dockerCertPath,
	}).Debug("Connecting to Docker via TLS")
	DockerClient, err = docker.NewTLSClient(dockerHost, filepath.Join(dockerCertPath, "cert.pem"), filepath.Join(dockerCertPath, "key.pem"), filepath.Join(dockerCertPath, "ca.pem"))
	if err != nil {
		return err
	}

	log.Debug("Connected via TLS")
	return nil
}

func popHostAndPath() (string, string) {
	return os.Getenv("DOCKER_HOST"), os.Getenv("DOCKER_CERT_PATH")
}

func checkKeysAndCerts(dPath string) error {
	toCheck := []string{"cert.pem", "key.pem", "ca.pem"}
	for _, f := range toCheck {
		f = filepath.Join(dPath, f)
		if _, err := os.Stat(f); err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("The marmots could not find a file that was required to connect to Docker.\nThey get a file does not exist error from the OS.\nFile needed:\t%s\n", f)
			} else if os.IsNotExist(err) {
				return fmt.Errorf("The marmots could not find a file that was required to connect to Docker.\nThey get a permissions error for the file.\nPlease check your file permissions.\nFile needed:\t%s\n", f)
			} else {
				return fmt.Errorf("The marmots could not find a file that was required to connect to Docker.\nThe file exists and the user has the right permissions.\nColor the marmots confused.\nFile needed:\t%s\nError:\t%v\n", f, err)
			}
		}
	}
	return nil
}

func mustInstallError() error {
	errBase := "The marmots cannot connect to Docker.\nDo you have docker installed?\nIf not please visit here:\t"
	dInst := "https://docs.docker.com/installation/"

	switch runtime.GOOS {
	case "linux":
		return fmt.Errorf("%s%s\nDo you have docker installed and running?\nIf not please [sudo services start docker] on Ubuntu.\nAlso check that your user is in the docker group (or rerun with sudo).\nTo fix this please run [sudo usermod -a -G docker $USER] on Ubuntu with your user substituted.", errBase, dInst)
	case "darwin":
		return fmt.Errorf("%s%s\n", errBase, (dInst + "mac/"))
	case "windows":
		return fmt.Errorf("%s%s\n", errBase, (dInst + "windows/"))
	default:
		return fmt.Errorf("%s%s\n", errBase, dInst)
	}

	return nil
}

// need to add ssh.exe to PATH, it resides in GIT dir.
// see: https://docs.docker.com/installation/windows/#from-your-shell
func prepWin() error {
	// note this is for running from cmd.exe ... watch out for powershell....
	cmd := exec.Command("set", `PATH=%PATH%;"c:\Program Files (x86)\Git\bin"`)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func setIPFSHostViaDockerHost(dockerHost string) {
	u, err := url.Parse(dockerHost)
	if err != nil {
		IfExit(fmt.Errorf("The marmots could not parse the URL for the DockerHost to populate the IPFS Host.\nPlease check that your docker-machine VM is running with [docker-machine ls]\nError:\t%v\n", err))
	}
	dIP, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		IfExit(fmt.Errorf("The marmots could not split the host and port for the DockerHost to populate the IPFS Host.\nPlease check that your docker-machine VM is running with [docker-machine ls]\nError:\t%v\n", err))

	}
	dockerIP := fmt.Sprintf("%s%s", "http://", dIP)

	log.WithField("url", dockerIP).Debug("Setting ERIS_IPFS_HOST")
	os.Setenv("ERIS_IPFS_HOST", dockerIP)
}
