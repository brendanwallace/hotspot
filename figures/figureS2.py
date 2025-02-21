import seaborn as sns
import pandas as pd
import matplotlib.pyplot as plt

import util
import settings


def plot_S2(D1, D2, D4, D8):


	for data in [D1, D2, D4, D8]:

		outcome_plot = sns.relplot(
		    data=data,
		    row="variable",
		    y="Peak size difference",
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


