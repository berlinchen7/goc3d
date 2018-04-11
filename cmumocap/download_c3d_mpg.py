#!/usr/bin/env python

#Parsing command-line arguments
import argparse

parser = argparse.ArgumentParser()
parser.add_argument('-s', dest='subject', default='38', type=str, help='subject number.')
parser.add_argument('-t', dest='trial', default='4', type=str, help='trial number.')
parser.add_argument('-mpg', dest='outputMPGPath', default= 'mpg_data', type=str, help='indicate the path of downloaded mpg file.')
parser.add_argument('-c3d', dest='outputC3DPath', default= 'c3d_data', type=str, help='indicate the path of downloaded c3d file.')

args = parser.parse_args()

subject = args.subject
trial = args.trial

outputMPGFile = args.outputMPGPath + '/' + subject + '_' + trial + '.mpg'
outputC3DFile = args.outputC3DPath + '/' + subject + '_' + trial + '.c3d'

subject = subject.zfill(2) if len(subject) < 3 else subject.zfill(3) #accomodating for the perks of the filename formatting used by the database
trial = trial.zfill(2)


#Downloading c3d and mpg files (some code taken from https://www.codementor.io/aviaryan/downloading-files-from-urls-in-python-77q3bs0un)
import requests
import os.path

def is_downloadable(url):
    """
    Does the url contain a downloadable resource
    """
    h = requests.head(url, allow_redirects=True)
    header = h.headers
    content_type = header.get('content-type')
    if 'text' in content_type.lower():
        return False
    if 'html' in content_type.lower():
        return False
    return True

if os.path.isfile(outputC3DFile):
    print(outputC3DFile + ' already exists')
else:
    print('Downloading .c3d file to ' + outputC3DFile)
    c3d_url = 'http://mocap.cs.cmu.edu/subjects/' + subject + '/' + subject + '_' + trial + '.c3d' 
    assert is_downloadable(c3d_url) == True, "not able to fetch specified c3d file!"
    c3d_r = requests.get(c3d_url, allow_redirects=True)
    open(outputC3DFile, 'wb').write(c3d_r.content)
    print('c3d file downloaded!')

if os.path.isfile(outputMPGFile):
    print(outputMPGFile + ' already exists')
else:
    print('Downloading .mpg file to ' + outputMPGFile)
    mpg_url = 'http://mocap.cs.cmu.edu/subjects/' + subject + '/' + subject + '_' + trial + '.mpg'
    assert is_downloadable(mpg_url) == True, "not able to fetch specified mpg file!\n  If .avi file exists, then manually download it, convert it to .mpg file, and run pipeline.sh again."
    mpg_r = requests.get(mpg_url, allow_redirects=True)
    open(outputMPGFile, 'wb').write(mpg_r.content)
    print('mpg file downloaded!')

