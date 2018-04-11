#!/bin/bash

##Set paremters:
subject=3
trial=3
labels="C7,T10,LFWT,LKNE,LANK,LSHO,LWRA,RFWT,RKNE,RANK,RSHO,RWRA,LFHD,LBHD,RBHD,RFHD,LHEE,RHEE"
# labels="LANK,RANK"
# labels="RANK,RKNE,RTHI"
# labels="RANK,RKNE,RTHI,LANK,LKNE,LTHI"

# NOTE:"LEBL" and "REBL" are not present in 38_4 and 85_5

size_of_moving_average=20
plot_x_axis_range=200

##Start execution
set -e #stop executing in case one of the commands results in an error

#Download specified data:
python download_c3d_mpg.py -s "$subject" -t "$trial"

#Run data with M_CW and generate csv file:
go run cmumocap.go -i /Users/berlin/go/src/github.com/berlin/goc3d/cmumocap/c3d_data/"$subject"_"$trial".c3d \
  -o /Users/berlin/go/src/github.com/berlin/goc3d/cmumocap/csv_data/"$subject"_"$trial"_"$labels".csv \
  -head /Users/berlin/go/src/github.com/berlin/goc3d/cmumocap/csv_data/headerInfo_"$subject"_"$trial"_"$labels".csv \
  -l "$labels" 

#Convert csv file to animation (where most of the time is spent):
python csv_to_mov.py -csvdata /Users/berlin/go/src/github.com/berlin/goc3d/cmumocap/csv_data/"$subject"_"$trial"_"$labels".csv \
  -csvheader /Users/berlin/go/src/github.com/berlin/goc3d/cmumocap/csv_data/headerInfo_"$subject"_"$trial"_"$labels".csv \
  -o /Users/berlin/go/src/github.com/berlin/goc3d/cmumocap/auxillary_mov_data/"$subject"_"$trial".mov \
  -ma $size_of_moving_average \
  -figsize $plot_x_axis_range

#Finally merge animation with the mpg/mov file:
ffmpeg \
  -i mpg_data/"$subject"_"$trial".mpg \
  -i auxillary_mov_data/"$subject"_"$trial".mov \
  -filter_complex '[0:v]pad=iw*2:ih[int];[int][1:v]overlay=W/2:0[vid]' \
  -map [vid] \
  -c:v libx264 \
  -crf 23 \
  -preset fast \
  "$subject"_"$trial"_"$labels".mov

#Generate slow version of a video:
ffmpeg -i "$subject"_"$trial"_"$labels".mov -filter:v "setpts=8.0*PTS" slow_"$subject"_"$trial"_"$labels".mov #1/8.0 times as fast

echo done


# Potentially useful commands:

# To see the specifications of a given video:
# ffprobe -v quiet -print_format json -show_format -show_streams mpg_data/38_04.mpg 

# To convert images to video:
# ffmpeg -r 120 -i png_data/%d.png -pix_fmt yuv420p -vf scale=352:240  test.mpg
