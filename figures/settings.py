import numpy as np
import json
from scipy import stats
import matplotlib
import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd


CM = 1/2.54
DPI = 600


EXTINCTION_CUTOFF = 50
TRIALS = 1000
N = 1000

R0_END = 8.0

COLUMN_WIDTH = 3.42 # inches
FULL_WIDTH = 7.0 # inches
HEIGHT = 1.8 # inches


## Plotting setup
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

