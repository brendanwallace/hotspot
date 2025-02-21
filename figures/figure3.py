"""
Adds a row for peak time & for duration to PNAS figure 3.

Adds a column of the homogeneous case.
"""

import seaborn as sns
import pandas as pd
import matplotlib.pyplot as plt

import util
import settings


	
def plot_3(data, timing_values=[], undiffed=False, set_width=True):
	"""
	Timing values can be: {"Duration", or "Peak time"}
	"""

	
	row_values=['Peak size', 'Extent'] + timing_values
	if not undiffed:
		row_values = [r + " difference" for r in row_values]


	# only consider the ones with outbreaks
	_data = data[data["Outbreak probability"] == 1]

	# drop the control, unless this is undiffed
	if not undiffed:
		_data = _data[_data["Hotspot fraction"] > 0]


	melt = pd.melt(
		_data[_data["R0"] <= settings.R0_END],
		id_vars=['R0', 'Hotspot fraction', 'Risk tolerance mean', 'Risk tolerance variance'], 
		value_vars=row_values,
		var_name="variable", value_name="value",
	)

	outcome_plot = sns.relplot(
	    data=melt,
	    row="variable",
	    y="value",
	    x="R0",
	    kind="line",
	    col="Hotspot fraction",
	    errorbar=None,
	    hue="Risk tolerance mean",
	    style="Risk tolerance variance",
	    facet_kws={"sharey":"row"},
	    height=settings.HEIGHT,
	    #aspect=1.0,
	)


	#outcome_plot.figure.tight_layout()

	for i, row_value in enumerate(row_values):
		outcome_plot.axes[i][0].set_ylabel(row_value)


	hs_fractions=["0.25", "0.5", "0.75"]

	if undiffed:
		hs_fractions = ["0"] + hs_fractions

	for i, ax in enumerate(outcome_plot.axes[0]):
	    ax.set_title("Hotspot fraction = {}".format(hs_fractions[i]))
	for i in range(1, len(row_values)):
		for ax in outcome_plot.axes[i]:
		    ax.set_title("")

	if set_width:
		util.set_width(outcome_plot)

	return outcome_plot

