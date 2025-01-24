import numpy as np
from scipy import stats
import matplotlib.pyplot as plt
import seaborn as sns
from seaborn import objects as so

import util
import extinction
import settings

#D1_SIMULATION = "simulation,D=1,T=1000.json"
#data = util.load_data(D1_SIMULATION)


def base_plot(data, alpha=1.0):
    rel = sns.relplot(
        data=data[data["R0"] <= settings.R0_END],
        kind="line",
        errorbar=None,
        n_boot=None,
        alpha=alpha,
        y="Outbreak probability",
        x="R0",
        col="Hotspot fraction",
        style="Risk tolerance variance",
        hue="Risk tolerance mean",
        height=settings.HEIGHT,
        #aspect=settings.DEFAULT_WIDTH/settings.DEFAULT_HEIGHT,
    )
    w, h = rel.figure.get_size_inches()
    rel.figure.set_size_inches(settings.FULL_WIDTH, h)
    return rel


def plot_2_A(data):

    #so.Plot().layout(size=(settings.FULL_WIDTH, settings.DEFAULT_HEIGHT))
    rp = base_plot(data, alpha=1.0)

    for ax in rp.axes[0]:
        sns.lineplot(
            x=extinction.X,
            y=extinction.HOMOGENEOUS,
            ax=ax,
            color="black",
            alpha=0.7,
            legend=None,
        )
    #rp.figure.set_size_inches(settings.FULL_WIDTH, settings.HEIGHT)
    # if save_figs:
    #     rp.savefig(settings.IMAGE_LOCATION + "figure2A.pdf", dpi=settings.DPI)

    return rp

#plot_2_A(data)
        



def plot_2_B(data, theoretical_risk_means=[0.125, 0.25, 0.5]):
    #so.Plot().layout(size=(settings.FULL_WIDTH, settings.DEFAULT_HEIGHT))
    rp = base_plot(data, alpha=0.4)

    # # add the homogeneous ABM case
    # for ax in rp.axes[0]:
    #     sns.lineplot(
    #         data=control,
    #         y="Extinction Probability",
    #         x="R0",
    #         color="black",
    #         ax=ax,
    #         legend=None,
    #         alpha=0.4,
    #     )


    for hs, hotspot in enumerate([0.25, 0.5, 0.75]):

        # Adds the theoretical line for homogeneous case
        sns.lineplot(
            x=extinction.X,
            y=extinction.HOMOGENEOUS,
            ax=rp.axes[0][hs],
            color="black",
            alpha=0.7,
            legend=None,
        )
        for rm, risk_mean in enumerate(theoretical_risk_means): 
            Y = []
            for R0 in extinction.X:
                beta_c = R0 / settings.N * (1-hotspot)
                beta_h = R0 / settings.N * hotspot / risk_mean / risk_mean
                Y.append(1.0 - extinction.theoretical_extinction_poisson(beta_c, beta_h, risk_mean))
            sns.lineplot(
                x=extinction.X,
                y=Y,
                ax=rp.axes[0][hs],
                color=(settings.COLORS)[rm],
                legend=None,
            )

    # if save_figs:
    #     rp.savefig(settings.IMAGE_LOCATION + "figure2B.pdf", dpi=settings.DPI)

    return rp


#plot_2_A(data)
#plot_2_B(data)