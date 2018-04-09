#!/bin/bash

#Paremters:
subject=38
trial=4
# labels="C7,T10,LFWT,LKNE,LANK,LSHO,LWRA,RFWT,RKNE,RANK,RSHO,RWRA,LFHD,LBHD,RBHD,RFHD,LHEE,RHEE"
# "LEBL" and "REBL" are not present in subject 38, trial 4
labels="LANK,RANK"
size_of_moving_average=20
plot_x_axis_range=200

#Download specified data:
python download_c3d_mpg.py -s "$subject" -t "$trial"

#Run data with M_CW and generate csv file:
go run cmumocap.go -i /Users/berlin/go/src/github.com/berlin/goc3d/cmumocap/c3d_data/"$subject"_"$trial".c3d -o /Users/berlin/go/src/github.com/berlin/goc3d/cmumocap/csv_data/"$subject"_"$trial"_"$labels".csv -l "$labels" 

#Convert csv file to animation (where most of the time is spent):
python csv_to_mov.py -csv /Users/berlin/go/src/github.com/berlin/goc3d/cmumocap/csv_data/"$subject"_"$trial"_"$labels".csv \
  -o /Users/berlin/go/src/github.com/berlin/goc3d/cmumocap/auxillary_mpg_data/"$subject"_"$trial".mov \
  -ma $size_of_moving_average \
  -figsize $plot_x_axis_range

#Finally merge animation with the mpg/mov file:
ffmpeg \
  -i mpg_data/"$subject"_"$trial".mpg \
  -i auxillary_mpg_data/"$subject"_"$trial".mov \
  -filter_complex '[0:v]pad=iw*2:ih[int];[int][1:v]overlay=W/2:0[vid]' \
  -map [vid] \
  -c:v libx264 \
  -crf 23 \
  -preset fast \
  "$labels".mov

echo done


# Potentially useful commands:

# To see the specifications of a given video:
# ffprobe -v quiet -print_format json -show_format -show_streams mpg_data/38_04.mpg 

# To slow down/speed up a video:
ffmpeg -i "$labels".mov -filter:v "setpts=8.0*PTS" slow_"$labels".mov #1/8.0 times as fast

# To convert images to video:
# ffmpeg -r 120 -i png_data/%d.png -pix_fmt yuv420p -vf scale=352:240  test.mpg
