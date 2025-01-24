import json
import numpy as np
import pandas as pd

import settings

DATA_LOCATION =  '../data/'

COLUMN_NAMES = {
    "RunSets.Parameters.R0": "R0",
    "HotspotFraction": "Hotspot fraction",
    "RiskVariance": "Risk tolerance variance",
    # make sure RiskMean is made categorical first
    # "RiskMean": "Risk tolerance mean",
    "MaxIDiff": "Peak size difference",
    "FinalRDiff": "Extent difference",
    "OutbreakProbability": "Outbreak probability",
    "PeakTimeDiff": "Peak time difference",
    "DurationDiff": "Duration difference",
    "FinalR": "Extent",
    "PeakTime": "Peak time",
    "MaxI": "Peak size",
}


def load_data(filename, drop_control=True, risk_means=None):

    with open(DATA_LOCATION + filename) as file:
            json_file = json.load(file, parse_float=lambda f: round(float(f), 3))
            
    data = process(pd.json_normalize(
        json_file,
        record_path=["RunSets", "Runs"],
        meta=[
            "RiskMean",
            "RiskVariance",
            "HotspotFraction",
            ["RunSets", "Parameters", "R0"],
            ["RunSets", "Parameters", "RunType"]
        ],
    ), drop_control, risk_means)

    return data


# Process data
def process(data, drop_control, risk_means):

    if risk_means is not None:
        data = data[data["RiskMean"].isin(risk_means)]
    
    data["OutbreakProbability"] = 1.0 - 1.0*(data["FinalR"] < settings.EXTINCTION_CUTOFF)
    data["Risk tolerance mean"] = pd.Categorical(data["RiskMean"])
    
    # Uses np.tile to replicate the control series.
    def tile(column):
        num_ppf = len(data["HotspotFraction"].unique())
        return pd.Series(np.tile(data[data["HotspotFraction"] == 0][column], num_ppf))

    # Add a <value>Diff column, which is <value> minus <value> in the homogeneous case.
    # This is used for figure 3.
    for column_to_diff in ["MaxI", "FinalR", "PeakTime", "Duration"]:
        data[column_to_diff + "Control"] = tile(column_to_diff)
        data[column_to_diff + "Diff"] = data[column_to_diff] - data[column_to_diff + "Control"]

    if drop_control:
        data = data[data["HotspotFraction"] != 0]
    
    data = data.rename(columns=COLUMN_NAMES)
    return data


def set_width(plot):
    w, h = plot.figure.get_size_inches()
    plot.figure.set_size_inches(settings.FULL_WIDTH, h)


def make_original_copy(data):
    original = data[data["RiskMean"].isin([0.5, 0.25, 0.125])].copy(deep=True)
    original["Risk tolerance mean"] = pd.Categorical(original["RiskMean"])
    original["Risk tolerance mean"].unique()
    return original