package main

type Command = []string

// func executeCommand(c Command) {
// 	prog := c[0]
// 	args := c[1:]
// 	cmd := exec.Command(prog, args...)
// 	if err := cmd.Run(); err != nil {
// 		log.Fatal(err)
// 	}
// }

var GitInit = Command{"git", "init"}

// func runPreFilesCreationCommands() {
// 	commands := []Command{
// 		// Example command
// 		// GitInit,
// 	}
// 	for _, cmd := range commands {
// 		executeCommand(cmd)
// 	}
// }

// func runPostFilesCreationCommands() {
// }
