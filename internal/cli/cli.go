package cli

type Command struct {
    Name        string
    Description string
    Callback    func(args []string) error
}
