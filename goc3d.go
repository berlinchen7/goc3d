package goc3d

import (
	"bufio"
	"os"

	pb "gopkg.in/cheggaaa/pb.v1"
)

func ReadC3D(filename string, eta bool) {
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		panic(err)
	}

	stats, statsErr := f.Stat()
	if statsErr != nil {
		panic(statsErr)
	}

	var size int = int(stats.Size())
	var bar *pb.ProgressBar

	if eta == true {
		bar = pb.StartNew(size)
	}

	bufr := bufio.NewReader(f)

	for i := 0; i < size; i++ {
		buf, rerr := bufr.ReadByte()
		if rerr != nil {
			panic(rerr)
		}

		// fmt.Println(buf)

		if eta == true {
			bar.Increment()
		}
	}

	if eta == true {
		bar.Finish()
	}
}
