import btk

# See: http://biomechanical-toolkit.github.io/docs/Wrapping/Python/_getting_started.html
reader = btk.btkAcquisitionFileReader() # build a btk reader object
reader.SetFilename("/Users/berlin/Downloads/38_04.c3d")#/Users/berlin/go/src/github.com/berlin/goc3d/cmumocap/Eb015pi.c3d") # set a filename to the reader
reader.Update()
acq = reader.GetOutput() # acq is the btk aquisition object
#print('Marker labels:')
#for i in range(0, acq.GetPoints().GetItemNumber()):
#    print(acq.GetPoint(i).GetLabel())   
for i, a in enumerate(acq.GetPoint(1).GetValues()):
	if i < 20:
		print a
