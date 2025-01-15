import seaborn as sns
import pandas as pd
import matplotlib.pyplot as plt

import util
import settings

#D8_DATA = "simulation,D=8,T=1000.json"
#data = util.load_data(D8_DATA)

def plot_3(data, save_figs=settings.SAVE_FIGS):


	melt = pd.melt(
		data[data["R0"] <= settings.R0_END],
		id_vars=['R0', 'Hotspot fraction', 'Risk tolerance mean', 'Risk tolerance variance'], 
		value_vars=['Extent difference', 'Peak size difference'],
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

	outcome_plot.figure.set_figwidth(settings.FULL_WIDTH)
	#outcome_plot.figure.tight_layout()


	outcome_plot.axes[0][0].set_ylabel("Extent difference")
	outcome_plot.axes[1][0].set_ylabel("Peak size difference")

	hs_fractions=["0.25", "0.5", "0.75"]
	for i, ax in enumerate(outcome_plot.axes[0]):
	    ax.set_title("Hotspot fraction = {}".format(hs_fractions[i]))
	for ax in outcome_plot.axes[1]:
	    ax.set_title("")


	util.set_width(outcome_plot)
	    
	#outcome_plot.fig
	if save_figs:
		outcome_plot.savefig(settings.IMAGE_LOCATION + 'figure3.pdf', format='pdf', dpi=settings.DPI)
	
	return outcome_plot

