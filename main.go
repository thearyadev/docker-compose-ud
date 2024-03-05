package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func hasComposeFile(directory string) bool {
	_, errYML := os.Stat(filepath.Join(directory, "docker-compose.yml"))
	_, errYAML := os.Stat(filepath.Join(directory, "docker-compose.yaml"))
	_, errCompose := os.Stat(filepath.Join(directory, "compose.yml"))
	_, errComposeYAML := os.Stat(filepath.Join(directory, "compose.yaml"))
	return !os.IsNotExist(errYML) || !os.IsNotExist(errYAML) || !os.IsNotExist(errCompose) || !os.IsNotExist(errComposeYAML)
}
func get_action(action string) func(string) {

	switch action {
	case "up":
		return func(file string) {
			command := exec.Command("docker", "compose", "up", "-d")
			err := command.Run()
			if err != nil {
				log.Fatalf("[ERROR] %s failed: %s", file, err)
			}

		}
	case "down":
		return func(file string) {
			command := exec.Command("docker", "compose", "down")
			err := command.Run()
			if err != nil {
				log.Printf("[ERROR] %s failed: %s", file, err)
			}
		}
	case "update":
		return func(file string) {
			command_down := exec.Command("docker", "compose", "down")
			err := command_down.Run()
			if err != nil {
				log.Printf("[ERROR] %s failed: %s", file, err)
			}

			command_pull := exec.Command("docker", "compose", "pull")

			err = command_pull.Run()
			if err != nil {
				log.Printf("[ERROR] %s failed: %s", file, err)
			}		
			command_up := exec.Command("docker", "compose", "up", "-d")
			err = command_up.Run()
			if err != nil {
				log.Printf("[ERROR] %s failed: %s", file, err)
			}
		}

	}
	return nil
}

func main() {
	var dir = flag.String("d", "", "Directory to run this program")
	var mode = flag.String("m", "", "up, down, update")

	flag.Parse()
	action := get_action(*mode)

	if action == nil {
		log.Fatal("Not a valid action")
	}

	files, err := os.ReadDir(*dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if !file.IsDir() {
			log.Printf("[WARN] %s is not a directory", file.Name())
			continue
		}

		if !hasComposeFile(filepath.Join(*dir, file.Name())) {
			log.Printf("[WARN] %s does not have a compose file", file.Name())
			continue
		}
		os.Chdir(filepath.Join(*dir, file.Name()))
		log.Printf("[INFO] %s", file.Name())
		action(file.Name())
		os.Chdir("../")
	}
}
