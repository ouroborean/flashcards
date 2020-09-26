package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"math/rand"
	"bufio"
	"time"
)

type Flashcard struct {
	side1		string
	side2		string
	cardType	string
	wordType	string
}

type Incision struct {
	startIndex1		int
	stopIndex1		int
	startIndex2		int
	stopIndex2		int
	wordType		string
	wordTypeArray 	[]rune
	wordSide1		string
	wordSide2		string
}

func flashCardReg (card Flashcard) (string){
	var correctAnswer string
	randomSide := rand.Intn(2)

	if randomSide == 0 {
		fmt.Println(card.side1)
		correctAnswer = card.side2
	}
	if randomSide == 1 {
		fmt.Println(card.side2)
		correctAnswer = card.side1
	}
	return correctAnswer
}

func findRandomFlashCard (wordType string, allCards []Flashcard) (string, string){

	
	var validCards []Flashcard
	var chosenSide1 string
	var chosenSide2 string

	
	
	for _, card := range allCards {
		if card.wordType == wordType {
			validCards = append(validCards, card)
		}
	}

	

	var chosenCard = rand.Intn(len(validCards))

	chosenSide1 = validCards[chosenCard].side1
	chosenSide2 = validCards[chosenCard].side2
	return chosenSide1, chosenSide2
	
}

func checkForBlankSide2 (blank Incision, side string) (Incision) {

	recording := false
	var tempStartIndex int
	var tempStopIndex int
	var tempWordTypeArray []rune

	for i, v := range side {

		if string(v) == "}" {

			recording = false
			tempStopIndex = i
			if (string(tempWordTypeArray) == blank.wordType){
				blank.startIndex2 = tempStartIndex
				blank.stopIndex2 = tempStopIndex
			} else {
				tempWordTypeArray = tempWordTypeArray[:0]
			}

		}

		if recording == true {
			tempWordTypeArray = append(tempWordTypeArray, v)
		}

		if string(v) == "{" {
			recording = true
			tempStartIndex = i
		}

	}

	return blank

}

func checkForBlankSide1 (side string) (Incision, bool) {
	
	var blank Incision
	loopContinue := false
	recording := false
	for i, v := range side {


		if string(v) == "}" {
			blank.stopIndex1 = i
			blank.wordType = string(blank.wordTypeArray)

			

			blank.wordTypeArray = blank.wordTypeArray[:0]
			recording = false
			loopContinue = true
			break
		}

		if recording == true {
			blank.wordTypeArray = append(blank.wordTypeArray, v)
		}

		if string(v) == "{" {

			blank.startIndex1 = i
			recording = true

		}

	}

	return blank, loopContinue
	
}

func flashCardBlank (card Flashcard, allCards []Flashcard) (string){
	
	//Store the blankcard's sides so we can actually fucking
	//work with them
	tempSide1 := card.side1
	tempSide2 := card.side2
	var blank Incision
	//Bool to stop loop when we're finally done
	
	var correctAnswer string
	randomSide := rand.Intn(2)

	//Bool to start the wordType recording process
	
	for {
		mightBeBlanks:= true

		blank, mightBeBlanks = checkForBlankSide1(tempSide1)
		
		blank = checkForBlankSide2(blank, tempSide2)

		if mightBeBlanks == false {
			break
		}

		blank.wordSide1, blank.wordSide2 = findRandomFlashCard(blank.wordType, allCards)

		tempSide1 = string(tempSide1[:blank.startIndex1]) + blank.wordSide1 + string(tempSide1[blank.stopIndex1+1:])
		tempSide2 = string(tempSide2[:blank.startIndex2]) + blank.wordSide2 + string(tempSide2[blank.stopIndex2+1:])

		

	}

	if randomSide == 0 {
		fmt.Println(tempSide1)
		correctAnswer = tempSide2
	}
	if randomSide == 1 {
		fmt.Println(tempSide2)
		correctAnswer = tempSide1
	}

	return correctAnswer

}

func main(){
	var lastCard int
	counter := 0
	var x Flashcard
	var correctAnswer string
	for _, filename := range os.Args[1:] {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dup3: %v\n", err)
			continue
		}
		
		var allCards []Flashcard

		allStrings := strings.Split(string(data), "\n")
		
		for _, line := range allStrings {
			
	
			
			if counter == 0 {
				x.cardType = strings.TrimSpace(line)
			}
			if counter == 1 {
				x.side1 = strings.TrimSpace(line)
			}
			if counter == 2 {
				x.side2 = strings.TrimSpace(line)
			}
			if counter == 3 {
				x.wordType = strings.TrimSpace(line)
				allCards = append(allCards, x)
				counter = -1
			}
			counter++
		}

		
		for {
			rand.Seed(time.Now().UnixNano())
			randomCard := rand.Intn(len(allCards))
			
			if randomCard == lastCard {

				if (lastCard == len(allCards)-1){
					randomCard--
				} else {
					randomCard++
				}

			}

			lastCard = randomCard

			if allCards[randomCard].cardType == "flashcard" {
				correctAnswer = flashCardReg(allCards[randomCard])
			}
			if allCards[randomCard].cardType == "blankcard" {
				correctAnswer = flashCardBlank(allCards[randomCard], allCards)
			}


			input := bufio.NewScanner(os.Stdin)
			input.Scan()
			guess := input.Text()

			if guess == "exit" {
				break
			}
			if guess == correctAnswer {
				fmt.Println("Correct!\n")
			} else {
				fmt.Println("Incorrect :(")
				fmt.Println("The correct answer was :", correctAnswer + "\n")
			}
		}



	}

	

}