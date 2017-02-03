// Facedown is when you showdown a run or a set to your opponent.
// Facedowns from both players are part of the game state.
// Here lies the rules for valid facedowns.
// Strike on a prior run and strikes on a strike is allowed.

package main

import (
    "fmt"
    "sort"
)

type Facedown struct {
    cards Cards
    isRun bool
    points int
}

func(fd *Facedown) Print() {
    for _, card := range fd.cards {
        fmt.Print(card.Display(), " ")
    }
    fmt.Println("Points ",fd.getPoints())
}

func (fd *Facedown) isValid() bool {
    if len(fd.cards) < 3 {
        return false
    }
    sort.Sort(fd.cards)
//    fd.Print()
    sameval, samesuit := true, true
    val := fd.cards[0].getValue()
    suit := fd.cards[0].getSuit()
    for _, card := range fd.cards {
        if card.getValue() != val {
            sameval = false
        }
        if card.getSuit() != suit {
            samesuit = false
        }
    }
    if samesuit {
        prev := fd.cards[0].getValue()
        i := 1
        if prev == 1 {
            if fd.cards[1].getValue() == 2 || fd.cards[len(fd.cards) - 1] % 13 == 0 {
                prev = fd.cards[1].getValue()
                i = 2
            }
        }
        for i < len(fd.cards) && fd.cards[i].getValue() == (prev + 1) % 13 {
            prev = fd.cards[i].getValue()
            i++
        }
        fd.isRun = (i == len(fd.cards))
        return (i == len(fd.cards))
    }
    if sameval {
        return true
    }
    return false
}

func (fd *Facedown) getPoints() int {
    var pts, fc, ace int
    for _, card := range fd.cards {
        val := card.getValue()
        if val == 0 || val > 10 {
            val = 10
            fc += 1
        } else if val == 1 {
            ace += 1
            fc += 1
        }
        pts = pts + val
    }
    if fc > 1 {
        pts = pts + (ace * 14)
    }
    return pts
}

func (fd *Facedown) addCard(card *Card) bool {
    if card.getSuit() == fd.cards[0].getSuit() {
        if (card.getValue() == fd.cards[0].getValue() - 1) ||
            (card.getValue() == fd.cards[len(fd.cards) - 1].getValue() + 1) {
            return true
        } else {
            return false
        }
    } else if card.getValue() == fd.cards[0].getValue() {
        return true
    }
    return false
}

func (fd *Facedown) addCards(cards []Card) (int, bool) {
    cards = append(cards, fd.cards...)
    nfd := Facedown{cards: cards}
    if nfd.isValid() {
        points := nfd.getPoints() - fd.getPoints()
        fd = &nfd
        return points, true
    } else {
        return 0, false
    }
}
