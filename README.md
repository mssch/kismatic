# Kismatic Platform

## Using Cobra
`go install github.com/spf13/cobra/cobra` installs the cobra CLI  
`cobra add command1 --config cobra.yaml` creates `command1` using the specified config file   
`cobra add command2 -p 'command1Cmd' --config cobra.yaml` creates `command2` as a *child* of `command1`
