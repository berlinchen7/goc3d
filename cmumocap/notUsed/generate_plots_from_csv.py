import csv
import os
import argparse
import matplotlib
matplotlib.use('TkAgg')
import matplotlib.pyplot as plt
import math

# from multiprocessing import Pool

parser = argparse.ArgumentParser()
parser.add_argument('-csv', dest='csvFilePath', default='csv_data/38_04_Liu.csv', type=str, help='input csv file path')
parser.add_argument('-o', dest='destDirPath', default='png_data', type=str, help='output png directory path')

args = parser.parse_args()

csvFilePath = args.csvFilePath
destDirPath = args.destDirPath

with open(csvFilePath, 'rb') as csvFile:
    r = csv.reader(csvFile, delimiter=' ')
    data = [row for row in r][0][1:]
    data[0], data[-1] = data[0][1:], data[-1][:-1] # getting rid of the '[' and ']'

data = [float(d) for d in data]

print(len(data))
for i in range(len(data)):
    if not bool(i%1000):
        print('First ' + str(i) + ' frames processed.')
    ymin, ymax = min(data), max(data) 
    plt.plot(range(len(data[:i+1])), data[:i+1], 'g')
    plt.ylim(float(math.floor(ymin)), float(math.ceil(ymax)))
    plt.xlim(0, len(data))
    plt.savefig(destDirPath + '/' + str(i) + '.png')

#p = Pool(5)
# def generatePlot(i):
#p.map_async(generatePlot, range(len(data)))
#p.close()  
#p.join()
