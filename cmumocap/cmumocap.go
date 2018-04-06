package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/kniren/gota/dataframe"
	"github.com/kniren/gota/series"
	"github.com/berlin/goc3d"
	"github.com/kzahedi/goent/continuous/state"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func getData(index int, label string, data goc3d.C3DData, useJerk bool) ([]string, [][]float64) {
	if index >= len(data.Points) {
		return []string{}, [][]float64{}
	}

	var labels []string

	labels = append(labels, strings.Trim(label, " ")+".X")
	labels = append(labels, strings.Trim(label, " ")+".Y")
	labels = append(labels, strings.Trim(label, " ")+".Z")
	labels = append(labels, strings.Trim(label, " ")+".A")

	//fmt.Println("Reading trajectory", index)

	points := data.Points[index]

	rdata := make([][]float64, 4, 4)
	for i := 0; i < 4; i++ {
		rdata[i] = make([]float64, len(points), len(points))
	}

	for i, p := range points {
		rdata[0][i] = float64(p.X)
		rdata[1][i] = float64(p.Y)
		rdata[2][i] = float64(p.Z)
	}

	// velocity
	for i := 1; i < len(points); i++ {
		xdist := (rdata[0][i] - rdata[0][i-1])
		ydist := (rdata[1][i] - rdata[1][i-1])
		zdist := (rdata[2][i] - rdata[2][i-1])
		dist := math.Sqrt(xdist*xdist + ydist*ydist + zdist*zdist)

		rdata[3][i] = dist
	}

	// acceleration
	for i := 1; i < len(points); i++ {
		rdata[3][i] = rdata[3][i] - rdata[3][i-1]
	}

	if useJerk == true {
		// third derivative
		for i := 1; i < len(points); i++ {
			rdata[3][i] = rdata[3][i] - rdata[3][i-1]
		}

		// jerk derivative
		for i := 1; i < len(points); i++ {
			rdata[3][i] = rdata[3][i] - rdata[3][i-1]
		}
	}

	return labels, rdata
}

func MinMax(data []float64) (float64, float64) {
	min := data[0]
	max := data[0]

	for _, v := range data {
		if v > max {
			max = v
		}
		if v < min {
			v = min
		}
	}
	return min, max
}

func normaliseDataFrame(df dataframe.DataFrame) dataframe.DataFrame {
	names := df.Names()

	var r dataframe.DataFrame

	fmt.Println("Normalising data")
	bar := pb.StartNew(len(names))
	for i, name := range names {
		s := df.Col(name)
		f := s.Float()
		min, max := MinMax(f)
		if math.Abs(max-min) > 0.0001 {
			for i := range f {
				f[i] = (f[i] - min) / (max - min)
			}
		}
		d := dataframe.New(
			series.New(f, series.Float, name),
		)
		if i == 0 {
			r = d
		} else {
			r = r.CBind(d)
		}
		bar.Increment()
	}
	bar.Finish()
	return r
}

func contains(lst []string, l string) bool {
	for _, s := range lst {
		if strings.Compare(strings.Trim(s, " "), strings.Trim(l, " ")) == 0 {
			return true
		}
	}
	return false
}

func main() {
	inputPtr := flag.String("i", "/Users/zahedi/projects/habil/data/c3d/c3d/38_04.c3d", "input filename")
	outputPtr := flag.String("o", "/Users/zahedi/projects/habil/data/c3d/mcw/38_04.txt", "input filename")
	labelsPtr := flag.String("l", "", "labels")
	minPtr := flag.Int("n", 1000, "minimum number of data points")
	printLabelsPtr := flag.Bool("L", false, "print all the labels")
	exportToCsv := flag.Bool("csv", false, "export to csv")
	useJerk := flag.Bool("j", false, "use jerk instead of curvature")

	flag.Parse()

	if _, err := os.Stat(*outputPtr); err == nil {
  		fmt.Println(fmt.Sprintf("\n%s already exists", *outputPtr))
  		os.Exit(0)
	}


	fmt.Println(fmt.Sprintf("Working on %s", *inputPtr))

	selectedLabels := strings.Split(*labelsPtr, ",")
	fmt.Println(*labelsPtr)
	header, info, data := goc3d.ReadC3D(*inputPtr)
	fmt.Println(header)

	if len(data.Points[0]) < *minPtr {
		fmt.Println(fmt.Sprintf("skipping because not enough data points (%d<%d)", len(data.Points[0]), *minPtr))
		os.Exit(0)
	}
	var labels []string

	prefix := ""

	for _, p := range info.Parameters {
		if p.Name == "LABEL_PREFIXES" {
			prefix = strings.Trim(p.StringData[0], " ")
		}
	}

	for _, p := range info.Parameters {
		if p.Name == "LABELS" {
			for _, s := range p.StringData {
				labels = append(labels, s)
			}
		}
	}

	for _, p := range info.Parameters {
		if p.Name == "LABELS2" {
			for _, s := range p.StringData {
				labels = append(labels, s)
			}
		}
	}

	if *printLabelsPtr == true {
		for _, l := range labels {
			fmt.Println(l)
		}
		os.Exit(0)
	}
	
	for i := range selectedLabels {
		selectedLabels[i] = prefix + selectedLabels[i]
	}
	var df dataframe.DataFrame

	if len(selectedLabels) > 0 { // check if all given labels are found
		var notfound []string
		for _, s := range selectedLabels {
			if contains(labels, s) == false {
				notfound = append(notfound, s)
			}
		}
		if len(notfound) > 0 {
			//fmt.Println("Parameters:", info.Parameters)
			fmt.Println("Labels not found:")
			for _, l := range notfound {
				fmt.Println(" ", l)
			}
			os.Exit(0)
		}
	}

	fmt.Println("Extracting data")
	bar := pb.StartNew(len(selectedLabels))
	var indices []int

	for i, l := range labels {
		if len(selectedLabels) > 0 && contains(selectedLabels, l) == true {
			indices = append(indices, i)
	//		fmt.Println(i, l)
		}
	}

	//fmt.Println(indices)
	
	for _, i := range indices {
		l := labels[i]
	//	fmt.Println("Reading data from", strings.Trim(l, " "), "index", i)
		ls, ds := getData(i, l, data, *useJerk)
		for j := range ls {
			d := dataframe.New(
				series.New(ds[j], series.Float, ls[j]),
			)
			if i == 0 && j == 0 {
				df = d
			} else {
				df = df.CBind(d)
			}
		}
		bar.Increment()
	}
	bar.Finish()

	// Not sure what is the functionality:
	if *exportToCsv == true {
		csvFilename := strings.Replace(*inputPtr, "c3d", "csv", -1)
		csvFilename = strings.Replace(csvFilename, "csv", "c3d", 1) // first one must be c3d not csv 
		fmt.Println("Exporting to", csvFilename)
		file, _ := os.Create(csvFilename)
		defer file.Close()
		df.WriteCSV(file)
		os.Exit(0)
	}

	df = normaliseDataFrame(df)

	m := make([][]float64, df.Nrow(), df.Nrow())
	for r := 0; r < df.Nrow(); r++ {
		m[r] = make([]float64, df.Ncol(), df.Ncol())
	}

	colIndex := 0
	for _, name := range df.Names() {
		isX := strings.HasSuffix(name, ".X") == true
		isY := strings.HasSuffix(name, ".Y") == true
		isZ := strings.HasSuffix(name, ".Z") == true
		if isX || isY || isZ {
			for row := 0; row < df.Nrow(); row++ {
				m[row][colIndex] = df.Elem(row, colIndex).Float()
			}
			colIndex++
		}
	}

	for _, name := range df.Names() {
		isA := strings.HasSuffix(name, ".A") == true
		if isA {
			for row := 0; row < df.Nrow(); row++ {
				m[row][colIndex] = df.Elem(row, colIndex).Float()
			}
			colIndex++
		}
	}

	nrOfLabels := df.Ncol() / 4

	if nrOfLabels == 0 {
		nrOfLabels = len(labels)
	}

	fmt.Println("Number of labels:", nrOfLabels)

	w2w1a1 := make([][]float64, df.Nrow()-1, df.Nrow()-1)
	for i := 0; i < df.Nrow()-1; i++ {
		w2w1a1[i] = make([]float64, nrOfLabels*7, nrOfLabels*7) // x', y', z', x, y, z, a
	}

	// w2
	for row := 0; row < df.Nrow()-1; row++ {
		for col := 0; col < nrOfLabels*3; col++ { // x, y, z
			w2w1a1[row][col] = m[row+1][col]
		}
	}

	// w1, a1
	for row := 0; row < df.Nrow()-1; row++ {
		for col := 0; col < nrOfLabels*4; col++ { // x, y, z, a
			w2w1a1[row][nrOfLabels*3+col] = m[row][col]
		}
	}

	w2Indices := make([]int, nrOfLabels*3, nrOfLabels*3) // x, y, z
	w1Indices := make([]int, nrOfLabels*3, nrOfLabels*3) // x, y, z
	a1Indices := make([]int, nrOfLabels, nrOfLabels)     // a

	index := 0
	for i := 0; i < nrOfLabels*3; i++ {
		w2Indices[i] = index
		index++
	}
	for i := 0; i < nrOfLabels*3; i++ {
		w1Indices[i] = index
		index++
	}
	for i := 0; i < nrOfLabels; i++ {
		a1Indices[i] = index
		index++
	}

	mcw := state.MorphologicalComputationW(w2w1a1, w2Indices, w1Indices, a1Indices, 40, true)

	// fmt.Println(mcw)
	fmt.Println("Result written to", *outputPtr)
	file, _ := os.Create(*outputPtr)
	defer file.Close()
	file.WriteString(fmt.Sprintf("MC_w: %f\n", mcw))
	file.WriteString(fmt.Sprintf("Number of data points: %d\n", len(data.Points[0])))
}
