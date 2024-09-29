package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	todos := []Todo{{1, "Finish this app", "high", false}}
	currentId := 1

	printTable(todos)

	for {
		fmt.Print("Enter command: ")
		handleInput(reader, &todos, &currentId)
		printTable(todos)
	}
}

func printTable(todos []Todo) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Task", "Priority", "Done"})
	table.SetAlignment(tablewriter.ALIGN_CENTER)
	for _, v := range todos {
		table.Append([]string{
			strconv.Itoa(v.ID),
			v.Task,
			v.Priority,
			formatAndColorDoneStatus(v.Done),
		})
	}
	table.Render()
	fmt.Println("Write 'info' into the terminal for more information.")
	fmt.Println()
}

func printInfo() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"cmd", "Operation", "Format"})
	table.SetRowLine(true)
	table.SetAutoWrapText(false)
	data := [][]string{
		{"a", "adding new todo", "\033[34ma;yourTodo;priority.\033[0m"},
		{"d", "making your todo done", "\033[34md;id\033[0m"},
		{"p", "changing priority of your todo", "\033[34mp;id;newPriority\033[0m"},
		{"r", "removing todo", "\033[34mr;id\033[0m"},
		{"ra", "removing all your todos", "\033[34mra\033[0m"},
	}
	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	fmt.Println("Example: " + ColorBlue + "'a;Finish this app;high'" + ColorReset)
	fmt.Println()
}

func handleInput(reader *bufio.Reader, todos *[]Todo, currentId *int) {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	slice := strings.Split(input, ";")
	operation := slice[0]

	switch operation {
	case CommandAdd:
		err := validateLength(slice, 3, "Command format for add: a;title;priority")
		if err != nil {
			printError(err)
			return
		}

		newTodoItem := newTodo(slice[1], slice[2], currentId)
		*todos = append(*todos, *newTodoItem)
		clearTerminal()

	case CommandPriority:
		inputId, err := validatLengthAndID(slice, 3, "Command format for changing priority: p;id;priority")
		if err != nil {
			return
		}

		for i := 0; i < len(*todos); i++ {
			if (*todos)[i].ID == inputId {
				(*todos)[i].Priority = slice[2]
				clearTerminal()
				return
			}
		}
		idDoesNotExist()
		return

	case CommandDone:
		err := validateLength(slice, 2, "Command format for done: d;id")
		if err != nil {
			printError(err)
			return
		}

		inputId, err := validateIdFormat(slice)
		if err != nil {
			printError(err)
			return
		}

		for i := 0; i < len(*todos); i++ {
			if (*todos)[i].ID == inputId {
				(*todos)[i].makeDone()
				clearTerminal()
				return
			}
		}
		idDoesNotExist()
		return

	case CommandRemove:
		err := validateLength(slice, 2, "Command format for remove: r;id")
		if err != nil {
			printError(err)
			return
		}

		inputId, err := validateIdFormat(slice)
		if err != nil {
			printError(err)
			return
		}

		for i := 0; i < len(*todos); i++ {
			if (*todos)[i].ID == inputId {
				*todos = append((*todos)[:i], (*todos)[i+1:]...)
				clearTerminal()
				return
			}
		}
		idDoesNotExist()
		return

	case CommandRemoveAll:
		err := validateLength(slice, 1, "Command format for deleting all tasks: ra")
		if err != nil {
			printError(err)
			return
		}
		*todos = []Todo{}
		clearTerminal()
		return

	case CommandInfo:
		err := validateLength(slice, 1, "Write 'info' for more information.")
		if err != nil {
			printError(err)
			return
		}
		clearTerminal()
		printInfo()
		return

	default:
		err := fmt.Errorf("ERROR: Oops, something went wrong. Try again or write 'info' for more information.")
		printError(err)
		return
	}
}

func validatLengthAndID(slice []string, expectedLength int, formatErrorMessage string) (int, error) {
	err := validateLength(slice, expectedLength, formatErrorMessage)
	if err != nil {
		printError(err)
		return 0, err
	}

	id, err := validateIdFormat(slice)
	if err != nil {
		printError(err)
		return 0, err
	}
	return id, nil
}

func validateIdFormat(slice []string) (int, error) {
	id, err := strconv.Atoi(slice[1])
	if err != nil {
		clearTerminal()
		return 0, fmt.Errorf("ERROR: Invalid ID format. ID should be a number.")
	}
	return id, nil
}

func validateLength(slice []string, expectedLength int, formatErrorMessage string) error {
	errorStart := "ERROR: Invalid input. "
	if len(slice) != expectedLength {
		return fmt.Errorf(errorStart + formatErrorMessage)
	}
	return nil
}

func printError(err error) {
	clearTerminal()
	fmt.Println(ColorRed + err.Error() + ColorReset)
}

func clearTerminal() {
	var c *exec.Cmd
	if strings.Contains(runtime.GOOS, "windows") {
		c = exec.Command("cmd", "/c", "cls")
	} else {
		c = exec.Command("clear")
	}
	c.Stdout = os.Stdout
	c.Run()
}

func idDoesNotExist() {
	err := fmt.Errorf("ERROR: Todo with this ID doesn't exist.")
	printError(err)
	return
}

func formatAndColorDoneStatus(done bool) string {
	formatedDone := strconv.FormatBool(done)

	if !done {
		return ColorYellow + formatedDone + ColorReset
	}
	return ColorGreen + formatedDone + ColorReset
}

func (t *Todo) makeDone() {
	t.Done = true
	return
}

func newTodo(task string, prio string, id *int) *Todo {
	*id++
	return &Todo{
		ID:       *id,
		Task:     task,
		Priority: prio,
		Done:     false,
	}
}

type Todo struct {
	ID       int
	Task     string
	Priority string
	Done     bool
}

const (
	CommandAdd       = "a"
	CommandDone      = "d"
	CommandPriority  = "p"
	CommandRemove    = "r"
	CommandRemoveAll = "ra"
	CommandInfo      = "info"

	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorReset  = "\033[0m"
)
