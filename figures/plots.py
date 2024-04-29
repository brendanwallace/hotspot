import numpy as np
import json
from scipy import stats
import matplotlib
import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd

from . import util


cm = 1/2.54
DPI = 600


IMAGE_LOCATION = '/Users/brendan/Documents/projects/hotspot/writeup/images/'
D1_DATA = "simulation,D=1,T=1000.json"
D8_DATA = "simulation,D=8,T=1000.json"
DIFEQ_DATA = "difeq,D=1,T=1.json"



EXTINCTION_CUTOFF = 50
TRIALS = 1000
N = 1000
SAVE_FIGS = False

## Plotting setup
blue_red_purple = ["#466be3","#be89f0", "#f74a64"]
blue_red_purple_corrected = ["#466be3","#bb5ea3", "#f74a64"]
brp3 = ["#67bdff", "#be89f0", "#d32349"]
brp4 = ['#006eff', '#cf78f5', '#ff5759']
brp5 = ["#006eff", "#7640f5", "#f74a75"] # Currently I like this the best!


colors = brp5
sns.set_style("ticks")
sns.set_palette(colors)
sns.color_palette()
sns.set_context("paper", rc={"font.size":8,"axes.titlesize":8,"axes.labelsize":8})


data = util.load_data("simulation,D=1,T=1000.json")
