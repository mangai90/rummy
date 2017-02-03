// 500 rummy game play.
// Contains Player info and the game state.
// Code that publishes and subscribes the game play at each turn.
// Still in progress..
// TODO:full run through of a hand.

package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
    "github.com/pubnub/go/messaging"
)

const (
    sub_key = "sub-c-4a943d42-d4cb-11e6-8b9b-02ee2ddab7fe"
    pub_key = "pub-c-d57052d6-2b7f-4488-8504-5e3fd36a5fe2"
)

var (
    deck Deck
    player *Player
)

type Player struct {
    Name string
    FirstTurn bool
    Hand []Card
    Facedowns []Facedown
}

type Message struct {
    Type string
    DnD Deck
    NextTurn bool
    Hand Cards
    Facedowns []Facedown
    HandPoints int
}

type GameState struct {
    Hand Cards
    Facedowns []Facedown
    GamePoints int
    HandPoints int
    OppFacedowns []Facedown
    OppGamePoints int
    OppHandPoints int
    Hands int
}

func (state *GameState) updateMyturn(fds []Facedown, handpts int) {
    state.Facedowns = append(state.Facedowns, fds...)
    state.HandPoints += handpts
}

func (state *GameState) updateOppturn(fds []Facedown, handpts int) {
    state.OppFacedowns = append(state.OppFacedowns, fds...)
    state.OppHandPoints += handpts
}

func (state *GameState) updateHand(handpts, opppts int) {
    state.GamePoints += state.HandPoints + handpts
    state.HandPoints = 0
    state.OppGamePoints += state.OppHandPoints + opppts
    state.OppHandPoints = 0
    state.Facedowns = []Facedown{}
    state.OppFacedowns = []Facedown{}
    state.Hands += 1
    state.Hand = []Card{}
}

func NewPlayer(name string, turn bool) *Player {
    return &Player{Name: name, FirstTurn: turn}
}

func (p *Player)  playTurn() {
    // assuming this works.
    // Pick Up:
    // Deck or Discard
    // update state
    // publish deck
    
    // Facedowns: optional
    // update state
    // publish facedowns
    
    // Discard:
    // update state
    // publish deck
}

func main() {
    deck.Shuffle()

    pubnub := messaging.NewPubnub(pub_key, sub_key, "", "", false, "")
    channel := "rummy"

    successChannel := make(chan []byte)
    errorChannel := make(chan []byte)
    pubsuccessChannel := make(chan []byte)
    puberrorChannel := make(chan []byte)

    // check presence
    // find a player to start game
    // oppentent agrees to play
    // optional not implemented yet.
    // Default channel "Rummy".
    name := "X"
    firstturn := true
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter your name (X) or Y: ")
    if n, err := reader.ReadString('\n'); err != nil && n != "" {
        if n != "X" && n != "x" {
            firstturn = false
            name = n
        }
    }
    player = NewPlayer(name, firstturn)

    fmt.Println("Welcome to Rummy,", name)

    fmt.Print("Enter Channel name that you want to join (rummy)")
    if ch, err := reader.ReadString('\n'); err != nil && ch != "" {
        channel = ch
    }
    fmt.Println(name, "has joined Channel", channel)
    // check presence
    msg := fmt.Sprint("Player", name, "joined")
    // 500 pts loop

    // First turn player calls Deal()
    // initialise own state.
    // publish Opp's Hand as msg type: deal

    // If not first turn player
    // wait for msg type : deal on Subcribe
    // initialise state with the Hand in msg
    if player.FirstTurn {
        myhand, opphand := deck.Deal()
        game := GameState{Hand: myhand}
        fmt.Print(game.Hand.Display())
        fmt.Print(opphand.Display())

        // publish opphand to channel
    } else {
        // wait for deal msg
    }
    // loop through till one player has no cards left
    // each turn wait for subs msg with nextturn == myturn
    // playturn()
    // If Gameover update state, publish state
    // updateHand()
    go pubnub.Subscribe(channel, "", successChannel, false, errorChannel)
    go handlesSubResponce(successChannel, errorChannel)
    go pubnub.Publish(channel, msg, pubsuccessChannel, puberrorChannel)
    go handlesPubResponce(pubsuccessChannel, puberrorChannel)
    //player.playTurn(reader, state)
}

func handlesSubResponce(successChannel, errorChannel chan []byte) {
    // On sucess check msg type
    // type = deal, initiate player's hand.
    // type = deck, update state and display
    // type = turn, update state, call playturn()

    // On Error display error
    for {
        select {
        case response := <-successChannel:
            var msg []interface{}
            err := json.Unmarshal(response, &msg)
            if err != nil {
                fmt.Println(err)
                return
            }
            switch m := msg[0].(type) {
            case float64:
                fmt.Println(msg[1].(string))
            case []interface{}:
                fmt.Printf("Received message '%s' on channel '%s'\n", m[0], msg[2])
                return
            default:
                panic(fmt.Sprintf("Unknown type: %T", m))
            }
        case err := <-errorChannel:
            fmt.Println(string(err))
        case <-messaging.SubscribeTimeout():
            fmt.Println("Subscribe() timeout")
        }
    }
}

func handlesPubResponce(successChannel, errorChannel chan []byte) {
    // On success echo "Published turn"
    // On error echo error
    select {
    case response := <-successChannel:
        fmt.Println(string(response))
    case err := <-errorChannel:
        fmt.Println(string(err))
    case <-messaging.Timeout():
        fmt.Println("Publish() timeout")
    }
}
