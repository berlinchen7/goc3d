#!/usr/bin/env python

import matplotlib
matplotlib.use('TKAgg')
import numpy as np
from matplotlib import pyplot as plt
from matplotlib import animation

import csv
import os
import argparse
import math
import numpy as np

#parsing command-line arguments
parser = argparse.ArgumentParser()
parser.add_argument('-csvdata', dest='csvDataFilePath', default='csv_data/38_04_Liu.csv', type=str, help='input csv data file path')
parser.add_argument('-csvheader', dest='csvHeaderFilePath', default='csv_data/headerInfo_default.csv', type=str, help='input csv header file path')
parser.add_argument('-ma', dest='movingAverage', default='5', type=int, help='size of moving average applied to smooth the output curve')
parser.add_argument('-o', dest='destFile', default='auxillary_mpg_data/tmp.mov', type=str, help='output animtion file path')
parser.add_argument('-figsize', dest='figSize', default=800, type=int, help='the size of the range of displayed x-axis')

args = parser.parse_args()

csvDataFilePath = args.csvDataFilePath
csvHeaderFilePath = args.csvHeaderFilePath
destFile = args.destFile
movingAverage = args.movingAverage
figSize = args.figSize

#determining the offset of c3d data
with open(csvHeaderFilePath, 'rb') as csvFile:
    r = csv.reader(csvFile, delimiter='\t')
    for i in r:
    	item = i[0].replace(" ","").split("=")
    	if item[0] == 'Firstframe':
    		first_frame = int(item[1])
    		break

#extracting the data from csv file
with open(csvDataFilePath, 'rb') as csvFile:
    r = csv.reader(csvFile, delimiter=' ')
    data = [row for row in r][0][1:]
    data[0], data[-1] = data[0][1:], data[-1][:-1] # getting rid of the '[' and ']'
data = [float(d) for d in data]
ymin, ymax = min(data), max(data)

#applying moving average
smooth_data = np.convolve(data, np.ones((movingAverage,))/movingAverage, mode='same')
data = smooth_data.tolist()

#Prepending zeros to data so the index matches with the actual frame number, which is needed for side-by-side comparison with footage
zeros = [0 for i in range(first_frame-1)]
data = zeros + data

#Setting up the figure, the axis, and the plot element we want to animate
fig = plt.figure(figsize=(3.5, 2.2))#6, 3.5)) # choice of figsize is such that the combined plot look reasonable on my laptop, so may need to be changed to be more robust
ax = plt.axes(ylim=(float(math.floor(ymin)), float(math.ceil(ymax))))
line, = ax.plot([], [], lw=2)

def init():
    line.set_data([], [])
    return line,

#animation function.  This is called sequentially
#NOTE: I reinitialize the plot every time animate() is called just so I could reconfigure the x-axis during the animation; however, this slows down anim.save() significantly.
def animate(i):
	x_lim = (i - figSize + 1, i) if i - figSize >= 0 else (0, figSize-1)
	ax = plt.axes(xlim= x_lim, ylim=(float(math.floor(ymin)), float(math.ceil(ymax)))) #xlim = ((i/figSize)*figSize, ((i/figSize)+1)*figSize),
	line, = ax.plot([], [], lw=2, color='green')
	x = np.array(range(max(i - figSize + 1, 0), i+1)) #np.array(range((i/figSize)*figSize, i))
	y = np.array(data[max(i - figSize + 1, 0): i+1]) # np.array(data[(i/figSize)*figSize:i+1])
	line.set_data(x, y)
	return line,

#call the animator.  blit=True means only re-draw the parts that have changed.
anim = animation.FuncAnimation(fig, animate, init_func=init, frames=len(data), interval=83.3, blit=False)

print("\nSaving plot animation from csv's...")

anim.save(destFile, fps=120, extra_args=['-vcodec', 'libx264'])

#To inspect the animation:
# plt.show()
