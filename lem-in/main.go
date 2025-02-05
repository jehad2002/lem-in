package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AntRoom struct {
	nbr  int
	path []string
	Ant  []string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <filename>")
		return
	}

	fileName := os.Args[1]
	data, err := readInputFile(fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
// ********data analysis ******************in parsegraphdata 
	graph, startNode, endNode, numAnts := parseGraphData(data)
	if startNode == endNode {
		checkError("ERROR: Start room is equal to End room")
	}
	if numAnts == 0 {
		checkError("ERROR: No ants specified")
	}
// *************var to genarateAntsname************** var for utilpath to get them for way*******************
	antNames := generateAntNames(numAnts)
	utilsPaths := utilPaths(graph, startNode, endNode)

	if len(utilsPaths) == 0 {
		checkError("ERROR: No valid paths found")
	}
// change path to strct and put in the app
	paths := matriceToStruct(utilsPaths)
	resolve(paths, antNames)
	displayAnt(paths, antNames)
}

func readInputFile(fileName string) (string, error) {
	f, err := os.Open("./file/" + fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		content.WriteString(scanner.Text())
		content.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return content.String(), nil
}
// get the data  and data it to []string
func parseGraphData(data string) (map[string][]string, string, string, int) {
	lines := strings.Split(strings.TrimSpace(data), "\n")

	graph := make(map[string][]string)
	var startNode, endNode string
	var numAnts int
// for to check the line and read the numants in line not nil
	for i, line := range lines {
		if line != "" && i == 0 {
			numAnts, _ = strconv.Atoi(line)
			continue
		}
// read the line after ##start/end to get the rooms
		if line == "##start" {
			if i+1 < len(lines) {
				startNode = strings.Fields(lines[i+1])[0]
			}
		} else if line == "##end" {
			if i+1 < len(lines) {
				endNode = strings.Fields(lines[i+1])[0]
			}
		// if theres error - in line put it in tow par
		} else if strings.Contains(line, "-") {
			parts := strings.Split(line, "-")
			node1 := parts[0]
			node2 := parts[1]

			graph[node1] = append(graph[node1], node2)
			graph[node2] = append(graph[node2], node1)
		}
	}

	return graph, startNode, endNode, numAnts
}

func utilPaths(graph map[string][]string, startNode string, endNode string) [][]string {
	allPaths := findPaths(graph, startNode, endNode)
	if len(allPaths) == 0 {
		return [][]string{} // Return empty slice if no paths found
	}

	allCombinations := [][][]string{}

	for i := range allPaths {
		allPathsHelp := moveSliceToBeginning(allPaths, i)
		current := allPaths[i]

		validPaths := [][]string{}
		validPaths = append(validPaths, current)
		tab := [][]string{}
		tab = append(tab, current)
		for _, path := range allPathsHelp[1:] {
			isHere := false
			for _, v := range tab {
				if !compareSlices(path, v) {
					help := path[:len(path)-1]
					for _, node := range help[1:] {
						if contains(v, node) {
							isHere = true
						}
					}
				}
			}
			// if the path in the tab 
			if !isHere {
				tab = append(tab, path)
				validPaths = append(validPaths, path)
				allCombinations = append(allCombinations, validPaths)
			}
		}
	}

	if len(allCombinations) == 0 {
		return allPaths[:1] // return the first way
	}

	return maxLen(allCombinations) // max path
}

func findPaths(graph map[string][]string, startNode string, endNode string) [][]string {
	queue := [][]string{{startNode}}
	paths := [][]string{} // to save the all path available
	for len(queue) > 0 {
		path := queue[0]
		queue = queue[1:]
//delee this path
		lastNode := path[len(path)-1]
// get the last node in the path **************
		if lastNode == endNode {
			paths = append(paths, path)
		}

		for _, adjacent := range graph[lastNode] {
			if !contains(path, adjacent) {
				newPath := append(append([]string(nil), path...), adjacent)
				queue = append(queue, newPath)
			}
		}
	}
	return paths
	//***********************************if the way empty add it to the new path*********************************************************
}

func moveSliceToBeginning(slice [][]string, index int) [][]string {
	if index < 0 || index >= len(slice) {
		return slice
	}
// add list to save the result
	result := make([][]string, len(slice))
	copy(result[0:], slice[index:index+1]) 
	copy(result[1:], slice[0:index])
	copy(result[index+1:], slice[index+1:]) // copt to the end
	return result
}

func compareSlices(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	for i, val := range slice1 {
		if val != slice2[i] {
			return false
		}
	}

	return true
}

func maxLen(listPath [][][]string) [][]string {
	max := 0
	utils := [][]string{}
	for _, path := range listPath {
		if len(path) > max {
			utils = path
			max = len(path)
		}
	}

	return utils
}

func generateAntNames(nbrFourmi int) []string {
	antNames := make([]string, nbrFourmi)
	for i := 1; i <= nbrFourmi; i++ {
		antNames[i-1] = "L" + strconv.Itoa(i)
	}
	return antNames
}

func resolve(paths []*AntRoom, nameAnts []string) {
	if len(paths) == 0 {
		checkError("ERROR: No paths provided to resolve")
	}

	currentPath := paths[0]
	currentPath.nbr++
	currentPath.Ant = append(currentPath.Ant, nameAnts[0])
	index := 0
	for i := 1; i < len(nameAnts); i++ {
		if index+1 < len(paths) {
			placeAnt(paths, index, nameAnts[i])
			index++
		} else {
			index = 0
			placeAnt(paths, index, nameAnts[i])
			index++
		}
	}
}

func placeAnt(paths []*AntRoom, currentIndex int, nameAnt string) {
	if len(paths) == 1 {
		paths[currentIndex].Ant = append(paths[currentIndex].Ant, nameAnt)
		paths[currentIndex].nbr++
		return
	}
	currentPath := paths[currentIndex]
	nextPath := paths[currentIndex+1]
	initialPath := paths[0]

	if currentPath.nbr+1 > nextPath.nbr+1 {
		nextPath.nbr++
		nextPath.Ant = append(nextPath.Ant, nameAnt)
	} else if nextPath.nbr >= initialPath.nbr {
		initialPath.nbr++
		initialPath.Ant = append(initialPath.Ant, nameAnt)
	} else {
		currentPath.nbr++
		currentPath.Ant = append(currentPath.Ant, nameAnt)
	}
}

func matriceToStruct(paths [][]string) []*AntRoom {
	var pathStruct = make([]*AntRoom, len(paths))
	for i, path := range paths {
		pathStruct[i] = &AntRoom{
			nbr:  len(path),
			path: path[1:],
		}
	}
	return pathStruct
}

func displayAnt(paths []*AntRoom, nameAnt []string) {
	End := paths[0].path[len(paths[0].path)-1]
	var test, i int
	for {
		for _, room := range paths {
			l := i
			for k := 0; k <= i; k++ {
				if k < len(room.Ant) && l < len(room.path) {
					fmt.Print(room.Ant[k] + "-" + room.path[l] + " ")
					if room.path[l] == End && room.Ant[k] == nameAnt[len(nameAnt)-1] {
						test = 1
					}
				}
				l--
			}
		}
		i++
		fmt.Println()
		if test == 1 {
			break
		}
	}
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

func checkError(err string) {
	fmt.Println(err)
	os.Exit(1)
}
