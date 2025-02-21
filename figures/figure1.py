import numpy as np
import pandas as pd
import seaborn as sns
from scipy import stats
from matplotlib import pyplot as plt


import settings
import util

FIGURE1_B_HEIGHT = 1
FIGURE1_B_WIDTH = 1
LINEWIDTH = 1.5


def plot_1_B():

    x = np.linspace(0, 1, 100)
    bases = [(1, 7, 0.125), (1, 3, 0.25),  (1, 1, 0.5)]
    factors = [(2, "low"), (1, "medium"), (0.1, "high")]


    data = {
        "variance": [],
        "RiskMean": [],
        "ρ": [],
        "density": [],
    }

    for i, (a, b, base) in enumerate(bases):
        for j, (f, factor) in enumerate(factors):
            


            g = stats.beta(f*a, f*b)
            y = g.pdf(x)
            data["variance"].extend([factor]*len(y))
            data["RiskMean"].extend([base]*len(y))
            data["ρ"].extend(x)
            data["density"].extend(y)

    distributions = pd.DataFrame(data)
    distributions["mean"] = pd.Categorical(distributions["RiskMean"])
            

    rel = sns.relplot(
                data=distributions,
                kind="line",
                linewidth=LINEWIDTH,
                x="ρ", y="density", row="mean",
                height=FIGURE1_B_HEIGHT,
                aspect=FIGURE1_B_WIDTH/FIGURE1_B_HEIGHT,
                col="variance",
                hue="mean",
                style="variance")
    rel.set_titles("")
    rel.set(xlabel="Risk tolerance")
    rel.set(ylabel="Population density")
    rel.tight_layout()

    # don't label the y axis points (these are density and it's confusing)
    rel.axes[0][0].get_yaxis().set_ticklabels([])

    # removes redundent axis label text
    rel.axes[0][0].set_ylabel("")
    rel.axes[2][0].set_ylabel("")
    rel.axes[2][0].set_xlabel("")
    rel.axes[2][2].set_xlabel("")

    return rel

