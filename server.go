package main

type FundServer struct {
    commands chan TransactionCommand
    fund *Fund
}

func NewFundServer(initialBalance int) *FundServer {
    server := &FundServer{
        // make() creates builtins like channels, maps, and slices
        commands: make(chan TransactionCommand),
        fund: NewFund(initialBalance),
    }

    // Spawn off the server's main loop immediately
    go server.loop()
    return server
}

// Typedef the callback for readability
type Transactor func(fund interface{})

// Add a new command type with a callback and a semaphore channel
type TransactionCommand struct {
    Transactor Transactor
    Done chan bool
}

func (s *FundServer) Balance() int {
    var balance int
    s.Transact(func(managedValue interface{}) {
        f := managedValue.(*Fund)
        balance = f.Balance()
    })
    return balance
}

func (s *FundServer) Withdraw(amount int) {
    s.Transact(func (managedValue interface{}) {
        f := managedValue.(*Fund)
        f.Withdraw(amount)
    })
}

func (s *FundServer) Transact(transactor Transactor) {
    command := TransactionCommand{
        Transactor: transactor,
        Done: make(chan bool),
    }
    s.commands <- command
    <- command.Done
}

func ( s *FundServer) loop() {
     for transaction := range s.commands {
        // Now we don't need any type-switch mess
        transaction.Transactor(s.fund)
        transaction.Done <- true
     }
}

type WithdrawCommand struct {
     Amount int
}

type BalanceCommand struct {
     Response chan int
}