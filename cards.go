// Representation and functionality of playing cards.
// Struct representation of a Deck and Discard pile.
// Manipulating a deck.

package main

import (
    "fmt"
    "math/rand"
    "strings"
)

// A deck of cards is numbers 1 to 52
//   1 - 13 Spades      0
//  14 - 26 Clubs       1
//  27 - 39 Hearts      2
//  40 - 52 Diamonds    3

type Card int
type Cards []Card

const (
    Spades = iota
    Clubs
    Hearts
    Diamonds
)
var (
    symbol = [...]string{Spades:"S", Clubs:"C", Hearts:"H", Diamonds:"D"}
    facevalue = map[int]string{
        1: "A",
        2: "2",
        3: "3",
        4: "4",
        5: "5",
        6: "6",
        7: "7",
        8: "8",
        9: "9",
        10: "10",
        11: "J",
        12: "Q",
        0: "K",
    }
)

func (card *Card) isValid() bool {
    return (*card > 0 && *card < 53)
}

func (card *Card) getSuit() int {
    return int((*card - 1) / 13)
}

func (card *Card) getValue() int {
    return int(*card % 13)
}

func (card *Card) Display() string {
    val := card.getValue()
    fval, ok := facevalue[val]
    if !ok {
        return fmt.Sprint("Value not valid ", val)
    }
    suit := card.getSuit()
    return fval + symbol[suit]
}

func (cards Cards) Display() string {
    out := []string{}
    for _, c := range cards {
        out = append(out, c.Display())
    }
    return strings.Join(out, " ")
}

func (cards Cards) Len() int{
    return len(cards)
}

func (cards Cards) Less(i, j int) bool {
    return cards[i] < cards[j]
}

func (cards Cards) Swap(i, j int) {
    cards[i], cards[j] = cards[j], cards[i]
}

type Deck struct {
    Deck Cards
    Discard Cards
}

func (deck *Deck) Shuffle() {
    deck.Deck, deck.Discard = []Card{}, []Card{}
    for i := 1; i < 53; i++ {
        deck.Deck = append(deck.Deck, Card(i))
    }
    for i := range deck.Deck {
        j := rand.Intn(i + 1)
        deck.Deck[i], deck.Deck[j] = deck.Deck[j], deck.Deck[i]
    }
}

func (deck *Deck) Deal() (handX, handY Cards) {
    deck.Shuffle()
    var playerX, playerY Cards
    var card Card
    for i := 0; i <27; i++ {
        l := len(deck.Deck)
        card, deck.Deck = deck.Deck[l-1], deck.Deck[:l-1]
        if i % 2 == 0 {
            playerX = append(playerX, card)
        } else {
            playerY = append(playerY, card)
        }
    }
    return playerX, playerY
}

func (deck *Deck) PickupDeck() (*Card, error) {
    if l := len(deck.Deck); l > 0 {
        card := deck.Deck[l - 1]
        deck.Deck = deck.Deck[:l - 1]
        deck.PublishDeck()
        return &card, nil
    }
    return nil, fmt.Errorf("Deck empty")
}

func (deck *Deck) PickupDiscard(n int) ([]Card, error) {
    if l := len(deck.Discard); l >= n && n > 0 {
        cards := deck.Discard[l - n:]
        deck.Discard = deck.Discard[:l - n]
        deck.PublishDeck()
        return cards, nil
    }
    return nil, fmt.Errorf("Not enough cards in Discard")
}

func (deck *Deck) DiscardCard(card *Card) {
    deck.Discard = append(deck.Discard, *card)
    deck.PublishDeck()
}

func (deck *Deck) PublishDeck() {
    //pubnub.Publish()
/*
    json := { Decksize: len(deck.Deck),
        Discard: deck.Discard,
    }
*/
}
