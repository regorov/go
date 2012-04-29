#!/usr/bin/env python

import os
import shutil
import numpy
#import json
from musical.audio import source, save
from musical.theory import Note, Scale

NOTE_LENGTH = 0.3
ECHO_TIMES = 5
ECHO_DECAY = 0.5
ROOT = "F4"
SCALE = "minor"
NUM_NOTES = 18

def render_note(note):
    print "Rendering %r" % note
    note_data = source.sine(note.frequency(), NOTE_LENGTH) * 0.2
    
    for i in xrange(len(note_data)):
        note_data[i] *= 1.0 - (float(i) / len(note_data))
    
    factor = 1.0
    data = note_data
    
    for i in xrange(ECHO_TIMES):
        factor *= ECHO_DECAY
        data = numpy.append(data, note_data * factor)
    
    fname = "out/%s%d.wav" % (note.note, note.octave)
    save.save_wave(data, fname)
    return fname

if os.path.exists("out"):
    shutil.rmtree("out")
os.mkdir("out")

root = Note(ROOT)
scale = Scale(root, SCALE)

notes = [None] * NUM_NOTES

mid = NUM_NOTES / 2
notes[mid] = root

for i in range(mid):
    notes[i] = scale.transpose(root, i - mid)

for i in range(mid + 1, NUM_NOTES):
    notes[i] = scale.transpose(root, i - mid)

files = []

for note in notes:
    files.append(render_note(note))

#f = open("files.json", "w")
#json.dump(files, f)
#f.close()
