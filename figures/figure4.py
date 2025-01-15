import seaborn as sns
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

import util
import settings

# DIFEQ_DATA = "difeq,D=1,T=1.json"
# data = util.load_data(DIFEQ_DATA, drop_control=False)


PLOT_HEIGHT = 1
PLOT_WIDTH = 1
LINEWIDTH = 1.5


def plot_4(data, save_figs=settings.SAVE_FIGS):


    Study = data[
        (data["Risk tolerance mean"] == 0.25) &
        (data["Risk tolerance variance"] == "high")]


    Study = Study[(Study["R0"] == 1.2) | (Study["R0"] == 2.0) | (Study["R0"] == 3.0)]
    Study = Study[(Study["Hotspot fraction"] == 0.0) | (Study["Hotspot fraction"] == 0.5)]


    Study = Study.explode(["Ts", "Is", "Rs", "Rts"])
    mask = Study["Hotspot fraction"] == 0.5
    Study["Model type"] = pd.Categorical(np.where(mask, "Hotspot", "Homogeneous"))

    melt = pd.melt(
        Study,
        id_vars=['R0', 'Model type', 'Ts'], 
        value_vars=['Rts', 'Is', 'Rs'],
        var_name="variable", value_name="value",
    )

    g = sns.FacetGrid(
        data=melt, col="R0", row="variable", sharex="col", sharey="row",
        #height=11*cm*2/3,
        height=PLOT_HEIGHT,
    ) # model_type hue
    g.map_dataframe(
        sns.lineplot, x="Ts", y="value", hue="Model type",
        estimator=None, n_boot=0, style="Model type",
        dashes=[settings.LINE_DASHES[0], settings.LINE_DASHES[2]],
        palette=[settings.HOMOGENEOUS_COLOR, settings.PURPLE],
    )

    g.axes[0][0].set_ylabel("Rt")
    g.axes[1][0].set_ylabel("Infected")
    g.axes[2][0].set_ylabel("Recovered")

    g.set(xlabel="Days")
    #g.add_legend()
    #sns.move_legend(g, "upper right")

    R0s =["1.2", "2.0", "3.0"]
    for i, ax in enumerate(g.axes[0]):
        ax.set_title("R0 = {}".format(R0s[i]))
    for ax in g.axes[1]:
        ax.set_title("")
    for ax in g.axes[2]:
        ax.set_title("")


    #util.set_width(g)
    _, h = g.figure.get_size_inches()
    g.figure.set_size_inches(settings.COLUMN_WIDTH, h)



    if settings.SAVE_FIGS:
        g.savefig(settings.IMAGE_LOCATION + 'figure4.pdf', format='pdf', dpi=settings.DPI)
    
    return g
