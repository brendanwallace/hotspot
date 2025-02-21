import seaborn as sns
import pandas as pd
import matplotlib.pyplot as plt
import numpy as np

import util
import settings


PLOT_HEIGHT = 1
PLOT_WIDTH = 1
LINEWIDTH = 1.5

FIG4_COLORS = {
    0.125: settings.BLUE,
    0.25: settings.PURPLE,
    0.5: settings.RED,
}

FIG4_TEXTURES = {
    "low": settings.LINE_DASHES[0],
    "medium": settings.LINE_DASHES[1],
    "high": settings.LINE_DASHES[2],
}


def plot_4(data,
    risk_tolerance_mean=0.25,
    risk_tolerance_variance="high",
    R0s = [1.25, 2.0, 3.0],

    ):

    color = FIG4_COLORS[risk_tolerance_mean]
    texture = FIG4_TEXTURES[risk_tolerance_variance]

    Study = data[
        (data["Risk tolerance mean"] == risk_tolerance_mean) &
        (data["Risk tolerance variance"] == risk_tolerance_variance)]


    Study = Study[Study["R0"].isin(R0s)]
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
        dashes=[settings.LINE_DASHES[0], texture],
        palette=[settings.HOMOGENEOUS_COLOR, color],
    )

    g.axes[0][0].set_ylabel("Rt")
    g.axes[1][0].set_ylabel("Infected")
    g.axes[2][0].set_ylabel("Recovered")

    g.set(xlabel="Time")

    for i, ax in enumerate(g.axes[0]):
        ax.set_title("R0 = {}".format(R0s[i]))
    for ax in g.axes[1]:
        ax.set_title("")
    for ax in g.axes[2]:
        ax.set_title("")


    #util.set_width(g)
    _, h = g.figure.get_size_inches()
    g.figure.set_size_inches(settings.COLUMN_WIDTH, h)



    return g
