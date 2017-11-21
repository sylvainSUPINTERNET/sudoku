package main;

import (
	"fmt"
	_"math/rand"
	_"time"
	"math/rand"
	"time"
	"os"
	"io/ioutil"
	"io"
	"bufio"
	"strconv"
	_"reflect"
	_"sync"
	"runtime"

)

func main() {

	nbCore := runtime.NumCPU();
	fmt.Println(nbCore)

	var channel chan Sudoku
	channel = make(chan Sudoku)


	sudoku := Sudoku{}

	//manualy or random
	grid := sudoku.grid
	grid = fillDdefault(sudoku)
	//fmt.Println(grid); //grid filled auto

	/*
	gridRandom := sudoku.grid
	gridRandom = fillRandom(sudoku)
	fmt.Println(gridRandom); //grid filled rand
	*/

	start := time.Now()
	display(solve(grid, 1))
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed);

	//var wg sync.WaitGroup
	//into go devant display boule wg.Add(1) a la fin de la boucle wg.Wait()
	//dans display prendre en param wg *sync.WaitGroup wg.Done()


	for i := 0; i <= 5; i++ {
		grid_to_resolve := sudoku.grid; //empty grid to push value from .txt

		file := getSudokuGame("game", i);
		scanner := bufio.NewScanner(file)
		lineGrid := 0;
		for scanner.Scan() { // internally, it advances token based on sperator
			//fmt.Println(scanner.Text())  // token in unicode-char
			//fmt.Println("Length", len(scanner.Text()))

			retourScan := scanner.Text()
			retourScanConvRune := []rune(retourScan);
			lineGrid += 1
			//y == line number
			fmt.Println("game", i);
			fmt.Println("line", lineGrid, retourScanConvRune)
			fmt.Println("\n")
			//fmt.Println("value index :", 0 ,strconv.QuoteRune(retourScanConvRune[0]))
			for rowGrid := 0; rowGrid < 9; rowGrid++ {
				//fmt.Println(strconv.QuoteRune(retourScanConvRune[rowGrid]))
				//pas un 0 mais un nombre
				//push dans la grid la valeur de la rune
				numberStr := string(retourScanConvRune[rowGrid])
				number, err := strconv.Atoi(numberStr)

				fmt.Println("Le nombre ", number)
				grid_to_resolve[lineGrid-1][rowGrid] = number
				if err != nil {
					//fmt.Println("Le nombre ", number)
					//fmt.Println("line" , lineGrid)
					//fmt.Println("row", rowGrid)
					//grid_to_resolve[lineGrid-1][rowGrid] = 9

				}

			}


		}
		display(solve(grid, 1))


		b, err := ioutil.ReadAll(file)
		if err != nil {
			print(err)
		} else {
			fmt.Println("");
			fmt.Println(b)
		}

	}

}

type Sudoku struct {
	grid [9][9] int
}

//Remplissage en dur
func fillDdefault(sudoku Sudoku) [9][9]int {

	var grid = [9][9]int{}

	grid[0][1] = 7
	grid[0][4] = 2

	grid[1][0] = 5
	grid[1][8] = 4

	grid[2][2] = 2
	grid[2][5] = 6

	grid[3][0] = 3
	grid[3][3] = 6

	grid[4][1] = 5
	grid[4][7] = 8

	grid[5][5] = 3
	grid[5][8] = 7

	grid[6][2] = 5
	grid[6][4] = 3
	grid[6][6] = 2

	grid[7][4] = 6
	grid[7][8] = 8

	grid[8][2] = 9
	grid[8][7] = 3

	return grid
}

func fillRandom(sudoku Sudoku) [9][9]int {
	for row := 0; row < 9; row++ {
		for column := 0; column < 9; column++ {
			sudoku.grid[row][column] = randInt(0, 9) //0 == case vide
		}
	}
	return sudoku.grid
}

func randInt(min int, max int) int {
	return min + rand.Intn(max+1-min)
}

//Check cadran / column / rown

func checkColumn(num int, column int, grid [9][9]int) bool {
	for i := 0; i < 9; i++ {
		if grid[i][column] == num {
			return true
		}
	}
	return false
}

func checkRow(num int, row int, grid [9][9]int) bool {
	for i := 0; i < 9; i++ {
		if grid[row][i] == num {
			return true
		}
	}
	return false
}

func checkCadran(num int, row int, column int, grid [9][9]int) bool {

	var baseRow = row - (row % 3)
	var baseColumn = column - (column % 3)

	for r := baseRow; r < baseRow+3; r++ {
		for c := baseColumn; c < baseColumn+3; c++ {
			if grid[r][c] == num {
				return true
			}
		}
	}
	return false
}

func isValid(num int, row int, column int, grid [9][9]int) bool {

	inRow := checkRow(num, row, grid)
	inColumn := checkColumn(num, column, grid)
	inCadran := checkCadran(num, row, column, grid)

	if (inRow == false && inColumn == false && inCadran == false) {
		return true; //all checking rigth
	}

	return false;
}

//si on trouve un 0 dans les raw == remplace par une valeur
// sinon on tente la possibilité par la valeur déjà donné
func setPossibility(grid [9][9]int) (int, int) {
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			if grid[r][c] == 0 {
				return r, c
			}
		}
	}
	return 200, 200
}

func isSolved(grid [9][9]int) bool {
	for r := 0; r < 9; r++ {
		for i := 0; i < 9; i++ {
			if !checkRow(i+1, r, grid) {
				return false
			}
		}
	}
	return true
}

// solve => recursivement va appeler tout d'abord replaceEmpty
// setPossibility ==> check si la value est 0 ou pas / si c'est 0 alors il lui attribut une valeur entre 1 et 9 SINON return la valeur actuel (grid)
// ENSUITE va etre appeler isValid (check cadran / column / row) SI cela return true alors on attribut a la grid la vlauer aux index (param) et rappeler la method
//ENSUITE isSolved qui va check si a la row la valeur qu'on via attribuer n'existe pas SI elle existe ALORS on reset la valeur à 0 et ce qui va donc au moment du recall de solve
// nous faire rerentrer dans replaceEmpty (cette raw fraichement remis à 0) et nous set une nouvelle valeur à tester entre 1 et 9
// etc jusqua isSolved return true, et stop la boucle

func solve(grid [9][9]int, num int) [9][9]int {
	// ************************ for empty value in grid
	//on cherche une possibilité hors de 0 et on l'utilise
	var r, c int = setPossibility(grid)
	//200 => pas de 0 donc return the initial grid (perfect)
	if r == 200 && c == 200 {
		return grid
	}
	// ************************

	//recursive, try for each number between 1 - 9

	for num <= 9 {
		//if for the num, row / column / cadran are fine => place the value at this index et recall the method herself solve  and reset num at 1s
		if isValid(num, r, c, grid) {
			grid[r][c] = num
			grid = solve(grid, 1)
		}
		num++
	}

	//erreur return sur le checking de la row (value) already used (DUPLICATION)
	if !isSolved(grid) {
		grid[r][c] = 0 // => Si une duplication est noté, alors on réatribut la valeur ou il y une duplication à 0
		// et on rappel ainsi setPossibility, qui va nous setter une nouvelle possibilité entre 1 / 9
	}

	return grid //for restart all once again if necessary
}

func display(grid [9][9]int) {
	fmt.Println("");
	for i := 0; i < 9; i++ {
		if i == 3 || i == 6 {
			fmt.Println(" __________________ ")
		}
		fmt.Println(grid[i])
	}
	fmt.Println("");
}

//creer des sudoku dans les txt
//creer la methode qui les ouvrent

func getSudokuGame(nameFile string, counter int) io.Reader {
	file, err := os.Open("./sudoku_test/" + nameFile + strconv.Itoa(counter) + ".txt")
	if err != nil {
		fmt.Println(err);
	}
	//fmt.Println(file)
	return file
}
