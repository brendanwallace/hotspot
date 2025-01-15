import numpy as np
import json
from scipy import stats
import matplotlib
import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd


CM = 1/2.54
DPI = 600


IMAGE_LOCATION = '/Users/brendan/Documents/projects/hotspot/figures/pnas/'
D1_DATA = "simulation,D=1,T=1000.json"
D8_DATA = "simulation,D=8,T=1000.json"
DIFEQ_DATA = "difeq,D=1,T=1.json"



EXTINCTION_CUTOFF = 50
TRIALS = 1000
N = 1000
SAVE_FIGS = False

R0_END = 8.0

COLUMN_WIDTH = 3.42 # inches
FULL_WIDTH = 7.0 # inches
HEIGHT = 1.8 # inches
#DEFAULT_ASPECT = 1.0 # inches


## Plotting setup
# blue_red_purple = ["#466be3","#be89f0", "#f74a64"]
# blue_red_purple_corrected = ["#466be3","#bb5ea3", "#f74a64"]
# brp3 = ["#67bdff", "#be89f0", "#d32349"]
# brp4 = ['#006eff', '#cf78f5', '#ff5759']
# brp5 = ["#006eff", "#7640f5", "#f74a75"] # Currently I like this the best.


LINE_DASHES = ["", (4, 1.5), (1, 1)]
COLORS = [BLUE, PURPLE, RED] = ["#006eff", "#7640f5", "#f74a75"]
HOMOGENEOUS_COLOR = "black"
sns.set_style("ticks")
sns.set_palette(COLORS)
sns.color_palette()




sns.set_context("paper", rc={
	"font.size":8,
	"axes.titlesize":8,
	"axes.labelsize":8,
	'lines.linewidth': 1.0,
	#'lines.linewidth': 2.25,
})

